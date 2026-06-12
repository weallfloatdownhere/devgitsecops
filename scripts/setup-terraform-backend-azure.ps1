#
# Azure Terraform Backend Setup Script (PowerShell)
# This script creates all resources needed for storing Terraform state in Azure
#
# Prerequisites: 
# - Azure CLI installed and logged in (az login)
# - Appropriate permissions to create resources
#
# Usage:
#   .\setup-terraform-backend-azure.ps1 -Environment "production"
#   .\setup-terraform-backend-azure.ps1 -Environment "dev" -Location "eastus"

param(
    [Parameter(Mandatory=$false)]
    [string]$Environment = "dev",
    
    [Parameter(Mandatory=$false)]
    [string]$Location = "eastus",
    
    [Parameter(Mandatory=$false)]
    [string]$Prefix = "terraform"
)

$ErrorActionPreference = "Stop"

# Generate unique suffix for globally unique names
$uniqueSuffix = Get-Random -Minimum 10000 -Maximum 99999

# Resource names (customize as needed)
$ResourceGroupName = "${Prefix}-${Environment}-rg"
$StorageAccountName = "${Prefix}${Environment}sa${uniqueSuffix}".ToLower() -replace '[^a-z0-9]', ''
$StorageAccountName = $StorageAccountName.Substring(0, [Math]::Min(24, $StorageAccountName.Length))
$ContainerName = "tfstate"
$KeyVaultName = "${Prefix}-${Environment}-kv-${uniqueSuffix}"

# Ensure Key Vault name is valid (max 24 chars, alphanumeric and hyphens)
if ($KeyVaultName.Length -gt 24) {
    $KeyVaultName = $KeyVaultName.Substring(0, 24)
}

# Tags
$Tags = "Environment=${Environment} Purpose=TerraformBackend ManagedBy=Script"

Write-Host "================================================" -ForegroundColor Green
Write-Host "  Azure Terraform Backend Setup" -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Green
Write-Host ""
Write-Host "Environment:       $Environment"
Write-Host "Location:          $Location"
Write-Host "Resource Group:    $ResourceGroupName"
Write-Host "Storage Account:   $StorageAccountName"
Write-Host "Container:         $ContainerName"
Write-Host "Key Vault:         $KeyVaultName"
Write-Host ""

# Confirm before proceeding
$confirmation = Read-Host "Do you want to continue? (yes/no)"
if ($confirmation -ne "yes") {
    Write-Host "Aborted." -ForegroundColor Yellow
    exit 0
}

# Get current user information
Write-Host "[1/7] Getting current user information..." -ForegroundColor Green
$CurrentUser = az account show --query user.name -o tsv
$CurrentUserObjectId = az ad signed-in-user show --query id -o tsv 2>$null
Write-Host "Current user: $CurrentUser"

# Create Resource Group
Write-Host "[2/7] Creating resource group: $ResourceGroupName..." -ForegroundColor Green
az group create `
  --name $ResourceGroupName `
  --location $Location `
  --tags $Tags

# Create Storage Account
Write-Host "[3/7] Creating storage account: $StorageAccountName..." -ForegroundColor Green
az storage account create `
  --name $StorageAccountName `
  --resource-group $ResourceGroupName `
  --location $Location `
  --sku Standard_LRS `
  --encryption-services blob `
  --https-only true `
  --min-tls-version TLS1_2 `
  --allow-blob-public-access false `
  --tags $Tags

# Enable versioning (recommended for state file protection)
Write-Host "[3a/7] Enabling blob versioning..." -ForegroundColor Green
az storage account blob-service-properties update `
  --account-name $StorageAccountName `
  --resource-group $ResourceGroupName `
  --enable-versioning true

# Create Blob Container
Write-Host "[4/7] Creating blob container: $ContainerName..." -ForegroundColor Green
az storage container create `
  --name $ContainerName `
  --account-name $StorageAccountName `
  --auth-mode login

