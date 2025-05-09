package crossbalance_test

import (
	"testing"

	keepertest "github.com/eni-chain/go-eni/testutil/keeper"
	"github.com/eni-chain/go-eni/testutil/nullify"
	crossbalance "github.com/eni-chain/go-eni/x/crossbalance/module"
	"github.com/eni-chain/go-eni/x/crossbalance/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.CrossbalanceKeeper(t)
	crossbalance.InitGenesis(ctx, k, genesisState)
	got := crossbalance.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
