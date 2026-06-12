package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display version, commit, and build information for devgitsecops.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("devgitsecops version %s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
