package evmrpc

import (
	"context"
	"errors"
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
	"github.com/eni-chain/go-eni/evmrpc/ethapi"
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
	//write(txbz)

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
	if err != nil {
		s.logger.Error("failed to broadcast tx", "err", err)
		return
	}

	return
}

// var lock sync.Mutex
//
//	func write(tx []byte) {
//		lock.Lock()
//		defer lock.Unlock()
//		// 1. 打开文件（不存在则创建，以追加模式写入）
//		file, err := os.OpenFile("./tx.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
//		if err != nil {
//			panic(err)
//		}
//		defer file.Close() // 确保函数退出时关闭文件
//
//		// 2. 写入一行字符串（自动追加到文件末尾）
//		file.WriteString(hex.EncodeToString(tx))
//		_, err = file.WriteString("\n") // 注意换行符 \n
//		if err != nil {
//			panic(err)
//		}
//	}
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
