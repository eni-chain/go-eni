package sdk_test

import (
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/evm/keeper"
	"github.com/cosmos/cosmos-sdk/x/evm/types"
	sdkmod "github.com/eni-chain/go-eni/syscontract/genesis/sdk"
	"testing"

	"cosmossdk.io/log"
	cosmossdk_io_math "cosmossdk.io/math"
	"cosmossdk.io/store/metrics"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func GetEVMKeeper() (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	transientStoreKey := storetypes.NewKVStoreKey(types.TransientStoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeIAVL, db)

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	k := keeper.NewKeeper(
		storeKey,
		transientStoreKey,
		nil,
		nil,
		nil,
		nil,
		cdc,
		log.NewNopLogger(),
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx

}

// Mock or real EVMKeeper should be initialized here
func getTestEVMKeeper() (*sdkmod.Hub, sdk.Context) {
	// This should be replaced with your test setup code
	evmKeeper, ctx := GetEVMKeeper()

	if evmKeeper == nil {
		panic("Failed to get EVM keeper")
	}
	hub, err := sdkmod.NewHub(evmKeeper)
	if err != nil {
		panic(err)
	}
	return hub, ctx
}

func TestApplyForValidator(t *testing.T) {
	hub, ctx := getTestEVMKeeper()

	caller := common.HexToAddress("0x123")
	node := common.HexToAddress("0x456")
	agent := common.HexToAddress("0x789")
	name := "TestValidator"
	description := "Testing validator application"
	pubKey := []byte("pubkey")
	value := cosmossdk_io_math.NewInt(10000)

	_, err := hub.ApplyForValidator(ctx, caller, node, agent, name, description, pubKey, &value)
	require.NoError(t, err)
}

func TestAuditPass(t *testing.T) {
	hub, ctx := getTestEVMKeeper()

	caller := common.HexToAddress("0xabc")
	operator := common.HexToAddress("0xdef")

	_, err := hub.AuditPass(ctx, caller, operator)
	require.NoError(t, err)
}

func TestBlockReward(t *testing.T) {
	hub, ctx := getTestEVMKeeper()

	caller := common.HexToAddress("0xaaa")
	node := common.HexToAddress("0xbbb")

	operator, amount, err := hub.BlockReward(ctx, caller, node)
	require.NoError(t, err)
	require.NotNil(t, operator)
	require.NotNil(t, amount)
}

func TestUpdateAdmin(t *testing.T) {
	hub, ctx := getTestEVMKeeper()

	caller := common.HexToAddress("0xaaa")
	newAdmin := common.HexToAddress("0xbbb")

	_, err := hub.UpdateAdmin(ctx, caller, newAdmin)
	require.NoError(t, err)
}
