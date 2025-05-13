package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/utils/helpers"
	"github.com/eni-chain/go-eni/x/crossbalance/types"
	"github.com/ethereum/go-ethereum/common"
)

func (k msgServer) TransferEniToEvm(goCtx context.Context, msg *types.MsgTransferEniToEvm) (*types.MsgTransferEniToEvmResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert sender address
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	// Convert EVM address
	evmAddr := common.HexToAddress(msg.To)
	if evmAddr == (common.Address{}) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid EVM address")
	}

	// Step 1: Transfer from ENI address to burn address
	burnAddr := sdk.AccAddress(common.HexToAddress("0x0000000000000000000000000000000000000000").Bytes())
	if err := k.BankKeeper.SendCoins(ctx, senderAddr, burnAddr, sdk.NewCoins(msg.Amount)); err != nil {
		return nil, errorsmod.Wrapf(err, "failed to burn coins")
	}

	// Step 2: Add balance to evm address
	if err = k.EvmKeeper.BankKeeper().AddCoins(ctx, nil, sdk.Coins{msg.Amount}); err != nil { // ToDo: x/evm evm addr should be sdk.AccAddress
		return nil, errorsmod.Wrapf(err, "failed to add balance to EVM address")
	}
	return &types.MsgTransferEniToEvmResponse{}, nil
}

func (k msgServer) TransferEvmToEni(goCtx context.Context, msg *types.MsgTransferEvmToEni) (*types.MsgTransferEvmToEniResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert EVM address
	evmAddr := common.HexToAddress(msg.From)
	if evmAddr == (common.Address{}) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid EVM address")
	}

	// Convert public key to Cosmos address
	evmAddrFromPubkey, _, _, err := helpers.GetAddressesFromPubkeyBytes(msg.Ak)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid EVM address (%s)", err)
	}

	// Check if the EVM address is the same as the one in the public key
	if evmAddr.Hex() != evmAddrFromPubkey.Hex() {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid EVM address in public key")
	}

	// Step 1: Burn coins from EVM address
	evmAccAddr := sdk.AccAddress(evmAddr.Bytes())              // todo evm address
	if !k.BankKeeper.HasBalance(ctx, evmAccAddr, msg.Amount) { // todo evmkeeper bankkeeper add hasbalance function
		return nil, errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient balance in EVM")
	}

	if err := k.BankKeeper.SendCoins(ctx, evmAccAddr, sdk.AccAddress(common.HexToAddress("0x0000000000000000000000000000000000000000").Bytes()), sdk.NewCoins(msg.Amount)); err != nil {
		return nil, errorsmod.Wrapf(err, "failed to burn coins from EVM address")
	}

	// Step 2: Mint coins to ENI address
	recipientAddr, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid destination address (%s)", err)
	}

	if err := k.BankKeeper.AddCoins(ctx, recipientAddr, sdk.NewCoins(msg.Amount)); err != nil {
		return nil, errorsmod.Wrapf(err, "failed to add coins")
	}

	return &types.MsgTransferEvmToEniResponse{}, nil
}
