package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/eni-chain/go-eni/x/binding/types"
)

func (k msgServer) CreateBinding(goCtx context.Context, msg *types.MsgCreateBinding) (*types.MsgCreateBindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetBinding(
		ctx,
		msg.Index,
	)
	if isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var binding = types.Binding{
		Creator:       msg.Creator,
		Index:         msg.Index,
		EvmAddress:    msg.EvmAddress,
		CosmosAddress: msg.CosmosAddress,
	}

	k.SetBinding(
		ctx,
		binding,
	)
	return &types.MsgCreateBindingResponse{}, nil
}

func (k msgServer) UpdateBinding(goCtx context.Context, msg *types.MsgUpdateBinding) (*types.MsgUpdateBindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetBinding(
		ctx,
		msg.Index,
	)
	if !isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var binding = types.Binding{
		Creator:       msg.Creator,
		Index:         msg.Index,
		EvmAddress:    msg.EvmAddress,
		CosmosAddress: msg.CosmosAddress,
	}

	k.SetBinding(ctx, binding)

	return &types.MsgUpdateBindingResponse{}, nil
}

func (k msgServer) DeleteBinding(goCtx context.Context, msg *types.MsgDeleteBinding) (*types.MsgDeleteBindingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetBinding(
		ctx,
		msg.Index,
	)
	if !isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveBinding(
		ctx,
		msg.Index,
	)

	return &types.MsgDeleteBindingResponse{}, nil
}
