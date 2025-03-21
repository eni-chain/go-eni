package syscontract

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/eni-chain/go-eni/syscontract/genesis"
	"github.com/ethereum/go-ethereum/common"
	"strings"

	evmKeeper "github.com/eni-chain/go-eni/x/evm/keeper"
)

var contracts *contractsConfig

type Config struct {
	Addr common.Address
	Code string
}

type contractsConfig struct {
	Name    string
	Configs []*Config
}

func init() {
	contracts = &contractsConfig{
		Name: "syscontract",
		Configs: []*Config{
			{
				Addr: common.HexToAddress(HubAddr),
				Code: genesis.HubContract,
			},
			{
				Addr: common.HexToAddress(ValidatorManagerAddr),
				Code: genesis.ValidatorManagerContract,
			},
			{
				Addr: common.HexToAddress(VRFAddr),
				Code: genesis.VRFContract,
			},
		},
	}
}

func SetupSystemContracts(ctx sdk.Context, evmKeeper evmKeeper.Keeper) {
	if contracts == nil {
		evmKeeper.Logger().Info("empty contracts config", "height", ctx.BlockHeight())
		return
	}

	evmKeeper.Logger().Info(fmt.Sprintf("apply contracts %s at height %d", contracts.Name, ctx.BlockHeight()))
	for _, cfg := range contracts.Configs {
		evmKeeper.Logger().Info(fmt.Sprintf("contractsConfig contract %s", cfg.Addr.String()))

		newContractCode, err := hex.DecodeString(strings.TrimSpace(cfg.Code))
		if err != nil {
			panic(fmt.Errorf("failed to decode new contract code: %s", err.Error()))
		}
		caller := evmKeeper.AccountKeeper().GetModuleAddress(authtypes.FeeCollectorName)

		//addr0 := common.Address{0}
		body, err := evmKeeper.CallEVM(ctx, common.Address(caller), nil, nil, newContractCode)
		if err != nil {
			panic(err)
			panic(fmt.Errorf("failed to execute contract constructor: %s", err.Error()))
		}

		evmKeeper.SetCode(ctx, cfg.Addr, body)
	}
}
