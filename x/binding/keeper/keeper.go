package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/eni-chain/go-eni/x/binding/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger
		authority    string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// HasBinding checks if a binding relationship exists
func (k Keeper) HasBinding(ctx sdk.Context, cosmosAddr sdk.AccAddress) (bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetBindingKey(cosmosAddr)
	return store.Has(key)
}

// DeleteBinding deletes the binding relationship
func (k Keeper) DeleteBinding(ctx sdk.Context, cosmosAddr sdk.AccAddress) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetBindingKey(cosmosAddr)
	return store.Delete(key)
}