# Create Key Vault
Write-Host "[5/7] Creating Key Vault: $KeyVaultName..." -ForegroundColor Green
az keyvault create `
  --name $KeyVaultName `
  --resource-group $ResourceGroupName `
  --location $Location `
  --enabled-for-deployment true `
  --enabled-for-template-deployment true `
  --tags $Tags

# Set Key Vault access policy for current user
if ($CurrentUserObjectId) {
    Write-Host "[5a/7] Setting Key Vault access policy for current user..." -ForegroundColor Green
    az keyvault set-policy `
      --name $KeyVaultName `
      --object-id $CurrentUserObjectId `
      --secret-permissions get list set delete `
      --key-permissions get list create delete `
      --certificate-permissions get list create delete
}

# Get storage account key and store in Key Vault
Write-Host "[6/7] Storing storage account key in Key Vault..." -ForegroundColor Green
$StorageAccountKey = az storage account keys list `
  --resource-group $ResourceGroupName `
  --account-name $StorageAccountName `
  --query '[0].value' -o tsv

az keyvault secret set `
  --vault-name $KeyVaultName `
  --name "terraform-backend-storage-key" `
  --value $StorageAccountKey `
  --description "Storage account key for Terraform backend"

# Store configuration in Key Vault
Write-Host "[7/7] Storing configuration in Key Vault..." -ForegroundColor Green
az keyvault secret set `
  --vault-name $KeyVaultName `
  --name "terraform-backend-storage-account" `
  --value $StorageAccountName

az keyvault secret set `
  --vault-name $KeyVaultName `
  --name "terraform-backend-container" `
  --value $ContainerName

az keyvault secret set `
  --vault-name $KeyVaultName `
  --name "terraform-backend-resource-group" `
  --value $ResourceGroupName

# Get subscription ID
$SubscriptionId = az account show --query id -o tsv

Write-Host ""
Write-Host "================================================" -ForegroundColor Green
Write-Host "  Setup Complete!" -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Green
Write-Host ""
Write-Host "Terraform Backend Configuration:" -ForegroundColor Yellow
Write-Host ""
Write-Host "Add this to your Terraform configuration (backend.tf):"
Write-Host ""
Write-Host "---------------------------------------------------"
Write-Host @"
terraform {
  backend "azurerm" {
    resource_group_name  = "$ResourceGroupName"
    storage_account_name = "$StorageAccountName"
    container_name       = "$ContainerName"
    key                  = "terraform.tfstate"
  }
}
"@
Write-Host "---------------------------------------------------"
Write-Host ""
Write-Host "Environment Variables (optional):" -ForegroundColor Yellow
Write-Host ""
Write-Host "`$env:ARM_SUBSCRIPTION_ID = `"$SubscriptionId`""
Write-Host "`$env:ARM_ACCESS_KEY = `"(az keyvault secret show --vault-name $KeyVaultName --name terraform-backend-storage-key --query value -o tsv)`""
Write-Host ""
Write-Host "Or use Azure CLI authentication:" -ForegroundColor Yellow
Write-Host ""
Write-Host "az login"
Write-Host "terraform init"
Write-Host ""
Write-Host "Key Vault Information:" -ForegroundColor Yellow
Write-Host "Key Vault Name: $KeyVaultName"
Write-Host "Secrets stored:"
Write-Host "  - terraform-backend-storage-key"
Write-Host "  - terraform-backend-storage-account"
Write-Host "  - terraform-backend-container"
Write-Host "  - terraform-backend-resource-group"
Write-Host ""
Write-Host "All resources have been created successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Initialize Terraform: terraform init"
Write-Host "2. The state file will be stored in: $StorageAccountName/$ContainerName/terraform.tfstate"
Write-Host "3. Access storage key from Key Vault if needed:"
Write-Host "   az keyvault secret show --vault-name $KeyVaultName --name terraform-backend-storage-key --query value -o tsv"
Write-Host ""
