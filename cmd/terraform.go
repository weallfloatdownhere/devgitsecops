package cmd

import (
	"github.com/spf13/cobra"
)

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Terraform automation commands",
	Long: `Commands to automate common Terraform operations and infrastructure setup.

Available subcommands:
  setup-backend - Setup Azure backend infrastructure for Terraform state`,
}

func init() {
	rootCmd.AddCommand(terraformCmd)
}
