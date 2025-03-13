package keeper

import (
	"context"

	"errors"
	"fmt"
	"math"
	"math/big"
	"runtime/debug"

	cosmossdk_io_math "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/eni-chain/go-eni/x/evm/state"
	"github.com/eni-chain/go-eni/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/hashicorp/go-metrics"
)

var (
	ErrReadEstimate = errors.New("multiversion store value contains estimate, cannot read, aborting")
)

type msgServer struct {
	*Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (server msgServer) EVMTransaction(goCtx context.Context, msg *types.MsgEVMTransaction) (serverRes *types.MsgEVMTransactionResponse, err error) {
	if msg.IsAssociateTx() {
		// no-op in msg server for associate tx; all the work have been done in ante handler
		return &types.MsgEVMTransactionResponse{}, nil
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	tx, _ := msg.AsTransaction()

	// EVM has a special case here, mainly because for an EVM transaction the gas limit is set on EVM payload level, not on top-level GasWanted field
	// as normal transactions (because existing eth client can't). As a result EVM has its own dedicated ante handler chain. The full sequence is:

	// 	1. At the beginning of the ante handler chain, gas meter is set to infinite so that the ante processing itself won't run out of gas (EVM ante is pretty light but it does read a parameter or two)
	// 	2. At the end of the ante handler chain, gas meter is set based on the gas limit specified in the EVM payload; this is only to provide a GasWanted return value to tendermint mempool when CheckTx returns, and not used for anything else.
	// 	3. At the beginning of message server (here), gas meter is set to infinite again, because EVM internal logic will then take over and manage out-of-gas scenarios.
	// 	4. At the end of message server, gas consumed by EVM is adjusted to Eni's unit and counted in the original gas meter, because that original gas meter will be used to count towards block gas after message server returns
	originalGasMeter := ctx.GasMeter()
	ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())

	stateDB := state.NewDBImpl(ctx, &server, false)
	emsg := server.GetEVMMessage(ctx, tx, msg.Derived.SenderEVMAddr)
	gp := server.GetGasPool()

	defer func() {
		defer stateDB.Cleanup()
		if pe := recover(); pe != nil {
			//if !strings.Contains(fmt.Sprintf("%s", pe), occtypes.ErrReadEstimate.Error()) {
			debug.PrintStack()
			ctx.Logger().Error(fmt.Sprintf("EVM PANIC: %s", pe))
			telemetry.IncrCounter(1, types.ModuleName, "panics")
			//}
			panic(pe)
		}
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Got EVM state transition error (not VM error): %s", err))

			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, "errors", "state_transition"},
				1,
				[]metrics.Label{
					telemetry.NewLabel("type", err.Error()),
				},
			)
			return
		}
		extraSurplus := cosmossdk_io_math.ZeroInt()
		surplus, ferr := stateDB.Finalize()
		if ferr != nil {
			err = ferr
			ctx.Logger().Error(fmt.Sprintf("failed to finalize EVM stateDB: %s", err))

			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, "errors", "stateDB_finalize"},
				1,
				[]metrics.Label{
					telemetry.NewLabel("type", err.Error()),
				},
			)
			return
		}

		//
		receipt, rerr := server.WriteReceipt(ctx, stateDB, emsg, uint32(tx.Type()), tx.Hash(), serverRes.GasUsed, serverRes.VmError)
		if rerr != nil {
			err = rerr
			ctx.Logger().Error(fmt.Sprintf("failed to write EVM receipt: %s", err))

			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, "errors", "write_receipt"},
				1,
				[]metrics.Label{
					telemetry.NewLabel("type", err.Error()),
				},
			)
			return
		}

		// Add metrics for receipt status
		if receipt.Status == uint32(ethtypes.ReceiptStatusFailed) {
			telemetry.IncrCounter(1, "receipt", "status", "failed")
		} else {
			telemetry.IncrCounter(1, "receipt", "status", "success")
		}

		surplus = surplus.Add(extraSurplus)
		bloom := ethtypes.Bloom{}
		bloom.SetBytes(receipt.LogsBloom)
		server.AppendToEvmTxDeferredInfo(ctx, bloom, tx.Hash(), surplus)

		// GasUsed in serverRes is in EVM's gas unit, not Eni's gas unit.
		// PriorityNormalizer is the coefficient that's used to adjust EVM
		// transactions' priority, which is based on gas limit in EVM unit,
		// to Eni transactions' priority, which is based on gas limit in
		// Eni unit, so we use the same coefficient to convert gas unit here.
		adjustedGasUsed := server.GetPriorityNormalizer(ctx).MulInt64(int64(serverRes.GasUsed))
		originalGasMeter.ConsumeGas(adjustedGasUsed.TruncateInt().Uint64(), "evm transaction")
	}()

	res, applyErr := server.applyEVMMessage(ctx, emsg, stateDB, gp)
	serverRes = &types.MsgEVMTransactionResponse{
		Hash: tx.Hash().Hex(),
	}
	if applyErr != nil {
		// This should not happen, as anything that could cause applyErr is supposed to
		// be checked in CheckTx first
		err = applyErr

		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "errors", "apply_message"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("type", err.Error()),
			},
		)

		return
	}

	// if applyErr is nil then res must be non-nil
	if res.Err != nil {
		serverRes.VmError = res.Err.Error()

		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "errors", "vm_execution"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("type", serverRes.VmError),
			},
		)
	}

	serverRes.GasUsed = res.UsedGas
	serverRes.ReturnData = res.ReturnData
	serverRes.Logs = types.NewLogsFromEth(stateDB.GetAllLogs())

	return
}

