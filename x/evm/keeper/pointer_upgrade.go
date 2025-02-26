package keeper

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/utils"
	"github.com/eni-chain/go-eni/x/evm/artifacts"
	"github.com/eni-chain/go-eni/x/evm/state"
	"github.com/eni-chain/go-eni/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
)

func (k *Keeper) RunWithOneOffEVMInstance(
	ctx sdk.Context, runner func(*vm.EVM) error, logger func(string, string),
) error {
	stateDB := state.NewDBImpl(ctx, k, false)
	evmModuleAddress := k.GetEVMAddressOrDefault(ctx, k.AccountKeeper().GetModuleAddress(types.ModuleName))
	gp := core.GasPool(math.MaxUint64)
	blockCtx, err := k.GetVMBlockContext(ctx, gp)
	if err != nil {
		logger("get block context", err.Error())
		return err
	}
	cfg := types.DefaultChainConfig().EthereumConfig(k.ChainID(ctx))
	txCtx := core.NewEVMTxContext(&core.Message{From: evmModuleAddress, GasPrice: utils.Big0})
	evmInstance := vm.NewEVM(*blockCtx, stateDB, cfg, vm.Config{})
	evmInstance.SetTxContext(txCtx)
	err = runner(evmInstance)
	if err != nil {
		logger("upserting pointer", err.Error())
		return err
	}
	surplus, err := stateDB.Finalize()
	if err != nil {
		logger("finalizing", err.Error())
		return err
	}
	if !surplus.IsZero() {
		logger("non-zero surplus", surplus.String())
	}
	return nil
}

func (k *Keeper) UpsertERCNativePointer(
	ctx sdk.Context, evm *vm.EVM, token string, metadata utils.ERCMetadata,
) (contractAddr common.Address, err error) {
	return k.UpsertERCPointer(
		ctx, evm, "native", []interface{}{
			token, metadata.Name, metadata.Symbol, metadata.Decimals,
		}, k.GetERC20NativePointer, k.SetERC20NativePointer,
	)
}

func (k *Keeper) UpsertERCPointer(
	ctx sdk.Context, evm *vm.EVM, typ string, args []interface{}, getter PointerGetter, setter PointerSetter,
) (contractAddr common.Address, err error) {
	pointee := args[0].(string)
	evmModuleAddress := k.GetEVMAddressOrDefault(ctx, k.AccountKeeper().GetModuleAddress(types.ModuleName))

	var bin []byte
	bin, err = artifacts.GetParsedABI(typ).Pack("", args...)
	if err != nil {
		panic(err)
	}
	bin = append(artifacts.GetBin(typ), bin...)
	existingAddr, _, exists := getter(ctx, pointee)
	//suppliedGas := k.getEvmGasLimitFromCtx(ctx)
	suppliedGas := uint64(math.MaxUint64)
	var remainingGas uint64
	if exists {
		var ret []byte
		contractAddr = existingAddr
		ret, remainingGas, err = evm.GetDeploymentCode(vm.AccountRef(evmModuleAddress), bin, suppliedGas, utils.Uint2560, existingAddr)
		k.SetCode(ctx, contractAddr, ret)
	} else {
		_, contractAddr, remainingGas, err = evm.Create(vm.AccountRef(evmModuleAddress), bin, suppliedGas, utils.Uint2560)
	}
	if err != nil {
		return
	}
	ctx.GasMeter().ConsumeGas(k.GetCosmosGasLimitFromEVMGas(ctx, suppliedGas-remainingGas), "ERC pointer deployment")
	if err = setter(ctx, pointee, contractAddr); err != nil {
		return
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypePointerRegistered, sdk.NewAttribute(types.AttributeKeyPointerType, typ),
		sdk.NewAttribute(types.AttributeKeyPointerAddress, contractAddr.Hex()), sdk.NewAttribute(types.AttributeKeyPointee, pointee)))
	return
}
