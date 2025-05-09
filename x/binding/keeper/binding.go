package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/eni-chain/go-eni/x/binding/types"
)

// SetBinding set a specific binding in the store from its index
func (k Keeper) SetBinding(ctx context.Context, binding types.Binding) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.BindingKeyPrefix))
	b := k.cdc.MustMarshal(&binding)
	store.Set(types.BindingKey(
		binding.Index,
	), b)
}

// GetBinding returns a binding from its index
func (k Keeper) GetBinding(
	ctx context.Context,
	index string,

) (val types.Binding, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.BindingKeyPrefix))

	b := store.Get(types.BindingKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveBinding removes a binding from the store
func (k Keeper) RemoveBinding(
	ctx context.Context,
	index string,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.BindingKeyPrefix))
	store.Delete(types.BindingKey(
		index,
	))
}

// GetAllBinding returns all binding
func (k Keeper) GetAllBinding(ctx context.Context) (list []types.Binding) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.BindingKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Binding
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
