package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/crossbalance module sentinel errors
var (
	ErrInvalidSigner     = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample            = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrEmptyFromAddress  = sdkerrors.Register(ModuleName, 1, "from address cannot be empty")
	ErrEmptyToAddress    = sdkerrors.Register(ModuleName, 2, "to address cannot be empty")
	ErrZeroAmount        = sdkerrors.Register(ModuleName, 3, "amount cannot be zero")
	ErrInsufficientFunds = sdkerrors.Register(ModuleName, 4, "insufficient funds")
	ErrNoBinding         = sdkerrors.Register(ModuleName, 5, "no binding found for address")
)
