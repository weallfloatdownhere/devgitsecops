package cmd

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	setupEnvironment string
	setupLocation    string
	setupPrefix      string
	setupAutoApprove bool
)

var setupTerraformBackendCmd = &cobra.Command{
	Use:   "setup-backend-azure",
	Short: "Setup Azure Terraform backend infrastructure",
	Long: `Creates all resources needed for storing Terraform state in Azure:
  - Resource Group
  - Storage Account (with encryption, versioning, HTTPS-only)
  - Blob Container (tfstate)
  - Key Vault (for secure credential storage)

The command automatically checks your Azure login status and prompts for
authentication if needed. All resources follow Azure best practices for
security and reliability.`,
	Example: `  devgitsecops terraform setup-backend-azure
  devgitsecops terraform setup-backend-azure --environment production
  devgitsecops terraform setup-backend-azure --environment staging --location westus2
  devgitsecops terraform setup-backend-azure --environment prod --auto-approve`,
	RunE: runSetupTerraformBackend,
}

func init() {
	terraformCmd.AddCommand(setupTerraformBackendCmd)
	
	setupTerraformBackendCmd.Flags().StringVarP(&setupEnvironment, "environment", "e", "dev", "Environment name (dev, staging, production)")
	setupTerraformBackendCmd.Flags().StringVarP(&setupLocation, "location", "l", "eastus", "Azure region")
	setupTerraformBackendCmd.Flags().StringVarP(&setupPrefix, "prefix", "p", "terraform", "Resource name prefix")
	setupTerraformBackendCmd.Flags().BoolVarP(&setupAutoApprove, "auto-approve", "y", false, "Skip confirmation prompt")
}

