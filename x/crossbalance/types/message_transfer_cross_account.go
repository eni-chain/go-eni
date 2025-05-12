package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgTransferCrossAccount{}

func NewMsgTransferCrossAccount(creator string, fromAddress string, toAddress string, amount string) *MsgTransferCrossAccount {
	coin, err := sdk.ParseCoinNormalized(amount)
	if err != nil {
		panic(err) // TODO: handle error
	}

	return &MsgTransferCrossAccount{
		Creator:     creator,
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      &coin,
	}
}

// Route implements the sdk.Msg interface
func (msg *MsgTransferCrossAccount) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface
func (msg *MsgTransferCrossAccount) Type() string {
	return TypeMsgTransferCrossAccount
}

// GetSigners implements the sdk.Msg interface
func (msg *MsgTransferCrossAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgTransferCrossAccount) GetSignBytes() []byte {
	bz, err := msg.Marshal()
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgTransferCrossAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.FromAddress); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.ToAddress); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid to address (%s)", err)
	}
	if msg.Amount == nil || !msg.Amount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid amount")
	}
	return nil
}
