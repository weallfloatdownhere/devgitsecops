package cmd

import (
	"github.com/h4ckb/devgitsecops/internal/executor"
	"github.com/spf13/cobra"
)

var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Run AWS CLI commands",
	Long:  `Execute aws (AWS CLI) commands through devgitsecops.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		exec := executor.NewToolExecutor("aws")
		return exec.Execute(args)
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(awsCmd)
}
