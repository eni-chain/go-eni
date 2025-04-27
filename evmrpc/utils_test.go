package evmrpc

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	signingv2 "cosmossdk.io/x/tx/signing"
	"github.com/cometbft/cometbft/libs/bytes"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	signing "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

// type alias for missing types if needed
type (
	TxDecoder       = func([]byte) (interface{}, error)
	TxEncoder       = func(interface{}) ([]byte, error)
	SignModeHandler = interface{}
)

// --- Mock types for testing ---
type mockTMClient struct{ rpcclient.Client }

func (m *mockTMClient) Genesis(ctx context.Context) (*coretypes.ResultGenesis, error) {
	return &coretypes.ResultGenesis{Genesis: &tmtypes.GenesisDoc{InitialHeight: 1}}, nil
}
func (m *mockTMClient) Block(ctx context.Context, height *int64) (*coretypes.ResultBlock, error) {
	return &coretypes.ResultBlock{Block: &tmtypes.Block{Header: tmtypes.Header{Height: 10}}}, nil
}
func (m *mockTMClient) BlockByHash(ctx context.Context, hash []byte) (*coretypes.ResultBlock, error) {
	return &coretypes.ResultBlock{Block: &tmtypes.Block{Header: tmtypes.Header{Height: 20}}}, nil
}
func (m *mockTMClient) BlockResults(ctx context.Context, height *int64) (*coretypes.ResultBlockResults, error) {
	return &coretypes.ResultBlockResults{}, nil
}

type emptyKeyring struct{ keyring.Keyring }

func (e *emptyKeyring) List() ([]*keyring.Record, error) { return nil, errors.New("no keys") }

// --- mockTxConfig implements all methods of client.TxConfig ---
type mockTxConfig struct{}

func (m *mockTxConfig) TxEncoder() sdk.TxEncoder {
	//TODO implement me
	panic("implement me")
}

func (m *mockTxConfig) TxDecoder() sdk.TxDecoder {
	//TODO implement me
	panic("implement me")
}

func (m *mockTxConfig) TxJSONEncoder() sdk.TxEncoder {
	//TODO implement me
	panic("implement me")
}

func (m *mockTxConfig) TxJSONDecoder() sdk.TxDecoder {
	//TODO implement me
	panic("implement me")
}

func (m *mockTxConfig) UnmarshalSignatureJSON(i []byte) ([]signing.SignatureV2, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockTxConfig) WrapTxBuilder(tx sdk.Tx) (client.TxBuilder, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockTxConfig) SigningContext() *signingv2.Context {
	//TODO implement me
	panic("implement me")
}

func (m *mockTxConfig) MarshalSignatureJSON(sigs []signing.SignatureV2) ([]byte, error) {
	return nil, nil
}
func (m *mockTxConfig) NewTxBuilder() client.TxBuilder {
	return &mockTxBuilder{}
}

func (m *mockTxConfig) SignModeHandler() *signingv2.HandlerMap               { return nil }
func (m *mockTxConfig) DefaultSignModes() []string                           { return nil }
func (m *mockTxConfig) SignModeHandlerMap() map[string]*signingv2.HandlerMap { return nil }
func (m *mockTxConfig) GetTxType() interface{}                               { return nil }

type mockTx struct{}

func (m *mockTx) GetMsgs() []interface{} { return nil }

// mockTxBuilder implements client.TxBuilder (empty for test)
type mockTxBuilder struct{}

func (m *mockTxBuilder) SetMsgs(msgs ...sdk.Msg) error {

	//TODO implement me
	panic("implement me")
}

func (m *mockTxBuilder) SetSignatures(signatures ...signing.SignatureV2) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockTxBuilder) SetFeePayer(feePayer sdk.AccAddress) {
	//TODO implement me
	panic("implement me")
}

func (m *mockTxBuilder) SetFeeGranter(feeGranter sdk.AccAddress) {
	//TODO implement me
	panic("implement me")
}

func (m *mockTxBuilder) SetMemo(string)                                    {}
func (m *mockTxBuilder) SetFeeAmount(coins sdk.Coins)                      {}
func (m *mockTxBuilder) SetGasLimit(uint64)                                {}
func (m *mockTxBuilder) SetTimeoutHeight(uint64)                           {}
func (m *mockTxBuilder) GetTx() authsigning.Tx                             { return nil }
func (m *mockTxBuilder) AddAuxSignerData(data txtypes.AuxSignerData) error { return nil }

func TestCheckVersion(t *testing.T) {
	err := CheckVersion(nil, nil)
	assert.NoError(t, err)
}

