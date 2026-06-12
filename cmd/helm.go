package cmd

import (
	"github.com/h4ckb/devgitsecops/internal/executor"
	"github.com/spf13/cobra"
)

var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Run helm commands",
	Long:  `Execute helm (Kubernetes package manager) commands through devgitsecops.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		exec := executor.NewToolExecutor("helm")
		return exec.Execute(args)
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(helmCmd)
}
