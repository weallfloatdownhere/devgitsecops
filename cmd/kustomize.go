package cmd

import (
	"github.com/h4ckb/devgitsecops/internal/executor"
	"github.com/spf13/cobra"
)

var kustomizeCmd = &cobra.Command{
	Use:   "kustomize",
	Short: "Run kustomize commands",
	Long:  `Execute kustomize (Kubernetes configuration management) commands through devgitsecops.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		exec := executor.NewToolExecutor("kustomize")
		return exec.Execute(args)
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(kustomizeCmd)
}
