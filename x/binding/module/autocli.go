package binding

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/eni-chain/go-eni/api/goeni/binding"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "BindingAll",
					Use:       "list-binding",
					Short:     "List all binding",
				},
				{
					RpcMethod:      "Binding",
					Use:            "show-binding [id]",
					Short:          "Shows a binding",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateBinding",
					Use:            "create-binding [index] [evmAddress] [cosmosAddress]",
					Short:          "Create a new binding",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "evmAddress"}, {ProtoField: "cosmosAddress"}},
				},
				{
					RpcMethod:      "UpdateBinding",
					Use:            "update-binding [index] [evmAddress] [cosmosAddress]",
					Short:          "Update binding",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "evmAddress"}, {ProtoField: "cosmosAddress"}},
				},
				{
					RpcMethod:      "DeleteBinding",
					Use:            "delete-binding [index]",
					Short:          "Delete binding",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
