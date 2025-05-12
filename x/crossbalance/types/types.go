package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CrossBalanceTransfer defines the cross-account transfer message
type CrossBalanceTransfer struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Amount      sdk.Coins      `json:"amount"`
}

// NewCrossBalanceTransfer creates a new cross-account transfer message
func NewCrossBalanceTransfer(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins) CrossBalanceTransfer {
	return CrossBalanceTransfer{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
}

// Validate validates the cross-account transfer message
func (m CrossBalanceTransfer) Validate() error {
	if m.FromAddress.Empty() {
		return ErrEmptyFromAddress
	}
	if m.ToAddress.Empty() {
		return ErrEmptyToAddress
	}
	if m.Amount.IsZero() {
		return ErrZeroAmount
	}
	return nil
}
