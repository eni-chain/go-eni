package keeper

import (
	cosmossdk_io_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/evm/types"
)

// modified eip-1559 adjustment using target gas used
func (k *Keeper) AdjustDynamicBaseFeePerGas(ctx sdk.Context, blockGasUsed uint64) *cosmossdk_io_math.Dec {
	//if ctx.ConsensusParams().Block == nil || ctx.ConsensusParams().Block == nil {
	//	return nil
	//}
	//prevBaseFee := k.GetNextBaseFeePerGas(ctx)
	//// set the resulting base fee for block n-1 on block n
	//k.SetCurrBaseFeePerGas(ctx, prevBaseFee)
	//targetGasUsed := cosmossdk_io_math.NewDecFromInt64(int64(k.GetTargetGasUsedPerBlock(ctx)))
	//if targetGasUsed.IsZero() { // avoid division by zero
	//	return &prevBaseFee // return the previous base fee as is
	//}
	//minimumFeePerGas := k.GetParams(ctx).MinimumFeePerGas
	//maximumFeePerGas := k.GetParams(ctx).MaximumFeePerGas
	//blockGasLimit := cosmossdk_io_math.NewDecFromInt64(ctx.ConsensusParams().Block.MaxGas)
	//blockGasUsedDec := cosmossdk_io_math.NewDecFromInt64(int64(blockGasUsed))
	//
	//// cap block gas used to block gas limit
	//
	////if blockGasUsedDec.GT(blockGasLimit) {  //todo check if this is correct
	//if blockGasUsedDec.Cmp(blockGasLimit) > 0 {
	//	blockGasUsedDec = blockGasLimit
	//}

	//var newBaseFee cosmossdk_io_math.Dec
	////if blockGasUsedDec.GT(targetGasUsed) { //todo check if this is correct
	//if blockGasUsedDec.Cmp(targetGasUsed) > 0 {
	//	// upward adjustment
	//	numerator, _ := blockGasUsedDec.Sub(targetGasUsed)
	//	denominator, _ := blockGasLimit.Sub(targetGasUsed)
	//	percentageFull, _ := numerator.Quo(denominator)
	//	adjustmentFactor := k.GetMaxDynamicBaseFeeUpwardAdjustment(ctx).Mul(percentageFull)
	//	newBaseFee, _ = prevBaseFee.Mul(sdk.NewDec(1).Add(adjustmentFactor))
	//} else {
	//	// downward adjustment
	//	numerator, _ := targetGasUsed.Sub(blockGasUsedDec)
	//	denominator := targetGasUsed
	//	percentageEmpty, _ := numerator.Quo(denominator)
	//	adjustmentFactor := k.GetMaxDynamicBaseFeeDownwardAdjustment(ctx).Mul(percentageEmpty)
	//	newBaseFee = prevBaseFee.Mul(sdk.NewDec(1).Sub(adjustmentFactor))
	//}
	//
	//// Ensure the new base fee is not lower than the minimum fee
	//if newBaseFee.LT(minimumFeePerGas) {
	//	newBaseFee = minimumFeePerGas
	//}
	//
	//// Ensure the new base fee is not higher than the maximum fee
	//if newBaseFee.GT(maximumFeePerGas) {
	//	newBaseFee = maximumFeePerGas
	//}

	// Set the new base fee for the next height
	//k.SetNextBaseFeePerGas(ctx, newBaseFee)
	//
	//return &newBaseFee
	return nil
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

func (k *Keeper) SetCurrBaseFeePerGas(ctx sdk.Context, baseFeePerGas cosmossdk_io_math.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz, err := baseFeePerGas.MarshalJSON()
	if err != nil {
		panic(err)
	}
	store.Set(types.BaseFeePerGasPrefix, bz)
}

func (k *Keeper) SetNextBaseFeePerGas(ctx sdk.Context, baseFeePerGas cosmossdk_io_math.Dec) {
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
