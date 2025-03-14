package sdk

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/eni-chain/go-eni/syscontract"
	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// ValidatorManager is the Go interface for the ValidatorManager contract
type ValidatorManager struct {
	abi abi.ABI
}

// NewValidatorManager creates a new instance of ValidatorManager
func NewValidatorManager() (*ValidatorManager, error) {
	parsedABI, err := abi.JSON(strings.NewReader(validatorManagerABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return &ValidatorManager{
		abi: parsedABI,
	}, nil
}

// Directly call EVM to execute contract methods
func callEVM(evm *vm.EVM, caller common.Address, contractAddr common.Address, input []byte, value *uint256.Int) ([]byte, error) {

	// Execute the call
	ret, _, err := evm.Call(
		vm.AccountRef(caller),
		contractAddr,
		input,
		uint64(10000000), // gas limit
		value,
	)

	return ret, err
}

// GetPubkey gets the public key of a validator
func (vm *ValidatorManager) GetPubkey(evm *vm.EVM, caller common.Address, validator common.Address) ([]byte, error) {
	input, err := vm.abi.Pack("getPubkey", validator)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	ret, err := callEVM(evm, caller, common.HexToAddress(syscontract.ValidatorManagerAddr), input, uint256.NewInt(0))
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	var pubkey []byte
	err = vm.abi.UnpackIntoInterface(&pubkey, "getPubkey", ret)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return pubkey, nil
}

// GetValidatorSet gets the set of validators
func (vm *ValidatorManager) GetValidatorSet(evm *vm.EVM, caller common.Address) ([]common.Address, error) {
	input, err := vm.abi.Pack("getValidatorSet")
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	ret, err := callEVM(evm, caller, common.HexToAddress(syscontract.ValidatorManagerAddr), input, uint256.NewInt(0))
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	var validators []common.Address
	err = vm.abi.UnpackIntoInterface(&validators, "getValidatorSet", ret)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return validators, nil
}

// AddValidator adds a new validator
func (vm *ValidatorManager) AddValidator(
	evm *vm.EVM,
	caller common.Address,
	operator common.Address,
	node common.Address,
	agent common.Address,
	amount *big.Int,
	enterTime *big.Int,
	name string,
	description string,
	pubKey []byte,
) error {
	input, err := vm.abi.Pack("addValidator", operator, node, agent, amount, enterTime, name, description, pubKey)
	if err != nil {
		return fmt.Errorf("failed to pack ABI: %v", err)
	}

	_, err = callEVM(evm, caller, common.HexToAddress(syscontract.ValidatorManagerAddr), input, uint256.NewInt(0))
	if err != nil {
		return fmt.Errorf("EVM call failed: %v", err)
	}

	return nil
}

// UpdateConsensus updates the consensus node set
func (vm *ValidatorManager) UpdateConsensus(evm *vm.EVM, caller common.Address, nodes []common.Address) error {
	input, err := vm.abi.Pack("undateConsensus", nodes)
	if err != nil {
		return fmt.Errorf("failed to pack ABI: %v", err)
	}

	_, err = callEVM(evm, caller, common.HexToAddress(syscontract.ValidatorManagerAddr), input, uint256.NewInt(0))
	if err != nil {
		return fmt.Errorf("EVM call failed: %v", err)
	}

	return nil
}

// GetPledgeAmount gets the pledge amount of a node
func (vm *ValidatorManager) GetPledgeAmount(evm *vm.EVM, caller common.Address, node common.Address) (*big.Int, error) {
	input, err := vm.abi.Pack("getPledgeAmount", node)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	ret, err := callEVM(evm, caller, common.HexToAddress(syscontract.ValidatorManagerAddr), input, uint256.NewInt(0))
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	var amount *big.Int
	err = vm.abi.UnpackIntoInterface(&amount, "getPledgeAmount", ret)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return amount, nil
}

// ExampleUsage example usage
func ExampleUsage(evm *vm.EVM, caller common.Address) {
	validatorManager, err := NewValidatorManager()
	if err != nil {
		fmt.Printf("failed to create ValidatorManager: %v\n", err)
		return
	}

	// Get the set of validators
	validators, err := validatorManager.GetValidatorSet(evm, caller)
	if err != nil {
		fmt.Printf("failed to get validator set: %v\n", err)
		return
	}
	fmt.Printf("validator set: %v\n", validators)

	// Get the public key of a validator
	if len(validators) > 0 {
		pubkey, err := validatorManager.GetPubkey(evm, caller, validators[0])
		if err != nil {
			fmt.Printf("failed to get public key: %v\n", err)
			return
		}
		fmt.Printf("validator public key: %s\n", hex.EncodeToString(pubkey))
	}
}
