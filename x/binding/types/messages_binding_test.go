package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/eni-chain/go-eni/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateBinding_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateBinding
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateBinding{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateBinding{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateBinding_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateBinding
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateBinding{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateBinding{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgDeleteBinding_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteBinding
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteBinding{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteBinding{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
