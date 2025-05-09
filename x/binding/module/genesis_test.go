package binding_test

import (
	"testing"

	keepertest "github.com/eni-chain/go-eni/testutil/keeper"
	"github.com/eni-chain/go-eni/testutil/nullify"
	binding "github.com/eni-chain/go-eni/x/binding/module"
	"github.com/eni-chain/go-eni/x/binding/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		BindingList: []types.Binding{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BindingKeeper(t)
	binding.InitGenesis(ctx, k, genesisState)
	got := binding.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.BindingList, got.BindingList)
	// this line is used by starport scaffolding # genesis/test/assert
}
