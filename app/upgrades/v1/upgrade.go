package v1

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/app/upgrades"
	"github.com/ethereum/go-ethereum/common"
)

// Testnet constants (ChainID 174)
const (
	TestnetUpgradeName = "v1-bytecode-redirect-testnet"
	TestnetChainID     = "ENI Testnet"
	TestnetHeight      int64 = 28845000
	TestnetOldContract = "0x16E7a7cD80fE7c61311Ce079Bcea7B993FDD02B4"
	TestnetNewContract = "0xfA165C2C58A8A22bA0508a1CF34cA79E1d891F0b"
)

// Mainnet constants (ChainID 173)
const (
	MainnetUpgradeName = "v1-bytecode-redirect-mainnet"
	MainnetChainID     = "ENI Mainnet"
	MainnetHeight      int64 = 30028500
	MainnetOldContract = "0x3ba1da3C4ab549B1816D8050E4368440b92d92A0"
	MainnetNewContract = "0x5e63338AfCE7d5a4d97Da1714f6a6bfA5A26F0f9"
)

// EVMCodeKeeper defines the interface for reading/writing contract bytecode.
type EVMCodeKeeper interface {
	GetCode(ctx sdk.Context, addr common.Address) []byte
	SetCode(ctx sdk.Context, addr common.Address, code []byte)
}

// HardForkUpgradeHandler redirects bytecode from a source contract address to a target contract address.
type HardForkUpgradeHandler struct {
	name            string
	targetHeight    int64
	targetChainID   string
	oldContractAddr string
	newContractAddr string
	evmKeeper       EVMCodeKeeper
}

// NewTestnetHandler creates a hard fork handler for testnet (ChainID 174, height 27820000).
func NewTestnetHandler(evmKeeper EVMCodeKeeper) upgrades.HardForkHandler {
	return &HardForkUpgradeHandler{
		name:            TestnetUpgradeName,
		targetHeight:    TestnetHeight,
		targetChainID:   TestnetChainID,
		oldContractAddr: TestnetOldContract,
		newContractAddr: TestnetNewContract,
		evmKeeper:       evmKeeper,
	}
}

// NewMainnetHandler creates a hard fork handler for mainnet (ChainID 173, height 29270000).
func NewMainnetHandler(evmKeeper EVMCodeKeeper) upgrades.HardForkHandler {
	return &HardForkUpgradeHandler{
		name:            MainnetUpgradeName,
		targetHeight:    MainnetHeight,
		targetChainID:   MainnetChainID,
		oldContractAddr: MainnetOldContract,
		newContractAddr: MainnetNewContract,
		evmKeeper:       evmKeeper,
	}
}

func (h *HardForkUpgradeHandler) GetName() string {
	return h.name
}

func (h *HardForkUpgradeHandler) GetTargetChainID() string {
	return h.targetChainID
}

func (h *HardForkUpgradeHandler) GetTargetHeight() int64 {
	return h.targetHeight
}

func (h *HardForkUpgradeHandler) ExecuteHandler(ctx sdk.Context) error {
	logger := ctx.Logger()
	oldAddr := common.HexToAddress(h.oldContractAddr)
	newAddr := common.HexToAddress(h.newContractAddr)

	logger.Info("========== HARD FORK BEGIN ==========")
	logger.Info(fmt.Sprintf("[HardFork] handler=%s chainID=%s height=%d",
		h.name, h.targetChainID, ctx.BlockHeight()))
	logger.Info(fmt.Sprintf("[HardFork] old contract (target): %s", h.oldContractAddr))
	logger.Info(fmt.Sprintf("[HardFork] new contract (source): %s", h.newContractAddr))

	// Read existing bytecode of the old contract before replacement
	oldCode := h.evmKeeper.GetCode(ctx, oldAddr)
	logger.Info(fmt.Sprintf("[HardFork] old contract current bytecode size: %d bytes", len(oldCode)))

	// Read the correct bytecode from the new contract
	code := h.evmKeeper.GetCode(ctx, newAddr)
	if len(code) == 0 {
		logger.Error(fmt.Sprintf("[HardFork] FAILED: source contract %s has no bytecode", h.newContractAddr))
		return fmt.Errorf("source contract %s has no bytecode, cannot perform bytecode redirect", h.newContractAddr)
	}
	logger.Info(fmt.Sprintf("[HardFork] new contract bytecode size: %d bytes", len(code)))

	// Write the correct bytecode to the old contract address
	h.evmKeeper.SetCode(ctx, oldAddr, code)
	logger.Info(fmt.Sprintf("[HardFork] SetCode executed: wrote %d bytes to %s", len(code), h.oldContractAddr))

	// Verify the write was successful
	verifyCode := h.evmKeeper.GetCode(ctx, oldAddr)
	if len(verifyCode) != len(code) {
		logger.Error(fmt.Sprintf("[HardFork] VERIFICATION FAILED: expected %d bytes, got %d bytes at %s",
			len(code), len(verifyCode), h.oldContractAddr))
		return fmt.Errorf("bytecode verification failed for %s: expected %d bytes, got %d bytes",
			h.oldContractAddr, len(code), len(verifyCode))
	}
	logger.Info(fmt.Sprintf("[HardFork] VERIFICATION PASSED: %s now has %d bytes of bytecode", h.oldContractAddr, len(verifyCode)))
	logger.Info("========== HARD FORK COMPLETE ==========")

	return nil
}
