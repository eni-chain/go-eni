package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/eni-chain/go-eni/x/evm/keeper"
	"github.com/eni-chain/go-eni/x/evm/types"
)

// This migration is to fix total supply mismatch caused by mishandled
// ante surplus
func FixTotalSupply(ctx sdk.Context, k *keeper.Keeper) error {
	balances := k.BankKeeper().GetAccountsBalances(ctx)
	correctSupply := sdk.ZeroInt()
	for _, balance := range balances {
		correctSupply = correctSupply.Add(balance.Coins.AmountOf(sdk.MustGetBaseDenom()))
	}
	totalWeiBalance := sdk.ZeroInt()
	k.BankKeeper().IterateAllWeiBalances(ctx, func(aa sdk.AccAddress, i sdk.Int) bool {
		totalWeiBalance = totalWeiBalance.Add(i)
		return false
	})
	weiInUeni, weiRemainder := bankkeeper.SplitUeniWeiAmount(totalWeiBalance)
	if !weiRemainder.IsZero() {
		ctx.Logger().Error("wei total supply has been compromised as well; rounding up and adding to reserve")
		if err := k.BankKeeper().AddWei(ctx, k.AccountKeeper().GetModuleAddress(types.ModuleName), bankkeeper.OneUeniInWei.Sub(weiRemainder)); err != nil {
			return err
		}
		weiInUeni = weiInUeni.Add(sdk.OneInt())
	}
	correctSupply = correctSupply.Add(weiInUeni)
	currentSupply := k.BankKeeper().GetSupply(ctx, sdk.MustGetBaseDenom()).Amount
	if !currentSupply.Equal(correctSupply) {
		k.BankKeeper().SetSupply(ctx, sdk.NewCoin(sdk.MustGetBaseDenom(), correctSupply))
	}
	return nil
}
