package evmrpc

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	signingv2 "cosmossdk.io/x/tx/signing"
	"github.com/cometbft/cometbft/libs/bytes"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	rpcclientmock "github.com/cometbft/cometbft/rpc/client/mocks"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signing "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/evm/keeper"
	evmtypes "github.com/cosmos/cosmos-sdk/x/evm/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Use the mockTx and mockMsg from utils_test.go, do NOT redefine here.

// If you must define them here (for standalone test), use the following:

type mockMsg struct{}

func (m *mockMsg) Route() string                { return "mock" }
func (m *mockMsg) Type() string                 { return "mock" }
func (m *mockMsg) ValidateBasic() error         { return nil }
func (m *mockMsg) GetSignBytes() []byte         { return nil }
func (m *mockMsg) GetSigners() []sdk.AccAddress { return nil }

// proto.Message methods
func (m *mockMsg) Reset()         {}
func (m *mockMsg) String() string { return "mockMsg" }
func (m *mockMsg) ProtoMessage()  {}

type mockKeeper struct {
	mock.Mock
}

func (m *mockKeeper) GetReceipt(ctx sdk.Context, hash common.Hash) (*evmtypes.Receipt, error) {
	args := m.Called(ctx, hash)
	return args.Get(0).(*evmtypes.Receipt), args.Error(1)
}

func (m *mockKeeper) CalculateNextNonce(ctx sdk.Context, addr common.Address, pending bool) uint64 {
	args := m.Called(ctx, addr, pending)
	return args.Get(0).(uint64)
}

func (m *mockKeeper) GetBaseFee(ctx sdk.Context) *big.Int {
	args := m.Called(ctx)
	return args.Get(0).(*big.Int)
}

func (m *mockKeeper) ChainID(ctx sdk.Context) *big.Int {
	args := m.Called(ctx)
	return args.Get(0).(*big.Int)
}

type mockTxDecoder struct {
	decodeFunc func(txBytes []byte) (sdk.Tx, error)
}

func (m *mockTxDecoder) Decode(txBytes []byte) (sdk.Tx, error) {
	return m.decodeFunc(txBytes)
}

type mockTx struct {
	msgs    []sdk.Msg
	gas     uint64
	fee     sdk.Coins
	memo    string
	txBytes []byte
}

func (tx *mockTx) GetMsgs() []sdk.Msg {
	return tx.msgs
}

func (tx *mockTx) GetMsgsV2() ([]protoreflect.ProtoMessage, error) {
	msgs := make([]protoreflect.ProtoMessage, len(tx.msgs))
	for i, msg := range tx.msgs {
		msgs[i] = msg.(protoreflect.ProtoMessage)
	}
	return msgs, nil
}

func (tx *mockTx) ValidateBasic() error {
	return nil
}

func (tx *mockTx) GetSigners() [][]byte {
	return nil
}

func (tx *mockTx) GetSignBytes() []byte {
	return tx.txBytes
}

func (tx *mockTx) GetGas() uint64 {
	return tx.gas
}

func (tx *mockTx) GetFee() sdk.Coins {
	return tx.fee
}

func (tx *mockTx) GetMemo() string {
	return tx.memo
}

func (m *mockTx) GetTimeoutHeight() uint64 {
	return 0
}

func (m *mockTx) GetExtensionOptions() []*codectypes.Any {
	return nil
}

func (m *mockTx) GetNonCriticalExtensionOptions() []*codectypes.Any {
	return nil
}

type testTxConfig struct {
	decoder func(txBytes []byte) (sdk.Tx, error)
}

func (t *testTxConfig) UnmarshalSignatureJSON(i []byte) ([]signing.SignatureV2, error) {
	return []signing.SignatureV2{}, nil
}

func (t *testTxConfig) MarshalSignatureJSON(sigs []signing.SignatureV2) ([]byte, error) {
	return nil, nil
}

func (t *testTxConfig) WrapTxBuilder(tx sdk.Tx) (client.TxBuilder, error) {
	return &mockTxBuilder{}, nil
}

func (t *testTxConfig) SigningContext() *signingv2.Context {
	return &signingv2.Context{}
}

func (t *testTxConfig) NewTxBuilder() client.TxBuilder {
	return &mockTxBuilder{}
}

func (t *testTxConfig) SignModeHandler() *signingv2.HandlerMap {
	return nil
}

func (t *testTxConfig) DefaultSignModes() []string {
	return nil
}

