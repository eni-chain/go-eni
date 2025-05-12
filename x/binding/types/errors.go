package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/binding module sentinel errors
var (
	ErrInvalidSigner      = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample             = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrEmptyCosmosAddress = sdkerrors.Register(ModuleName, 1, "cosmos address cannot be empty")
	ErrEmptyEvmAddress    = sdkerrors.Register(ModuleName, 2, "evm address cannot be empty")
	ErrBindingNotFound    = sdkerrors.Register(ModuleName, 3, "binding not found")
	ErrBindingExists      = sdkerrors.Register(ModuleName, 4, "binding already exists")
)
