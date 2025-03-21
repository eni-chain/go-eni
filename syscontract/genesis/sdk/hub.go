package sdk

import (
	"fmt"
	"math/big"
	"strings"

	cosmossdk_io_math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmKeeper "github.com/cosmos/cosmos-sdk/x/evm/keeper"
	"github.com/eni-chain/go-eni/syscontract"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// Hub is the Go interface for the Hub contract
type Hub struct {
	abi       abi.ABI
	evmKeeper *evmKeeper.Keeper
}

// NewHub creates a new instance of Hub
func NewHub(evmKeeper *evmKeeper.Keeper) (*Hub, error) {
	parsedABI, err := abi.JSON(strings.NewReader(HubABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return &Hub{
		abi:       parsedABI,
		evmKeeper: evmKeeper,
	}, nil
}

// ApplyForValidator applies to become a validator
func (h *Hub) ApplyForValidator(
	ctx sdk.Context,
	caller common.Address,
	node common.Address,
	agent common.Address,
	name string,
	description string,
	pubKey []byte,
	value *cosmossdk_io_math.Int,
) ([]byte, error) {
	input, err := h.abi.Pack("applyForValidator", node, agent, name, description, pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.HubAddr)
	to := &address
	retData, err := h.evmKeeper.CallEVM(ctx, caller, to, value, input)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	return retData, nil
}

// AuditPass approves a validator application
func (h *Hub) AuditPass(
	ctx sdk.Context,
	caller common.Address,
	operator common.Address,
) ([]byte, error) {
	input, err := h.abi.Pack("auditPass", operator)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.HubAddr)
	to := &address
	retData, err := h.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	return retData, nil
}

// BlockReward distributes block rewards to a validator
func (h *Hub) BlockReward(
	ctx sdk.Context,
	caller common.Address,
	node common.Address,
) (common.Address, *big.Int, error) {
	input, err := h.abi.Pack("blockReward", node)
	if err != nil {
		return common.Address{0}, nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.HubAddr)
	to := &address
	retData, err := h.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return common.Address{0}, big.NewInt(0), nil
	}
	if retData == nil {
		return common.Address{0}, big.NewInt(0), nil
	}

	rets, err := h.abi.Unpack("blockReward", retData)
	if err != nil {
		return common.Address{0}, big.NewInt(0), nil
	}

	operator := rets[0].(common.Address)
	amount := rets[1].(*big.Int)

	return operator, amount, nil
}

// UpdateAdmin updates the admin address
func (h *Hub) UpdateAdmin(
	ctx sdk.Context,
	caller common.Address,
	admin common.Address,
) ([]byte, error) {
	input, err := h.abi.Pack("updateAdmin", admin)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	address := common.HexToAddress(syscontract.HubAddr)
	to := &address
	retData, err := h.evmKeeper.CallEVM(ctx, caller, to, nil, input)
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	return retData, nil
}

// ExampleHub Example usage of Hub
func ExampleHub() {
	evmKeeper, ctx := GetEVMKeeper()
	if evmKeeper == nil {
		panic("Failed to get EVM keeper")
	}
	// Create a new Hub instance
	hub, err := NewHub(evmKeeper)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Hub: %v", err))
	}

	// Example context and addresses
	caller := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// Apply to become a validator
	nodeAddr := common.HexToAddress("0x2345678901234567890123456789012345678901")
	agentAddr := common.HexToAddress("0x3456789012345678901234567890123456789012")
	name := "MyValidator"
	description := "A reliable validator for the network"
	pubKey := []byte("validator-pubkey-bytes")

	// Convert pledge amount to cosmossdk_io_math.Int
	pledgeAmount := cosmossdk_io_math.NewInt(10000)

	_, err = hub.ApplyForValidator(ctx, caller, nodeAddr, agentAddr, name, description, pubKey, &pledgeAmount)
	if err != nil {
		panic(fmt.Sprintf("Failed to apply for validator: %v", err))
	}

	// Approve a validator application
	operatorAddr := common.HexToAddress("0x4567890123456789012345678901234567890123")
	_, err = hub.AuditPass(ctx, caller, operatorAddr)
	if err != nil {
		panic(fmt.Sprintf("Failed to approve validator: %v", err))
	}

	// Distribute block rewards
	validatorNodeAddr := common.HexToAddress("0x5678901234567890123456789012345678901234")
	operator, reward, err := hub.BlockReward(ctx, caller, validatorNodeAddr)
	if err != nil {
		panic(fmt.Sprintf("Failed to distribute block reward: %v", err))
	}
	fmt.Printf("Block reward: %s, for: %x\n", reward.String(), operator)

	// Update admin address
	newAdminAddr := common.HexToAddress("0x6789012345678901234567890123456789012345")
	_, err = hub.UpdateAdmin(ctx, caller, newAdminAddr)
	if err != nil {
		panic(fmt.Sprintf("Failed to update admin: %v", err))
	}
}
