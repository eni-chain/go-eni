package types

import (
	"strings"

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
	//if msg.Creator == "" {
	//	return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "creator address cannot be empty")
	//}
	//if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
	//	return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	//}

	if msg.Index == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index cannot be empty")
	}

	if msg.EvmAddress == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "evm address cannot be empty")
	}
	// Validate EVM address format (should be 0x followed by 40 hex characters)
	evmAddr := strings.ToLower(msg.EvmAddress)
	if !strings.HasPrefix(evmAddr, "0x") || len(evmAddr) != 42 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid evm address format")
	}

	if msg.CosmosAddress == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "cosmos address cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(msg.CosmosAddress); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid cosmos address (%s)", err)
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
	if msg.Creator == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "creator address cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Index == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index cannot be empty")
	}

	if msg.EvmAddress == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "evm address cannot be empty")
	}
	// Validate EVM address format (should be 0x followed by 40 hex characters)
	evmAddr := strings.ToLower(msg.EvmAddress)
	if !strings.HasPrefix(evmAddr, "0x") || len(evmAddr) != 42 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid evm address format")
	}

	if msg.CosmosAddress == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "cosmos address cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(msg.CosmosAddress); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid cosmos address (%s)", err)
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
	if msg.Creator == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "creator address cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Index == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index cannot be empty")
	}

	return nil
}
