package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devgitsecops",
	Short: "A unified CLI tool for DevOps operations",
	Long: `devgitsecops is a unified command-line interface that embeds multiple 
DevOps tools into a single binary, making it easier to manage your 
infrastructure and deployments.

Embedded tools:
  - kubectl    (Kubernetes CLI)
  - kustomize  (Kubernetes configuration management)
  - helm       (Kubernetes package manager)
  - terraform  (Infrastructure as Code)
  - az         (Azure CLI)
  - aws        (AWS CLI)`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.devgitsecops.yaml)")
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".devgitsecops")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
