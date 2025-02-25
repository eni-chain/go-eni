package keeper

import (
	"cosmossdk.io/store/prefix"
	//"github.com/cosmos/cosmos-sdk/store/prefix"
	cosmossdk_io_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
)

func (k *Keeper) AddAnteSurplus(ctx sdk.Context, txHash common.Hash, surplus cosmossdk_io_math.Int) error {
	store := prefix.NewStore(ctx.TransientStore(k.transientStoreKey), types.AnteSurplusPrefix)
	bz, err := surplus.Marshal()
	if err != nil {
		return err
	}
	store.Set(txHash[:], bz)
	return nil
}

func (k *Keeper) GetAnteSurplusSum(ctx sdk.Context) cosmossdk_io_math.Int {
	iter := prefix.NewStore(ctx.TransientStore(k.transientStoreKey), types.AnteSurplusPrefix).Iterator(nil, nil)
	defer iter.Close()
	res := cosmossdk_io_math.ZeroInt()
	for ; iter.Valid(); iter.Next() {
		surplus := cosmossdk_io_math.Int{}
		_ = surplus.Unmarshal(iter.Value())
		res = res.Add(surplus)
	}
	return res
}