func (k *Keeper) GetGasPool() core.GasPool {
	return math.MaxUint64
}

func (k *Keeper) GetEVMMessage(ctx sdk.Context, tx *ethtypes.Transaction, sender common.Address) *core.Message {
	msg := &core.Message{
		Nonce:            tx.Nonce(),
		GasLimit:         tx.Gas(),
		GasPrice:         new(big.Int).Set(tx.GasPrice()),
		GasFeeCap:        new(big.Int).Set(tx.GasFeeCap()),
		GasTipCap:        new(big.Int).Set(tx.GasTipCap()),
		To:               tx.To(),
		Value:            tx.Value(),
		Data:             tx.Data(),
		AccessList:       tx.AccessList(),
		SkipNonceChecks:  false,
		SkipFromEOACheck: false,
		BlobHashes:       tx.BlobHashes(),
		BlobGasFeeCap:    tx.BlobGasFeeCap(),
		From:             sender,
	}
	// If baseFee provided, set gasPrice to effectiveGasPrice.
	baseFee := k.GetBaseFee(ctx)
	if baseFee != nil {
		msg.GasPrice = BigMin(msg.GasPrice.Add(msg.GasTipCap, baseFee), msg.GasFeeCap)
	}
	return msg
}

// BigMin returns the smaller of x or y.
func BigMin(x, y *big.Int) *big.Int {
	if x.Cmp(y) > 0 {
		return y
	}
	return x
}
func (k Keeper) applyEVMMessage(ctx sdk.Context, msg *core.Message, stateDB *state.DBImpl, gp core.GasPool) (*core.ExecutionResult, error) {
	blockCtx, err := k.GetVMBlockContext(ctx, gp)
	if err != nil {
		return nil, err
	}
	cfg := types.DefaultChainConfig().EthereumConfig(k.ChainID(ctx))
	txCtx := core.NewEVMTxContext(msg)
	evmInstance := vm.NewEVM(*blockCtx, stateDB, cfg, vm.Config{})
	evmInstance.SetTxContext(txCtx)
	st := core.NewStateTransition(evmInstance, msg, &gp, true) // fee already charged in ante handler
	return st.Execute()
}

func (server msgServer) Send(goCtx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	recipient := server.GetEniAddressOrDefault(ctx, common.HexToAddress(msg.ToAddress))
	_, err := bankkeeper.NewMsgServerImpl(server.BankKeeper()).Send(goCtx, &banktypes.MsgSend{
		FromAddress: msg.FromAddress,
		ToAddress:   recipient.String(),
		Amount:      msg.Amount,
	})
	if err != nil {
		return nil, err
	}
	return &types.MsgSendResponse{}, nil
}

func (server msgServer) RegisterPointer(goCtx context.Context, msg *types.MsgRegisterPointer) (*types.MsgRegisterPointerResponse, error) {
	panic("unknown pointer type")
}

func (server msgServer) AssociateContractAddress(goCtx context.Context, msg *types.MsgAssociateContractAddress) (*types.MsgAssociateContractAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr := sdk.MustAccAddressFromBech32(msg.Address) // already validated
	// check if address is for a contract

	evmAddr := common.BytesToAddress(addr)
	existingEvmAddr, ok := server.GetEVMAddress(ctx, addr)
	if ok {
		if existingEvmAddr.Cmp(evmAddr) != 0 {
			ctx.Logger().Error(fmt.Sprintf("unexpected associated EVM address %s exists for contract %s: expecting %s", existingEvmAddr.Hex(), addr.String(), evmAddr.Hex()))
		}
		return nil, errors.New("contract already has an associated address")
	}
	server.SetAddressMapping(ctx, addr, evmAddr)
	return &types.MsgAssociateContractAddressResponse{}, nil
}

func (server msgServer) Associate(context.Context, *types.MsgAssociate) (*types.MsgAssociateResponse, error) {
	return &types.MsgAssociateResponse{}, nil
}
