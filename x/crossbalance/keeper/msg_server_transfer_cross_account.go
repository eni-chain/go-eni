package keeper

import (
	"context"
	"errors"

	sdkerrors "cosmossdk.io/errors"

	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/crossbalance/types"
)

func (k msgServer) TransferCrossAccount(goCtx context.Context, msg *types.MsgTransferCrossAccount) (*types.MsgTransferCrossAccountResponse, error) {
	// get the shared internal account of from and to
	fromRealAddr, err := k.resolveRealAddress(goCtx, msg.FromAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid from address")
	}

	toRealAddr, err := k.resolveRealAddress(goCtx, msg.ToAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid to address")
	}

	// query the balance of fromRealAddr
	coin := *msg.Amount
	if !k.BankKeeper.HasBalance(goCtx, fromRealAddr, coin) {
		return nil, sdkerrors.Wrapf(err, "insufficient funds for %s", fromRealAddr.String())
	}

	// execute the transfer (internal account to internal account)
	err = k.BankKeeper.SendCoins(goCtx, fromRealAddr, toRealAddr, sdk.NewCoins(coin))
	if err != nil {
		return nil, sdkerrors.Wrap(errors.New("failed to send coins"), "failed to send coins")
	}

	return &types.MsgTransferCrossAccountResponse{}, nil
}

func (k msgServer) resolveRealAddress(goCtx context.Context, inputAddr string) (sdk.AccAddress, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// if inputAddr is a cosmos address (bech32)
	if accAddr, err := sdk.AccAddressFromBech32(inputAddr); err == nil {
		return accAddr, nil
	}

	// try to find the evm address binding
	evmAddr := strings.ToLower(inputAddr)
	hasBinding, err := k.BindingKeeper.HasBinding(ctx, sdk.AccAddress(evmAddr))
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to check binding for evm address %s", evmAddr)
	}
	if !hasBinding {
		return nil, sdkerrors.Wrapf(errors.New("no binding for evm address"), "no binding for evm address %s", evmAddr)
	}

	// Since we can't get the binding directly, we'll need to use the EVM address as the real address
	return sdk.AccAddress(evmAddr), nil
}
