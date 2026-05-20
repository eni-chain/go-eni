package v1_test

import (
	"testing"

	"cosmossdk.io/log"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/app/upgrades"
	v1 "github.com/eni-chain/go-eni/app/upgrades/v1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// mockEVMCodeKeeper implements v1.EVMCodeKeeper for testing.
type mockEVMCodeKeeper struct {
	codes map[common.Address][]byte
}

func newMockEVMCodeKeeper() *mockEVMCodeKeeper {
	return &mockEVMCodeKeeper{
		codes: make(map[common.Address][]byte),
	}
}

func (m *mockEVMCodeKeeper) GetCode(_ sdk.Context, addr common.Address) []byte {
	return m.codes[addr]
}

func (m *mockEVMCodeKeeper) SetCode(_ sdk.Context, addr common.Address, code []byte) {
	m.codes[addr] = code
}

func newTestContext(chainID string, height int64) sdk.Context {
	return sdk.NewContext(nil, cmtproto.Header{
		ChainID: chainID,
		Height:  height,
	}, false, log.NewNopLogger())
}

// --- Testnet Tests ---

func TestTestnet_BytecodeRedirect_Success(t *testing.T) {
	mock := newMockEVMCodeKeeper()

	newAddr := common.HexToAddress(v1.TestnetNewContract)
	expectedCode := []byte("correct-bytecode-for-testing")
	mock.SetCode(sdk.Context{}, newAddr, expectedCode)

	oldAddr := common.HexToAddress(v1.TestnetOldContract)
	mock.SetCode(sdk.Context{}, oldAddr, []byte("incorrect-old-bytecode"))

	handler := v1.NewTestnetHandler(mock)
	ctx := newTestContext(v1.TestnetChainID, v1.TestnetHeight)

	err := handler.ExecuteHandler(ctx)
	require.NoError(t, err)

	actualCode := mock.GetCode(ctx, oldAddr)
	require.Equal(t, expectedCode, actualCode)
}

func TestTestnet_BytecodeRedirect_SourceEmpty(t *testing.T) {
	mock := newMockEVMCodeKeeper()

	handler := v1.NewTestnetHandler(mock)
	ctx := newTestContext(v1.TestnetChainID, v1.TestnetHeight)

	err := handler.ExecuteHandler(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "has no bytecode")
}

func TestTestnet_ChainIDMatch(t *testing.T) {
	mock := newMockEVMCodeKeeper()

	newAddr := common.HexToAddress(v1.TestnetNewContract)
	expectedCode := []byte("correct-bytecode")
	mock.SetCode(sdk.Context{}, newAddr, expectedCode)

	oldAddr := common.HexToAddress(v1.TestnetOldContract)
	mock.SetCode(sdk.Context{}, oldAddr, []byte("old-incorrect-bytecode"))

	manager := upgrades.NewHardForkManager(v1.TestnetChainID)
	manager.RegisterHandler(v1.NewTestnetHandler(mock))

	ctx := newTestContext(v1.TestnetChainID, v1.TestnetHeight)
	require.True(t, manager.TargetHeightReached(ctx))
	manager.ExecuteForTargetHeight(ctx)

	actualCode := mock.GetCode(ctx, oldAddr)
	require.Equal(t, expectedCode, actualCode)
}

func TestTestnet_ChainIDMismatch(t *testing.T) {
	mock := newMockEVMCodeKeeper()

	oldAddr := common.HexToAddress(v1.TestnetOldContract)
	oldCode := []byte("old-incorrect-bytecode")
	mock.SetCode(sdk.Context{}, oldAddr, oldCode)

	// Register testnet handler on mainnet manager - should be filtered out
	manager := upgrades.NewHardForkManager(v1.MainnetChainID)
	manager.RegisterHandler(v1.NewTestnetHandler(mock))

	ctx := newTestContext(v1.MainnetChainID, v1.TestnetHeight)
	require.False(t, manager.TargetHeightReached(ctx))

	actualCode := mock.GetCode(ctx, oldAddr)
	require.Equal(t, oldCode, actualCode)
}

func TestTestnet_BeforeTargetHeight(t *testing.T) {
	mock := newMockEVMCodeKeeper()

	oldAddr := common.HexToAddress(v1.TestnetOldContract)
	oldCode := []byte("old-incorrect-bytecode")
	mock.SetCode(sdk.Context{}, oldAddr, oldCode)

	newAddr := common.HexToAddress(v1.TestnetNewContract)
	mock.SetCode(sdk.Context{}, newAddr, []byte("correct-bytecode"))

	manager := upgrades.NewHardForkManager(v1.TestnetChainID)
	manager.RegisterHandler(v1.NewTestnetHandler(mock))

	ctx := newTestContext(v1.TestnetChainID, v1.TestnetHeight-1)
	require.False(t, manager.TargetHeightReached(ctx))

	actualCode := mock.GetCode(ctx, oldAddr)
	require.Equal(t, oldCode, actualCode)
}

func TestTestnet_SourceEmptyPanics(t *testing.T) {
	mock := newMockEVMCodeKeeper()

	manager := upgrades.NewHardForkManager(v1.TestnetChainID)
	manager.RegisterHandler(v1.NewTestnetHandler(mock))

	ctx := newTestContext(v1.TestnetChainID, v1.TestnetHeight)
	require.Panics(t, func() {
		manager.ExecuteForTargetHeight(ctx)
	})
}

// --- Mainnet Tests ---

func TestMainnet_BytecodeRedirect_Success(t *testing.T) {
	mock := newMockEVMCodeKeeper()

	newAddr := common.HexToAddress(v1.MainnetNewContract)
	expectedCode := []byte("correct-mainnet-bytecode")
	mock.SetCode(sdk.Context{}, newAddr, expectedCode)

	oldAddr := common.HexToAddress(v1.MainnetOldContract)
	mock.SetCode(sdk.Context{}, oldAddr, []byte("incorrect-mainnet-bytecode"))

	handler := v1.NewMainnetHandler(mock)
	ctx := newTestContext(v1.MainnetChainID, v1.MainnetHeight)

	err := handler.ExecuteHandler(ctx)
	require.NoError(t, err)

	actualCode := mock.GetCode(ctx, oldAddr)
	require.Equal(t, expectedCode, actualCode)
}

func TestMainnet_BytecodeRedirect_SourceEmpty(t *testing.T) {
	mock := newMockEVMCodeKeeper()

	handler := v1.NewMainnetHandler(mock)
	ctx := newTestContext(v1.MainnetChainID, v1.MainnetHeight)

	err := handler.ExecuteHandler(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "has no bytecode")
}

func TestMainnet_ChainIDMatch(t *testing.T) {
	mock := newMockEVMCodeKeeper()

	newAddr := common.HexToAddress(v1.MainnetNewContract)
	expectedCode := []byte("correct-mainnet-bytecode")
	mock.SetCode(sdk.Context{}, newAddr, expectedCode)

	oldAddr := common.HexToAddress(v1.MainnetOldContract)
	mock.SetCode(sdk.Context{}, oldAddr, []byte("old-mainnet-bytecode"))

	manager := upgrades.NewHardForkManager(v1.MainnetChainID)
	manager.RegisterHandler(v1.NewMainnetHandler(mock))

	ctx := newTestContext(v1.MainnetChainID, v1.MainnetHeight)
	require.True(t, manager.TargetHeightReached(ctx))
	manager.ExecuteForTargetHeight(ctx)

	actualCode := mock.GetCode(ctx, oldAddr)
	require.Equal(t, expectedCode, actualCode)
}

func TestMainnet_ChainIDMismatch(t *testing.T) {
	mock := newMockEVMCodeKeeper()

	oldAddr := common.HexToAddress(v1.MainnetOldContract)
	oldCode := []byte("old-mainnet-bytecode")
	mock.SetCode(sdk.Context{}, oldAddr, oldCode)

	// Register mainnet handler on testnet manager - should be filtered out
	manager := upgrades.NewHardForkManager(v1.TestnetChainID)
	manager.RegisterHandler(v1.NewMainnetHandler(mock))

	ctx := newTestContext(v1.TestnetChainID, v1.MainnetHeight)
	require.False(t, manager.TargetHeightReached(ctx))

	actualCode := mock.GetCode(ctx, oldAddr)
	require.Equal(t, oldCode, actualCode)
}

func TestMainnet_BeforeTargetHeight(t *testing.T) {
	mock := newMockEVMCodeKeeper()

	oldAddr := common.HexToAddress(v1.MainnetOldContract)
	oldCode := []byte("old-mainnet-bytecode")
	mock.SetCode(sdk.Context{}, oldAddr, oldCode)

	newAddr := common.HexToAddress(v1.MainnetNewContract)
	mock.SetCode(sdk.Context{}, newAddr, []byte("correct-mainnet-bytecode"))

	manager := upgrades.NewHardForkManager(v1.MainnetChainID)
	manager.RegisterHandler(v1.NewMainnetHandler(mock))

	ctx := newTestContext(v1.MainnetChainID, v1.MainnetHeight-1)
	require.False(t, manager.TargetHeightReached(ctx))

	actualCode := mock.GetCode(ctx, oldAddr)
	require.Equal(t, oldCode, actualCode)
}

func TestMainnet_SourceEmptyPanics(t *testing.T) {
	mock := newMockEVMCodeKeeper()

	manager := upgrades.NewHardForkManager(v1.MainnetChainID)
	manager.RegisterHandler(v1.NewMainnetHandler(mock))

	ctx := newTestContext(v1.MainnetChainID, v1.MainnetHeight)
	require.Panics(t, func() {
		manager.ExecuteForTargetHeight(ctx)
	})
}

// --- Handler Metadata Tests ---

func TestTestnetHandler_Metadata(t *testing.T) {
	mock := newMockEVMCodeKeeper()
	handler := v1.NewTestnetHandler(mock)

	require.Equal(t, v1.TestnetUpgradeName, handler.GetName())
	require.Equal(t, v1.TestnetChainID, handler.GetTargetChainID())
	require.Equal(t, v1.TestnetHeight, handler.GetTargetHeight())
}

func TestMainnetHandler_Metadata(t *testing.T) {
	mock := newMockEVMCodeKeeper()
	handler := v1.NewMainnetHandler(mock)

	require.Equal(t, v1.MainnetUpgradeName, handler.GetName())
	require.Equal(t, v1.MainnetChainID, handler.GetTargetChainID())
	require.Equal(t, v1.MainnetHeight, handler.GetTargetHeight())
}
