package epoch

import (
	"context"
	"encoding/json"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	"github.com/cosmos/cosmos-sdk/telemetry"
	evmKeeper "github.com/cosmos/cosmos-sdk/x/evm/keeper"
	"github.com/eni-chain/go-eni/syscontract"
	syscontractSdk "github.com/eni-chain/go-eni/syscontract/genesis/sdk"
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	// this line is used by starport scaffolding # 1

	modulev1 "github.com/eni-chain/go-eni/api/goeni/epoch/module"
	"github.com/eni-chain/go-eni/x/epoch/keeper"
	"github.com/eni-chain/go-eni/x/epoch/types"
)

var (
	_ module.AppModuleBasic      = (*AppModule)(nil)
	_ module.AppModuleSimulation = (*AppModule)(nil)
	_ module.HasGenesis          = (*AppModule)(nil)
	_ module.HasInvariants       = (*AppModule)(nil)
	_ module.HasConsensusVersion = (*AppModule)(nil)
	_ module.HasABCIEndBlock     = (*AppModule)(nil)

	_ appmodule.AppModule       = (*AppModule)(nil)
	_ appmodule.HasBeginBlocker = (*AppModule)(nil)
	//_ appmodule.HasEndBlocker   = (*AppModule)(nil)
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface that defines the
// independent methods a Cosmos SDK module needs to implement.
type AppModuleBasic struct {
	cdc codec.BinaryCodec
}

func NewAppModuleBasic(cdc codec.BinaryCodec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the name of the module as a string.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the amino codec for the module, which is used
// to marshal and unmarshal structs to/from []byte in order to persist them in the module's KVStore.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// RegisterInterfaces registers a module's interface types and their concrete implementations as proto.Message.
func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// DefaultGenesis returns a default GenesisState for the module, marshalled to json.RawMessage.
// The default GenesisState need to be defined by the module developer and is primarily used for testing.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis used to validate the GenesisState, given in its json.RawMessage form.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface that defines the inter-dependent methods that modules need to implement
type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	EvmKeeper     *evmKeeper.Keeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	EvmKeeper *evmKeeper.Keeper,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		EvmKeeper:      EvmKeeper,
	}
}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// RegisterInvariants registers the invariants of the module. If an invariant deviates from its predicted value, the InvariantRegistry triggers appropriate logic (most often the chain will be halted)
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the module's genesis initialization. It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) {

	/******************init system contract*************************/
	if ctx.BlockHeight() == 0 {
		syscontract.SetupSystemContracts(ctx, am.EvmKeeper)
	}

	/******************init epoch****************************/
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	cdc.MustUnmarshalJSON(gs, &genState)

	if genState.GetEpoch() == nil {
		epoch := types.Epoch{
			GenesisTime:             ctx.BlockTime(),
			EpochInterval:           50,
			CurrentEpoch:            1,
			CurrentEpochStartHeight: 1,
			CurrentEpochHeight:      1,
		}
		am.keeper.SetEpoch(ctx, epoch)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeNewEpoch,
				sdk.NewAttribute(types.AttributeEpochNumber, fmt.Sprint(epoch.CurrentEpoch)),
				sdk.NewAttribute(types.AttributeEpochTime, fmt.Sprint(epoch.CurrentEpochStartHeight)),
				sdk.NewAttribute(types.AttributeEpochHeight, fmt.Sprint(epoch.CurrentEpochHeight)),
			),
		)

		telemetry.SetGauge(float32(epoch.CurrentEpoch), "epoch", "current")
	} else {
		InitGenesis(ctx, am.keeper, genState)
	}
}

// ExportGenesis returns the module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(genState)
}

