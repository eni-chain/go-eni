package ethtx

import sdk "cosmossdk.io/math"

func (tx *LegacyTx) SetTo(v string) {
	tx.To = v
}

func (tx *LegacyTx) SetAmount(v sdk.Int) {
	tx.Amount = &v
}

func (tx *LegacyTx) SetGasPrice(v sdk.Int) {
	tx.GasPrice = &v
}
