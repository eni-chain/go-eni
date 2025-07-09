package evmrpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"

	"time"

	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	coserrors "github.com/cosmos/cosmos-sdk/types/errors"
	evmante "github.com/cosmos/cosmos-sdk/x/evm/ante"
	"github.com/cosmos/cosmos-sdk/x/evm/keeper"
	"github.com/cosmos/cosmos-sdk/x/evm/types"
	"github.com/cosmos/cosmos-sdk/x/evm/types/ethtx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethapi"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	//rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
)

type SendAPI struct {
	tmClient       rpcclient.Client
	txConfig       client.TxConfig
	sendConfig     *SendConfig
	keeper         *keeper.Keeper
	ctxProvider    func(int64) sdk.Context
	homeDir        string
	backend        *Backend
	connectionType ConnectionType
	logger         log.Logger
}

type SendConfig struct {
	slow bool
}

func NewSendAPI(tmClient rpcclient.Client, txConfig client.TxConfig, sendConfig *SendConfig, k *keeper.Keeper,
	ctxProvider func(int64) sdk.Context, homeDir string, simulateConfig *SimulateConfig,
	connectionType ConnectionType, logger log.Logger) *SendAPI {
	return &SendAPI{
		tmClient:       tmClient,
		txConfig:       txConfig,
		sendConfig:     sendConfig,
		keeper:         k,
		ctxProvider:    ctxProvider,
		homeDir:        homeDir,
		backend:        NewBackend(ctxProvider, k, txConfig.TxDecoder(), tmClient, simulateConfig),
		connectionType: connectionType,
		logger:         logger,
	}
}

func (s *SendAPI) SendRawTransaction(ctx context.Context, input hexutil.Bytes) (hash common.Hash, err error) {
	startTime := time.Now()
	defer recordMetrics("eth_sendRawTransaction", s.connectionType, startTime, err == nil)
	tx := new(ethtypes.Transaction)
	if err = tx.UnmarshalBinary(input); err != nil {
		return
	}
	hash = tx.Hash()
	txData, err := ethtx.NewTxDataFromTx(tx)
	if err != nil {
		s.logger.Error("failed to convert tx to tx data", "err", err)
		return
	}
	msg, err := types.NewMsgEVMTransaction(txData)
	if err != nil {
		s.logger.Error("failed to convert tx to MsgEVMTransaction", "err", err)
		return
	}
	err = evmante.PreprocessMsgSender(msg)
	if err != nil {
		s.logger.Error("failed to convert MsgEVMTransaction to evmante.PreprocessMsgSender", "err", err)
		return
	}
	txBuilder := s.txConfig.NewTxBuilder()
	if err = txBuilder.SetMsgs(msg); err != nil {
		return
	}
	txbz, encodeErr := s.txConfig.TxEncoder()(txBuilder.GetTx())
	if encodeErr != nil {
		return hash, encodeErr
	}

	if s.sendConfig.slow {
		res, broadcastError := s.tmClient.BroadcastTxCommit(ctx, txbz)
		if broadcastError != nil {
			err = broadcastError
		} else if res == nil {
			err = errors.New("missing broadcast response")
		} else if res.CheckTx.Code != 0 {
			err = sdkerrors.ABCIError(coserrors.RootCodespace, res.CheckTx.Code, "")
		}
	} else {
		res, broadcastError := s.tmClient.BroadcastTxSync(ctx, txbz)
		if broadcastError != nil {
			err = broadcastError
		} else if res == nil {
			err = errors.New("missing broadcast response")
		} else if res.Code != 0 {
			err = sdkerrors.ABCIError(coserrors.RootCodespace, res.Code, "")
		}
	}
	return
}

func (s *SendAPI) SignTransaction(_ context.Context, args apitypes.SendTxArgs, _ *string) (result *ethapi.SignTransactionResult, returnErr error) {
	startTime := time.Now()
	defer recordMetrics("eth_signTransaction", s.connectionType, startTime, returnErr == nil)
	var unsignedTx, err = args.ToTransaction()
	if err != nil {
		return nil, err
	}
	signedTx, err := s.signTransaction(unsignedTx, args.From.Address().Hex())
	if err != nil {
		return nil, err
	}
	data, err := signedTx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &ethapi.SignTransactionResult{Raw: data, Tx: signedTx}, nil
}

type TxArgs struct {
	ethapi.TransactionArgs
}

// data retrieves the transaction calldata. Input field is preferred.
func (args *TxArgs) data() []byte {
	if args.Input != nil {
		return *args.Input
	}
	if args.Data != nil {
		return *args.Data
	}
	return nil
}

// SetDefaults fills in default values for unspecified tx fields.
func (args *TxArgs) SetDefaults(ctx context.Context, b *Backend) error {
	if err := args.setFeeDefaults(ctx, b); err != nil {
		return err
	}
	if args.Value == nil {
		args.Value = new(hexutil.Big)
	}
	if args.Nonce == nil {
		nonce, err := b.GetPoolNonce(ctx, *args.From)
		if err != nil {
			return err
		}
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}
	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
		return errors.New(`both "data" and "input" are set and not equal. Please use "input" to pass transaction call data`)
	}
	if args.To == nil && len(args.data()) == 0 {
		return errors.New(`contract creation without any data provided`)
	}
	// Estimate the gas usage if necessary.
	if args.Gas == nil {
		// These fields are immutable during the estimation, safe to
		// pass the pointer directly.
		data := args.data()
		callArgs := TxArgs{
			ethapi.TransactionArgs{
				From:                 args.From,
				To:                   args.To,
				GasPrice:             args.GasPrice,
				MaxFeePerGas:         args.MaxFeePerGas,
				MaxPriorityFeePerGas: args.MaxPriorityFeePerGas,
				Value:                args.Value,
				Data:                 (*hexutil.Bytes)(&data),
				AccessList:           args.AccessList,
			},
		}
		pendingBlockNr := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
		estimated, err := ethapi.DoEstimateGas(ctx, b, callArgs.TransactionArgs, pendingBlockNr, nil, nil, b.RPCGasCap())
		if err != nil {
			return err
		}
		args.Gas = &estimated
		//log.Trace("Estimate gas usage automatically", "gas", args.Gas)
	}
	// If chain id is provided, ensure it matches the local chain id. Otherwise, set the local
	// chain id as the default.
	want := b.ChainConfig().ChainID
	if args.ChainID != nil {
		if have := (*big.Int)(args.ChainID); have.Cmp(want) != 0 {
			return fmt.Errorf("chainId does not match node's (have=%v, want=%v)", have, want)
		}
	} else {
		args.ChainID = (*hexutil.Big)(want)
	}
	return nil
}

