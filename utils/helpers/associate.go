package helpers

//type AssociationHelper struct {
//	evmKeeper     pcommon.EVMKeeper
//	bankKeeper    pcommon.BankKeeper
//	accountKeeper pcommon.AccountKeeper
//}
//
//func NewAssociationHelper(evmKeeper pcommon.EVMKeeper, bankKeeper pcommon.BankKeeper, accountKeeper pcommon.AccountKeeper) *AssociationHelper {
//	return &AssociationHelper{evmKeeper: evmKeeper, bankKeeper: bankKeeper, accountKeeper: accountKeeper}
//}
//
//func (p AssociationHelper) AssociateAddresses(ctx sdk.Context, eniAddr sdk.AccAddress, evmAddr common.Address, pubkey cryptotypes.PubKey) error {
//	p.evmKeeper.SetAddressMapping(ctx, eniAddr, evmAddr)
//	if acc := p.accountKeeper.GetAccount(ctx, eniAddr); acc.GetPubKey() == nil {
//		if err := acc.SetPubKey(pubkey); err != nil {
//			return err
//		}
//		p.accountKeeper.SetAccount(ctx, acc)
//	}
//	return p.MigrateBalance(ctx, evmAddr, eniAddr)
//}
//
//func (p AssociationHelper) MigrateBalance(ctx sdk.Context, evmAddr common.Address, eniAddr sdk.AccAddress) error {
//	castAddr := sdk.AccAddress(evmAddr[:])
//	castAddrBalances := p.bankKeeper.SpendableCoins(ctx, castAddr)
//	if !castAddrBalances.IsZero() {
//		if err := p.bankKeeper.SendCoins(ctx, castAddr, eniAddr, castAddrBalances); err != nil {
//			return err
//		}
//	}
//	castAddrWei := p.bankKeeper.GetWeiBalance(ctx, castAddr)
//	if !castAddrWei.IsZero() {
//		if err := p.bankKeeper.SendCoinsAndWei(ctx, castAddr, eniAddr, sdk.ZeroInt(), castAddrWei); err != nil {
//			return err
//		}
//	}
//	if p.bankKeeper.LockedCoins(ctx, castAddr).IsZero() {
//		p.accountKeeper.RemoveAccount(ctx, authtypes.NewBaseAccountWithAddress(castAddr))
//	}
//	return nil
//}
