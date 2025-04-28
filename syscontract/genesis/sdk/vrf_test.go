package sdk

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/evm/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEVMKeeper is a mock implementation of the EVM keeper.
type MockEVMKeeper struct {
	mock.Mock
}

func (m *MockEVMKeeper) CallEVM(ctx sdk.Context, from, to common.Address, value *big.Int, input []byte) ([]byte, error) {
	args := m.Called(ctx, from, to, value, input)
	return args.Get(0).([]byte), args.Error(1)
}

// Constants for testing
const (
	TestVRFAddr = "0x1234567890123456789012345678901234567890"
	TestABI     = `[
		{"name":"getRandomSeed","inputs":[{"type":"uint256"}],"outputs":[{"type":"bytes"}]},
		{"name":"init","inputs":[{"type":"bytes"}],"outputs":[{"type":"bytes"}]},
		{"name":"sendRandom","inputs":[{"type":"bytes"},{"type":"uint256"}],"outputs":[{"type":"bool"}]},
		{"name":"updateAdmin","inputs":[{"type":"address"}],"outputs":[{"type":"bytes"}]},
		{"name":"updateConsensusSet","inputs":[{"type":"uint256"}],"outputs":[{"type":"address[]"}]}
	]`
)

// Setup function to create a VRF instance with mocked dependencies
func setupVRF(t *testing.T) (*VRF, *MockEVMKeeper, sdk.Context) {
	// Parse the test ABI
	parsedABI, err := abi.JSON(strings.NewReader(TestABI))
	if err != nil {
		t.Fatalf("Failed to parse test ABI: %v", err)
	}

	evmKeeper := &MockEVMKeeper{}
	ctx := sdk.Context{}
	vrf := &VRF{
		abi:       parsedABI,
		evmKeeper: &keeper.Keeper{},
	}
	return vrf, evmKeeper, ctx
}

func TestNewVRF(t *testing.T) {
	// Test case: Successful creation
	t.Run("Success", func(t *testing.T) {
		vrf, err := NewVRF(&keeper.Keeper{})
		assert.NoError(t, err)
		assert.NotNil(t, vrf)
		assert.NotNil(t, vrf.abi)
		assert.NotNil(t, vrf.evmKeeper)
	})

	// Test case: Invalid ABI
	t.Run("InvalidABI", func(t *testing.T) {
		// Create a VRF with an invalid ABI directly
		invalidABI := "invalid json"
		_, err := abi.JSON(strings.NewReader(invalidABI))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid character")
	})
}

func TestGetRandomSeed(t *testing.T) {
	vrf, evmKeeper, ctx := setupVRF(t)
	caller := common.HexToAddress("0x123")
	epoch := big.NewInt(1234)

	// Test case: Successful call
	t.Run("Success", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("getRandomSeed", epoch)
		returnData := []byte("random-seed")
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vrf.GetRandomSeed(ctx, caller, epoch)
		assert.NoError(t, err)
		assert.Equal(t, returnData, result)
	})

	// Test case: ABI pack failure
	t.Run("ABIPackFailure", func(t *testing.T) {
		// Invalid epoch type to cause pack failure
		result, err := vrf.GetRandomSeed(ctx, caller, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to pack ABI")
		assert.Nil(t, result)
	})

	// Test case: EVM call failure
	t.Run("EVMCallFailure", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("getRandomSeed", epoch)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return([]byte{}, fmt.Errorf("evm error")).Once()
		result, err := vrf.GetRandomSeed(ctx, caller, epoch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EVM call failed")
		assert.Nil(t, result)
	})

	// Test case: Unpack failure
	t.Run("UnpackFailure", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("getRandomSeed", epoch)
		// Invalid return data to cause unpack failure
		returnData := []byte{0x00} // Invalid for bytes type
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vrf.GetRandomSeed(ctx, caller, epoch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unpack return value")
		assert.Nil(t, result)
	})
}

func TestInit(t *testing.T) {
	vrf, evmKeeper, ctx := setupVRF(t)
	caller := common.HexToAddress("0x123")
	rnd := []byte("random-seed")

	// Test case: Successful initialization
	t.Run("Success", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("init", rnd)
		returnData := []byte("init-result")
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vrf.Init(ctx, caller, rnd)
		assert.NoError(t, err)
		assert.Equal(t, returnData, result)
	})

	// Test case: ABI pack failure
	t.Run("ABIPackFailure", func(t *testing.T) {
		// Invalid rnd type to cause pack failure
		result, err := vrf.Init(ctx, caller, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to pack ABI")
		assert.Nil(t, result)
	})

	// Test case: EVM call failure
	t.Run("EVMCallFailure", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("init", rnd)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return([]byte{}, fmt.Errorf("evm error")).Once()
		result, err := vrf.Init(ctx, caller, rnd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EVM call failed")
		assert.Nil(t, result)
	})
}