func runSetupTerraformBackend(cmd *cobra.Command, args []string) error {
	// Check if az cli is available
	if !isAzCliAvailable() {
		return fmt.Errorf("Azure CLI is not installed or not found. Install it with: devgitsecops install az --auto")
	}

	// Check if user is logged in to Azure
	if err := ensureAzureLogin(); err != nil {
		return fmt.Errorf("Azure login failed: %w", err)
	}

	// Generate unique suffix
	rand.Seed(time.Now().UnixNano())
	uniqueSuffix := rand.Intn(90000) + 10000

	// Calculate resource names
	resourceGroupName := fmt.Sprintf("%s-%s-rg", setupPrefix, setupEnvironment)
	storageAccountName := fmt.Sprintf("%s%ssa%d", setupPrefix, setupEnvironment, uniqueSuffix)
	storageAccountName = strings.ToLower(strings.ReplaceAll(storageAccountName, "-", ""))
	if len(storageAccountName) > 24 {
		storageAccountName = storageAccountName[:24]
	}
	containerName := "tfstate"
	keyVaultName := fmt.Sprintf("%s-%s-kv-%d", setupPrefix, setupEnvironment, uniqueSuffix)
	if len(keyVaultName) > 24 {
		keyVaultName = keyVaultName[:24]
	}

	tags := fmt.Sprintf("Environment=%s Purpose=TerraformBackend ManagedBy=devgitsecops", setupEnvironment)

	// Display configuration
	printHeader("Azure Terraform Backend Setup")
	fmt.Println()
	fmt.Printf("Environment:       %s\n", setupEnvironment)
	fmt.Printf("Location:          %s\n", setupLocation)
	fmt.Printf("Resource Group:    %s\n", resourceGroupName)
	fmt.Printf("Storage Account:   %s\n", storageAccountName)
	fmt.Printf("Container:         %s\n", containerName)
	fmt.Printf("Key Vault:         %s\n", keyVaultName)
	fmt.Println()

	// Confirm before proceeding
	if !setupAutoApprove {
		if !confirmAction("Do you want to continue? (yes/no): ") {
			fmt.Println("Aborted.")
			return nil
		}
	}

	// Get current user information
	printStep("[1/7] Getting current user information...")
	currentUser, err := runAzCommand("account", "show", "--query", "user.name", "-o", "tsv")
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	fmt.Printf("Current user: %s\n", currentUser)

	currentUserObjectID, _ := runAzCommand("ad", "signed-in-user", "show", "--query", "id", "-o", "tsv")

	// Create Resource Group
	printStep("[2/7] Creating resource group: " + resourceGroupName + "...")
	if err := runAzCommandWithOutput("group", "create",
		"--name", resourceGroupName,
		"--location", setupLocation,
		"--tags", tags); err != nil {
		return fmt.Errorf("failed to create resource group: %w", err)
	}

	// Create Storage Account
	printStep("[3/7] Creating storage account: " + storageAccountName + "...")
	if err := runAzCommandWithOutput("storage", "account", "create",
		"--name", storageAccountName,
		"--resource-group", resourceGroupName,
		"--location", setupLocation,
		"--sku", "Standard_LRS",
		"--encryption-services", "blob",
		"--https-only", "true",
		"--min-tls-version", "TLS1_2",
		"--allow-blob-public-access", "false",
		"--tags", tags); err != nil {
		return fmt.Errorf("failed to create storage account: %w", err)
	}

	// Enable versioning
	printStep("[3a/7] Enabling blob versioning...")
	if err := runAzCommandWithOutput("storage", "account", "blob-service-properties", "update",
		"--account-name", storageAccountName,
		"--resource-group", resourceGroupName,
		"--enable-versioning", "true"); err != nil {
		return fmt.Errorf("failed to enable versioning: %w", err)
	}

	// Create Blob Container
	printStep("[4/7] Creating blob container: " + containerName + "...")
	if err := runAzCommandWithOutput("storage", "container", "create",
		"--name", containerName,
		"--account-name", storageAccountName,
		"--auth-mode", "login"); err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Create Key Vault
	printStep("[5/7] Creating Key Vault: " + keyVaultName + "...")
	if err := runAzCommandWithOutput("keyvault", "create",
		"--name", keyVaultName,
		"--resource-group", resourceGroupName,
		"--location", setupLocation,
		"--enabled-for-deployment", "true",
		"--enabled-for-template-deployment", "true",
		"--tags", tags); err != nil {
		return fmt.Errorf("failed to create key vault: %w", err)
	}

	// Set Key Vault access policy
	if currentUserObjectID != "" {
		printStep("[5a/7] Setting Key Vault access policy for current user...")
		if err := runAzCommandWithOutput("keyvault", "set-policy",
			"--name", keyVaultName,
			"--object-id", currentUserObjectID,
			"--secret-permissions", "get", "list", "set", "delete",
			"--key-permissions", "get", "list", "create", "delete",
			"--certificate-permissions", "get", "list", "create", "delete"); err != nil {
			fmt.Printf("Warning: Failed to set access policy: %v\n", err)
		}
	}

	// Get storage account key
	printStep("[6/7] Storing storage account key in Key Vault...")
	storageKey, err := runAzCommand("storage", "account", "keys", "list",
		"--resource-group", resourceGroupName,
		"--account-name", storageAccountName,
		"--query", "[0].value",
		"-o", "tsv")
	if err != nil {
		return fmt.Errorf("failed to get storage key: %w", err)
	}

	// Store secrets in Key Vault
	if err := runAzCommandWithOutput("keyvault", "secret", "set",
		"--vault-name", keyVaultName,
		"--name", "terraform-backend-storage-key",
		"--value", storageKey,
		"--description", "Storage account key for Terraform backend"); err != nil {
		return fmt.Errorf("failed to store storage key: %w", err)
	}

	printStep("[7/7] Storing configuration in Key Vault...")
	secrets := map[string]string{
		"terraform-backend-storage-account":   storageAccountName,
		"terraform-backend-container":         containerName,
		"terraform-backend-resource-group":    resourceGroupName,
	}

	for name, value := range secrets {
		if err := runAzCommandWithOutput("keyvault", "secret", "set",
			"--vault-name", keyVaultName,
			"--name", name,
			"--value", value); err != nil {
			fmt.Printf("Warning: Failed to store secret %s: %v\n", name, err)
		}
	}

	// Get subscription ID
	subscriptionID, _ := runAzCommand("account", "show", "--query", "id", "-o", "tsv")

	// Print success message and configuration
	fmt.Println()
	printHeader("Setup Complete!")
	fmt.Println()
	printSuccess("Terraform Backend Configuration:")
	fmt.Println()
	fmt.Println("Add this to your Terraform configuration (backend.tf):")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Printf(`terraform {
  backend "azurerm" {
    resource_group_name  = "%s"
    storage_account_name = "%s"
    container_name       = "%s"
    key                  = "terraform.tfstate"
  }
}
`, resourceGroupName, storageAccountName, containerName)
	fmt.Println("---------------------------------------------------")
	fmt.Println()
	
	printSuccess("Environment Variables (optional):")
	fmt.Println()
	fmt.Printf("export ARM_SUBSCRIPTION_ID=\"%s\"\n", subscriptionID)
	fmt.Printf("export ARM_ACCESS_KEY=\"$(az keyvault secret show --vault-name %s --name terraform-backend-storage-key --query value -o tsv)\"\n", keyVaultName)
	fmt.Println()
	
	printSuccess("Or use Azure CLI authentication:")
	fmt.Println()
	fmt.Println("az login")
	fmt.Println("terraform init")
	fmt.Println()
	
	printSuccess("Key Vault Information:")
	fmt.Printf("Key Vault Name: %s\n", keyVaultName)
	fmt.Println("Secrets stored:")
	fmt.Println("  - terraform-backend-storage-key")
	fmt.Println("  - terraform-backend-storage-account")
	fmt.Println("  - terraform-backend-container")
	fmt.Println("  - terraform-backend-resource-group")
	fmt.Println()
	
	printSuccess("All resources have been created successfully!")
	fmt.Println()
	
	printSuccess("Next steps:")
	fmt.Println("1. Initialize Terraform: terraform init")
	fmt.Printf("2. The state file will be stored in: %s/%s/terraform.tfstate\n", storageAccountName, containerName)
	fmt.Printf("3. Access storage key from Key Vault if needed:\n")
	fmt.Printf("   az keyvault secret show --vault-name %s --name terraform-backend-storage-key --query value -o tsv\n", keyVaultName)
	fmt.Println()

	return nil
}

