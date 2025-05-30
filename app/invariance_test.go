package app_test

//import (
//	"fmt"
//	"testing"
//	"time"
//
//	"cosmossdk.io/math"
//	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/cosmos/cosmos-sdk/x/evm/types"
//	app "github.com/eni-chain/go-eni/app"
//	"github.com/stretchr/testify/require"
//)
//
//func TestLightInvarianceChecks(t *testing.T) {
//	tm := time.Now().UTC()
//	valPub := secp256k1.GenPrivKey().PubKey()
//	accounts := []sdk.AccAddress{
//		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
//		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
//	}
//	ueniCoin := func(i int64) sdk.Coin { return sdk.NewCoin("ueni", math.NewInt(i)) }
//	ueniCoins := func(i int64) sdk.Coins { return sdk.NewCoins(ueniCoin(i)) }
//	for i, tt := range []struct {
//		preUeni    []int64
//		preSupply  int64
//		postUeni   []int64
//		postWei    []int64
//		postSupply int64
//		success    bool
//	}{
//		{
//			preUeni:    []int64{0, 0},
//			preSupply:  5,
//			postUeni:   []int64{1, 2},
//			postWei:    []int64{0, 0},
//			postSupply: 8,
//			success:    true,
//		},
//		{
//			preUeni:    []int64{2, 1},
//			preSupply:  3,
//			postUeni:   []int64{0, 0},
//			postWei:    []int64{0, 0},
//			postSupply: 0,
//			success:    true,
//		},
//		{
//			preUeni:    []int64{1, 0},
//			preSupply:  10,
//			postUeni:   []int64{0, 1},
//			postWei:    []int64{0, 0},
//			postSupply: 10,
//			success:    true,
//		},
//		{
//			preUeni:    []int64{1, 0},
//			preSupply:  10,
//			postUeni:   []int64{0, 0},
//			postWei:    []int64{500_000_000_000, 500_000_000_000},
//			postSupply: 10,
//			success:    true,
//		},
//		{
//			preUeni:    []int64{0, 0},
//			preSupply:  10,
//			postUeni:   []int64{1, 0},
//			postWei:    []int64{0, 0},
//			postSupply: 10,
//			success:    true,
//		},
//		{
//			preUeni:    []int64{0, 0},
//			preSupply:  10,
//			postUeni:   []int64{0, 0},
//			postWei:    []int64{2, 1},
//			postSupply: 10,
//			success:    true,
//		},
//		{
//			preUeni:    []int64{1, 0},
//			preSupply:  10,
//			postUeni:   []int64{1, 1},
//			postWei:    []int64{0, 0},
//			postSupply: 10,
//			success:    false,
//		},
//		{
//			preUeni:    []int64{1, 0},
//			preSupply:  10,
//			postUeni:   []int64{0, 0},
//			postWei:    []int64{0, 0},
//			postSupply: 10,
//			success:    false,
//		},
//		{
//			preUeni:    []int64{1, 0},
//			preSupply:  10,
//			postUeni:   []int64{0, 1},
//			postWei:    []int64{500_000_000_000, 500_000_000_000},
//			postSupply: 10,
//			success:    false,
//		},
//		{
//			preUeni:    []int64{1, 0},
//			preSupply:  10,
//			postUeni:   []int64{0, 1},
//			postWei:    []int64{0, 0},
//			postSupply: 10,
//			success:    false,
//		},
//		{
//			preUeni:    []int64{0, 0},
//			preSupply:  10,
//			postUeni:   []int64{0, 0},
//			postWei:    []int64{2, 2},
//			postSupply: 10,
//			success:    false,
//		},
//		{
//			preUeni:    []int64{0, 0},
//			preSupply:  10,
//			postUeni:   []int64{0, 0},
//			postWei:    []int64{1, 1},
//			postSupply: 10,
//			success:    false,
//		},
//	} {
//		fmt.Printf("Running test %d\n", i)
//		testWrapper := app.NewTestWrapperWithSc(t, tm, valPub, false)
//		a, ctx := testWrapper.App, testWrapper.Ctx
//		for i := range tt.preUeni {
//			if tt.preUeni[i] > 0 {
//				a.BankKeeper.AddCoins(ctx, accounts[i], ueniCoins(tt.preUeni[i]))
//			}
//		}
//		if tt.preSupply > 0 {
//			a.BankKeeper.MintCoins(ctx, types.ModuleName, ueniCoins(tt.preSupply))
//		}
//		a.Commit()
//		for i := range tt.postUeni {
//			ueniDiff := tt.postUeni[i] - tt.preUeni[i]
//			if ueniDiff > 0 {
//				a.BankKeeper.AddCoins(ctx, accounts[i], ueniCoins(ueniDiff))
//			} else if ueniDiff < 0 {
//				a.BankKeeper.SubUnlockedCoins(ctx, accounts[i], ueniCoins(-ueniDiff))
//			}
//
//		}
//		a.BankKeeper.MintCoins(ctx, types.ModuleName, ueniCoins(tt.postSupply))
//		a.SetDeliverStateToCommit()
//		f := func() { a.LightInvarianceChecks(a.WriteState(), app.LightInvarianceConfig{SupplyEnabled: true}) }
//		if tt.success {
//			require.NotPanics(t, f)
//		} else {
//			require.Panics(t, f)
//		}
//		safeClose(a)
//	}
//}
//
//// TODO: remove once snapshot manager can be closed gracefully in tests
//func safeClose(a *app.App) {
//	defer func() {
//		_ = recover()
//	}()
//	a.Close()
//}
