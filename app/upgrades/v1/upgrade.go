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
	TestnetHeight      int64 = 27820000
	TestnetOldContract = "0x642Ad927f251cE7c4796A1dA56caa25F8FEdB5eE"
	TestnetNewContract = "0x9e2EB483aB9873cbb84c47720f8fed61af5b3d8F"
)

// Mainnet constants (ChainID 173)
const (
	MainnetUpgradeName = "v1-bytecode-redirect-mainnet"
	MainnetChainID     = "ENI Mainnet"
	MainnetHeight      int64 = 29270000
	MainnetOldContract = "0x002Ad927f251cE7c4796A1dA56caa25F8FEdB5eE"
	MainnetNewContract = "0xFF2EB483aB9873cbb84c47720f8fed61af5b3d8F"
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
	oldAddr := common.HexToAddress(h.oldContractAddr)
	newAddr := common.HexToAddress(h.newContractAddr)

	// Read the correct bytecode from the new contract
	code := h.evmKeeper.GetCode(ctx, newAddr)
	if len(code) == 0 {
		return fmt.Errorf("source contract %s has no bytecode, cannot perform bytecode redirect", h.newContractAddr)
	}

	ctx.Logger().Info(fmt.Sprintf(
		"Hard fork v1: redirecting bytecode from %s to %s (code size: %d bytes)",
		h.newContractAddr, h.oldContractAddr, len(code),
	))

	// Write the correct bytecode to the old contract address
	h.evmKeeper.SetCode(ctx, oldAddr, code)

	ctx.Logger().Info(fmt.Sprintf(
		"Hard fork v1: bytecode redirect completed for %s",
		h.oldContractAddr,
	))

	return nil
}
