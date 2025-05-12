package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BindingKeeper defines expected keeper from the binding module
type BindingKeeper interface {
	HasBinding(ctx context.Context, cosmosAddr sdk.AccAddress) (bool, error)
	DeleteBinding(ctx context.Context, cosmosAddr sdk.AccAddress) error
}

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	HasBalance(ctx context.Context, addr sdk.AccAddress, amt sdk.Coin) bool
	AddCoins(ctx context.Context, addr sdk.AccAddress, amt sdk.Coins) error
	SubUnlockedCoins(ctx context.Context, addr sdk.AccAddress, amt sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
