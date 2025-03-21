package sdk

//
//import (
//	"math/big"
//	"testing"
//
//	cosmossdk_io_math "cosmossdk.io/math"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/eni-chain/go-eni/syscontract"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//)
//
//func TestNewHub(t *testing.T) {
//	mockKeeper := new(MockEvmKeeper)
//	hub, err := NewHub(mockKeeper)
//
//	assert.NoError(t, err)
//	assert.NotNil(t, hub)
//	assert.Equal(t, mockKeeper, hub.evmKeeper)
//}
//
//func TestApplyForValidator(t *testing.T) {
//	mockKeeper := new(MockEvmKeeper)
//	hub, _ := NewHub(mockKeeper)
//
//	ctx := sdk.Context{}
//	caller := common.HexToAddress("0x1234567890123456789012345678901234567890")
//	node := common.HexToAddress("0x2345678901234567890123456789012345678901")
//	agent := common.HexToAddress("0x3456789012345678901234567890123456789012")
//	name := "validator-1"
//	description := "Test validator"
//	pubKey := []byte("test-pubkey")
//	value := cosmossdk_io_math.NewInt(1000000)
//
//	// Mock the return data
//	returnData := []byte{0x01} // Simplified for test
//
//	// Set up expectations
//	hubAddr := common.HexToAddress(syscontract.HubAddr)
//	mockKeeper.On("CallEVM", ctx, caller, &hubAddr, &value, mock.Anything).Return(returnData, nil)
//
//	// Call the method
//	result, err := hub.ApplyForValidator(ctx, caller, node, agent, name, description, pubKey, &value)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Equal(t, returnData, result)
//	mockKeeper.AssertExpectations(t)
//}
//
//func TestAuditPass(t *testing.T) {
//	mockKeeper := new(MockEvmKeeper)
//	hub, _ := NewHub(mockKeeper)
//
//	ctx := sdk.Context{}
//	caller := common.HexToAddress("0x1234567890123456789012345678901234567890")
//	operator := common.HexToAddress("0x2345678901234567890123456789012345678901")
//
//	// Mock the return data
//	returnData := []byte{0x01} // Simplified for test
//
//	// Set up expectations
//	hubAddr := common.HexToAddress(syscontract.HubAddr)
//	mockKeeper.On("CallEVM", ctx, caller, &hubAddr, (*cosmossdk_io_math.Int)(nil), mock.Anything).Return(returnData, nil)
//
//	// Call the method
//	result, err := hub.AuditPass(ctx, caller, operator)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Equal(t, returnData, result)
//	mockKeeper.AssertExpectations(t)
//}
//
//func TestBlockReward(t *testing.T) {
//	mockKeeper := new(MockEvmKeeper)
//	hub, _ := NewHub(mockKeeper)
//
//	ctx := sdk.Context{}
//	caller := common.HexToAddress("0x1234567890123456789012345678901234567890")
//	node := common.HexToAddress("0x2345678901234567890123456789012345678901")
//
//	// Expected reward
//	expectedReward := big.NewInt(1000)
//
//	// Mock the return data that would be unpacked to the reward
//	returnData := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xe8}
//
//	// Set up expectations
//	hubAddr := common.HexToAddress(syscontract.HubAddr)
//	mockKeeper.On("CallEVM", ctx, caller, &hubAddr, (*cosmossdk_io_math.Int)(nil), mock.Anything).Return(returnData, nil)
//
//	// Call the method
//	reward, err := hub.BlockReward(ctx, caller, node)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Equal(t, expectedReward, reward)
//	mockKeeper.AssertExpectations(t)
//}
//
//func TestUpdateAdmin(t *testing.T) {
//	mockKeeper := new(MockEvmKeeper)
//	hub, _ := NewHub(mockKeeper)
//
//	ctx := sdk.Context{}
//	caller := common.HexToAddress("0x1234567890123456789012345678901234567890")
//	admin := common.HexToAddress("0x2345678901234567890123456789012345678901")
//
//	// Mock the return data
//	returnData := []byte{0x01} // Simplified for test
//
//	// Set up expectations
//	hubAddr := common.HexToAddress(syscontract.HubAddr)
//	mockKeeper.On("CallEVM", ctx, caller, &hubAddr, (*cosmossdk_io_math.Int)(nil), mock.Anything).Return(returnData, nil)
//
//	// Call the method
//	result, err := hub.UpdateAdmin(ctx, caller, admin)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Equal(t, returnData, result)
//	mockKeeper.AssertExpectations(t)
//}
