package migrations

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/evm/keeper"
	"github.com/ethereum/go-ethereum/common"
)

func MigrateCastAddressBalances(ctx sdk.Context, k *keeper.Keeper) (rerr error) {
	k.IterateEniAddressMapping(ctx, func(evmAddr common.Address, eniAddr sdk.AccAddress) bool {
		castAddr := sdk.AccAddress(evmAddr[:])
		if balances := k.BankKeeper().SpendableCoins(ctx, castAddr); !balances.IsZero() {
			if err := k.BankKeeper().SendCoins(ctx, castAddr, eniAddr, balances); err != nil {
				ctx.Logger().Error(fmt.Sprintf("error migrating balances from cast address to real address for %s due to %s", evmAddr.Hex(), err))
				rerr = err
				return true
			}
		}
		if wei := k.BankKeeper().GetWeiBalance(ctx, castAddr); !wei.IsZero() {
			if err := k.BankKeeper().SendCoinsAndWei(ctx, castAddr, eniAddr, sdk.ZeroInt(), wei); err != nil {
				ctx.Logger().Error(fmt.Sprintf("error migrating wei from cast address to real address for %s due to %s", evmAddr.Hex(), err))
				rerr = err
				return true
			}
		}
		return false
	})
	return
}
