package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/h4ckb/devgitsecops/internal/downloader"
	"github.com/h4ckb/devgitsecops/internal/executor"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [tool]",
	Short: "Install or link embedded tools",
	Long: `Install or link the embedded tools to devgitsecops.

This command can automatically download tools from their official sources,
link to existing installations on your system, or accept manual paths.

Examples:
  devgitsecops install kubectl                           # Download kubectl automatically
  devgitsecops install --all --auto                      # Download all tools
  devgitsecops install kubectl --from /path/to/kubectl   # Link existing binary`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fromPath, _ := cmd.Flags().GetString("from")
		all, _ := cmd.Flags().GetBool("all")
		auto, _ := cmd.Flags().GetBool("auto")

		if all {
			return installAll(auto)
		}

		if len(args) == 0 {
			return fmt.Errorf("please specify a tool name or use --all flag")
		}

		toolName := args[0]
		validTools := map[string]bool{
			"kubectl":   true,
			"kustomize": true,
			"helm":      true,
			"terraform": true,
			"az":        true,
			"aws":       true,
		}

		if !validTools[toolName] {
			return fmt.Errorf("invalid tool name: %s", toolName)
		}

		// If --from is specified, use manual installation
		if fromPath != "" {
			return installFromPath(toolName, fromPath)
		}

		// Otherwise, auto-download (this is the new default behavior!)
		return autoInstall(toolName)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringP("from", "f", "", "Path to existing binary to link (optional)")
	installCmd.Flags().BoolP("all", "a", false, "Attempt to download/install all tools")
	installCmd.Flags().Bool("auto", false, "Search system paths first before downloading")
}

func installFromPath(toolName, srcPath string) error {
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("source binary does not exist: %s", srcPath)
	}

	if err := executor.CopyBinary(srcPath, toolName); err != nil {
		return fmt.Errorf("failed to install %s: %w", toolName, err)
	}

	fmt.Printf("✓ Successfully installed %s\n", toolName)
	return nil
}

func autoInstall(toolName string) error {
	// First try to find the tool in system PATH
	searchPaths := []string{}
	
	// Add common paths based on OS
	if runtime.GOOS == "windows" {
		searchPaths = append(searchPaths,
			filepath.Join(os.Getenv("ProgramFiles"), "bin"),
			filepath.Join(os.Getenv("ProgramFiles(x86)"), "bin"),
			"C:\\tools\\bin",
		)
	} else {
		searchPaths = append(searchPaths,
			"/usr/local/bin",
			"/usr/bin",
			"/opt/homebrew/bin",
			filepath.Join(os.Getenv("HOME"), ".local", "bin"),
		)
	}

	// Add PATH directories
	pathEnv := os.Getenv("PATH")
	if pathEnv != "" {
		pathDirs := filepath.SplitList(pathEnv)
		searchPaths = append(searchPaths, pathDirs...)
	}

	binaryName := toolName
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	// Search for the binary in system
	for _, dir := range searchPaths {
		fullPath := filepath.Join(dir, binaryName)
		if _, err := os.Stat(fullPath); err == nil {
			fmt.Printf("Found %s at: %s\n", toolName, fullPath)
			return installFromPath(toolName, fullPath)
		}
	}

	// If not found, try to download it automatically
	fmt.Printf("%s not found in system. Attempting automatic download...\n", toolName)
	return downloadTool(toolName)
}

func downloadTool(toolName string) error {
	binDir := executor.GetBinDir()
	
	progressCallback := func(msg string) {
		fmt.Println(msg)
	}
	
	if err := downloader.DownloadTool(toolName, binDir, progressCallback); err != nil {
		return fmt.Errorf("failed to download %s: %w\nYou can manually install it and use: devgitsecops install %s --from /path/to/%s", toolName, err, toolName, toolName)
	}
	
	return nil
}

func installAll(auto bool) error {
	tools := []string{"kubectl", "kustomize", "helm", "terraform"}
	
	fmt.Println("Installing tools...")
	fmt.Println("Note: az and aws cli must be installed manually (large Python-based tools)")
	fmt.Println()
	
	success := 0
	failed := 0
	
	for _, tool := range tools {
		fmt.Printf("Installing %s... ", tool)
		
		exec := executor.NewToolExecutor(tool)
		if exec.IsInstalled() {
			fmt.Println("✓ Already installed")
			success++
			continue
		}
		
		if auto {
			// Try system first, then download
			if err := autoInstall(tool); err != nil {
				fmt.Printf("✗ Failed: %v\n", err)
				failed++
			} else {
				success++
			}
		} else {
			// Just download
			if err := downloadTool(tool); err != nil {
				fmt.Printf("✗ Failed: %v\n", err)
				failed++
			} else {
				success++
			}
		}
	}
	
	fmt.Printf("\n✓ Installed: %d  ✗ Failed: %d\n", success, failed)
	return nil
}

func printInstallInstructions(toolName string) error {
	fmt.Printf("To install %s, simply run:\n\n", toolName)
	fmt.Printf("  devgitsecops install %s\n\n", toolName)
	fmt.Println("This will automatically download and install the tool.")
	fmt.Println()
	
	instructions := map[string]string{
		"kubectl": `kubectl will be downloaded from the official Kubernetes release repository.`,
		"kustomize": `kustomize will be downloaded from the official GitHub releases.`,
		"helm": `helm will be downloaded from the official Helm release repository.`,
		"terraform": `terraform will be downloaded from the official HashiCorp releases.`,
		"az": `az cli must be installed manually from:
  https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-windows
  
  After installation, link it with:
    devgitsecops install az --auto`,
		"aws": `aws cli must be installed manually from:
  https://aws.amazon.com/cli/
  
  After installation, link it with:
    devgitsecops install aws --auto`,
	}
	
	if inst, ok := instructions[toolName]; ok {
		fmt.Println(inst)
	}
	
	return nil
}
