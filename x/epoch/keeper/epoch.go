package keeper

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/epoch/types"
)

const EpochKey = "epoch"

func (k Keeper) SetEpoch(ctx sdk.Context, epoch types.Epoch) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	value, err := k.cdc.Marshal(&epoch)
	if err != nil {
		return err
	}
	store.Set([]byte(EpochKey), value)
	return nil
}

func (k Keeper) GetEpoch(ctx sdk.Context) (epoch types.Epoch) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := adapter.Get([]byte(EpochKey))
	if b != nil {
		k.cdc.MustUnmarshal(b, &epoch)
		return epoch
	}
	return epoch
}
