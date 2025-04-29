package sdk

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/evm/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var TestValidatorManagerAddr = "0x1234567890123456789012345678901234567890"

// Setup function to create a ValidatorManager instance with mocked dependencies
func setupValidatorManager(t *testing.T) (*ValidatorManager, *MockEVMKeeper, sdk.Context) {
	parsedABI, err := abi.JSON(strings.NewReader(TestABI))
	if err != nil {
		t.Fatalf("Failed to parse test ABI: %v", err)
	}

	evmKeeper := &MockEVMKeeper{}
	ctx := sdk.Context{}
	vm := &ValidatorManager{
		abi:       parsedABI,
		evmKeeper: &keeper.Keeper{},
	}
	return vm, evmKeeper, ctx
}

func TestNewValidatorManager(t *testing.T) {
	// Test case: Successful creation
	t.Run("Success", func(t *testing.T) {
		vm, err := NewValidatorManager(&keeper.Keeper{})
		assert.NoError(t, err)
		assert.NotNil(t, vm)
		assert.NotNil(t, vm.abi)
		assert.NotNil(t, vm.evmKeeper)
	})

	// Test case: Invalid ABI
	t.Run("InvalidABI", func(t *testing.T) {
		invalidABI := "invalid json"
		_, err := abi.JSON(strings.NewReader(invalidABI))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid character")
	})
}

