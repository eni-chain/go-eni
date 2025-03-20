package keeper

import (
	cosmossdk_io_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/evm/types"
)

// modified eip-1559 adjustment using target gas used
func (k *Keeper) AdjustDynamicBaseFeePerGas(ctx sdk.Context, blockGasUsed uint64) *cosmossdk_io_math.LegacyDec {
	if ctx.ConsensusParams().Block == nil {
		return nil
	} //TODO remove:ctx.ConsensusParams() == nil ||
	prevBaseFee := k.GetNextBaseFeePerGas(ctx)
	// set the resulting base fee for block n-1 on block n
	k.SetCurrBaseFeePerGas(ctx, prevBaseFee)
	targetGasUsed := cosmossdk_io_math.LegacyNewDec(int64(k.GetTargetGasUsedPerBlock(ctx)))
	if targetGasUsed.IsZero() { // avoid division by zero
		return &prevBaseFee // return the previous base fee as is
	}
	minimumFeePerGas := k.GetParams(ctx).MinimumFeePerGas
	maximumFeePerGas := k.GetParams(ctx).MaximumFeePerGas
	blockGasLimit := cosmossdk_io_math.LegacyNewDec(ctx.ConsensusParams().Block.MaxGas)
	blockGasUsedDec := cosmossdk_io_math.LegacyNewDec(int64(blockGasUsed))

	// cap block gas used to block gas limit
	if blockGasUsedDec.GT(blockGasLimit) {
		blockGasUsedDec = blockGasLimit
	}

	var newBaseFee cosmossdk_io_math.LegacyDec
	if blockGasUsedDec.GT(targetGasUsed) {
		// upward adjustment
		numerator := blockGasUsedDec.Sub(targetGasUsed)
		denominator := blockGasLimit.Sub(targetGasUsed)
		percentageFull := numerator.Quo(denominator)
		adjustmentFactor := k.GetMaxDynamicBaseFeeUpwardAdjustment(ctx).Mul(percentageFull)
		newBaseFee = prevBaseFee.Mul(cosmossdk_io_math.LegacyNewDec(1).Add(adjustmentFactor))
	} else {
		// downward adjustment
		numerator := targetGasUsed.Sub(blockGasUsedDec)
		denominator := targetGasUsed
		percentageEmpty := numerator.Quo(denominator)
		adjustmentFactor := k.GetMaxDynamicBaseFeeDownwardAdjustment(ctx).Mul(percentageEmpty)
		newBaseFee = prevBaseFee.Mul(cosmossdk_io_math.LegacyNewDec(1).Sub(adjustmentFactor))
	}

	// Ensure the new base fee is not lower than the minimum fee
	if newBaseFee.LT(minimumFeePerGas) {
		newBaseFee = minimumFeePerGas
	}

	// Ensure the new base fee is not higher than the maximum fee
	if newBaseFee.GT(maximumFeePerGas) {
		newBaseFee = maximumFeePerGas
	}

	// Set the new base fee for the next height
	k.SetNextBaseFeePerGas(ctx, newBaseFee)

	return &newBaseFee
}

// dont have height be a prefix, just store the current base fee directly
func (k *Keeper) GetCurrBaseFeePerGas(ctx sdk.Context) cosmossdk_io_math.LegacyDec {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BaseFeePerGasPrefix)
	if bz == nil {
		minFeePerGas := k.GetMinimumFeePerGas(ctx)
		if minFeePerGas.IsNil() {
			minFeePerGas = types.DefaultParams().MinimumFeePerGas
		}
		return minFeePerGas
	}
	d := cosmossdk_io_math.LegacyDec{}
	err := d.UnmarshalJSON(bz)
	if err != nil {
		panic(err)
	}
	return d
}

func (k *Keeper) SetCurrBaseFeePerGas(ctx sdk.Context, baseFeePerGas cosmossdk_io_math.LegacyDec) {
	store := ctx.KVStore(k.storeKey)
	bz, err := baseFeePerGas.MarshalJSON()
	if err != nil {
		panic(err)
	}
	store.Set(types.BaseFeePerGasPrefix, bz)
}

func (k *Keeper) SetNextBaseFeePerGas(ctx sdk.Context, baseFeePerGas cosmossdk_io_math.LegacyDec) {
	store := ctx.KVStore(k.storeKey)
	bz, err := baseFeePerGas.MarshalJSON()
	if err != nil {
		panic(err)
	}
	store.Set(types.NextBaseFeePerGasPrefix, bz)
}

func (k *Keeper) GetNextBaseFeePerGas(ctx sdk.Context) cosmossdk_io_math.LegacyDec {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextBaseFeePerGasPrefix)
	if bz == nil {
		minFeePerGas := k.GetMinimumFeePerGas(ctx)
		if minFeePerGas.IsNil() || minFeePerGas.IsZero() {
			minFeePerGas = types.DefaultParams().MinimumFeePerGas
		}
		return minFeePerGas
	}
	d := cosmossdk_io_math.LegacyDec{}
	err := d.UnmarshalJSON(bz)
	if err != nil {
		panic(err)
	}
	return d
}
