package syscontract

import (
	"cosmossdk.io/log"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	evmKeeper "github.com/cosmos/cosmos-sdk/x/evm/keeper"
	"github.com/eni-chain/go-eni/syscontract/genesis"
	syscontractSdk "github.com/eni-chain/go-eni/syscontract/genesis/sdk"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"os"
	"strings"
)

var contracts *contractsConfig

var logger = log.NewLogger(os.Stdout)

type Config struct {
	Addr common.Address
	Code string
	Abi  abi.ABI
}

type contractsConfig struct {
	Name    string
	Configs []*Config
}

func init() {
	proxyABI, err := abi.JSON(strings.NewReader(syscontractSdk.PROXYABI))
	if err != nil {
		logger.Error(fmt.Sprintf("parse proxy contract abi failed:%v", err.Error()))
	}

	hubABI, err := abi.JSON(strings.NewReader(syscontractSdk.HubABI))
	if err != nil {
		logger.Error(fmt.Sprintf("parse hub abi contract failed:%v", err.Error()))
	}

	vrfABI, err := abi.JSON(strings.NewReader(syscontractSdk.VRFABI))
	if err != nil {
		logger.Error(fmt.Sprintf("parse vrf contract abi failed:%v", err.Error()))
	}

	validatorManagerABI, err := abi.JSON(strings.NewReader(syscontractSdk.ValidatorManagerABI))
	if err != nil {
		logger.Error(fmt.Sprintf("parse validator manager abi failed:%v", err.Error()))
	}

	contracts = &contractsConfig{
		Name: "syscontract",
		Configs: []*Config{
			{
				Addr: common.HexToAddress(syscontractSdk.ProxyAddr),
				Code: genesis.ProxyContract,
				Abi:  proxyABI,
			},
			{
				Addr: common.HexToAddress(syscontractSdk.HubAddr),
				Code: genesis.HubContract,
				Abi:  hubABI,
			},
			{
				Addr: common.HexToAddress(syscontractSdk.ValidatorManagerAddr),
				Code: genesis.ValidatorManagerContract,
				Abi:  validatorManagerABI,
			},
			{
				Addr: common.HexToAddress(syscontractSdk.VRFAddr),
				Code: genesis.VRFContract,
				Abi:  vrfABI,
			},
		},
	}
}

func SetupSystemContracts(ctx sdk.Context, evmKeeper *evmKeeper.Keeper) {
	if contracts == nil {
		evmKeeper.Logger().Info("empty contracts config", "height", ctx.BlockHeight())
		return
	}

	evmKeeper.Logger().Info(fmt.Sprintf("apply contracts %s at height %d", contracts.Name, ctx.BlockHeight()))

	for _, cfg := range contracts.Configs {
		if cfg.Addr == common.HexToAddress(syscontractSdk.ProxyAddr) {
			continue
		}

		evmKeeper.Logger().Info(fmt.Sprintf("contractsConfig contract %s", cfg.Addr.String()))

		newContractCode, err := hex.DecodeString(strings.TrimSpace(cfg.Code))
		if err != nil {
			panic(fmt.Errorf("failed to decode new contract code: %s", err.Error()))
		}
		caller := evmKeeper.AccountKeeper().GetModuleAddress(authtypes.FeeCollectorName)

		//addr0 := common.Address{0}
		body, err := evmKeeper.CallEVM(ctx, common.Address(caller), nil, nil, newContractCode)
		if err != nil {
			panic(fmt.Errorf("failed to execute contract constructor: %s", err.Error()))
		}

		evmKeeper.SetCode(ctx, cfg.Addr, body)
		calldata, err := cfg.Abi.Pack("init")
		if err != nil {
			panic(fmt.Errorf("failed to pack calldata: %s", err.Error()))
		}

		_, err = evmKeeper.CallEVM(ctx, common.Address(caller), &cfg.Addr, nil, calldata)
		if err != nil {
			panic(fmt.Errorf("failed to execute contract init: %s", err.Error()))
		}
	}
}
