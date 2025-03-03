package keeper

import (
	cosmossdk_io_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"math/big"

	"github.com/eni-chain/go-eni/utils"
	"github.com/eni-chain/go-eni/x/evm/config"

	"github.com/eni-chain/go-eni/x/evm/types"
)

const BaseDenom = "ueni"

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	return nil
}

func (k *Keeper) GetParamsIfExists(ctx sdk.Context) types.Params {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := types.Params{}
	k.Paramstore.GetParamSetIfExists(sdkCtx, &params)
	return params
}

func (k *Keeper) GetBaseDenom(ctx sdk.Context) string {
	return BaseDenom
}

func (k *Keeper) GetPriorityNormalizer(ctx sdk.Context) cosmossdk_io_math.LegacyDec {
	return k.GetParams(ctx).PriorityNormalizer
}

func (k *Keeper) GetBaseFeePerGas(ctx sdk.Context) cosmossdk_io_math.LegacyDec {
	return k.GetParams(ctx).BaseFeePerGas
}

func (k *Keeper) GetMaxDynamicBaseFeeUpwardAdjustment(ctx sdk.Context) cosmossdk_io_math.LegacyDec {
	return k.GetParams(ctx).MaxDynamicBaseFeeUpwardAdjustment
}

func (k *Keeper) GetMaxDynamicBaseFeeDownwardAdjustment(ctx sdk.Context) cosmossdk_io_math.LegacyDec {
	return k.GetParams(ctx).MaxDynamicBaseFeeDownwardAdjustment
}

func (k *Keeper) GetMinimumFeePerGas(ctx sdk.Context) cosmossdk_io_math.LegacyDec {
	return k.GetParams(ctx).MinimumFeePerGas
}

func (k *Keeper) GetMaximumFeePerGas(ctx sdk.Context) cosmossdk_io_math.LegacyDec {
	return k.GetParams(ctx).MaximumFeePerGas
}

func (k *Keeper) GetTargetGasUsedPerBlock(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).TargetGasUsedPerBlock
}

func (k *Keeper) GetDeliverTxHookWasmGasLimit(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).DeliverTxHookWasmGasLimit
}

func (k *Keeper) ChainID(ctx sdk.Context) *big.Int {
	if k.EthBlockTestConfig.Enabled {
		// replay is for eth mainnet so always return 1
		return utils.Big1
	}
	// return mapped chain ID
	return config.GetEVMChainID(ctx.ChainID())

}

/*
*
eni gas = evm gas * multiplier
eni gas price = fee / eni gas = fee / (evm gas * multiplier) = evm gas / multiplier
*/
func (k *Keeper) GetEVMGasLimitFromCtx(ctx sdk.Context) uint64 {
	//return k.getEvmGasLimitFromCtx(ctx)
	return 0
}

func (k *Keeper) GetCosmosGasLimitFromEVMGas(ctx sdk.Context, evmGas uint64) uint64 {
	//gasMultipler := k.GetPriorityNormalizer(ctx)
	//gasLimitBigInt := sdk.NewDecFromInt(sdk.NewIntFromUint64(evmGas)).Mul(gasMultipler).TruncateInt().BigInt()
	//if gasLimitBigInt.Cmp(utils.BigMaxU64) > 0 {
	//	gasLimitBigInt = utils.BigMaxU64
	//}
	//return gasLimitBigInt.Uint64()
	return evmGas / k.GetPriorityNormalizer(ctx).BigInt().Uint64() // todo: fix this
}
