package keeper_test

import (
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	keepertest "github.com/eni-chain/go-eni/testutil/keeper"
	"github.com/eni-chain/go-eni/x/binding/keeper"
	"github.com/eni-chain/go-eni/x/binding/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestBindingMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.BindingKeeper(t)
	srv := keeper.NewMsgServerImpl(k)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateBinding{Creator: creator,
			Index: strconv.Itoa(i),
		}
		_, err := srv.CreateBinding(ctx, expected)
		require.NoError(t, err)
		rst, found := k.GetBinding(ctx,
			expected.Index,
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestBindingMsgServerUpdate(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgUpdateBinding
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateBinding{Creator: creator,
				Index: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateBinding{Creator: "B",
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateBinding{Creator: creator,
				Index: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.BindingKeeper(t)
			srv := keeper.NewMsgServerImpl(k)
			expected := &types.MsgCreateBinding{Creator: creator,
				Index: strconv.Itoa(0),
			}
			_, err := srv.CreateBinding(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateBinding(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetBinding(ctx,
					expected.Index,
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestBindingMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeleteBinding
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteBinding{Creator: creator,
				Index: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteBinding{Creator: "B",
				Index: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteBinding{Creator: creator,
				Index: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.BindingKeeper(t)
			srv := keeper.NewMsgServerImpl(k)

			_, err := srv.CreateBinding(ctx, &types.MsgCreateBinding{Creator: creator,
				Index: strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteBinding(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetBinding(ctx,
					tc.request.Index,
				)
				require.False(t, found)
			}
		})
	}
}
