package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ToolExecutor handles execution of embedded tools
type ToolExecutor struct {
	ToolName   string
	BinaryPath string
}

// NewToolExecutor creates a new tool executor
func NewToolExecutor(toolName string) *ToolExecutor {
	return &ToolExecutor{
		ToolName:   toolName,
		BinaryPath: getToolPath(toolName),
	}
}

// Execute runs the tool with the provided arguments
func (te *ToolExecutor) Execute(args []string) error {
	// Check if tool binary exists
	if !te.IsInstalled() {
		return fmt.Errorf("%s is not installed. Run 'devgitsecops install %s --auto' to download it automatically", te.ToolName, te.ToolName)
	}

	cmd := exec.Command(te.BinaryPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	return cmd.Run()
}

// IsInstalled checks if the tool binary exists
func (te *ToolExecutor) IsInstalled() bool {
	_, err := os.Stat(te.BinaryPath)
	return err == nil
}

// GetVersion returns the version of the installed tool
func (te *ToolExecutor) GetVersion() (string, error) {
	if !te.IsInstalled() {
		return "", fmt.Errorf("%s is not installed", te.ToolName)
	}

	cmd := exec.Command(te.BinaryPath, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Try --version flag
		cmd = exec.Command(te.BinaryPath, "--version")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return "", err
		}
	}

	return string(output), nil
}

// getToolPath returns the path where the tool binary should be stored
func getToolPath(toolName string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	binDir := filepath.Join(homeDir, ".devgitsecops", "bin")
	binaryName := toolName
	
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	return filepath.Join(binDir, binaryName)
}

// GetBinDir returns the directory where tool binaries are stored
func GetBinDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".devgitsecops", "bin")
}

// EnsureBinDir creates the bin directory if it doesn't exist
func EnsureBinDir() error {
	binDir := GetBinDir()
	return os.MkdirAll(binDir, 0755)
}

// CopyBinary copies a binary file to the bin directory
func CopyBinary(srcPath, toolName string) error {
	if err := EnsureBinDir(); err != nil {
		return err
	}

	destPath := getToolPath(toolName)
	
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source binary: %w", err)
	}
	defer src.Close()

	dest, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination binary: %w", err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return fmt.Errorf("failed to copy binary: %w", err)
	}

	return nil
}
