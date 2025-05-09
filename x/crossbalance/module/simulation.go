package crossbalance

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/eni-chain/go-eni/testutil/sample"
	crossbalancesimulation "github.com/eni-chain/go-eni/x/crossbalance/simulation"
	"github.com/eni-chain/go-eni/x/crossbalance/types"
)

// avoid unused import issue
var (
	_ = crossbalancesimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgTransferCrossAccount = "op_weight_msg_transfer_cross_account"
	// TODO: Determine the simulation weight value
	defaultWeightMsgTransferCrossAccount int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	crossbalanceGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&crossbalanceGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgTransferCrossAccount int
	simState.AppParams.GetOrGenerate(opWeightMsgTransferCrossAccount, &weightMsgTransferCrossAccount, nil,
		func(_ *rand.Rand) {
			weightMsgTransferCrossAccount = defaultWeightMsgTransferCrossAccount
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTransferCrossAccount,
		crossbalancesimulation.SimulateMsgTransferCrossAccount(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgTransferCrossAccount,
			defaultWeightMsgTransferCrossAccount,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				crossbalancesimulation.SimulateMsgTransferCrossAccount(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
