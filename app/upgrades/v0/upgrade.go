package v0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eni-chain/go-eni/app/upgrades"
)

const (
	UpgradeName = "v0"
)

// HardForkUpgradeHandler defines an example hard fork handler that will be
// executed during BeginBlock at a target height and chain-ID.
type HardForkUpgradeHandler struct {
	TargetHeight  int64
	TargetChainID string
	//WasmKeeper    wasm.Keeper
}

func NewHardForkUpgradeHandler(height int64, chainID string) upgrades.HardForkHandler {
	return HardForkUpgradeHandler{
		TargetHeight:  height,
		TargetChainID: chainID,
		//WasmKeeper:    wk,
	}
}

func (h HardForkUpgradeHandler) GetName() string {
	return UpgradeName
}

func (h HardForkUpgradeHandler) GetTargetChainID() string {
	return h.TargetChainID
}

func (h HardForkUpgradeHandler) GetTargetHeight() int64 {
	return h.TargetHeight
}

func (h HardForkUpgradeHandler) ExecuteHandler(ctx sdk.Context) error {
	//govKeeper := wasmkeeper.NewGovPermissionKeeper(h.WasmKeeper)
	// If other contract need to be migrated, create functions for them and pass
	// the govKeeper to them.
	//return h.migrateGringotts(ctx, govKeeper)
	return nil
}
