package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{}
}

var DefaultPriorityNormalizer = cosmossdk_io_math.LegacyNewDec(1)

// DefaultBaseFeePerGas determines how much ueni per gas spent is
// burnt rather than go to validators (similar to base fee on
// Ethereum).
var DefaultBaseFeePerGas = cosmossdk_io_math.LegacyNewDec(0)         // used for static base fee, deprecated in favor of dynamic base fee
var DefaultMinFeePerGas = cosmossdk_io_math.LegacyNewDec(1000000000) // 1gwei
var DefaultDeliverTxHookWasmGasLimit = uint64(300000)

var DefaultWhitelistedCwCodeHashesForDelegateCall = [][]byte(nil)

var DefaultMaxDynamicBaseFeeUpwardAdjustment = cosmossdk_io_math.LegacyNewDecWithPrec(189, 4)  // 1.89%
var DefaultMaxDynamicBaseFeeDownwardAdjustment = cosmossdk_io_math.LegacyNewDecWithPrec(39, 4) // .39%
var DefaultTargetGasUsedPerBlock = uint64(250000)                                              // 250k
var DefaultMaxFeePerGas = cosmossdk_io_math.LegacyNewDec(1000000000000)                        // 1,000gwei

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		PriorityNormalizer:                     DefaultPriorityNormalizer,
		BaseFeePerGas:                          DefaultBaseFeePerGas,
		MaxDynamicBaseFeeUpwardAdjustment:      DefaultMaxDynamicBaseFeeUpwardAdjustment,
		MaxDynamicBaseFeeDownwardAdjustment:    DefaultMaxDynamicBaseFeeDownwardAdjustment,
		MinimumFeePerGas:                       DefaultMinFeePerGas,
		DeliverTxHookWasmGasLimit:              DefaultDeliverTxHookWasmGasLimit,
		WhitelistedCwCodeHashesForDelegateCall: DefaultWhitelistedCwCodeHashesForDelegateCall,
		TargetGasUsedPerBlock:                  DefaultTargetGasUsedPerBlock,
		MaximumFeePerGas:                       DefaultMaxFeePerGas,
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}
