package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/evm/keeper"
	"github.com/eni-chain/go-eni/x/evm/types"
)

func MigrateEip1559MaxFeePerGas(ctx sdk.Context, k *keeper.Keeper) error {
	keeperParams := k.GetParamsIfExists(ctx)
	keeperParams.MaximumFeePerGas = types.DefaultParams().MaximumFeePerGas
	k.SetParams(ctx, keeperParams)
	return nil
}
