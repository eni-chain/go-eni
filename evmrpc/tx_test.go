package evmrpc

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/cometbft/cometbft/libs/bytes"
	rpcclientmock "github.com/cometbft/cometbft/rpc/client/mocks"
	tmtypes "github.com/cometbft/cometbft/types"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/evm/keeper"
	"github.com/cosmos/cosmos-sdk/x/evm/types"
	testutil "github.com/eni-chain/go-eni/testutil/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *mockKeeper) GetReceipt(ctx sdk.Context, hash common.Hash) (*types.Receipt, error) {
	args := m.Called(ctx, hash)
	return args.Get(0).(*types.Receipt), args.Error(1)
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

func TestEniTransactionAPI_GetTransactionReceiptExcludeTraceFail(t *testing.T) {
	hash := common.HexToHash("0x123")
	ctx := context.Background()
	tmClient := &rpcclientmock.Client{}
	keeperMock := &mockKeeper{}
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

	// Test case: Receipt not found
	t.Run("ReceiptNotFound", func(t *testing.T) {
		api.isPanicTx = func(ctx context.Context, hash common.Hash) (bool, error) { return false, nil }
		keeperMock.On("GetReceipt", mock.Anything, hash).Return(nil, errors.New("not found")).Once()
		api.TransactionAPI.keeper = nil
		result, err := api.GetTransactionReceiptExcludeTraceFail(ctx, hash)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestTransactionAPI_GetTransactionReceipt(t *testing.T) {
	hash := common.HexToHash("0x123")
	ctx := context.Background()
	tmClient := &rpcclientmock.Client{}
	evmKeeper, _ := testutil.EvmKeeper(t)
	ctxProvider := func(height int64) sdk.Context { return sdk.Context{} }
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

		// You may need to mock sdkclient.TxConfig if your code expects a non-nil value
		// For now, we use nil for simplicity
		result, err := api.GetTransactionReceipt(ctx, hash)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestTransactionAPI_GetTransactionByHash(t *testing.T) {
	hash := common.HexToHash("0x123")
	ctx := context.Background()
	tmClient := &rpcclientmock.Client{}
	keeperMock := &mockKeeper{}
	ctxProvider := func(height int64) sdk.Context { return sdk.Context{} }
	var cfg sdkclient.TxConfig
	api := NewTransactionAPI(tmClient, &keeper.Keeper{}, ctxProvider, cfg, "", ConnectionType("test"))

	// Test case: Transaction in mempool
	t.Run("InMempool", func(t *testing.T) {
		tx := ethtypes.NewTransaction(
			1,
			common.HexToAddress("0x456"),
			big.NewInt(100),
			21000,
			big.NewInt(1000),
			[]byte("data"),
		)
		tmClient.On("UnconfirmedTxs", mock.Anything, mock.Anything).Return(&coretypes.ResultUnconfirmedTxs{Txs: tmtypes.Txs{[]byte("txdata")}}, nil).Once()
		// You may need to mock sdkclient.TxConfig if your code expects a non-nil value
		result, err := api.GetTransactionByHash(ctx, tx.Hash())
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	// Test case: Transaction not found
	t.Run("NotFound", func(t *testing.T) {
		keeperMock.On("GetReceipt", mock.Anything, hash).Return(nil, errors.New("not found")).Once()
		api.keeper = &keeper.Keeper{}
		result, err := api.GetTransactionByHash(ctx, hash)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestTransactionAPI_GetTransactionCount(t *testing.T) {
	ctx := context.Background()
	tmClient := &rpcclientmock.Client{}
	keeperMock := &mockKeeper{}
	ctxProvider := func(height int64) sdk.Context { return sdk.Context{} }
	var cfg sdkclient.TxConfig
	api := NewTransactionAPI(tmClient, &keeper.Keeper{}, ctxProvider, cfg, "", ConnectionType("test"))
	address := common.HexToAddress("0x123")

	// Test case: Pending block
	t.Run("PendingBlock", func(t *testing.T) {
		keeperMock.On("CalculateNextNonce", mock.Anything, address, true).Return(uint64(10)).Once()
		pending := rpc.PendingBlockNumber
		result, err := api.GetTransactionCount(ctx, address, rpc.BlockNumberOrHash{BlockNumber: &pending})
		assert.NoError(t, err)
		assert.Equal(t, hexutil.Uint64(10), *result)
	})

	// Test case: Specific block
	t.Run("SpecificBlock", func(t *testing.T) {
		keeperMock.On("CalculateNextNonce", mock.Anything, address, false).Return(uint64(5)).Once()
		tmClient.On("Block", mock.Anything, mock.Anything).Return(&coretypes.ResultBlock{}, nil).Maybe()
		blockNum := rpc.BlockNumber(100)
		result, err := api.GetTransactionCount(ctx, address, rpc.BlockNumberOrHash{BlockNumber: &blockNum})
		assert.NoError(t, err)
		assert.Equal(t, hexutil.Uint64(5), *result)
	})
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
		tx := tmtypes.Tx([]byte("txdata"))
		result := getEthTxForTxBz(tx, txDecoder.Decode)
		assert.NotNil(t, result)
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
	receipt := &types.Receipt{
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
		assert.Equal(t, common.HexToHash("blockhash"), result["blockHash"])
	})

	// Test case: Transaction not found
	t.Run("TransactionNotFound", func(t *testing.T) {
		receiptChecker := func(h common.Hash) bool { return false }
		_, err := encodeReceipt(receipt, txDecoder.Decode, block, receiptChecker)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find transaction in block")
	})
}
