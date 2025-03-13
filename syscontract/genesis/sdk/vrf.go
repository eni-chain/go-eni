package sdk

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/eni-chain/go-eni/syscontract"
	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// VRF is the Go interface for the VRF contract
type VRF struct {
	abi abi.ABI
}

// NewVRF creates a new instance of VRF
func NewVRF() (*VRF, error) {
	parsedABI, err := abi.JSON(strings.NewReader(VRFABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return &VRF{
		abi: parsedABI,
	}, nil
}

// GetRandomSeed gets the random seed for a specific epoch
func (v *VRF) GetRandomSeed(evm *vm.EVM, caller common.Address, epoch *big.Int) ([]byte, error) {
	input, err := v.abi.Pack("getRandomSeed", epoch)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	ret, err := callEVM(evm, caller, common.HexToAddress(syscontract.VRFAddr), input, uint256.NewInt(0))
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	var randomSeed []byte
	err = v.abi.UnpackIntoInterface(&randomSeed, "getRandomSeed", ret)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return randomSeed, nil
}

// Init initializes the VRF contract with a random seed
func (v *VRF) Init(evm *vm.EVM, caller common.Address, rnd []byte) error {
	input, err := v.abi.Pack("init", rnd)
	if err != nil {
		return fmt.Errorf("failed to pack ABI: %v", err)
	}

	_, err = callEVM(evm, caller, common.HexToAddress(syscontract.VRFAddr), input, uint256.NewInt(0))
	if err != nil {
		return fmt.Errorf("EVM call failed: %v", err)
	}

	return nil
}

// SendRandom sends a random value for a specific epoch
func (v *VRF) SendRandom(evm *vm.EVM, caller common.Address, rnd []byte, epoch *big.Int) (bool, error) {
	input, err := v.abi.Pack("sendRandom", rnd, epoch)
	if err != nil {
		return false, fmt.Errorf("failed to pack ABI: %v", err)
	}

	ret, err := callEVM(evm, caller, common.HexToAddress(syscontract.VRFAddr), input, uint256.NewInt(0))
	if err != nil {
		return false, fmt.Errorf("EVM call failed: %v", err)
	}

	var success bool
	err = v.abi.UnpackIntoInterface(&success, "sendRandom", ret)
	if err != nil {
		return false, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return success, nil
}

// UpdateAdmin updates the admin address
func (v *VRF) UpdateAdmin(evm *vm.EVM, caller common.Address, admin common.Address) error {
	input, err := v.abi.Pack("updateAdmin", admin)
	if err != nil {
		return fmt.Errorf("failed to pack ABI: %v", err)
	}

	_, err = callEVM(evm, caller, common.HexToAddress(syscontract.VRFAddr), input, uint256.NewInt(0))
	if err != nil {
		return fmt.Errorf("EVM call failed: %v", err)
	}

	return nil
}

// UpdateConsensusSet updates the consensus set for a specific epoch
func (v *VRF) UpdateConsensusSet(evm *vm.EVM, caller common.Address, epoch *big.Int) ([]common.Address, error) {
	input, err := v.abi.Pack("updateConsensusSet", epoch)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	ret, err := callEVM(evm, caller, common.HexToAddress(syscontract.VRFAddr), input, uint256.NewInt(0))
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	var nodes []common.Address
	err = v.abi.UnpackIntoInterface(&nodes, "updateConsensusSet", ret)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return nodes, nil
}
