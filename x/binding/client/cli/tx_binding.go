package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/eni-chain/go-eni/x/binding/types"
	"github.com/spf13/cobra"
)

func CmdCreateBinding() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-binding [index] [evm-address] [cosmos-address]",
		Short: "Create a new binding between EVM and Cosmos addresses",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			index := args[0]
			evmAddress := args[1]
			cosmosAddress := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateBinding(
				clientCtx.GetFromAddress().String(),
				index,
				evmAddress,
				cosmosAddress,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func CmdUpdateBinding() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-binding [index] [evm-address] [cosmos-address]",
		Short: "Update an existing binding between EVM and Cosmos addresses",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			index := args[0]
			evmAddress := args[1]
			cosmosAddress := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateBinding(
				clientCtx.GetFromAddress().String(),
				index,
				evmAddress,
				cosmosAddress,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func CmdDeleteBinding() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-binding [index]",
		Short: "Delete an existing binding",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			index := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteBinding(
				clientCtx.GetFromAddress().String(),
				index,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
