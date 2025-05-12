package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewBinding creates a new binding relationship
func NewBinding(cosmosAddr sdk.AccAddress, evmAddr string, creator sdk.AccAddress) Binding {
	return Binding{
		Index:         cosmosAddr.String(),
		CosmosAddress: cosmosAddr.String(),
		EvmAddress:    evmAddr,
		Creator:       creator.String(),
	}
}

// Validate validates the binding relationship
func (b Binding) Validate() error {
	if b.CosmosAddress == "" {
		return ErrEmptyCosmosAddress
	}
	if b.EvmAddress == "" {
		return ErrEmptyEvmAddress
	}
	return nil
}
