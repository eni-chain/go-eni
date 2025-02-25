package evm

import (
	"context"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	cosmossdk_io_math "cosmossdk.io/math"
	"encoding/json"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/eni-chain/go-eni/utils"
	"github.com/eni-chain/go-eni/x/evm/state"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	// this line is used by starport scaffolding # 1

	modulev1 "github.com/eni-chain/go-eni/api/goeni/evm/module"
	"github.com/eni-chain/go-eni/x/evm/keeper"
	"github.com/eni-chain/go-eni/x/evm/types"
)

var (
	_ module.AppModuleBasic      = (*AppModule)(nil)
	_ module.AppModuleSimulation = (*AppModule)(nil)
	_ module.HasGenesis          = (*AppModule)(nil)
	_ module.HasInvariants       = (*AppModule)(nil)
	_ module.HasConsensusVersion = (*AppModule)(nil)

	_ appmodule.AppModule       = (*AppModule)(nil)
	_ appmodule.HasBeginBlocker = (*AppModule)(nil)
	_ appmodule.HasEndBlocker   = (*AppModule)(nil)
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
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

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
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
	}
}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(&am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// RegisterInvariants registers the invariants of the module. If an invariant deviates from its predicted value, the InvariantRegistry triggers appropriate logic (most often the chain will be halted)
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the module's genesis initialization. It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) {
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	cdc.MustUnmarshalJSON(gs, &genState)

	InitGenesis(ctx, am.keeper, genState)
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
func (am AppModule) BeginBlock(_ context.Context) error {
	// clear tx/tx responses from last block
	am.keeper.SetMsgs([]*types.MsgEVMTransaction{})
	am.keeper.SetTxResults([]*abci.ExecTxResult{})
	return nil
}

// EndBlock contains the logic that is automatically triggered at the end of each block.
// The end block implementation is optional.
func (am AppModule) EndBlock(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	//newBaseFee := am.keeper.AdjustDynamicBaseFeePerGas(ctx, uint64(req.BlockGasUsed))
	//if newBaseFee != nil {
	//	metrics.GaugeEvmBlockBaseFee(newBaseFee.TruncateInt().BigInt(), req.Height)
	//}
	var _ sdk.AccAddress // to avoid unused error
	if am.keeper.EthBlockTestConfig.Enabled {
		blocks := am.keeper.BlockTest.Json.Blocks
		block, err := blocks[ctx.BlockHeight()-1].Decode()
		if err != nil {
			panic(err)
		}
		_ = am.keeper.GetEniAddressOrDefault(ctx, block.Header().Coinbase)
	} else {
		_ = am.keeper.AccountKeeper().GetModuleAddress(authtypes.FeeCollectorName)
	}
	evmTxDeferredInfoList := am.keeper.GetAllEVMTxDeferredInfo(ctx)
	denom := am.keeper.GetBaseDenom(ctx)
	surplus := am.keeper.GetAnteSurplusSum(ctx)
	for _, deferredInfo := range evmTxDeferredInfoList {
		txHash := common.BytesToHash(deferredInfo.TxHash)
		if deferredInfo.Error != "" && txHash.Cmp(ethtypes.EmptyTxsHash) != 0 {
			_ = am.keeper.SetTransientReceipt(ctx, txHash, &types.Receipt{
				TxHashHex:        txHash.Hex(),
				TransactionIndex: deferredInfo.TxIndex,
				VmError:          deferredInfo.Error,
				BlockNumber:      uint64(ctx.BlockHeight()),
			})
			continue
		}
		idx := int(deferredInfo.TxIndex)
		coinbaseAddress := state.GetCoinbaseAddress(idx)
		balance := am.keeper.BankKeeper().SpendableCoins(ctx, coinbaseAddress).AmountOf(denom)
		//weiBalance := am.keeper.BankKeeper().GetWeiBalance(ctx, coinbaseAddress)
		weiBalance := am.keeper.BankKeeper().GetBalance(ctx, coinbaseAddress, denom)
		if !balance.IsZero() || !weiBalance.IsZero() {
			// todo  check code correct
			//if err := am.keeper.BankKeeper().SendCoinsAndWei(ctx, coinbaseAddress, coinbase, balance, weiBalance); err != nil {
			//	ctx.Logger().Error(fmt.Sprintf("failed to send ueni surplus from %s to coinbase account due to %s", coinbaseAddress.String(), err))
			//}
		}
		surplus = surplus.Add(deferredInfo.Surplus)
	}
	if surplus.IsPositive() {
		surplusUeni, surplusWei := state.SplitUeniWeiAmount(surplus.BigInt())
		if surplusUeni.GT(cosmossdk_io_math.ZeroInt()) {
			// todo  check code correct
			//if err := am.keeper.BankKeeper().AddCoins(ctx, am.keeper.AccountKeeper().GetModuleAddress(types.ModuleName), sdk.NewCoins(sdk.NewCoin(am.keeper.GetBaseDenom(ctx), surplusUeni)), true); err != nil {
			//	ctx.Logger().Error("failed to send ueni surplus of %s to EVM module account", surplusUeni)
			//}
		}
		if surplusWei.GT(cosmossdk_io_math.ZeroInt()) {
			// todo  check code correct
			//if err := am.keeper.BankKeeper().AddWei(ctx, am.keeper.AccountKeeper().GetModuleAddress(types.ModuleName), surplusWei); err != nil {
			//	ctx.Logger().Error("failed to send wei surplus of %s to EVM module account", surplusWei)
			//}
		}
	}
	am.keeper.SetBlockBloom(ctx, utils.Map(evmTxDeferredInfoList, func(i *types.DeferredInfo) ethtypes.Bloom { return ethtypes.BytesToBloom(i.TxBloom) }))
	return nil
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

	// todo  check code correct
	AccountKeeper *authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	StakingKeeper *stakingkeeper.Keeper
}

type ModuleOutputs struct {
	depinject.Out

	EvmKeeper keeper.Keeper
	Module    appmodule.AppModule
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
		in.AccountKeeper,
		in.BankKeeper,
		in.StakingKeeper,
	)
	m := NewAppModule(
		in.Cdc,
		k,
		in.AccountKeeper,
		in.BankKeeper,
	)

	return ModuleOutputs{EvmKeeper: k, Module: m}
}
