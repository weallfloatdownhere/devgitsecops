package downloader

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ToolInfo contains download information for a tool
type ToolInfo struct {
	Name        string
	Version     string
	DownloadURL string
	BinaryName  string
	IsArchive   bool
}

// GetToolInfo returns download information for supported tools
func GetToolInfo(toolName string) (*ToolInfo, error) {
	// Detect latest versions and construct download URLs
	switch toolName {
	case "kubectl":
		return &ToolInfo{
			Name:        "kubectl",
			Version:     "v1.36.0",
			DownloadURL: fmt.Sprintf("https://dl.k8s.io/release/v1.36.0/bin/%s/%s/kubectl%s", runtime.GOOS, runtime.GOARCH, getExeSuffix()),
			BinaryName:  "kubectl" + getExeSuffix(),
			IsArchive:   false,
		}, nil

	case "kustomize":
		// Kustomize releases are in archives
		version := "v5.4.3"
		os := runtime.GOOS
		arch := runtime.GOARCH
		if arch == "amd64" {
			arch = "amd64"
		}
		if os == "windows" {
			return &ToolInfo{
				Name:        "kustomize",
				Version:     version,
				DownloadURL: fmt.Sprintf("https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/%s/kustomize_%s_windows_amd64.zip", version, version),
				BinaryName:  "kustomize" + getExeSuffix(),
				IsArchive:   true,
			}, nil
		}
		return &ToolInfo{
			Name:        "kustomize",
			Version:     version,
			DownloadURL: fmt.Sprintf("https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/%s/kustomize_%s_%s_%s.tar.gz", version, version, os, arch),
			BinaryName:  "kustomize" + getExeSuffix(),
			IsArchive:   true,
		}, nil

	case "helm":
		version := "v3.16.2"
		os := runtime.GOOS
		arch := runtime.GOARCH
		ext := ".tar.gz"
		if os == "windows" {
			ext = ".zip"
		}
		return &ToolInfo{
			Name:        "helm",
			Version:     version,
			DownloadURL: fmt.Sprintf("https://get.helm.sh/helm-%s-%s-%s%s", version, os, arch, ext),
			BinaryName:  "helm" + getExeSuffix(),
			IsArchive:   true,
		}, nil

	case "terraform":
		version := "1.9.8"
		os := runtime.GOOS
		arch := runtime.GOARCH
		return &ToolInfo{
			Name:        "terraform",
			Version:     version,
			DownloadURL: fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_%s.zip", version, version, os, arch),
			BinaryName:  "terraform" + getExeSuffix(),
			IsArchive:   true,
		}, nil

	case "az":
		return nil, fmt.Errorf("az cli requires manual installation from https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-windows (large Python-based tool)")

	case "aws":
		return nil, fmt.Errorf("aws cli requires manual installation from https://aws.amazon.com/cli/ (large Python-based tool)")

	default:
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}
}

// DownloadTool downloads and installs a tool
func DownloadTool(toolName, destDir string, progressCallback func(string)) error {
	info, err := GetToolInfo(toolName)
	if err != nil {
		return err
	}

	if progressCallback != nil {
		progressCallback(fmt.Sprintf("Downloading %s %s...", info.Name, info.Version))
	}

	// Download the file
	resp, err := http.Get(info.DownloadURL)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: HTTP %d", resp.StatusCode)
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("%s-*", toolName))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Download to temp file
	written, err := io.Copy(tmpFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save download: %w", err)
	}
	tmpFile.Close()

	if progressCallback != nil {
		progressCallback(fmt.Sprintf("Downloaded %d bytes", written))
	}

	// Ensure destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Extract if archive, otherwise just move
	if info.IsArchive {
		if progressCallback != nil {
			progressCallback("Extracting archive...")
		}
		if err := extractBinary(tmpFile.Name(), info.BinaryName, destDir); err != nil {
			return fmt.Errorf("failed to extract: %w", err)
		}
	} else {
		// Direct binary, just move it
		destPath := filepath.Join(destDir, info.BinaryName)
		if err := copyFile(tmpFile.Name(), destPath); err != nil {
			return fmt.Errorf("failed to copy binary: %w", err)
		}
		if err := os.Chmod(destPath, 0755); err != nil {
			return fmt.Errorf("failed to set permissions: %w", err)
		}
	}

	if progressCallback != nil {
		progressCallback(fmt.Sprintf("✓ %s installed successfully", info.Name))
	}

	return nil
}

// extractBinary extracts a binary from an archive
func extractBinary(archivePath, binaryName, destDir string) error {
	// Handle zip files
	if strings.HasSuffix(archivePath, ".zip") || runtime.GOOS == "windows" {
		r, err := zip.OpenReader(archivePath)
		if err != nil {
			return fmt.Errorf("failed to open zip: %w", err)
		}
		defer r.Close()

		// Find the binary in the archive
		for _, f := range r.File {
			// Check if this is the binary we want (might be in a subdirectory)
			baseName := filepath.Base(f.Name)
			if baseName == binaryName || strings.HasSuffix(f.Name, binaryName) {
				rc, err := f.Open()
				if err != nil {
					return fmt.Errorf("failed to open file in zip: %w", err)
				}
				defer rc.Close()

				destPath := filepath.Join(destDir, binaryName)
				outFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
				if err != nil {
					return fmt.Errorf("failed to create output file: %w", err)
				}
				defer outFile.Close()

				_, err = io.Copy(outFile, rc)
				if err != nil {
					return fmt.Errorf("failed to copy file: %w", err)
				}
				return nil
			}
		}

		return fmt.Errorf("binary %s not found in archive", binaryName)
	}

	// For non-Windows or .tar.gz files, return error for now
	return fmt.Errorf("unsupported archive format for %s", archivePath)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// getExeSuffix returns the executable suffix for the current OS
func getExeSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}
