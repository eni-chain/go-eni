package binding

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/eni-chain/go-eni/testutil/sample"
	bindingsimulation "github.com/eni-chain/go-eni/x/binding/simulation"
	"github.com/eni-chain/go-eni/x/binding/types"
)

// avoid unused import issue
var (
	_ = bindingsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgCreateBinding = "op_weight_msg_binding"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateBinding int = 100

	opWeightMsgUpdateBinding = "op_weight_msg_binding"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateBinding int = 100

	opWeightMsgDeleteBinding = "op_weight_msg_binding"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteBinding int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	bindingGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		BindingList: []types.Binding{
			{
				Creator: sample.AccAddress(),
				Index:   "0",
			},
			{
				Creator: sample.AccAddress(),
				Index:   "1",
			},
		},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&bindingGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateBinding int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateBinding, &weightMsgCreateBinding, nil,
		func(_ *rand.Rand) {
			weightMsgCreateBinding = defaultWeightMsgCreateBinding
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateBinding,
		bindingsimulation.SimulateMsgCreateBinding(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateBinding int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateBinding, &weightMsgUpdateBinding, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateBinding = defaultWeightMsgUpdateBinding
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateBinding,
		bindingsimulation.SimulateMsgUpdateBinding(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteBinding int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteBinding, &weightMsgDeleteBinding, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteBinding = defaultWeightMsgDeleteBinding
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteBinding,
		bindingsimulation.SimulateMsgDeleteBinding(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateBinding,
			defaultWeightMsgCreateBinding,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				bindingsimulation.SimulateMsgCreateBinding(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateBinding,
			defaultWeightMsgUpdateBinding,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				bindingsimulation.SimulateMsgUpdateBinding(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteBinding,
			defaultWeightMsgDeleteBinding,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				bindingsimulation.SimulateMsgDeleteBinding(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