// ConsensusVersion is a sequence number for state-breaking change of the module.
// It should be incremented on each consensus-breaking change introduced by the module.
// To avoid wrong/empty versions, the initial version should be set to 1.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock contains the logic that is automatically triggered at the beginning of each block.
// The begin block implementation is optional.
func (am AppModule) BeginBlock(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	lastEpoch := am.keeper.GetEpoch(sdkCtx)

	if uint64(sdkCtx.BlockHeight())-(lastEpoch.CurrentEpochStartHeight) >= lastEpoch.EpochInterval {
		//am.keeper.AfterEpochEnd(ctx, lastEpoch)
		newEpoch := types.Epoch{
			GenesisTime:             lastEpoch.GenesisTime,
			EpochInterval:           lastEpoch.EpochInterval,
			CurrentEpoch:            lastEpoch.CurrentEpoch + 1,
			CurrentEpochStartHeight: uint64(sdkCtx.BlockHeight()),
			CurrentEpochHeight:      sdkCtx.BlockHeight(),
		}
		am.keeper.SetEpoch(sdkCtx, newEpoch)

		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeNewEpoch,
				sdk.NewAttribute(types.AttributeEpochNumber, fmt.Sprint(newEpoch.CurrentEpoch)),
				sdk.NewAttribute(types.AttributeEpochTime, fmt.Sprint(newEpoch.CurrentEpochStartHeight)),
				sdk.NewAttribute(types.AttributeEpochHeight, fmt.Sprint(newEpoch.CurrentEpochHeight)),
			),
		)

		telemetry.SetGauge(float32(newEpoch.CurrentEpoch), "epoch", "current")
		//am.keeper.BeforeEpochStart(ctx, newEpoch)
	}
	return nil
}

// EndBlock contains the logic that is automatically triggered at the end of each block.
// The end block implementation is optional.
func (am AppModule) EndBlock(goCtx context.Context) ([]abci.ValidatorUpdate, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//The last block of the epoch updates the consensus set for the next epoch
	epoch := am.keeper.GetEpoch(ctx)
	if epoch.EpochInterval == 0 {
		return nil, nil
	}

	if uint64(ctx.BlockHeight())%epoch.EpochInterval != 0 {
		return nil, nil
	}

	vrf, err := syscontractSdk.NewVRF(am.EvmKeeper)
	if err != nil {
		return nil, err
	}

	addr := am.EvmKeeper.AccountKeeper().GetModuleAddress(authtypes.FeeCollectorName)
	caller := common.Address(addr)
	epochNum := big.NewInt(int64(epoch.CurrentEpoch))

	addrs, err := vrf.UpdateConsensusSet(ctx, caller, epochNum)
	if err != nil {
		return nil, err
	}

	valSet, err := syscontractSdk.NewValidatorManager(am.EvmKeeper)
	if err != nil {
		return nil, err
	}

	pubKeys, err := valSet.GetPubKeysBySequence(ctx, caller, addrs)
	if err != nil {
		return nil, err
	}

	validatorSet := make([]abci.ValidatorUpdate, len(pubKeys))
	for i := 0; i < len(pubKeys); i++ {
		//innerPk := crypto.PublicKey_Ed25519{Ed25519: pkBytes}
		//pubKey := crypto.PublicKey{Sum: &innerPk}
		pk := crypto.PublicKey_Ed25519{Ed25519: pubKeys[i]}
		pubKey := crypto.PublicKey{Sum: &pk}
		validatorSet[i].PubKey = pubKey
		validatorSet[i].Power = 1
	}

	return validatorSet, nil
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// ----------------------------------------------------------------------------
// App Wiring Setup
// ----------------------------------------------------------------------------

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	StoreService store.KVStoreService
	Cdc          codec.Codec
	Config       *modulev1.Module
	Logger       log.Logger

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
	EvmKeeper     *evmKeeper.Keeper
}

type ModuleOutputs struct {
	depinject.Out

	EpochKeeper keeper.Keeper
	Module      appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}
	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.Logger,
		authority.String(),
	)
	m := NewAppModule(
		in.Cdc,
		k,
		in.AccountKeeper,
		in.BankKeeper,
		in.EvmKeeper,
	)

	return ModuleOutputs{EpochKeeper: k, Module: m}
}
