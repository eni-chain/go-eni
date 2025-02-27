package app

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	//"github.com/cosmos/cosmos-sdk/storev2/commitment"
	"github.com/spf13/cast"
)

type LightInvarianceConfig struct {
	SupplyEnabled bool `mapstructure:"supply_enabled"`
}

var DefaultLightInvarianceConfig = LightInvarianceConfig{
	SupplyEnabled: true,
}

const (
	flagSupplyEnabled = "light_invariance.supply_enabled"
)

func ReadLightInvarianceConfig(opts servertypes.AppOptions) (LightInvarianceConfig, error) {
	cfg := DefaultLightInvarianceConfig // copy
	var err error
	if v := opts.Get(flagSupplyEnabled); v != nil {
		if cfg.SupplyEnabled, err = cast.ToBoolE(v); err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}

//
//func (app *App) LightInvarianceChecks(cms sdk.CommitMultiStore, config LightInvarianceConfig) {
//	if config.SupplyEnabled {
//		app.LightInvarianceTotalSupply(cms)
//	}
//}
//
//func (app *App) LightInvarianceTotalSupply(cms sdk.CommitMultiStore) {
//	defer metrics.MeasureSince(
//		[]string{"eni", "lightinvariance_supply", "milliseconds"},
//		time.Now().UTC(),
//	)
//	ckv, ok := cms.GetStore(app.BankKeeper.GetStoreKey()).(*commitment.Store)
//	if !ok {
//		app.Logger().Error("bank store is not a memiavl store; cannot run light invariance check")
//		return
//	}
//	balanceChangePairs := ckv.GetChangedPairs(banktypes.BalancesPrefix)
//	ueniPostTotal := sdk.ZeroInt()
//	ueniChangedAddr := []sdk.AccAddress{}
//	for _, p := range balanceChangePairs {
//		if len(p.Key) < 2 {
//			// invalid key; ignore
//			metrics.IncrCounterWithLabels([]string{"eni", "lightinvariance_supply", "invalid_changed_key"}, 1, []metrics.Label{
//				{
//					Name:  "type",
//					Value: "eni",
//				},
//			})
//			app.Logger().Error(fmt.Sprintf("invalid changed pair key for ueni: %X", p.Key))
//			continue
//		}
//		addrLen := int(p.Key[1])
//		if len(p.Key) < addrLen+2 {
//			// invalid key length; ignore
//			metrics.IncrCounterWithLabels([]string{"eni", "lightinvariance_supply", "invalid_changed_key"}, 1, []metrics.Label{
//				{
//					Name:  "type",
//					Value: "eni",
//				},
//			})
//			app.Logger().Error(fmt.Sprintf("invalid changed pair key for ueni: %X", p.Key))
//			continue
//		}
//		addr := p.Key[2 : addrLen+2]
//		denom := p.Key[addrLen+2:]
//		if string(denom) != sdk.MustGetBaseDenom() {
//			continue
//		}
//		if !p.Delete {
//			var balance sdk.Coin
//			if err := balance.Unmarshal(p.Value); err != nil {
//				metrics.IncrCounterWithLabels([]string{"eni", "lightinvariance_supply", "unmarshal_failure"}, 1, []metrics.Label{
//					{
//						Name:  "type",
//						Value: "ueni",
//					}, {
//						Name:  "step",
//						Value: "post_block",
//					},
//				})
//				app.Logger().Error(fmt.Sprintf("failed to unmarshal balance: %s", err))
//				continue
//			}
//			if balance.Amount.IsNegative() {
//				panic(fmt.Sprintf("negative balance found for addr %s: %s", sdk.AccAddress(addr).String(), balance.String()))
//			}
//			ueniPostTotal = ueniPostTotal.Add(balance.Amount)
//		}
//		ueniChangedAddr = append(ueniChangedAddr, addr)
//	}
//	ueniPreTotal := sdk.ZeroInt()
//	for _, a := range ueniChangedAddr {
//		key := append(banktypes.CreateAccountBalancesPrefix(a), []byte(sdk.MustGetBaseDenom())...)
//		val := ckv.Get(key)
//		if val == nil {
//			continue
//		}
//		var balance sdk.Coin
//		if err := balance.Unmarshal(val); err != nil {
//			metrics.IncrCounterWithLabels([]string{"eni", "lightinvariance_supply", "unmarshal_failure"}, 1, []metrics.Label{
//				{
//					Name:  "type",
//					Value: "ueni",
//				}, {
//					Name:  "step",
//					Value: "pre_block",
//				},
//			})
//			app.Logger().Error(fmt.Sprintf("failed to unmarshal preblock balance: %s", err))
//			continue
//		}
//		ueniPreTotal = ueniPreTotal.Add(balance.Amount)
//	}
//	weiChangePairs := ckv.GetChangedPairs(banktypes.WeiBalancesPrefix)
//	weiPostTotal := sdk.ZeroInt()
//	weiChangedAddrs := []sdk.AccAddress{}
//	for _, p := range weiChangePairs {
//		var amt sdk.Int
//		if len(p.Key) < 1 {
//			metrics.IncrCounterWithLabels([]string{"eni", "lightinvariance_supply", "invalid_changed_key"}, 1, []metrics.Label{
//				{
//					Name:  "type",
//					Value: "wei",
//				},
//			})
//			app.Logger().Error(fmt.Sprintf("invalid changed pair key: %X", p.Key))
//			continue
//		}
//		if !p.Delete {
//			if err := amt.Unmarshal(p.Value); err != nil {
//				metrics.IncrCounterWithLabels([]string{"eni", "lightinvariance_supply", "unmarshal_failure"}, 1, []metrics.Label{
//					{
//						Name:  "type",
//						Value: "wei",
//					}, {
//						Name:  "step",
//						Value: "post_block",
//					},
//				})
//				app.Logger().Error(fmt.Sprintf("failed to unmarshal wei balance: %s", err))
//				continue
//			}
//			weiPostTotal = weiPostTotal.Add(amt)
//			if amt.IsNegative() {
//				panic(fmt.Sprintf("negative wei balance found for addr %s: %s", sdk.AccAddress(p.Key[1:]).String(), amt.String()))
//			}
//		}
//		weiChangedAddrs = append(weiChangedAddrs, p.Key[1:])
//	}
//	weiPreTotal := sdk.ZeroInt()
//	for _, a := range weiChangedAddrs {
//		key := append(banktypes.WeiBalancesPrefix, a...)
//		val := ckv.Get(key)
//		if val == nil {
//			continue
//		}
//		var amt sdk.Int
//		if err := amt.Unmarshal(val); err != nil {
//			metrics.IncrCounterWithLabels([]string{"eni", "lightinvariance_supply", "unmarshal_failure"}, 1, []metrics.Label{
//				{
//					Name:  "type",
//					Value: "wei",
//				}, {
//					Name:  "step",
//					Value: "pre_block",
//				},
//			})
//			app.Logger().Error(fmt.Sprintf("failed to unmarshal preblock wei balance: %s", err))
//			continue
//		}
//		weiPreTotal = weiPreTotal.Add(amt)
//	}
//	totalSupplyChangePairs := ckv.GetChangedPairs(banktypes.SupplyKey)
//	supplyChanged := sdk.ZeroInt()
//	preTotalSupply := sdk.ZeroInt()
//	if bz := ckv.Get(append(banktypes.SupplyKey, []byte(sdk.MustGetBaseDenom())...)); bz != nil {
//		var amt sdk.Int
//		if err := amt.Unmarshal(bz); err != nil {
//			metrics.IncrCounterWithLabels([]string{"eni", "lightinvariance_supply", "unmarshal_failure"}, 1, []metrics.Label{
//				{
//					Name:  "type",
//					Value: "total_supply",
//				}, {
//					Name:  "step",
//					Value: "pre_block",
//				},
//			})
//			app.Logger().Error(fmt.Sprintf("failed to unmarshal pre total supply: %s", err))
//			return
//		}
//		preTotalSupply = amt
//	}
//	for _, p := range totalSupplyChangePairs {
//		if string(p.Key[1:]) == sdk.MustGetBaseDenom() {
//			if p.Delete {
//				supplyChanged = preTotalSupply.Neg()
//			} else {
//				var amt sdk.Int
//				if err := amt.Unmarshal(p.Value); err != nil {
//					metrics.IncrCounterWithLabels([]string{"eni", "lightinvariance_supply", "unmarshal_failure"}, 1, []metrics.Label{
//						{
//							Name:  "type",
//							Value: "total_supply",
//						}, {
//							Name:  "step",
//							Value: "post_block",
//						},
//					})
//					app.Logger().Error(fmt.Sprintf("failed to unmarshal total supply: %s", err))
//				} else {
//					supplyChanged = amt.Sub(preTotalSupply)
//				}
//			}
//			break
//		}
//	}
//	weiDiff := weiPostTotal.Sub(weiPreTotal)
//	weiDiffInUeni, weiDiffRemainder := bankkeeper.SplitUeniWeiAmount(weiDiff)
//	if !weiDiffRemainder.IsZero() {
//		panic(fmt.Sprintf("non-zero wei diff found! Pre-block wei total %s, post-block wei total %s", weiPreTotal, weiPostTotal))
//	}
//	ueniDiff := ueniPreTotal.Sub(ueniPostTotal).Sub(weiDiffInUeni).Add(supplyChanged)
//	if !ueniDiff.IsZero() {
//		panic(fmt.Sprintf("unexpected ueni balance total found! Pre-block ueni total %s wei total %s total supply %s, post-block ueni total %s wei total %s total supply %s", ueniPreTotal, weiPreTotal, preTotalSupply, ueniPostTotal, weiPostTotal, preTotalSupply.Add(supplyChanged)))
//	}
//	app.Logger().Info("successfully verified supply light invariance")
//}
