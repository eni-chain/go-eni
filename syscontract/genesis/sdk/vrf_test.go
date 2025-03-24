package sdk

//
//import (
//	cosmossdk_io_math "cosmossdk.io/math"
//	"math/big"
//	"testing"
//
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/eni-chain/go-eni/syscontract"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//)
//
//func TestNewVRF(t *testing.T) {
//	mockKeeper := new(MockEvmKeeper)
//	vrf, err := NewVRF(mockKeeper)
//
//	assert.NoError(t, err)
//	assert.NotNil(t, vrf)
//	assert.Equal(t, mockKeeper, vrf.evmKeeper)
//}
//
//func TestGetRandomSeed(t *testing.T) {
//	mockKeeper := new(MockEvmKeeper)
//	vrf, _ := NewVRF(mockKeeper)
//
//	ctx := sdk.Context{}
//	caller := common.HexToAddress("0x1234567890123456789012345678901234567890")
//	epoch := big.NewInt(123)
//
//	// Expected random seed
//	expectedSeed := []byte("random-seed-data")
//
//	// Mock the return data that would be unpacked to the random seed
//	returnData := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x2d, 0x73, 0x65, 0x65, 0x64, 0x2d, 0x64, 0x61, 0x74, 0x61}
//
//	// Set up expectations
//	vrfAddr := common.HexToAddress(syscontract.VRFAddr)
//	mockKeeper.On("CallEVM", ctx, caller, &vrfAddr, (*cosmossdk_io_math.Int)(nil), mock.Anything).Return(returnData, nil)
//
//	// Call the method
//	seed, err := vrf.GetRandomSeed(ctx, caller, epoch)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Equal(t, expectedSeed, seed)
//	mockKeeper.AssertExpectations(t)
//}
//
//func TestInit(t *testing.T) {
//	mockKeeper := new(MockEvmKeeper)
//	vrf, _ := NewVRF(mockKeeper)
//
//	ctx := sdk.Context{}
//	caller := common.HexToAddress("0x1234567890123456789012345678901234567890")
//	rnd := []byte("initial-random-seed")
//
//	// Mock the return data
//	returnData := []byte{0x01} // Simplified for test
//
//	// Set up expectations
//	vrfAddr := common.HexToAddress(syscontract.VRFAddr)
//	mockKeeper.On("CallEVM", ctx, caller, &vrfAddr, (*cosmossdk_io_math.Int)(nil), mock.Anything).Return(returnData, nil)
//
//	// Call the method
//	result, err := vrf.Init(ctx, caller, rnd)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Equal(t, returnData, result)
//	mockKeeper.AssertExpectations(t)
//}
//
//func TestSendRandom(t *testing.T) {
//	mockKeeper := new(MockEvmKeeper)
//	vrf, _ := NewVRF(mockKeeper)
//
//	ctx := sdk.Context{}
//	caller := common.HexToAddress("0x1234567890123456789012345678901234567890")
//	rnd := []byte("random-data")
//	epoch := big.NewInt(123)
//
//	// Expected success value
//	expectedSuccess := true
//
//	// Mock the return data that would be unpacked to the success value
//	returnData := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
//
//	// Set up expectations
//	vrfAddr := common.HexToAddress(syscontract.VRFAddr)
//	mockKeeper.On("CallEVM", ctx, caller, &vrfAddr, (*cosmossdk_io_math.Int)(nil), mock.Anything).Return(returnData, nil)
//
//	// Call the method
//	success, err := vrf.SendRandom(ctx, caller, rnd, epoch)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Equal(t, expectedSuccess, success)
//	mockKeeper.AssertExpectations(t)
//}
//
//func TestUpdateAdmin_VRF(t *testing.T) {
//	mockKeeper := new(MockEvmKeeper)
//	vrf, _ := NewVRF(mockKeeper)
//
//	ctx := sdk.Context{}
//	caller := common.HexToAddress("0x1234567890123456789012345678901234567890")
//	admin := common.HexToAddress("0x2345678901234567890123456789012345678901")
//
//	// Mock the return data
//	returnData := []byte{0x01} // Simplified for test
//
//	// Set up expectations
//	vrfAddr := common.HexToAddress(syscontract.VRFAddr)
//	mockKeeper.On("CallEVM", ctx, caller, &vrfAddr, (*cosmossdk_io_math.Int)(nil), mock.Anything).Return(returnData, nil)
//
//	// Call the method
//	result, err := vrf.UpdateAdmin(ctx, caller, admin)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Equal(t, returnData, result)
//	mockKeeper.AssertExpectations(t)
//}
//
//func TestUpdateConsensusSet(t *testing.T) {
//	mockKeeper := new(MockEvmKeeper)
//	vrf, _ := NewVRF(mockKeeper)
//
//	ctx := sdk.Context{}
//	caller := common.HexToAddress("0x1234567890123456789012345678901234567890")
//	epoch := big.NewInt(123)
//
//	// Expected nodes
//	expectedNodes := []common.Address{
//		common.HexToAddress("0x1111111111111111111111111111111111111111"),
//		common.HexToAddress("0x2222222222222222222222222222222222222222"),
//	}
//
//	// Mock the return data that would be unpacked to the nodes
//	returnData := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22}
//
//	// Set up expectations
//	vrfAddr := common.HexToAddress(syscontract.VRFAddr)
//	mockKeeper.On("CallEVM", ctx, caller, &vrfAddr, (*cosmossdk_io_math.Int)(nil), mock.Anything).Return(returnData, nil)
//
//	// Call the method
//	nodes, err := vrf.UpdateConsensusSet(ctx, caller, epoch)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Equal(t, expectedNodes, nodes)
//	mockKeeper.AssertExpectations(t)
//}