func TestGetPubkey(t *testing.T) {
	vm, evmKeeper, ctx := setupValidatorManager(t)
	caller := common.HexToAddress("0x123")
	validator := common.HexToAddress("0x456")

	// Test case: Successful call
	t.Run("Success", func(t *testing.T) {
		inputData, err := vm.abi.Pack("getPubkey", validator)
		assert.NoError(t, err)
		returnData := []byte("pubkey-data")
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vm.GetPubkey(ctx, caller, validator)
		assert.NoError(t, err)
		assert.Equal(t, returnData, result)
	})

	// Test case: ABI pack failure
	t.Run("ABIPackFailure", func(t *testing.T) {
		result, err := vm.GetPubkey(ctx, caller, common.Address{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to pack ABI")
		assert.Nil(t, result)
	})

	// Test case: EVM call failure
	t.Run("EVMCallFailure", func(t *testing.T) {
		inputData, err := vm.abi.Pack("getPubkey", validator)
		assert.NoError(t, err)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return([]byte{}, fmt.Errorf("evm error")).Once()
		result, err := vm.GetPubkey(ctx, caller, validator)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EVM call failed")
		assert.Nil(t, result)
	})

	// Test case: Unpack failure
	t.Run("UnpackFailure", func(t *testing.T) {
		inputData, err := vm.abi.Pack("getPubkey", validator)
		assert.NoError(t, err)
		returnData := []byte{0x00} // Invalid for bytes type
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vm.GetPubkey(ctx, caller, validator)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unpack return value")
		assert.Nil(t, result)
	})
}

func TestGetPubKeysBySequence(t *testing.T) {
	vm, evmKeeper, ctx := setupValidatorManager(t)
	caller := common.HexToAddress("0x123")
	validators := []common.Address{
		common.HexToAddress("0x456"),
		common.HexToAddress("0x789"),
	}

	// Test case: Successful call
	t.Run("Success", func(t *testing.T) {
		inputData, err := vm.abi.Pack("getPubKeysBySequence", validators)
		assert.NoError(t, err)
		returnData := []byte{0x00, 0x01, 0x02} // Simplified bytes[] data
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vm.GetPubKeysBySequence(ctx, caller, validators)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	// Test case: Nil return data
	t.Run("NilReturnData", func(t *testing.T) {
		inputData, err := vm.abi.Pack("getPubKeysBySequence", validators)
		assert.NoError(t, err)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return([]byte{}, nil).Once()
		result, err := vm.GetPubKeysBySequence(ctx, caller, validators)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	// Test case: ABI pack failure
	t.Run("ABIPackFailure", func(t *testing.T) {
		result, err := vm.GetPubKeysBySequence(ctx, caller, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to pack ABI")
		assert.Nil(t, result)
	})

	// Test case: EVM call failure
	t.Run("EVMCallFailure", func(t *testing.T) {
		inputData, err := vm.abi.Pack("getPubKeysBySequence", validators)
		assert.NoError(t, err)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return([]byte{}, fmt.Errorf("evm error")).Once()
		result, err := vm.GetPubKeysBySequence(ctx, caller, validators)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EVM call failed")
		assert.Nil(t, result)
	})

	// Test case: Unpack failure
	t.Run("UnpackFailure", func(t *testing.T) {
		inputData, err := vm.abi.Pack("getPubKeysBySequence", validators)
		assert.NoError(t, err)
		returnData := []byte{0x00} // Invalid for bytes[] type
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vm.GetPubKeysBySequence(ctx, caller, validators)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unpack return value")
		assert.Nil(t, result)
	})
}

func TestGetValidatorSet(t *testing.T) {
	vm, evmKeeper, ctx := setupValidatorManager(t)
	caller := common.HexToAddress("0x123")

	// Test case: Successful call
	t.Run("Success", func(t *testing.T) {
		inputData, err := vm.abi.Pack("getValidatorSet")
		assert.NoError(t, err)
		validators := []common.Address{
			common.HexToAddress("0x456"),
			common.HexToAddress("0x789"),
		}
		returnData, err := vm.abi.Pack("", validators)
		assert.NoError(t, err)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vm.GetValidatorSet(ctx, caller)
		assert.NoError(t, err)
		assert.Equal(t, validators, result)
	})

	// Test case: ABI pack failure
	t.Run("ABIPackFailure", func(t *testing.T) {
		// Note: getValidatorSet has no inputs, so ABI pack failure is unlikely
		// Test by mocking an invalid ABI method
		vm.abi.Methods["getValidatorSet"] = abi.Method{} // Temporarily corrupt method
		result, err := vm.GetValidatorSet(ctx, caller)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to pack ABI")
		assert.Nil(t, result)
		// Restore ABI for other tests
		parsedABI, _ := abi.JSON(strings.NewReader(TestABI))
		vm.abi = parsedABI
	})

	// Test case: EVM call failure
	t.Run("EVMCallFailure", func(t *testing.T) {
		inputData, err := vm.abi.Pack("getValidatorSet")
		assert.NoError(t, err)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return([]byte{}, fmt.Errorf("evm error")).Once()
		result, err := vm.GetValidatorSet(ctx, caller)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EVM call failed")
		assert.Nil(t, result)
	})

	// Test case: Unpack failure
	t.Run("UnpackFailure", func(t *testing.T) {
		inputData, err := vm.abi.Pack("getValidatorSet")
		assert.NoError(t, err)
		returnData := []byte{0x00} // Invalid for address[] type
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vm.GetValidatorSet(ctx, caller)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unpack return value")
		assert.Nil(t, result)
	})
}

func TestAddValidator(t *testing.T) {
	vm, evmKeeper, ctx := setupValidatorManager(t)
	caller := common.HexToAddress("0x123")
	operator := common.HexToAddress("0x456")
	node := common.HexToAddress("0x789")
	agent := common.HexToAddress("0xabc")
	amount := big.NewInt(10000)
	enterTime := big.NewInt(time.Now().Unix())
	name := "Validator1"
	description := "Test validator"
	pubKey := []byte("pubkey-data")

	// Test case: Successful call
	t.Run("Success", func(t *testing.T) {
		inputData, err := vm.abi.Pack("addValidator", operator, node, agent, amount, enterTime, name, description, pubKey)
		assert.NoError(t, err)
		returnData := []byte("add-result")
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vm.AddValidator(ctx, caller, operator, node, agent, amount, enterTime, name, description, pubKey)
		assert.NoError(t, err)
		assert.Equal(t, returnData, result)
	})

	// Test case: ABI pack failure
	t.Run("ABIPackFailure", func(t *testing.T) {
		result, err := vm.AddValidator(ctx, caller, operator, node, agent, nil, enterTime, name, description, pubKey)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to pack ABI")
		assert.Nil(t, result)
	})

	// Test case: EVM call failure
	t.Run("EVMCallFailure", func(t *testing.T) {
		inputData, err := vm.abi.Pack("addValidator", operator, node, agent, amount, enterTime, name, description, pubKey)
		assert.NoError(t, err)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return([]byte{}, fmt.Errorf("evm error")).Once()
		result, err := vm.AddValidator(ctx, caller, operator, node, agent, amount, enterTime, name, description, pubKey)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EVM call failed")
		assert.Nil(t, result)
	})
}

func TestUpdateConsensus(t *testing.T) {
	vm, evmKeeper, ctx := setupValidatorManager(t)
	caller := common.HexToAddress("0x123")
	nodes := []common.Address{
		common.HexToAddress("0x456"),
		common.HexToAddress("0x789"),
	}

	// Test case: Successful call
	t.Run("Success", func(t *testing.T) {
		inputData, err := vm.abi.Pack("undateConsensus", nodes)
		assert.NoError(t, err)
		returnData := []byte("update-result")
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vm.UpdateConsensus(ctx, caller, nodes)
		assert.NoError(t, err)
		assert.Equal(t, returnData, result)
	})

	// Test case: ABI pack failure
	t.Run("ABIPackFailure", func(t *testing.T) {
		result, err := vm.UpdateConsensus(ctx, caller, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to pack ABI")
		assert.Nil(t, result)
	})

	// Test case: EVM call failure
	t.Run("EVMCallFailure", func(t *testing.T) {
		inputData, err := vm.abi.Pack("undateConsensus", nodes)
		assert.NoError(t, err)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return([]byte{}, fmt.Errorf("evm error")).Once()
		result, err := vm.UpdateConsensus(ctx, caller, nodes)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EVM call failed")
		assert.Nil(t, result)
	})
}

func TestGetPledgeAmount(t *testing.T) {
	vm, evmKeeper, ctx := setupValidatorManager(t)
	caller := common.HexToAddress("0x123")
	node := common.HexToAddress("0x456")

	// Test case: Successful call
	t.Run("Success", func(t *testing.T) {
		inputData, err := vm.abi.Pack("getPledgeAmount", node)
		assert.NoError(t, err)
		amount := big.NewInt(10000)
		returnData, err := vm.abi.Pack("", amount)
		assert.NoError(t, err)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vm.GetPledgeAmount(ctx, caller, node)
		assert.NoError(t, err)
		assert.Equal(t, amount, result)
	})

	// Test case: ABI pack failure
	t.Run("ABIPackFailure", func(t *testing.T) {
		result, err := vm.GetPledgeAmount(ctx, caller, common.Address{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to pack ABI")
		assert.Nil(t, result)
	})

	// Test case: EVM call failure
	t.Run("EVMCallFailure", func(t *testing.T) {
		inputData, err := vm.abi.Pack("getPledgeAmount", node)
		assert.NoError(t, err)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return([]byte{}, fmt.Errorf("evm error")).Once()
		result, err := vm.GetPledgeAmount(ctx, caller, node)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EVM call failed")
		assert.Nil(t, result)
	})

	// Test case: Unpack failure
	t.Run("UnpackFailure", func(t *testing.T) {
		inputData, err := vm.abi.Pack("getPledgeAmount", node)
		assert.NoError(t, err)
		returnData := []byte{0x00} // Invalid for uint256 type
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestValidatorManagerAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vm.GetPledgeAmount(ctx, caller, node)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unpack return value")
		assert.Nil(t, result)
	})
}