func TestSendRandom(t *testing.T) {
	vrf, evmKeeper, ctx := setupVRF(t)
	caller := common.HexToAddress("0x123")
	rnd := []byte("random-value")
	epoch := big.NewInt(1234)

	// Test case: Successful send
	t.Run("Success", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("sendRandom", rnd, epoch)
		returnData, _ := vrf.abi.Pack("", true) // ABI-encoded true
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vrf.SendRandom(ctx, caller, rnd, epoch)
		assert.NoError(t, err)
		assert.True(t, result)
	})

	// Test case: ABI pack failure
	t.Run("ABIPackFailure", func(t *testing.T) {
		result, err := vrf.SendRandom(ctx, caller, nil, epoch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to pack ABI")
		assert.False(t, result)
	})

	// Test case: EVM call failure
	t.Run("EVMCallFailure", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("sendRandom", rnd, epoch)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return([]byte{}, fmt.Errorf("evm error")).Once()
		result, err := vrf.SendRandom(ctx, caller, rnd, epoch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EVM call failed")
		assert.False(t, result)
	})

	// Test case: Unpack failure
	t.Run("UnpackFailure", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("sendRandom", rnd, epoch)
		returnData := []byte{0x00} // Invalid for bool type
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vrf.SendRandom(ctx, caller, rnd, epoch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unpack return value")
		assert.False(t, result)
	})
}

func TestUpdateAdmin(t *testing.T) {
	vrf, evmKeeper, ctx := setupVRF(t)
	caller := common.HexToAddress("0x123")
	admin := common.HexToAddress("0x456")

	// Test case: Successful update
	t.Run("Success", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("updateAdmin", admin)
		returnData := []byte("update-result")
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vrf.UpdateAdmin(ctx, caller, admin)
		assert.NoError(t, err)
		assert.Equal(t, returnData, result)
	})

	// Test case: ABI pack failure
	t.Run("ABIPackFailure", func(t *testing.T) {
		// Invalid admin type to cause pack failure
		result, err := vrf.UpdateAdmin(ctx, caller, common.Address{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to pack ABI")
		assert.Nil(t, result)
	})

	// Test case: EVM call failure
	t.Run("EVMCallFailure", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("updateAdmin", admin)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return([]byte{}, fmt.Errorf("evm error")).Once()
		result, err := vrf.UpdateAdmin(ctx, caller, admin)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EVM call failed")
		assert.Nil(t, result)
	})
}

func TestUpdateConsensusSet(t *testing.T) {
	vrf, evmKeeper, ctx := setupVRF(t)
	caller := common.HexToAddress("0x123")
	epoch := big.NewInt(1235)

	// Test case: Successful update
	t.Run("Success", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("updateConsensusSet", epoch)
		nodes := []common.Address{common.HexToAddress("0x456"), common.HexToAddress("0x789")}
		returnData, _ := vrf.abi.Pack("", nodes)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vrf.UpdateConsensusSet(ctx, caller, epoch)
		assert.NoError(t, err)
		assert.Equal(t, nodes, result)
	})

	// Test case: EVM call returns nil data and error
	t.Run("EVMCallNil", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("updateConsensusSet", epoch)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return([]byte{}, fmt.Errorf("evm error")).Once()
		result, err := vrf.UpdateConsensusSet(ctx, caller, epoch)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	// Test case: Nil return data
	t.Run("NilReturnData", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("updateConsensusSet", epoch)
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return([]byte{}, nil).Once()
		result, err := vrf.UpdateConsensusSet(ctx, caller, epoch)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	// Test case: ABI pack failure
	t.Run("ABIPackFailure", func(t *testing.T) {
		result, err := vrf.UpdateConsensusSet(ctx, caller, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to pack ABI")
		assert.Nil(t, result)
	})

	// Test case: Unpack failure
	t.Run("UnpackFailure", func(t *testing.T) {
		inputData, _ := vrf.abi.Pack("updateConsensusSet", epoch)
		returnData := []byte{0x00} // Invalid for address[] type
		evmKeeper.On("CallEVM", ctx, caller, common.HexToAddress(TestVRFAddr), (*big.Int)(nil), inputData).Return(returnData, nil).Once()
		result, err := vrf.UpdateConsensusSet(ctx, caller, epoch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unpack return value")
		assert.Nil(t, result)
	})
}

func TestParseRevert(t *testing.T) {
	// Test case: Valid revert data
	t.Run("ValidRevert", func(t *testing.T) {
		// Construct revert data: selector (08c379a0) + offset (32) + length (32) + message
		message := []byte("revert reason")
		length := big.NewInt(int64(len(message)))
		lengthBytes := make([]byte, 32)
		length.FillBytes(lengthBytes)
		data := append([]byte{0x08, 0xc3, 0x79, 0xa0}, make([]byte, 32)...) // Selector + offset
		data = append(data, lengthBytes...)
		data = append(data, message...)
		result, err := ParseRevert(data)
		assert.NoError(t, err)
		assert.Equal(t, message, result)
	})

	// Test case: Invalid selector
	t.Run("InvalidSelector", func(t *testing.T) {
		data := []byte{0x00, 0x00, 0x00, 0x00, 0x00} // Wrong selector
		result, err := ParseRevert(data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse revert data")
		assert.Nil(t, result)
	})

	// Test case: Data too short
	t.Run("DataTooShort", func(t *testing.T) {
		data := []byte{0x08, 0xc3, 0x79, 0xa0} // Only selector
		result, err := ParseRevert(data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse revert data")
		assert.Nil(t, result)
	})
}
