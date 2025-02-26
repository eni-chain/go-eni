package keeper_test

//
//import (
//	"bytes"
//	"testing"
//
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/eni-chain/go-eni/testutil/keeper"
//	evmkeeper "github.com/eni-chain/go-eni/x/evm/keeper"
//	"github.com/stretchr/testify/require"
//)
//
//func TestSetGetAddressMapping(t *testing.T) {
//	k := &keeper.EVMTestApp.EvmKeeper
//	ctx := keeper.EVMTestApp.GetContextForDeliverTx([]byte{})
//	eniAddr, evmAddr := keeper.MockAddressPair()
//	_, ok := k.GetEVMAddress(ctx, eniAddr)
//	require.False(t, ok)
//	_, ok = k.GetEniAddress(ctx, evmAddr)
//	require.False(t, ok)
//	k.SetAddressMapping(ctx, eniAddr, evmAddr)
//	foundEVM, ok := k.GetEVMAddress(ctx, eniAddr)
//	require.True(t, ok)
//	require.Equal(t, evmAddr, foundEVM)
//	foundEni, ok := k.GetEniAddress(ctx, evmAddr)
//	require.True(t, ok)
//	require.Equal(t, eniAddr, foundEni)
//	require.Equal(t, eniAddr, k.AccountKeeper().GetAccount(ctx, eniAddr).GetAddress())
//}
//
//func TestDeleteAddressMapping(t *testing.T) {
//	k := &keeper.EVMTestApp.EvmKeeper
//	ctx := keeper.EVMTestApp.GetContextForDeliverTx([]byte{})
//	eniAddr, evmAddr := keeper.MockAddressPair()
//	k.SetAddressMapping(ctx, eniAddr, evmAddr)
//	foundEVM, ok := k.GetEVMAddress(ctx, eniAddr)
//	require.True(t, ok)
//	require.Equal(t, evmAddr, foundEVM)
//	foundEni, ok := k.GetEniAddress(ctx, evmAddr)
//	require.True(t, ok)
//	require.Equal(t, eniAddr, foundEni)
//	k.DeleteAddressMapping(ctx, eniAddr, evmAddr)
//	_, ok = k.GetEVMAddress(ctx, eniAddr)
//	require.False(t, ok)
//	_, ok = k.GetEniAddress(ctx, evmAddr)
//	require.False(t, ok)
//}
//
//func TestGetAddressOrDefault(t *testing.T) {
//	k := &keeper.EVMTestApp.EvmKeeper
//	ctx := keeper.EVMTestApp.GetContextForDeliverTx([]byte{})
//	eniAddr, evmAddr := keeper.MockAddressPair()
//	defaultEvmAddr := k.GetEVMAddressOrDefault(ctx, eniAddr)
//	require.True(t, bytes.Equal(eniAddr, defaultEvmAddr[:]))
//	defaultEniAddr := k.GetEniAddressOrDefault(ctx, evmAddr)
//	require.True(t, bytes.Equal(defaultEniAddr, evmAddr[:]))
//}
//
//func TestSendingToCastAddress(t *testing.T) {
//	a := keeper.EVMTestApp
//	ctx := a.GetContextForDeliverTx([]byte{})
//	eniAddr, evmAddr := keeper.MockAddressPair()
//	castAddr := sdk.AccAddress(evmAddr[:])
//	sourceAddr, _ := keeper.MockAddressPair()
//	require.Nil(t, a.BankKeeper.MintCoins(ctx, "evm", sdk.NewCoins(sdk.NewCoin("ueni", sdk.NewInt(10)))))
//	require.Nil(t, a.BankKeeper.SendCoinsFromModuleToAccount(ctx, "evm", sourceAddr, sdk.NewCoins(sdk.NewCoin("ueni", sdk.NewInt(5)))))
//	amt := sdk.NewCoins(sdk.NewCoin("ueni", sdk.NewInt(1)))
//	require.Nil(t, a.BankKeeper.SendCoinsFromModuleToAccount(ctx, "evm", castAddr, amt))
//	require.Nil(t, a.BankKeeper.SendCoins(ctx, sourceAddr, castAddr, amt))
//	require.Nil(t, a.BankKeeper.SendCoinsAndWei(ctx, sourceAddr, castAddr, sdk.OneInt(), sdk.OneInt()))
//
//	a.EvmKeeper.SetAddressMapping(ctx, eniAddr, evmAddr)
//	require.NotNil(t, a.BankKeeper.SendCoinsFromModuleToAccount(ctx, "evm", castAddr, amt))
//	require.NotNil(t, a.BankKeeper.SendCoins(ctx, sourceAddr, castAddr, amt))
//	require.NotNil(t, a.BankKeeper.SendCoinsAndWei(ctx, sourceAddr, castAddr, sdk.OneInt(), sdk.OneInt()))
//}
//
//func TestEvmAddressHandler_GetEniAddressFromString(t *testing.T) {
//	a := keeper.EVMTestApp
//	ctx := a.GetContextForDeliverTx([]byte{})
//	eniAddr, evmAddr := keeper.MockAddressPair()
//	a.EvmKeeper.SetAddressMapping(ctx, eniAddr, evmAddr)
//
//	_, notAssociatedEvmAddr := keeper.MockAddressPair()
//	castAddr := sdk.AccAddress(notAssociatedEvmAddr[:])
//
//	type args struct {
//		ctx     sdk.Context
//		address string
//	}
//	tests := []struct {
//		name       string
//		args       args
//		want       sdk.AccAddress
//		wantErr    bool
//		wantErrMsg string
//	}{
//		{
//			name: "returns associated Eni address if input address is a valid 0x and associated",
//			args: args{
//				ctx:     ctx,
//				address: evmAddr.String(),
//			},
//			want: eniAddr,
//		},
//		{
//			name: "returns default Eni address if input address is a valid 0x not associated",
//			args: args{
//				ctx:     ctx,
//				address: notAssociatedEvmAddr.String(),
//			},
//			want: castAddr,
//		},
//		{
//			name: "returns Eni address if input address is a valid bech32 address",
//			args: args{
//				ctx:     ctx,
//				address: eniAddr.String(),
//			},
//			want: eniAddr,
//		},
//		{
//			name: "returns error if address is invalid",
//			args: args{
//				ctx:     ctx,
//				address: "invalid",
//			},
//			wantErr:    true,
//			wantErrMsg: "decoding bech32 failed: invalid bech32 string length 7",
//		}, {
//			name: "returns error if address is empty",
//			args: args{
//				ctx:     ctx,
//				address: "",
//			},
//			wantErr:    true,
//			wantErrMsg: "empty address string is not allowed",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := evmkeeper.NewEvmAddressHandler(&a.EvmKeeper)
//			got, err := h.GetEniAddressFromString(tt.args.ctx, tt.args.address)
//			if tt.wantErr {
//				require.NotNil(t, err)
//				require.Equal(t, tt.wantErrMsg, err.Error())
//				return
//			} else {
//				require.NoError(t, err)
//				require.Equal(t, tt.want, got)
//			}
//		})
//	}
//}
