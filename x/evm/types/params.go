package types

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyPriorityNormalizer                  = []byte("KeyPriorityNormalizer")
	KeyMinFeePerGas                        = []byte("KeyMinFeePerGas")
	KeyMaxFeePerGas                        = []byte("KeyMaximumFeePerGas")
	KeyMaxDynamicBaseFeeUpwardAdjustment   = []byte("KeyMaxDynamicBaseFeeUpwardAdjustment")
	KeyMaxDynamicBaseFeeDownwardAdjustment = []byte("KeyMaxDynamicBaseFeeDownwardAdjustment")
	KeyTargetGasUsedPerBlock               = []byte("KeyTargetGasUsedPerBlock")
	// deprecated
	KeyBaseFeePerGas = []byte("KeyBaseFeePerGas")
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

var DefaultPriorityNormalizer = math.LegacyNewDec(1)

// DefaultBaseFeePerGas determines how much ueni per gas spent is
// burnt rather than go to validators (similar to base fee on
// Ethereum).
var DefaultBaseFeePerGas = math.LegacyNewDec(0)         // used for static base fee, deprecated in favor of dynamic base fee
var DefaultMinFeePerGas = math.LegacyNewDec(1000000000) // 1gwei

var DefaultMaxDynamicBaseFeeUpwardAdjustment = math.LegacyNewDecWithPrec(189, 4)  // 1.89%
var DefaultMaxDynamicBaseFeeDownwardAdjustment = math.LegacyNewDecWithPrec(39, 4) // .39%
var DefaultTargetGasUsedPerBlock = uint64(250000)                                 // 250k
var DefaultMaxFeePerGas = math.LegacyNewDec(1000000000000)                        // 1,000gwei

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		PriorityNormalizer:                  DefaultPriorityNormalizer,
		BaseFeePerGas:                       DefaultBaseFeePerGas,
		MinimumFeePerGas:                    DefaultMinFeePerGas,
		MaxDynamicBaseFeeUpwardAdjustment:   DefaultMaxDynamicBaseFeeUpwardAdjustment,
		MaxDynamicBaseFeeDownwardAdjustment: DefaultMaxDynamicBaseFeeDownwardAdjustment,
		TargetGasUsedPerBlock:               DefaultTargetGasUsedPerBlock,
		MaximumFeePerGas:                    DefaultMaxFeePerGas,
		InitEniAmount:                       "",
		InitEniAddress:                      "",
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPriorityNormalizer, &p.PriorityNormalizer, validatePriorityNormalizer),
		paramtypes.NewParamSetPair(KeyBaseFeePerGas, &p.BaseFeePerGas, validateBaseFeePerGas),
		paramtypes.NewParamSetPair(KeyMaxDynamicBaseFeeUpwardAdjustment, &p.MaxDynamicBaseFeeUpwardAdjustment, validateBaseFeeAdjustment),
		paramtypes.NewParamSetPair(KeyMaxDynamicBaseFeeDownwardAdjustment, &p.MaxDynamicBaseFeeDownwardAdjustment, validateBaseFeeAdjustment),
		paramtypes.NewParamSetPair(KeyMinFeePerGas, &p.MinimumFeePerGas, validateMinFeePerGas),
		paramtypes.NewParamSetPair(KeyTargetGasUsedPerBlock, &p.TargetGasUsedPerBlock, func(i interface{}) error { return nil }),
		paramtypes.NewParamSetPair(KeyMaxFeePerGas, &p.MaximumFeePerGas, validateMaxFeePerGas),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validatePriorityNormalizer(p.PriorityNormalizer); err != nil {
		return err
	}
	if err := validateBaseFeePerGas(p.BaseFeePerGas); err != nil {
		return err
	}
	if err := validateMinFeePerGas(p.MinimumFeePerGas); err != nil {
		return err
	}
	if err := validateMaxFeePerGas(p.MaximumFeePerGas); err != nil {
		return err
	}
	if p.MinimumFeePerGas.LT(p.BaseFeePerGas) {
		return errors.New("minimum fee cannot be lower than base fee")
	}
	if err := validateBaseFeeAdjustment(p.MaxDynamicBaseFeeUpwardAdjustment); err != nil {
		return fmt.Errorf("invalid max dynamic base fee upward adjustment: %s, err: %s", p.MaxDynamicBaseFeeUpwardAdjustment, err)
	}
	if err := validateBaseFeeAdjustment(p.MaxDynamicBaseFeeDownwardAdjustment); err != nil {
		return fmt.Errorf("invalid max dynamic base fee downward adjustment: %s, err: %s", p.MaxDynamicBaseFeeDownwardAdjustment, err)
	}
	return nil
}

func validatePriorityNormalizer(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("nonpositive priority normalizer: %d", v)
	}

	return nil
}

func validateBaseFeePerGas(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("negative base fee per gas: %d", v)
	}

	return nil
}
func validateBaseFeeAdjustment(i interface{}) error {
	adjustment, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if adjustment.IsNegative() {
		return fmt.Errorf("negative base fee adjustment: %s", adjustment)
	}
	if adjustment.GT(math.LegacyOneDec()) {
		return fmt.Errorf("base fee adjustment must be less than or equal to 1: %s", adjustment)
	}
	return nil
}
func validateMinFeePerGas(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("negative min fee per gas: %d", v)
	}

	return nil
}

func validateMaxFeePerGas(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("negative max fee per gas: %d", v)
	}
	return nil
}

func validateDeliverTxHookWasmGasLimit(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("invalid deliver_tx_hook_wasm_gas_limit: must be greater than 0, got %d", v)
	}
	return nil
}

func validateWhitelistedCwHashesForDelegateCall(i interface{}) error {
	_, ok := i.([][]byte)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
