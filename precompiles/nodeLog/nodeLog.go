package ContractNodeLog

import (
	"cosmossdk.io/log"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"os"
)

// log level
const (
	lDebug = 1
	lInfo  = 2
	lWarn  = 3
	lError = 4
)

var logger = log.NewLogger(os.Stdout)

// nodeLog enables users to log in the contract
type nodeLog struct{}

//func (c *nodeLog)SetValue(v string){}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *nodeLog) RequiredGas(input []byte) uint64 {
	return params.BalanceGasFrontier
}

func (c *nodeLog) Run(_ *vm.EVM, sender common.Address, _ common.Address, input []byte, _ *big.Int, _ bool, _ bool) ([]byte, error) {
	return printLog(sender, input)
}

// printLog implements the nodeLog precompile
func printLog(sender common.Address, input []byte) ([]byte, error) {
	level := input[0] - '0'

	switch level {
	case lDebug:
		logger.Debug(fmt.Sprintf("EVM LOG>> contract[%s]- %s", sender, string(input[1:])))
	case lInfo:
		logger.Info(fmt.Sprintf("EVM LOG>> contract[%s]- %s", sender, string(input[1:])))
	case lWarn:
		logger.Warn(fmt.Sprintf("EVM LOG>> contract[%s]- %s", sender, string(input[1:])))
	case lError:
		logger.Error(fmt.Sprintf("EVM LOG>> contract[%s]- %s", sender, string(input[1:])))
	default:
		logger.Info(fmt.Sprintf("EVM LOG>> contract[%s]- %v", sender, level))
	}

	return nil, nil
}

func AddNodeLogToVM() bool {
	p := &nodeLog{}
	addr := common.BytesToAddress([]byte{0xa2})

	vm.PrecompiledContractsHomestead[addr] = p
	vm.PrecompiledContractsByzantium[addr] = p
	vm.PrecompiledContractsIstanbul[addr] = p
	vm.PrecompiledContractsBerlin[addr] = p
	vm.PrecompiledContractsCancun[addr] = p
	vm.PrecompiledContractsBLS[addr] = p
	vm.PrecompiledAddressesHomestead = append(vm.PrecompiledAddressesHomestead, addr)
	vm.PrecompiledAddressesByzantium = append(vm.PrecompiledAddressesByzantium, addr)
	vm.PrecompiledAddressesIstanbul = append(vm.PrecompiledAddressesIstanbul, addr)
	vm.PrecompiledAddressesBerlin = append(vm.PrecompiledAddressesBerlin, addr)
	vm.PrecompiledAddressesCancun = append(vm.PrecompiledAddressesCancun, addr)

	return true
}
