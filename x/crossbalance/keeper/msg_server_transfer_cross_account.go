package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/crossbalance/types"
)

func (k msgServer) TransferCrossAccount(goCtx context.Context, msg *types.MsgTransferCrossAccount) (*types.MsgTransferCrossAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgTransferCrossAccountResponse{}, nil
}
