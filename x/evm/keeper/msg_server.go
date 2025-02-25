package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/evm/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (server msgServer) EVMTransaction(goCtx context.Context, msg *types.MsgEVMTransaction) (serverRes *types.MsgEVMTransactionResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return
}

func (server msgServer) Send(goCtx context.Context, msg *types.MsgSend) (res *types.MsgSendResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return
}

func (server msgServer) RegisterPointer(goCtx context.Context, msg *types.MsgRegisterPointer) (res *types.MsgRegisterPointerResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return
}

func (server msgServer) Associate(goCtx context.Context, msg *types.MsgAssociate) (res *types.MsgAssociateResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return
}

func (server msgServer) AssociateContractAddress(goCtx context.Context, msg *types.MsgAssociateContractAddress) (res *types.MsgAssociateContractAddressResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return

}
