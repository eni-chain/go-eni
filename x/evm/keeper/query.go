package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/evm/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) EniAddressByEVMAddress(goCtx context.Context, req *types.QuerySeiAddressByEVMAddressRequest) (res *types.QuerySeiAddressByEVMAddressResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return
}
func (k Keeper) EVMAddressBySeiAddress(goCtx context.Context, req *types.QueryEVMAddressBySeiAddressRequest) (res *types.QueryEVMAddressBySeiAddressResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return
}
func (k Keeper) StaticCall(goCtx context.Context, req *types.QueryStaticCallRequest) (res *types.QueryStaticCallResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return
}
func (k Keeper) Pointer(goCtx context.Context, req *types.QueryPointerRequest) (res *types.QueryPointerResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return
}
func (k Keeper) PointerVersion(goCtx context.Context, req *types.QueryPointerVersionRequest) (res *types.QueryPointerVersionResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return
}
func (k Keeper) Pointee(goCtx context.Context, req *types.QueryPointeeRequest) (res *types.QueryPointeeResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return
}
