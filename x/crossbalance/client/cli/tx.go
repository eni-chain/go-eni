package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// NewTxCmd returns the root tx command for the crossbalance module.
func NewTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "crossbalance",
		Short:                      "Crossbalance transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		//CmdTransferCrossAccount(),
		CmdTransferEniToEvm(),
		CmdTransferEvmToEni(),
	)

	cmd.Flags().String(flags.FlagFrom, "", "Name of the sender (optional)")
	cmd.MarkFlagRequired(flags.FlagFrom)
	return cmd
}
