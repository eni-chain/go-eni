package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	crossbalanceTypes "github.com/eni-chain/go-eni/x/crossbalance/types"
	"github.com/spf13/cobra"
)

// NewTxCmd returns the root tx command for the binding module.
func NewTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "binding",
		Short:                      "Binding transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdTransferCrossAccount(),
		CmdCreateBinding(),
		CmdUpdateBinding(),
		CmdDeleteBinding(),
	)

	cmd.Flags().String(flags.FlagFrom, "", "Name of the sender (optional)")
	cmd.MarkFlagRequired(flags.FlagFrom)
	return cmd
}

func CmdTransferCrossAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-cross-account [fromAddress] [toAddress] [amount]",
		Short: "Transfer tokens between cosmos and evm address formats",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			from := args[0]
			to := args[1]
			amount := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := crossbalanceTypes.NewMsgTransferCrossAccount(
				clientCtx.GetFromAddress().String(),
				from,
				to,
				amount,
			)

			// generate or broadcast transaction
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
