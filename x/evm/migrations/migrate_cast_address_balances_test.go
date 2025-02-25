package migrations_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/eni-chain/go-eni/testutil/keeper"
	"github.com/eni-chain/go-eni/x/evm/migrations"
	"github.com/eni-chain/go-eni/x/evm/types"
	"github.com/stretchr/testify/require"
)

func TestMigrateCastAddressBalances(t *testing.T) {
	k := testkeeper.EVMTestApp.EvmKeeper
	ctx := testkeeper.EVMTestApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	require.Nil(t, k.BankKeeper().MintCoins(ctx, types.ModuleName, testkeeper.UeniCoins(100)))
	// unassociated account with funds
	eniAddr1, evmAddr1 := testkeeper.MockAddressPair()
	require.Nil(t, k.BankKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(evmAddr1[:]), testkeeper.UeniCoins(10)))
	// associated account without funds
	eniAddr2, evmAddr2 := testkeeper.MockAddressPair()
	k.SetAddressMapping(ctx, eniAddr2, evmAddr2)
	// associated account with funds
	eniAddr3, evmAddr3 := testkeeper.MockAddressPair()
	require.Nil(t, k.BankKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(evmAddr3[:]), testkeeper.UeniCoins(10)))
	k.SetAddressMapping(ctx, eniAddr3, evmAddr3)

	require.Nil(t, migrations.MigrateCastAddressBalances(ctx, &k))

	require.Equal(t, sdk.NewInt(10), k.BankKeeper().GetBalance(ctx, sdk.AccAddress(evmAddr1[:]), "ueni").Amount)
	require.Equal(t, sdk.ZeroInt(), k.BankKeeper().GetBalance(ctx, eniAddr1, "ueni").Amount)
	require.Equal(t, sdk.ZeroInt(), k.BankKeeper().GetBalance(ctx, sdk.AccAddress(evmAddr2[:]), "ueni").Amount)
	require.Equal(t, sdk.ZeroInt(), k.BankKeeper().GetBalance(ctx, eniAddr2, "ueni").Amount)
	require.Equal(t, sdk.ZeroInt(), k.BankKeeper().GetBalance(ctx, sdk.AccAddress(evmAddr3[:]), "ueni").Amount)
	require.Equal(t, sdk.NewInt(10), k.BankKeeper().GetBalance(ctx, eniAddr3, "ueni").Amount)
}
