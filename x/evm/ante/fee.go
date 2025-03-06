package ante

import (
	sdkerrors "cosmossdk.io/errors"
	coserrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/eni-chain/go-eni/app/antedecorators"
	"github.com/eni-chain/go-eni/utils"
	"github.com/eni-chain/go-eni/utils/metrics"
	"github.com/eni-chain/go-eni/x/evm/derived"
	"github.com/eni-chain/go-eni/x/evm/state"
	evmtypes "github.com/eni-chain/go-eni/x/evm/types"
	"github.com/eni-chain/go-eni/x/evm/types/ethtx"
	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"math/big"

	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmkeeper "github.com/eni-chain/go-eni/x/evm/keeper"
)

type EVMFeeCheckDecorator struct {
	evmKeeper *evmkeeper.Keeper
}

func NewEVMFeeCheckDecorator(evmKeeper *evmkeeper.Keeper) *EVMFeeCheckDecorator {
	return &EVMFeeCheckDecorator{
		evmKeeper: evmKeeper,
	}
}

func (fc EVMFeeCheckDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if simulate {
		return next(ctx, tx, simulate)
	}
	return next(ctx, tx, simulate)
	//todo gas related content, ignore for now
	msg := evmtypes.MustGetEVMTransactionMessage(tx)
	txData, err := evmtypes.UnpackTxData(msg.Data)
	if err != nil {
		return ctx, err
	}

	ver := msg.Derived.Version

	if txData.GetGasFeeCap().Cmp(fc.getBaseFee(ctx)) < 0 {
		return ctx, coserrors.ErrInsufficientFee
	}
	if txData.GetGasFeeCap().Cmp(fc.getMinimumFee(ctx)) < 0 {
		return ctx, coserrors.ErrInsufficientFee
	}
	if txData.GetGasTipCap().Sign() < 0 {
		return ctx, sdkerrors.Wrapf(coserrors.ErrInvalidRequest, "gas fee cap cannot be negative")
	}

	// if EVM version is Cancun or later, and the transaction contains at least one blob, we need to
	// make sure the transaction carries a non-zero blob fee cap.
	if ver >= derived.Cancun && len(txData.GetBlobHashes()) > 0 {
		// For now we are simply assuming excessive blob gas is 0. In the future we might change it to be
		// dynamic based on prior block usage.
		zero := uint64(0)
		if txData.GetBlobFeeCap().Cmp(eip4844.CalcBlobFee(&params.ChainConfig{CancunTime: &zero}, &ethtypes.Header{})) < 0 {
			return ctx, coserrors.ErrInsufficientFee
		}
	}

	// check if the sender has enough balance to cover fees
	etx, _ := msg.AsTransaction()
	emsg := fc.evmKeeper.GetEVMMessage(ctx, etx, msg.Derived.SenderEVMAddr)
	stateDB := state.NewDBImpl(ctx, fc.evmKeeper, false)
	gp := fc.evmKeeper.GetGasPool()
	blockCtx, err := fc.evmKeeper.GetVMBlockContext(ctx, gp)
	if err != nil {
		return ctx, err
	}
	cfg := evmtypes.DefaultChainConfig().EthereumConfig(fc.evmKeeper.ChainID(ctx))
	txCtx := core.NewEVMTxContext(emsg)
	evmInstance := vm.NewEVM(*blockCtx, stateDB, cfg, vm.Config{})
	evmInstance.SetTxContext(txCtx)
	//st, err := core.ApplyMessage(evmInstance, emsg, &gp)
	st := core.NewStateTransition(evmInstance, emsg, &gp, true)
	// run stateless checks before charging gas (mimicking Geth behavior)
	if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
		// we don't want to run nonce check here for CheckTx because we have special
		// logic for pending nonce during CheckTx in sig.go
		if err := st.StatelessChecks(); err != nil {
			return ctx, sdkerrors.Wrap(coserrors.ErrWrongSequence, err.Error())
		}
	}
	if err := st.BuyGas(); err != nil {
		return ctx, sdkerrors.Wrap(coserrors.ErrInsufficientFunds, err.Error())
	}
	if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
		surplus, err := stateDB.Finalize()
		if err != nil {
			return ctx, err
		}
		if err := fc.evmKeeper.AddAnteSurplus(ctx, etx.Hash(), surplus); err != nil {
			return ctx, err
		}
	}

	// calculate the priority by dividing the total fee with the native gas limit (i.e. the effective native gas price)
	priority := fc.CalculatePriority(ctx, txData)
	ctx = ctx.WithPriority(priority.Int64())

	return next(ctx, tx, simulate)
}

// minimum fee per gas required for a tx to be processed
func (fc EVMFeeCheckDecorator) getBaseFee(ctx sdk.Context) *big.Int {
	return fc.evmKeeper.GetCurrBaseFeePerGas(ctx).TruncateInt().BigInt()
}

// lowest allowed fee per gas, base fee will not be lower than this
func (fc EVMFeeCheckDecorator) getMinimumFee(ctx sdk.Context) *big.Int {
	return fc.evmKeeper.GetMinimumFeePerGas(ctx).TruncateInt().BigInt()
}

// CalculatePriority returns a priority based on the effective gas price of the transaction
func (fc EVMFeeCheckDecorator) CalculatePriority(ctx sdk.Context, txData ethtx.TxData) *big.Int {
	gp := txData.EffectiveGasPrice(utils.Big0)
	if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
		metrics.HistogramEvmEffectiveGasPrice(gp)
	}
	priority := math.LegacyNewDecFromBigInt(gp).Quo(fc.evmKeeper.GetPriorityNormalizer(ctx)).TruncateInt().BigInt()
	if priority.Cmp(big.NewInt(antedecorators.MaxPriority)) > 0 {
		priority = big.NewInt(antedecorators.MaxPriority)
	}
	return priority
	return nil
}
