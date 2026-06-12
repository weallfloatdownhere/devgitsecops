package cmd

import (
	"github.com/h4ckb/devgitsecops/internal/executor"
	"github.com/spf13/cobra"
)

var azCmd = &cobra.Command{
	Use:   "az",
	Short: "Run Azure CLI commands",
	Long:  `Execute az (Azure CLI) commands through devgitsecops.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		exec := executor.NewToolExecutor("az")
		return exec.Execute(args)
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(azCmd)
}