func TestGetBlockNumberByNrOrHash(t *testing.T) {
	ctx := context.Background()
	client := &mockTMClient{}
	// Test by block number
	bn := rpc.BlockNumber(5)
	bnoh := rpc.BlockNumberOrHash{BlockNumber: &bn}
	got, err := GetBlockNumberByNrOrHash(ctx, client, bnoh)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), *got)
	// Test by block hash
	hash := common.Hash{}
	bnoh = rpc.BlockNumberOrHash{BlockHash: &hash}
	got, err = GetBlockNumberByNrOrHash(ctx, client, bnoh)
	assert.NoError(t, err)
	assert.Equal(t, int64(20), *got)
}

func Test_bankExists(t *testing.T) {
	assert.True(t, bankExists(nil, nil))
}

func Test_blockByHash(t *testing.T) {
	client := &mockTMClient{}
	got, err := blockByHash(context.Background(), client, bytes.HexBytes("abc"))
	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func Test_blockByHashWithRetry(t *testing.T) {
	client := &mockTMClient{}
	got, err := blockByHashWithRetry(context.Background(), client, bytes.HexBytes("abc"), 1)
	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func Test_blockByNumber(t *testing.T) {
	client := &mockTMClient{}
	got, err := blockByNumber(context.Background(), client, nil)
	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func Test_blockByNumberWithRetry(t *testing.T) {
	client := &mockTMClient{}
	got, err := blockByNumberWithRetry(context.Background(), client, nil, 1)
	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func Test_blockResultsWithRetry(t *testing.T) {
	client := &mockTMClient{}
	got, err := blockResultsWithRetry(context.Background(), client, nil)
	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func Test_evmExists(t *testing.T) {
	assert.True(t, evmExists(nil, nil))
}

func Test_extractPrivKeyFromLocal(t *testing.T) {
	rl := &keyring.Record_Local{PrivKey: nil}
	got, err := extractPrivKeyFromLocal(rl)
	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrPrivKeyNotAvailable)
	// Only test nil branch, type assertion failure branch requires more complex mock
}

func Test_extractPrivKeyFromRecord(t *testing.T) {
	rec := &keyring.Record{}
	got, err := extractPrivKeyFromRecord(rec)
	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrPrivKeyExtr)
}

func Test_getAddressPrivKeyMap(t *testing.T) {
	kb := &emptyKeyring{}
	got := getAddressPrivKeyMap(kb)
	assert.NotNil(t, got)
	assert.Len(t, got, 0)
}

func Test_getBlockNumber(t *testing.T) {
	ctx := context.Background()
	client := &mockTMClient{}
	// Earliest block
	got, err := getBlockNumber(ctx, client, rpc.EarliestBlockNumber)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), *got)
	// Latest block
	got, err = getBlockNumber(ctx, client, rpc.LatestBlockNumber)
	assert.NoError(t, err)
	assert.Nil(t, got)
	// Custom block number
	got, err = getBlockNumber(ctx, client, rpc.BlockNumber(7))
	assert.NoError(t, err)
	assert.Equal(t, int64(7), *got)
}

func Test_getHeightFromBigIntBlockNumber(t *testing.T) {
	latest := int64(100)
	assert.Equal(t, latest, getHeightFromBigIntBlockNumber(latest, big.NewInt(rpc.LatestBlockNumber.Int64())))
	assert.Equal(t, int64(88), getHeightFromBigIntBlockNumber(latest, big.NewInt(88)))
}

func Test_getTestKeyring(t *testing.T) {
	kr, err := getTestKeyring("/tmp/nonexist")
	assert.NotNil(t, kr)
	assert.NoError(t, err)
}

func Test_getTxHashesFromBlock(t *testing.T) {
	block := &coretypes.ResultBlock{Block: &tmtypes.Block{Header: tmtypes.Header{Height: 1}}}
	txConfig := &mockTxConfig{}
	got := getTxHashesFromBlock(block, txConfig, false)
	assert.NotNil(t, got)
}

func Test_recordMetrics(t *testing.T) {
	recordMetrics("api", ConnectionType("ws"), time.Now(), true)
}

func Test_recordMetricsWithError(t *testing.T) {
	recordMetricsWithError("api", ConnectionType("ws"), time.Now(), errors.New("fail"))
}

func Test_shouldIncludeSynthetic(t *testing.T) {
	assert.False(t, shouldIncludeSynthetic("eth"))
	assert.True(t, shouldIncludeSynthetic("eni"))
	assert.Panics(t, func() { shouldIncludeSynthetic("foo") })
}
