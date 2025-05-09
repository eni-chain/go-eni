package binding

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/eni-chain/go-eni/x/binding/keeper"
	"github.com/eni-chain/go-eni/x/binding/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the binding
	for _, elem := range genState.BindingList {
		k.SetBinding(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.BindingList = k.GetAllBinding(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
