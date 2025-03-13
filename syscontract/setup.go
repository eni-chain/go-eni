package syscontract

import (
	"encoding/hex"
	"fmt"
	"github.com/eni-chain/go-eni/syscontract/genesis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/log"
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

func SetupSystemContracts(blockNumber *big.Int, statedb vm.StateDB, logger log.Logger) {
	if contracts == nil {
		logger.Info("Empty contracts config", "height", blockNumber.String())
		return
	}

	logger.Info(fmt.Sprintf("Apply contracts %s at height %d", contracts.Name, blockNumber.Int64()))
	for _, cfg := range contracts.Configs {
		logger.Info(fmt.Sprintf("contractsConfig contract %s", cfg.Addr.String()))

		newContractCode, err := hex.DecodeString(strings.TrimSpace(cfg.Code))
		if err != nil {
			panic(fmt.Errorf("failed to decode new contract code: %s", err.Error()))
		}
		statedb.SetCode(cfg.Addr, newContractCode)

	}
}
