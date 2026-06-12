package cmd

import (
	"github.com/h4ckb/devgitsecops/internal/executor"
	"github.com/spf13/cobra"
)

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Run terraform commands",
	Long:  `Execute terraform (Infrastructure as Code) commands through devgitsecops.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		exec := executor.NewToolExecutor("terraform")
		return exec.Execute(args)
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(terraformCmd)
}
