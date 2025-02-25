package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/eni-chain/go-eni/x/evm/artifacts/native"

	"github.com/eni-chain/go-eni/x/evm/types"
)

type PointerGetter func(sdk.Context, string) (common.Address, uint16, bool)
type PointerSetter func(sdk.Context, string, common.Address) error

var ErrorPointerToPointerNotAllowed = sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "cannot create a pointer to a pointer")

// ERC20 -> Native Token
func (k *Keeper) SetERC20NativePointer(ctx sdk.Context, token string, addr common.Address) error {
	return k.SetERC20NativePointerWithVersion(ctx, token, addr, native.CurrentVersion)
}

// ERC20 -> Native Token
func (k *Keeper) SetERC20NativePointerWithVersion(ctx sdk.Context, token string, addr common.Address, version uint16) error {
	if k.cwAddressIsPointer(ctx, token) {
		return ErrorPointerToPointerNotAllowed
	}
	err := k.setPointerInfo(ctx, types.PointerERC20NativeKey(token), addr[:], version)
	if err != nil {
		return err
	}
	return k.setPointerInfo(ctx, types.PointerReverseRegistryKey(addr), []byte(token), version)
}

// ERC20 -> Native Token
func (k *Keeper) GetERC20NativePointer(ctx sdk.Context, token string) (addr common.Address, version uint16, exists bool) {
	addrBz, version, exists := k.GetPointerInfo(ctx, types.PointerERC20NativeKey(token))
	if exists {
		addr = common.BytesToAddress(addrBz)
	}
	return
}

func (k *Keeper) evmAddressIsPointer(ctx sdk.Context, addr common.Address) bool {
	_, _, exists := k.GetPointerInfo(ctx, types.PointerReverseRegistryKey(addr))
	return exists
}

func (k *Keeper) cwAddressIsPointer(ctx sdk.Context, addr string) bool {
	_, _, exists := k.GetPointerInfo(ctx, types.PointerReverseRegistryKey(common.BytesToAddress([]byte(addr))))
	return exists
}

func (k *Keeper) GetPointerInfo(ctx sdk.Context, pref []byte) (addr []byte, version uint16, exists bool) {
	store := prefix.NewStore(ctx.KVStore(k.GetStoreKey()), pref)
	iter := store.ReverseIterator(nil, nil)
	defer iter.Close()
	exists = iter.Valid()
	if !exists {
		return
	}
	version = binary.BigEndian.Uint16(iter.Key())
	addr = iter.Value()
	return
}

func (k *Keeper) setPointerInfo(ctx sdk.Context, pref []byte, addr []byte, version uint16) error {
	store := prefix.NewStore(ctx.KVStore(k.GetStoreKey()), pref)
	versionBz := make([]byte, 2)
	binary.BigEndian.PutUint16(versionBz, version)
	store.Set(versionBz, addr)
	return nil
}

func (k *Keeper) deletePointerInfo(ctx sdk.Context, pref []byte, version uint16) {
	store := prefix.NewStore(ctx.KVStore(k.GetStoreKey()), pref)
	versionBz := make([]byte, 2)
	binary.BigEndian.PutUint16(versionBz, version)
	store.Delete(versionBz)
}

func (k *Keeper) GetNativePointee(ctx sdk.Context, erc20Address string) (token string, version uint16, exists bool) {
	// Ensure the key matches how it was set in SetERC20NativePointer
	key := types.PointerReverseRegistryKey(common.HexToAddress(erc20Address))
	addrBz, version, exists := k.GetPointerInfo(ctx, key)
	if exists {
		token = string(addrBz)
	}
	return
}
