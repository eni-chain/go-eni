package sdk

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmKeeper "github.com/cosmos/cosmos-sdk/x/evm/keeper"
	"github.com/eni-chain/go-eni/syscontract"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// ValidatorManager is the Go interface for the ValidatorManager contract
type ValidatorManager struct {
	abi       abi.ABI
	evmKeeper *evmKeeper.Keeper
}

// NewValidatorManager creates a new instance of ValidatorManager
func NewValidatorManager(evmKeeper *evmKeeper.Keeper) (*ValidatorManager, error) {
	parsedABI, err := abi.JSON(strings.NewReader(validatorManagerABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return &ValidatorManager{
		abi:       parsedABI,
		evmKeeper: evmKeeper,
	}, nil
}

// GetPubkey gets the public key of a validator
func (vm *ValidatorManager) GetPubkey(
	ctx sdk.Context,
	caller common.Address,
	validator common.Address,
) ([]byte, error) {
	input, err := vm.abi.Pack("getPubkey", validator)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.ValidatorManagerAddr)
	to := &address
	retData, err := vm.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	var pubkey []byte
	err = vm.abi.UnpackIntoInterface(&pubkey, "getPubkey", retData)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return pubkey, nil
}

// GetPubKeysBySequence gets the public keys of a validators by sequence
func (vm *ValidatorManager) GetPubKeysBySequence(
	ctx sdk.Context,
	caller common.Address,
	validators []common.Address,
) ([][]byte, error) {
	input, err := vm.abi.Pack("getPubKeysBySequence", validators)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.ValidatorManagerAddr)
	to := &address
	retData, err := vm.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	if retData == nil {
		return nil, nil
	}

	var pubkeys [][]byte
	err = vm.abi.UnpackIntoInterface(&pubkeys, "getPubKeysBySequence", retData)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return pubkeys, nil
}

// GetValidatorSet gets the set of validators
func (vm *ValidatorManager) GetValidatorSet(
	ctx sdk.Context,
	caller common.Address,
) ([]common.Address, error) {
	input, err := vm.abi.Pack("getValidatorSet")
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.ValidatorManagerAddr)
	to := &address
	retData, err := vm.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	var validators []common.Address
	err = vm.abi.UnpackIntoInterface(&validators, "getValidatorSet", retData)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return validators, nil
}

// AddValidator adds a new validator
func (vm *ValidatorManager) AddValidator(
	ctx sdk.Context,
	caller common.Address,
	operator common.Address,
	node common.Address,
	agent common.Address,
	amount *big.Int,
	enterTime *big.Int,
	name string,
	description string,
	pubKey []byte,
) ([]byte, error) {
	input, err := vm.abi.Pack("addValidator", operator, node, agent, amount, enterTime, name, description, pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.ValidatorManagerAddr)
	to := &address
	retData, err := vm.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	return retData, nil
}

// UpdateConsensus updates the consensus node set
func (vm *ValidatorManager) UpdateConsensus(
	ctx sdk.Context,
	caller common.Address,
	nodes []common.Address,
) ([]byte, error) {
	// Note: Fixed the typo in the function name from "undateConsensus" to "updateConsensus"
	// However, since the contract uses "undateConsensus", we need to keep that name in the ABI call
	input, err := vm.abi.Pack("undateConsensus", nodes)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.ValidatorManagerAddr)
	to := &address
	retData, err := vm.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	return retData, nil
}

// GetPledgeAmount gets the pledge amount of a node
func (vm *ValidatorManager) GetPledgeAmount(
	ctx sdk.Context,
	caller common.Address,
	node common.Address,
) (*big.Int, error) {
	input, err := vm.abi.Pack("getPledgeAmount", node)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.ValidatorManagerAddr)
	to := &address
	retData, err := vm.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	var amount *big.Int
	err = vm.abi.UnpackIntoInterface(&amount, "getPledgeAmount", retData)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return amount, nil
}

// ExampleValidatorManager Example usage of ValidatorManager
func ExampleValidatorManager() {
	evmKeeper, ctx := GetEVMKeeper()
	if evmKeeper == nil {
		panic("Failed to get EVM keeper")
	}
	// Create a new ValidatorManager instance
	vm, err := NewValidatorManager(evmKeeper)
	if err != nil {
		panic(fmt.Sprintf("Failed to create ValidatorManager: %v", err))
	}

	// Example context and addresses
	caller := common.HexToAddress("0x1234567890123456789012345678901234567890")
	validatorAddr := common.HexToAddress("0x2345678901234567890123456789012345678901")

	// Get validator's public key
	pubkey, err := vm.GetPubkey(ctx, caller, validatorAddr)
	if err != nil {
		panic(fmt.Sprintf("Failed to get validator pubkey: %v", err))
	}
	fmt.Printf("Validator pubkey: %x\n", pubkey)

	// Get the validator set
	validators, err := vm.GetValidatorSet(ctx, caller)
	if err != nil {
		panic(fmt.Sprintf("Failed to get validator set: %v", err))
	}
	fmt.Printf("Number of validators: %d\n", len(validators))

	// Add a new validator
	operatorAddr := common.HexToAddress("0x3456789012345678901234567890123456789012")
	nodeAddr := common.HexToAddress("0x4567890123456789012345678901234567890123")
	agentAddr := common.HexToAddress("0x5678901234567890123456789012345678901234")
	amount := big.NewInt(10000)
	enterTime := big.NewInt(time.Now().Unix())
	name := "Validator1"
	description := "A reliable validator"
	pubKey := []byte("validator-pubkey-bytes")

	_, err = vm.AddValidator(ctx, caller, operatorAddr, nodeAddr, agentAddr, amount, enterTime, name, description, pubKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to add validator: %v", err))
	}

	// Update consensus nodes
	nodes := []common.Address{
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
		common.HexToAddress("0x2222222222222222222222222222222222222222"),
		common.HexToAddress("0x3333333333333333333333333333333333333333"),
	}
	_, err = vm.UpdateConsensus(ctx, caller, nodes)
	if err != nil {
		panic(fmt.Sprintf("Failed to update consensus: %v", err))
	}

	// Get pledge amount
	pledgeAmount, err := vm.GetPledgeAmount(ctx, caller, nodeAddr)
	if err != nil {
		panic(fmt.Sprintf("Failed to get pledge amount: %v", err))
	}
	fmt.Printf("Pledge amount: %s\n", pledgeAmount.String())
}
