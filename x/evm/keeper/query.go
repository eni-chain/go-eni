package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) EniAddressByEVMAddress(goCtx context.Context, req *types.QueryEniAddressByEVMAddressRequest) (res *types.QueryEniAddressByEVMAddressResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr := common.HexToAddress(req.EvmAddress)
	eniAddr := k.GetEniAddressOrDefault(ctx, addr)
	return &types.QueryEniAddressByEVMAddressResponse{EniAddress: eniAddr.String()}, nil
}
func (k Keeper) EVMAddressByEniAddress(goCtx context.Context, req *types.QueryEVMAddressByEniAddressRequest) (res *types.QueryEVMAddressByEniAddressResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.EniAddress)
	if err != nil {
		return nil, err
	}
	evmAddr := k.GetEVMAddressOrDefault(ctx, addr)
	return &types.QueryEVMAddressByEniAddressResponse{EvmAddress: evmAddr.Hex()}, nil
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
