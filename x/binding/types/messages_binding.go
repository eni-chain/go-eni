package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateBinding{}

func NewMsgCreateBinding(
	creator string,
	index string,
	evmAddress string,
	cosmosAddress string,

) *MsgCreateBinding {
	return &MsgCreateBinding{
		Creator:       creator,
		Index:         index,
		EvmAddress:    evmAddress,
		CosmosAddress: cosmosAddress,
	}
}

func (msg *MsgCreateBinding) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateBinding{}

func NewMsgUpdateBinding(
	creator string,
	index string,
	evmAddress string,
	cosmosAddress string,

) *MsgUpdateBinding {
	return &MsgUpdateBinding{
		Creator:       creator,
		Index:         index,
		EvmAddress:    evmAddress,
		CosmosAddress: cosmosAddress,
	}
}

func (msg *MsgUpdateBinding) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteBinding{}

func NewMsgDeleteBinding(
	creator string,
	index string,

) *MsgDeleteBinding {
	return &MsgDeleteBinding{
		Creator: creator,
		Index:   index,
	}
}

func (msg *MsgDeleteBinding) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
