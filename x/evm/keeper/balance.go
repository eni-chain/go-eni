package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/evm/state"
)

func (k *Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress) *big.Int {
	denom := k.GetBaseDenom(ctx)
	allUeni := k.BankKeeper().GetBalance(ctx, addr, denom).Amount
	lockedUeni := k.BankKeeper().LockedCoins(ctx, addr).AmountOf(denom) // LockedCoins doesn't use iterators
	ueni := allUeni.Sub(lockedUeni)
	//wei := k.BankKeeper().GetWeiBalance(ctx, addr)
	wei := k.bankKeeper.GetBalance(ctx, addr, denom)
	return ueni.Mul(state.SdkUeniToSweiMultiplier).Add(wei.Amount).BigInt()
}
