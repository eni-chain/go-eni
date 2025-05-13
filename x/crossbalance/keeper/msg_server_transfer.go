package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/eni-chain/go-eni/x/crossbalance/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (k msgServer) TransferEniToEvm(goCtx context.Context, msg *types.MsgTransferEniToEvm) (*types.MsgTransferEniToEvmResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert sender address
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	// Convert EVM address
	evmAddr := common.HexToAddress(msg.To)

	// Check if sender has enough balance
	if !k.BankKeeper.HasBalance(ctx, sender, msg.Amount) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient balance")
	}

	// Deduct from sender's balance
	if err := k.BankKeeper.SendCoins(ctx, sender, evmAddr.Bytes(), sdk.NewCoins(msg.Amount)); err != nil {
		return nil, errorsmod.Wrapf(err, "failed to send coins")
	}

	return &types.MsgTransferEniToEvmResponse{}, nil
}

func (k msgServer) TransferEvmToEni(goCtx context.Context, msg *types.MsgTransferEvmToEni) (*types.MsgTransferEvmToEniResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert EVM address
	//evmAddr := common.HexToAddress(msg.From)

	// Convert public key to Cosmos address
	pubKey, err := crypto.UnmarshalPubkey(msg.Ak)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid public key (%s)", err)
	}

	// Get the Cosmos address from the public key
	cosmosFrom := sdk.AccAddress(crypto.PubkeyToAddress(*pubKey).Bytes())

	// Convert destination address
	to, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid destination address (%s)", err)
	}

	// Check if sender has enough balance
	if !k.BankKeeper.HasBalance(ctx, cosmosFrom, msg.Amount) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient balance")
	}

	// Deduct from sender's balance and add to recipient's balance
	if err := k.BankKeeper.SendCoins(ctx, cosmosFrom, to, sdk.NewCoins(msg.Amount)); err != nil {
		return nil, errorsmod.Wrapf(err, "failed to send coins")
	}

	return &types.MsgTransferEvmToEniResponse{}, nil
}
