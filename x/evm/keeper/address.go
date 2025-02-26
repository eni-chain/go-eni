package keeper

import (
	//"cosmossdk.io/store/prefix"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
)

func (k *Keeper) SetAddressMapping(ctx sdk.Context, eniAddress sdk.AccAddress, evmAddress common.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.EVMAddressToEniAddressKey(evmAddress), eniAddress)
	store.Set(types.EniAddressToEVMAddressKey(eniAddress), evmAddress[:])
	if !k.accountKeeper.HasAccount(ctx, eniAddress) {
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, eniAddress))
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeAddressAssociated,
		sdk.NewAttribute(types.AttributeKeyEniAddress, eniAddress.String()),
		sdk.NewAttribute(types.AttributeKeyEvmAddress, evmAddress.Hex()),
	))
}

func (k *Keeper) DeleteAddressMapping(ctx sdk.Context, eniAddress sdk.AccAddress, evmAddress common.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.EVMAddressToEniAddressKey(evmAddress))
	store.Delete(types.EniAddressToEVMAddressKey(eniAddress))
}

func (k *Keeper) GetEVMAddress(ctx sdk.Context, eniAddress sdk.AccAddress) (common.Address, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EniAddressToEVMAddressKey(eniAddress))
	addr := common.Address{}
	if bz == nil {
		return addr, false
	}
	copy(addr[:], bz)
	return addr, true
}

func (k *Keeper) GetEVMAddressOrDefault(ctx sdk.Context, eniAddress sdk.AccAddress) common.Address {
	addr, ok := k.GetEVMAddress(ctx, eniAddress)
	if ok {
		return addr
	}
	return common.BytesToAddress(eniAddress)
}

func (k *Keeper) GetEniAddress(ctx sdk.Context, evmAddress common.Address) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EVMAddressToEniAddressKey(evmAddress))
	if bz == nil {
		return []byte{}, false
	}
	return bz, true
}

func (k *Keeper) GetEniAddressOrDefault(ctx sdk.Context, evmAddress common.Address) sdk.AccAddress {
	addr, ok := k.GetEniAddress(ctx, evmAddress)
	if ok {
		return addr
	}
	return sdk.AccAddress(evmAddress[:])
}

func (k *Keeper) IterateEniAddressMapping(ctx sdk.Context, cb func(evmAddr common.Address, eniAddr sdk.AccAddress) bool) {
	iter := prefix.NewStore(ctx.KVStore(k.storeKey), types.EVMAddressToEniAddressKeyPrefix).Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		evmAddr := common.BytesToAddress(iter.Key())
		eniAddr := sdk.AccAddress(iter.Value())
		if cb(evmAddr, eniAddr) {
			break
		}
	}
}

// A sdk.AccAddress may not receive funds from bank if it's the result of direct-casting
// from an EVM address AND the originating EVM address has already been associated with
// a true (i.e. derived from the same pubkey) sdk.AccAddress.
func (k *Keeper) CanAddressReceive(ctx sdk.Context, addr sdk.AccAddress) bool {
	directCast := common.BytesToAddress(addr) // casting goes both directions since both address formats have 20 bytes
	associatedAddr, isAssociated := k.GetEniAddress(ctx, directCast)
	// if the associated address is the cast address itself, allow the address to receive (e.g. EVM contract addresses)
	return associatedAddr.Equals(addr) || !isAssociated // this means it's either a cast address that's not associated yet, or not a cast address at all.
}

type EvmAddressHandler struct {
	evmKeeper *Keeper
}

func NewEvmAddressHandler(evmKeeper *Keeper) EvmAddressHandler {
	return EvmAddressHandler{evmKeeper: evmKeeper}
}

func (h EvmAddressHandler) GetEniAddressFromString(ctx sdk.Context, address string) (sdk.AccAddress, error) {
	if common.IsHexAddress(address) {
		parsedAddress := common.HexToAddress(address)
		return h.evmKeeper.GetEniAddressOrDefault(ctx, parsedAddress), nil
	}
	parsedAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	return parsedAddress, nil
}
