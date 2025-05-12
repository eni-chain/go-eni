package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"

	sdkerrors "cosmossdk.io/errors"

	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/crossbalance/types"
)

// RootCodespace is the codespace for all errors defined in this package
const RootCodespace = "sdk"

var (
	// ErrInsufficientFunds is used when the account cannot pay requested amount.

	ErrInsufficientFunds = errorsmod.Register(RootCodespace, 5, "insufficient funds")

	// ErrUnknownAddress to doc
	ErrUnknownAddress = errorsmod.Register(RootCodespace, 9, "unknown address")
)

func (k msgServer) TransferCrossAccount(goCtx context.Context, msg *types.MsgTransferCrossAccount) (*types.MsgTransferCrossAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get the shared internal account of from and to
	fromRealAddr, err := k.resolveRealAddress(ctx, msg.FromAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid from address")
	}

	toRealAddr, err := k.resolveRealAddress(ctx, msg.ToAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid to address")
	}

	// query the balance of fromRealAddr
	coin := *msg.Amount
	if !k.BankKeeper.HasBalance(ctx, fromRealAddr, coin) {
		return nil, sdkerrors.Wrapf(ErrInsufficientFunds, "insufficient funds for %s", fromRealAddr.String())
	}

	// execute the transfer (internal account to internal account)
	err = k.BankKeeper.SendCoins(ctx, fromRealAddr, toRealAddr, sdk.NewCoins(coin))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to send coins")
	}

	return &types.MsgTransferCrossAccountResponse{}, nil
}

func (k msgServer) resolveRealAddress(ctx sdk.Context, inputAddr string) (sdk.AccAddress, error) {
	// if inputAddr is a cosmos address (bech32)
	if accAddr, err := sdk.AccAddressFromBech32(inputAddr); err == nil {
		return accAddr, nil
	}

	// try to find the evm address binding
	evmAddr := strings.ToLower(inputAddr)
	cosmosAddress, err := k.BindingKeeper.GetBinding(ctx, sdk.AccAddress(evmAddr))
	if err != nil {
		return nil, sdkerrors.Wrapf(ErrUnknownAddress, "no binding for evm address %s", evmAddr)
	}

	return sdk.AccAddressFromBech32(cosmosAddress.CosmosAddress)
}
