package cmd

import (
	"github.com/h4ckb/devgitsecops/internal/executor"
	"github.com/spf13/cobra"
)

var kubectlCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "Run kubectl commands",
	Long:  `Execute kubectl (Kubernetes CLI) commands through devgitsecops.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		exec := executor.NewToolExecutor("kubectl")
		return exec.Execute(args)
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(kubectlCmd)
}