// Helper functions

func ensureAzureLogin() error {
	// Check if user is already logged in
	cmd := exec.Command("az", "account", "show")
	err := cmd.Run()
	
	if err != nil {
		// User is not logged in, prompt to login
		fmt.Println()
		printStep("Azure login required. Opening browser for authentication...")
		fmt.Println()
		
		loginCmd := exec.Command("az", "login")
		loginCmd.Stdout = os.Stdout
		loginCmd.Stderr = os.Stderr
		loginCmd.Stdin = os.Stdin
		
		if err := loginCmd.Run(); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}
		
		fmt.Println()
		printSuccess("✓ Successfully logged in to Azure")
		fmt.Println()
	} else {
		printSuccess("✓ Already logged in to Azure")
	}
	
	return nil
}

func isAzCliAvailable() bool {
	_, err := exec.LookPath("az")
	return err == nil
}

func runAzCommand(args ...string) (string, error) {
	cmd := exec.Command("az", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func runAzCommandWithOutput(args ...string) error {
	cmd := exec.Command("az", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func confirmAction(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "yes" || response == "y"
}

func printHeader(text string) {
	fmt.Println("================================================")
	fmt.Printf("  %s\n", text)
	fmt.Println("================================================")
}

func printStep(text string) {
	fmt.Printf("\033[32m%s\033[0m\n", text)
}

func printSuccess(text string) {
	fmt.Printf("\033[33m%s\033[0m\n", text)
}
