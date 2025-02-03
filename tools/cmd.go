package tools

import (
	"github.com/spf13/cobra"

	migration "github.com/eni-chain/go-eni/tools/migration/cmd"
	scanner "github.com/eni-chain/go-eni/tools/tx-scanner/cmd"
)

func ToolCmd() *cobra.Command {
	toolsCmd := &cobra.Command{
		Use:   "tools",
		Short: "A set of useful tools for eni chain",
	}
	toolsCmd.AddCommand(scanner.ScanCmd())
	toolsCmd.AddCommand(migration.MigrateCmd())
	toolsCmd.AddCommand(migration.VerifyMigrationCmd())
	toolsCmd.AddCommand(migration.GenerateStats())
	return toolsCmd
}
