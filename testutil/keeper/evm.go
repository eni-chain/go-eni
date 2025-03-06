package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/eni-chain/go-eni/x/evm/keeper"
	"github.com/eni-chain/go-eni/x/evm/types"
)

func EvmKeeper(t testing.TB) (keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	//registry := codectypes.NewInterfaceRegistry()
	//cdc := codec.NewProtoCodec(registry)
	//authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	//k := keeper.NewKeeper(
	//	cdc,
	//	runtime.NewKVStoreService(storeKey),
	//	log.NewNopLogger(),
	//	authority.String(),
	//	nil,
	//	nil,
	//	nil,nil,
	//)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	//if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
	//	panic(err)
	//}

	//return k, ctx

	return keeper.Keeper{}, ctx
}