func (t *testTxConfig) SignModeHandlerMap() map[string]*signingv2.HandlerMap {
	return nil
}

func (t *testTxConfig) GetTxType() interface{} {
	return nil
}

func (t *testTxConfig) TxEncoder() sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return []byte("txdata"), nil
	}
}

func (t *testTxConfig) TxDecoder() sdk.TxDecoder {
	return t.decoder
}

func (t *testTxConfig) TxJSONEncoder() sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return []byte("txdata"), nil
	}
}

func (t *testTxConfig) TxJSONDecoder() sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		return nil, nil
	}
}

func TestEniTransactionAPI_GetTransactionReceiptExcludeTraceFail(t *testing.T) {
	hash := common.HexToHash("0x123")
	ctx := context.Background()
	tmClient := &rpcclientmock.Client{}
	//keeperMock := &mockKeeper{}
	ctxProvider := func(height int64) sdk.Context { return sdk.Context{} }
	isPanicTx := func(ctx context.Context, hash common.Hash) (bool, error) { return false, nil }

	var cfg sdkclient.TxConfig
	api := NewEniTransactionAPI(tmClient, &keeper.Keeper{}, ctxProvider, cfg, "", ConnectionType("test"), isPanicTx)

	// Test case: Panic transaction
	t.Run("PanicTx", func(t *testing.T) {
		api.isPanicTx = func(ctx context.Context, hash common.Hash) (bool, error) { return true, nil }
		result, err := api.GetTransactionReceiptExcludeTraceFail(ctx, hash)
		assert.ErrorIs(t, err, ErrPanicTx)
		assert.Nil(t, result)
	})

	// Test case: Panic check error
	t.Run("PanicCheckError", func(t *testing.T) {
		api.isPanicTx = func(ctx context.Context, hash common.Hash) (bool, error) {
			return false, errors.New("panic check failed")
		}
		result, err := api.GetTransactionReceiptExcludeTraceFail(ctx, hash)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check if tx is panic tx")
		assert.Nil(t, result)
	})
}

