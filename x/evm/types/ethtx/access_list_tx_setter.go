package ethtx

import math "cosmossdk.io/math"

func (tx *AccessListTx) SetTo(v string) {
	tx.To = v
}

func (tx *AccessListTx) SetAmount(v math.Int) {
	tx.Amount = &v
}

func (tx *AccessListTx) SetGasPrice(v math.Int) {
	tx.GasPrice = &v
}

func (tx *AccessListTx) SetAccesses(v AccessList) {
	tx.Accesses = v
}
