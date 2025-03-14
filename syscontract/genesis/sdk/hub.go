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

// Hub is the Go interface for the Hub contract
type Hub struct {
	abi abi.ABI
}

// NewHub creates a new instance of Hub
func NewHub() (*Hub, error) {
	parsedABI, err := abi.JSON(strings.NewReader(HubABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return &Hub{
		abi: parsedABI,
	}, nil
}

// ApplyForValidator applies to become a validator
func (h *Hub) ApplyForValidator(
	evm *vm.EVM,
	caller common.Address,
	node common.Address,
	agent common.Address,
	name string,
	description string,
	pubKey []byte,
	value *uint256.Int,
) error {
	input, err := h.abi.Pack("applyForValidator", node, agent, name, description, pubKey)
	if err != nil {
		return fmt.Errorf("failed to pack ABI: %v", err)
	}

	_, err = callEVM(evm, caller, common.HexToAddress(syscontract.HubAddr), input, value)
	if err != nil {
		return fmt.Errorf("EVM call failed: %v", err)
	}

	return nil
}

// AuditPass approves a validator application
func (h *Hub) AuditPass(evm *vm.EVM, caller common.Address, operator common.Address) error {
	input, err := h.abi.Pack("auditPass", operator)
	if err != nil {
		return fmt.Errorf("failed to pack ABI: %v", err)
	}

	_, err = callEVM(evm, caller, common.HexToAddress(syscontract.HubAddr), input, uint256.NewInt(0))
	if err != nil {
		return fmt.Errorf("EVM call failed: %v", err)
	}

	return nil
}

// BlockReward distributes block rewards to a validator
func (h *Hub) BlockReward(evm *vm.EVM, caller common.Address, node common.Address) (*big.Int, error) {
	input, err := h.abi.Pack("blockReward", node)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %v", err)
	}

	ret, err := callEVM(evm, caller, common.HexToAddress(syscontract.HubAddr), input, uint256.NewInt(0))
	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %v", err)
	}

	var amount *big.Int
	err = h.abi.UnpackIntoInterface(&amount, "blockReward", ret)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack return value: %v", err)
	}

	return amount, nil
}

// UpdateAdmin updates the admin address
func (h *Hub) UpdateAdmin(evm *vm.EVM, caller common.Address, admin common.Address) error {
	input, err := h.abi.Pack("updateAdmin", admin)
	if err != nil {
		return fmt.Errorf("failed to pack ABI: %v", err)
	}

	_, err = callEVM(evm, caller, common.HexToAddress(syscontract.HubAddr), input, uint256.NewInt(0))
	if err != nil {
		return fmt.Errorf("EVM call failed: %v", err)
	}

	return nil
}