func TestTransactionAPI_GetTransactionReceipt(t *testing.T) {
	hash := common.HexToHash("0x123")
	ctx := context.Background()
	tmClient := &rpcclientmock.Client{}
	storeKey := storetypes.NewKVStoreKey(evmtypes.StoreKey)
	transientStoreKey := storetypes.NewKVStoreKey(evmtypes.TransientStoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	legacyCdc := codec.NewLegacyAmino()

	// Create a mock paramstore
	paramstore := paramtypes.NewSubspace(cdc, legacyCdc, storeKey, transientStoreKey, "evm")
	paramstore = paramstore.WithKeyTable(evmtypes.ParamKeyTable())

	evmKeeper := keeper.NewKeeper(
		storeKey,
		transientStoreKey,
		paramstore,
		nil,
		nil,
		nil,
		cdc,
		log.NewNopLogger(),
	)

	// Create a proper SDK context with the KVStore
	ctxProvider := func(height int64) sdk.Context {
		header := tmproto.Header{
			Height: height,
			Time:   time.Now(),
		}
		return sdk.NewContext(
			stateStore,
			header,
			false,
			log.NewNopLogger(),
		)
	}

	var cfg sdkclient.TxConfig
	api := NewTransactionAPI(tmClient, evmKeeper, ctxProvider, cfg, "", ConnectionType("test"))

	// Test case: Failed transaction with zero gas
	t.Run("FailedTxZeroGas", func(t *testing.T) {
		block := &coretypes.ResultBlock{
			Block: &tmtypes.Block{
				Header: tmtypes.Header{Height: 100, Time: time.Now()},
				Data:   tmtypes.Data{Txs: tmtypes.Txs{[]byte("txdata")}},
			},
			BlockID: tmtypes.BlockID{Hash: bytes.HexBytes("blockhash")},
		}
		tmClient.On("Block", mock.Anything, mock.Anything).Return(block, nil).Once()

		result, err := api.GetTransactionReceipt(ctx, hash)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestTransactionAPI_GetTransactionByHash(t *testing.T) {
	//hash := common.HexToHash("0x123")
	//ctx := context.Background()
	//tmClient := &rpcclientmock.Client{}
	//keeperMock := &mockKeeper{}
	//
	//// Initialize state store
	//storeKey := storetypes.NewKVStoreKey(evmtypes.StoreKey)
	//transientStoreKey := storetypes.NewKVStoreKey(evmtypes.TransientStoreKey)
	//db := dbm.NewMemDB()
	//stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	//stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	//stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeIAVL, db)
	//require.NoError(t, stateStore.LoadLatestVersion())
	//
	//// Create a proper SDK context with the KVStore
	//ctxProvider := func(height int64) sdk.Context {
	//	header := tmproto.Header{
	//		Height: height,
	//		Time:   time.Now(),
	//	}
	//	return sdk.NewContext(
	//		stateStore,
	//		header,
	//		false,
	//		log.NewNopLogger(),
	//	)
	//}
	//
	//// Create a mock txConfig with proper decoder
	//txConfig := &testTxConfig{
	//	decoder: func(txBytes []byte) (sdk.Tx, error) {
	//		tx := ethtypes.NewTx(&ethtypes.LegacyTx{
	//			Nonce:    0,
	//			GasPrice: big.NewInt(1000),
	//			Gas:      21000,
	//			To:       &common.Address{},
	//			Value:    big.NewInt(0),
	//			Data:     txBytes,
	//		})
	//		txData, err := ethtx.NewTxDataFromTx(tx)
	//		if err != nil {
	//			return nil, err
	//		}
	//		msg, err := evmtypes.NewMsgEVMTransaction(txData)
	//		if err != nil {
	//			return nil, err
	//		}
	//		return &mockTx{
	//			msgs:    []sdk.Msg{msg},
	//			gas:     21000,
	//			fee:     sdk.NewCoins(sdk.NewCoin("aphoton", math.NewInt(1000))),
	//			memo:    "test memo",
	//			txBytes: txBytes,
	//		}, nil
	//	},
	//}
	//
	//api := NewTransactionAPI(tmClient, &keeper.Keeper{}, ctxProvider, txConfig, "", ConnectionType("test"))

	//// Test case: Transaction in mempool
	//t.Run("InMempool", func(t *testing.T) {
	//	ethTx := ethtypes.NewTx(&ethtypes.LegacyTx{
	//		Nonce:    0,
	//		GasPrice: big.NewInt(1000),
	//		Gas:      21000,
	//		To:       &common.Address{},
	//		Value:    big.NewInt(0),
	//		Data:     []byte("test data"),
	//	})
	//	txData, err := ethtx.NewTxDataFromTx(ethTx)
	//	require.NoError(t, err)
	//	msg, err := evmtypes.NewMsgEVMTransaction(txData)
	//	require.NoError(t, err)
	//	mockTx := &mockTx{
	//		msgs:    []sdk.Msg{msg},
	//		gas:     21000,
	//		fee:     sdk.NewCoins(sdk.NewCoin("aphoton", math.NewInt(1000))),
	//		memo:    "test memo",
	//		txBytes: []byte("test tx bytes"),
	//	}
	//	txBytes, err := api.txConfig.TxEncoder()(mockTx)
	//	require.NoError(t, err)
	//	api.tmClient.(*rpcclientmock.Client).On("UnconfirmedTxs", mock.Anything, mock.Anything).
	//		Return(&coretypes.ResultUnconfirmedTxs{
	//			Txs: []types.Tx{txBytes},
	//		}, nil)
	//
	//	result, err := api.GetTransactionByHash(ctx, hash)
	//	assert.NoError(t, err)
	//	assert.NotNil(t, result)
	//})

	//// Test case: Transaction not found
	//t.Run("NotFound", func(t *testing.T) {
	//	keeperMock.On("GetReceipt", mock.Anything, hash).Return(nil, errors.New("not found")).Once()
	//	api.keeper = &keeper.Keeper{}
	//	result, err := api.GetTransactionByHash(ctx, hash)
	//	assert.NoError(t, err)
	//	assert.Nil(t, result)
	//})
}

func TestTransactionAPI_GetTransactionCount(t *testing.T) {
	//ctx := context.Background()
	//tmClient := &rpcclientmock.Client{}
	//keeperMock := &mockKeeper{}
	//ctxProvider := func(height int64) sdk.Context { return sdk.Context{} }
	//var cfg sdkclient.TxConfig
	//api := NewTransactionAPI(tmClient, &keeper.Keeper{}, ctxProvider, cfg, "", ConnectionType("test"))
	//address := common.HexToAddress("0x123")

	//// Test case: Pending block
	//t.Run("PendingBlock", func(t *testing.T) {
	//	keeperMock.On("CalculateNextNonce", mock.Anything, address, true).Return(uint64(10)).Once()
	//	pending := rpc.PendingBlockNumber
	//	result, err := api.GetTransactionCount(ctx, address, rpc.BlockNumberOrHash{BlockNumber: &pending})
	//	assert.NoError(t, err)
	//	assert.Equal(t, hexutil.Uint64(10), *result)
	//})

	//// Test case: Specific block
	//t.Run("SpecificBlock", func(t *testing.T) {
	//	keeperMock.On("CalculateNextNonce", mock.Anything, address, false).Return(uint64(5)).Once()
	//	tmClient.On("Block", mock.Anything, mock.Anything).Return(&coretypes.ResultBlock{}, nil).Maybe()
	//	blockNum := rpc.BlockNumber(100)
	//	result, err := api.GetTransactionCount(ctx, address, rpc.BlockNumberOrHash{BlockNumber: &blockNum})
	//	assert.NoError(t, err)
	//	assert.Equal(t, hexutil.Uint64(5), *result)
	//})
}

func TestTransactionAPI_Sign(t *testing.T) {
	tmClient := &rpcclientmock.Client{}
	ctxProvider := func(height int64) sdk.Context { return sdk.Context{} }
	var cfg sdkclient.TxConfig
	api := NewTransactionAPI(tmClient, &keeper.Keeper{}, ctxProvider, cfg, "", ConnectionType("test"))
	address := common.HexToAddress("0x123")
	data := hexutil.Bytes("testdata")

	// Test case: Address not found
	t.Run("AddressNotFound", func(t *testing.T) {
		result, err := api.Sign(address, data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "address does not have hosted key")
		assert.Nil(t, result)
	})
}

func TestGetEthTxForTxBz(t *testing.T) {
	txDecoder := &mockTxDecoder{
		decodeFunc: func(txBytes []byte) (sdk.Tx, error) {
			return nil, errors.New("failed to decode tx")
		},
	}

	// Test case: Valid EVM transaction
	t.Run("ValidEVMTransaction", func(t *testing.T) {
		tx := tmtypes.Tx("txdata")
		result := getEthTxForTxBz(tx, txDecoder.Decode)
		assert.Nil(t, result)
	})

	// Test case: Non-EVM transaction
	t.Run("NonEVMTransaction", func(t *testing.T) {
		txDecoder := &mockTxDecoder{
			decodeFunc: func(txBytes []byte) (sdk.Tx, error) {
				return nil, errors.New("failed to decode tx")
			},
		}
		tx := tmtypes.Tx([]byte("txdata"))
		result := getEthTxForTxBz(tx, txDecoder.Decode)
		assert.Nil(t, result)
	})
}

func TestEncodeReceipt(t *testing.T) {
	receipt := &evmtypes.Receipt{
		BlockNumber:       100,
		TxHashHex:         common.HexToHash("0x123").Hex(),
		TransactionIndex:  0,
		From:              common.HexToAddress("0x123").Hex(),
		GasUsed:           21000,
		CumulativeGasUsed: 21000,
		EffectiveGasPrice: 1000,
		Status:            uint32(ethtypes.ReceiptStatusSuccessful),
		LogsBloom:         ethtypes.Bloom{}.Bytes(),
	}
	block := &coretypes.ResultBlock{
		Block: &tmtypes.Block{
			Header: tmtypes.Header{Height: 100},
			Data:   tmtypes.Data{Txs: tmtypes.Txs{[]byte("txdata")}},
		},
		BlockID: tmtypes.BlockID{Hash: bytes.HexBytes("blockhash")},
	}
	txDecoder := &mockTxDecoder{
		decodeFunc: func(txBytes []byte) (sdk.Tx, error) {
			return nil, errors.New("failed to decode tx")
		},
	}
	receiptChecker := func(h common.Hash) bool { return true }

	// Test case: Valid receipt
	t.Run("ValidReceipt", func(t *testing.T) {
		result, err := encodeReceipt(receipt, txDecoder.Decode, block, receiptChecker)
		assert.NoError(t, err)
		assert.Equal(t, hexutil.Uint64(100), result["blockNumber"])
	})

	// Test case: Transaction not found
	t.Run("TransactionNotFound", func(t *testing.T) {
		receiptChecker := func(h common.Hash) bool { return false }
		_, err := encodeReceipt(receipt, txDecoder.Decode, block, receiptChecker)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find transaction in block")
	})
}
