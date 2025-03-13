package evmrpc

import (
	"context"
	"errors"
	"github.com/eni-chain/go-eni/x/evm/ante"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	//sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "cosmossdk.io/errors"
	"github.com/eni-chain/go-eni/evmrpc/ethapi"
	"github.com/eni-chain/go-eni/x/evm/keeper"
	"github.com/eni-chain/go-eni/x/evm/types"
	"github.com/eni-chain/go-eni/x/evm/types/ethtx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
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
}

type SendConfig struct {
	slow bool
}

func NewSendAPI(tmClient rpcclient.Client, txConfig client.TxConfig, sendConfig *SendConfig, k *keeper.Keeper, ctxProvider func(int64) sdk.Context, homeDir string, simulateConfig *SimulateConfig, connectionType ConnectionType) *SendAPI {
	return &SendAPI{
		tmClient:       tmClient,
		txConfig:       txConfig,
		sendConfig:     sendConfig,
		keeper:         k,
		ctxProvider:    ctxProvider,
		homeDir:        homeDir,
		backend:        NewBackend(ctxProvider, k, txConfig.TxDecoder(), tmClient, simulateConfig),
		connectionType: connectionType,
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
		return
	}
	msg, err := types.NewMsgEVMTransaction(txData)
	if err != nil {
		return
	}
	ante.Preprocess2(msg)
	txBuilder := s.txConfig.NewTxBuilder()
	if err = txBuilder.SetMsgs(msg); err != nil {
		return
	}
	txbz, encodeErr := s.txConfig.TxEncoder()(txBuilder.GetTx())
	if encodeErr != nil {
		return hash, encodeErr
	}

	//h := bfttypes.Tx(txbz).Hash()
	//hash = common.BytesToHash(h)

	if s.sendConfig.slow {
		res, broadcastError := s.tmClient.BroadcastTxCommit(ctx, txbz)
		if broadcastError != nil {
			err = broadcastError
		} else if res == nil {
			err = errors.New("missing broadcast response")
		} else if res.CheckTx.Code != 0 {
			//err = sdkerrors.ABCIError(sdkerrors.RootCodespace, res.CheckTx.Code, "")
			//todo: need to confirm the codespace
			err = sdkerrors.ABCIError(sdkerrors.UndefinedCodespace, res.CheckTx.Code, "")
		}
	} else {
		res, broadcastError := s.tmClient.BroadcastTxSync(ctx, txbz)
		if broadcastError != nil {
			err = broadcastError
		} else if res == nil {
			err = errors.New("missing broadcast response")
		} else if res.Code != 0 {
			//err = sdkerrors.ABCIError(sdkerrors.RootCodespace, res.Code, "")
			//todo: need to confirm the codespace
			err = sdkerrors.ABCIError(sdkerrors.UndefinedCodespace, res.Code, "")
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

func (s *SendAPI) SendTransaction(ctx context.Context, args ethapi.TransactionArgs) (result common.Hash, returnErr error) {
	startTime := time.Now()
	defer recordMetrics("eth_sendTransaction", s.connectionType, startTime, returnErr == nil)
	//if err := args.SetDefaults(ctx, s.backend); err != nil {
	//	return common.Hash{}, err
	//}
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
	kb, err := getTestKeyring(s.homeDir)
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
