package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var _ sdk.Msg = &MsgTransferEniToEvm{}

func NewMsgTransferEniToEvm(
	sender string,
	to common.Address,
	amount sdk.Coin,
) *MsgTransferEniToEvm {
	return &MsgTransferEniToEvm{
		Sender: sender,
		To:     to.Hex(),
		Amount: amount,
	}
}

func (msg *MsgTransferEniToEvm) ValidateBasic() error {
	if msg.Sender == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "sender address cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.To == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "evm address cannot be empty")
	}
	if !common.IsHexAddress(msg.To) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid evm address format")
	}

	if msg.Amount.IsNil() || msg.Amount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount cannot be zero")
	}

	return nil
}

var _ sdk.Msg = &MsgTransferEvmToEni{}

func NewMsgTransferEvmToEni(
	from common.Address,
	ak []byte,
	to string,
	amount sdk.Coin,
) *MsgTransferEvmToEni {
	return &MsgTransferEvmToEni{
		From:   from.Hex(),
		Ak:     ak,
		To:     to,
		Amount: amount,
	}
}

func (msg *MsgTransferEvmToEni) ValidateBasic() error {
	if msg.From == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "evm address cannot be empty")
	}
	if !common.IsHexAddress(msg.From) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid evm address format")
	}

	if len(msg.Ak) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidPubKey, "public key cannot be empty")
	}

	if msg.To == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "cosmos address cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(msg.To); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid cosmos address (%s)", err)
	}

	if msg.Amount.IsNil() || msg.Amount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount cannot be zero")
	}

	return nil
}
