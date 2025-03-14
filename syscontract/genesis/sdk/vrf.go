package sdk

import (
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/syscontract"
	evmKeeper "github.com/eni-chain/go-eni/x/evm/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// VRF is the Go interface for the VRF contract
type VRF struct {
	abi       abi.ABI
	evmKeeper *evmKeeper.Keeper
}

// NewVRF creates a new instance of VRF
func NewVRF(evmKeeper *evmKeeper.Keeper) (*VRF, error) {
	parsedABI, err := abi.JSON(strings.NewReader(VRFABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return &VRF{
		abi:       parsedABI,
		evmKeeper: evmKeeper,
	}, nil
}

// GetRandomSeed gets the random seed for a specific epoch
func (v *VRF) GetRandomSeed(
	ctx sdk.Context,
	caller common.Address,
	epoch *big.Int,
) ([]byte, error) {
	input, err := v.abi.Pack("getRandomSeed", epoch)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.VRFAddr)
	to := &address
	retData, err := v.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	var randomSeed []byte
	err = v.abi.UnpackIntoInterface(&randomSeed, "getRandomSeed", retData)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return randomSeed, nil
}

// Init initializes the VRF contract with a random seed
func (v *VRF) Init(
	ctx sdk.Context,
	caller common.Address,
	rnd []byte,
) ([]byte, error) {
	input, err := v.abi.Pack("init", rnd)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.VRFAddr)
	to := &address
	retData, err := v.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	return retData, nil
}

// SendRandom sends a random value for a specific epoch
func (v *VRF) SendRandom(
	ctx sdk.Context,
	caller common.Address,
	rnd []byte,
	epoch *big.Int,
) (bool, error) {
	input, err := v.abi.Pack("sendRandom", rnd, epoch)
	if err != nil {
		return false, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.VRFAddr)
	to := &address
	retData, err := v.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return false, fmt.Errorf("EVM call failed: %v", err)
	}

	var success bool
	err = v.abi.UnpackIntoInterface(&success, "sendRandom", retData)
	if err != nil {
		return false, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return success, nil
}

// UpdateAdmin updates the admin address
func (v *VRF) UpdateAdmin(
	ctx sdk.Context,
	caller common.Address,
	admin common.Address,
) ([]byte, error) {
	input, err := v.abi.Pack("updateAdmin", admin)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.VRFAddr)
	to := &address
	retData, err := v.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	return retData, nil
}

// UpdateConsensusSet updates the consensus set for a specific epoch
func (v *VRF) UpdateConsensusSet(
	ctx sdk.Context,
	caller common.Address,
	epoch *big.Int,
) ([]common.Address, error) {
	input, err := v.abi.Pack("updateConsensusSet", epoch)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.VRFAddr)
	to := &address
	retData, err := v.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	if retData == nil {
		return nil, nil
	}

	var nodes []common.Address
	err = v.abi.UnpackIntoInterface(&nodes, "updateConsensusSet", retData)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return nodes, nil
}

// ExampleVRF Example usage of VRF
func ExampleVRF() {
	evmKeeper, ctx := GetEVMKeeper()
	if evmKeeper == nil {
		panic("Failed to get EVM keeper")
	}
	// Create a new VRF instance
	vrf, err := NewVRF(evmKeeper)
	if err != nil {
		panic(fmt.Sprintf("Failed to create VRF: %v", err))
	}

	caller := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// Initialize VRF with a random seed
	initialRandomSeed := []byte("initial-random-seed-bytes")
	_, err = vrf.Init(ctx, caller, initialRandomSeed)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize VRF: %v", err))
	}

	// Get random seed for a specific epoch
	epoch := big.NewInt(1234)
	randomSeed, err := vrf.GetRandomSeed(ctx, caller, epoch)
	if err != nil {
		panic(fmt.Sprintf("Failed to get random seed: %v", err))
	}
	fmt.Printf("Random seed for epoch %s: %x\n", epoch.String(), randomSeed)

	// Send a random value for a specific epoch
	newRandomValue := []byte("new-random-value-bytes")
	success, err := vrf.SendRandom(ctx, caller, newRandomValue, epoch)
	if err != nil {
		panic(fmt.Sprintf("Failed to send random value: %v", err))
	}
	fmt.Printf("Send random successful: %v\n", success)

	// Update admin address
	newAdminAddr := common.HexToAddress("0x2345678901234567890123456789012345678901")
	_, err = vrf.UpdateAdmin(ctx, caller, newAdminAddr)
	if err != nil {
		panic(fmt.Sprintf("Failed to update admin: %v", err))
	}

	// Update consensus set for a specific epoch
	newEpoch := big.NewInt(1235)
	consensusSet, err := vrf.UpdateConsensusSet(ctx, caller, newEpoch)
	if err != nil {
		panic(fmt.Sprintf("Failed to update consensus set: %v", err))
	}
	fmt.Printf("Number of nodes in consensus set: %d\n", len(consensusSet))
}
