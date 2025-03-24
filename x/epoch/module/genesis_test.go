package epoch_test

import (
	"testing"

	keepertest "github.com/eni-chain/go-eni/testutil/keeper"
	"github.com/eni-chain/go-eni/testutil/nullify"
	epoch "github.com/eni-chain/go-eni/x/epoch/module"
	"github.com/eni-chain/go-eni/x/epoch/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.EpochKeeper(t)
	epoch.InitGenesis(ctx, k, genesisState)
	got := epoch.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
