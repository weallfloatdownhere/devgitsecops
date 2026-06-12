package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/h4ckb/devgitsecops/internal/executor"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of embedded tools",
	Long:  `Display the installation status and versions of all embedded tools.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tools := []string{"kubectl", "kustomize", "helm", "terraform", "az", "aws"}
		
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "TOOL\tSTATUS\tVERSION")
		fmt.Fprintln(w, "----\t------\t-------")
		
		for _, tool := range tools {
			exec := executor.NewToolExecutor(tool)
			status := "Not Installed"
			version := "N/A"
			
			if exec.IsInstalled() {
				status = "✓ Installed"
				if v, err := exec.GetVersion(); err == nil {
					// Trim version output to first line
					lines := []rune(v)
					for i, r := range lines {
						if r == '\n' {
							version = string(lines[:i])
							break
						}
					}
					if version == "N/A" {
						version = v
					}
				}
			}
			
			fmt.Fprintf(w, "%s\t%s\t%s\n", tool, status, version)
		}
		
		w.Flush()
		fmt.Println("\nBinary location:", executor.GetBinDir())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
