package epoch

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/eni-chain/go-eni/x/epoch/keeper"
	"github.com/eni-chain/go-eni/x/epoch/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
	k.SetEpoch(
		ctx,
		*genState.Epoch,
	)
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	epoch := k.GetEpoch(ctx)
	genesis.Epoch = &epoch
	return genesis
}
