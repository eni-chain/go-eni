package migrations

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/utils"
	"github.com/eni-chain/go-eni/x/evm/keeper"
	"github.com/eni-chain/go-eni/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

func MigrateERCNativePointers(ctx sdk.Context, k *keeper.Keeper) error {
	iter := prefix.NewStore(ctx.KVStore(k.GetStoreKey()), append(types.PointerRegistryPrefix, types.PointerERC20NativePrefix...)).ReverseIterator(nil, nil)
	defer iter.Close()
	seen := map[string]struct{}{}
	for ; iter.Valid(); iter.Next() {
		token := string(iter.Key()[:len(iter.Key())-2]) // last two bytes are version
		if _, ok := seen[token]; ok {
			continue
		}
		seen[token] = struct{}{}
		addr := common.BytesToAddress(iter.Value())
		oName, err := k.QueryERCSingleOutput(ctx, "native", addr, "name")
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Failed to upgrade pointer for %s due to failed name query: %s", token, err))
			continue
		}
		oSymbol, err := k.QueryERCSingleOutput(ctx, "native", addr, "symbol")
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Failed to upgrade pointer for %s due to failed symbol query: %s", token, err))
			continue
		}
		oDecimals, err := k.QueryERCSingleOutput(ctx, "native", addr, "decimals")
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Failed to upgrade pointer for %s due to failed decimal query: %s", token, err))
			continue
		}
		_ = k.RunWithOneOffEVMInstance(ctx, func(e *vm.EVM) error {
			_, err := k.UpsertERCNativePointer(ctx.WithGasMeter(sdk.NewInfiniteGasMeterWithMultiplier(ctx)), e, token, utils.ERCMetadata{
				Name:     oName.(string),
				Symbol:   oSymbol.(string),
				Decimals: oDecimals.(uint8),
			})
			return err
		}, func(s1, s2 string) {
			ctx.Logger().Error(fmt.Sprintf("Failed to upgrade pointer for %s at step %s due to %s", token, s1, s2))
		})
	}
	return nil
}