// setFeeDefaults fills in default fee values for unspecified tx fields.
func (args *TxArgs) setFeeDefaults(ctx context.Context, b *Backend) error {
	// If both gasPrice and at least one of the EIP-1559 fee parameters are specified, error.
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	}
	// If the tx has completely specified a fee mechanism, no default is needed. This allows users
	// who are not yet synced past London to get defaults for other tx values. See
	// https://github.com/ethereum/go-ethereum/pull/23274 for more information.
	eip1559ParamsSet := args.MaxFeePerGas != nil && args.MaxPriorityFeePerGas != nil
	if (args.GasPrice != nil && !eip1559ParamsSet) || (args.GasPrice == nil && eip1559ParamsSet) {
		// Sanity check the EIP-1559 fee parameters if present.
		if args.GasPrice == nil && args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
			return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.MaxFeePerGas, args.MaxPriorityFeePerGas)
		}
		return nil
	}
	// Now attempt to fill in default value depending on whether London is active or not.
	head := b.CurrentHeader()
	if b.ChainConfig().IsLondon(head.Number) {
		// London is active, set maxPriorityFeePerGas and maxFeePerGas.
		if err := args.setLondonFeeDefaults(ctx, head, b); err != nil {
			return err
		}
	} else {
		if args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil {
			return errors.New("maxFeePerGas and maxPriorityFeePerGas are not valid before London is active")
		}
		// London not active, set gas price.
		price, err := b.SuggestGasTipCap(ctx)
		if err != nil {
			return err
		}
		args.GasPrice = (*hexutil.Big)(price)
	}
	return nil
}

// setLondonFeeDefaults fills in reasonable default fee values for unspecified fields.
func (args *TxArgs) setLondonFeeDefaults(ctx context.Context, head *ethtypes.Header, b *Backend) error {
	// Set maxPriorityFeePerGas if it is missing.
	if args.MaxPriorityFeePerGas == nil {
		tip, err := b.SuggestGasTipCap(ctx)
		if err != nil {
			return err
		}
		args.MaxPriorityFeePerGas = (*hexutil.Big)(tip)
	}
	// Set maxFeePerGas if it is missing.
	if args.MaxFeePerGas == nil {
		// Set the max fee to be 2 times larger than the previous block's base fee.
		// The additional slack allows the tx to not become invalidated if the base
		// fee is rising.
		val := new(big.Int).Add(
			args.MaxPriorityFeePerGas.ToInt(),
			new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
		)
		args.MaxFeePerGas = (*hexutil.Big)(val)
	}
	// Both EIP-1559 fee parameters are now set; sanity check them.
	if args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
		return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.MaxFeePerGas, args.MaxPriorityFeePerGas)
	}
	return nil
}

func (s *SendAPI) SendTransaction(ctx context.Context, arg ethapi.TransactionArgs) (result common.Hash, returnErr error) {
	startTime := time.Now()
	defer recordMetrics("eth_sendTransaction", s.connectionType, startTime, returnErr == nil)
	//if err := args.SetDefaults(ctx, s.backend); err != nil {
	//	return common.Hash{}, err
	//}
	args := TxArgs{arg}
	if err := args.SetDefaults(ctx, s.backend); err != nil {
		return common.Hash{}, err
	}

	var unsignedTx = args.ToTransaction(0)
	signedTx, err := s.signTransaction(unsignedTx, args.From.Hex())
	if err != nil {
		return common.Hash{}, err
	}
	data, err := signedTx.MarshalBinary()
	if err != nil {
		return common.Hash{}, err
	}
	return s.SendRawTransaction(ctx, data)
}

func (s *SendAPI) signTransaction(unsignedTx *ethtypes.Transaction, from string) (*ethtypes.Transaction, error) {
	//kb, err := getTestKeyring(s.homeDir)
	clientCtx := client.Context{}.WithViper("").WithHomeDir(s.homeDir)
	clientCtx, err := config.ReadFromClientConfig(clientCtx)
	if err != nil {
		return nil, err
	}

	cdc, ok := s.keeper.Codec().(*codec.ProtoCodec)
	if ok {
		clientCtx = clientCtx.WithCodec(cdc)
	}

	kb, err := client.NewKeyringFromBackend(clientCtx, keyring.BackendTest)
	if err != nil {
		return nil, err
	}
	privKey, ok := getAddressPrivKeyMap(kb)[from]
	if !ok {
		return nil, errors.New("from address does not have hosted key")
	}
	chainId := s.keeper.ChainID(s.ctxProvider(LatestCtxHeight))
	signer := ethtypes.LatestSignerForChainID(chainId)
	return ethtypes.SignTx(unsignedTx, signer, privKey)
}
