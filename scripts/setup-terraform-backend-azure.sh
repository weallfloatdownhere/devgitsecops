#!/bin/bash
#
# Azure Terraform Backend Setup Script
# This script creates all resources needed for storing Terraform state in Azure
#
# Prerequisites: 
# - Azure CLI installed and logged in (az login)
# - Appropriate permissions to create resources
#
# Usage:
#   ./setup-terraform-backend-azure.sh [environment]
#   Example: ./setup-terraform-backend-azure.sh production

set -e

# Configuration Variables
ENVIRONMENT="${1:-dev}"
LOCATION="${AZURE_LOCATION:-eastus}"
PREFIX="terraform"

# Resource names (customize as needed)
RESOURCE_GROUP_NAME="${PREFIX}-${ENVIRONMENT}-rg"
STORAGE_ACCOUNT_NAME="${PREFIX}${ENVIRONMENT}sa$(date +%s | tail -c 5)"  # Must be unique globally
CONTAINER_NAME="tfstate"
KEY_VAULT_NAME="${PREFIX}-${ENVIRONMENT}-kv-$(date +%s | tail -c 5)"     # Must be unique globally

# Tags
TAGS="Environment=${ENVIRONMENT} Purpose=TerraformBackend ManagedBy=Script"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}  Azure Terraform Backend Setup${NC}"
echo -e "${GREEN}================================================${NC}"
echo ""
echo "Environment:       ${ENVIRONMENT}"
echo "Location:          ${LOCATION}"
echo "Resource Group:    ${RESOURCE_GROUP_NAME}"
echo "Storage Account:   ${STORAGE_ACCOUNT_NAME}"
echo "Container:         ${CONTAINER_NAME}"
echo "Key Vault:         ${KEY_VAULT_NAME}"
echo ""

# Confirm before proceeding
read -p "Do you want to continue? (yes/no): " -r
echo
if [[ ! $REPLY =~ ^[Yy]es$ ]]; then
    echo -e "${YELLOW}Aborted.${NC}"
    exit 0
fi

# Get current user information
echo -e "${GREEN}[1/7] Getting current user information...${NC}"
CURRENT_USER=$(az account show --query user.name -o tsv)
CURRENT_USER_OBJECT_ID=$(az ad signed-in-user show --query id -o tsv 2>/dev/null || echo "")
echo "Current user: ${CURRENT_USER}"

# Create Resource Group
echo -e "${GREEN}[2/7] Creating resource group: ${RESOURCE_GROUP_NAME}...${NC}"
az group create \
  --name "${RESOURCE_GROUP_NAME}" \
  --location "${LOCATION}" \
  --tags ${TAGS}

# Create Storage Account
echo -e "${GREEN}[3/7] Creating storage account: ${STORAGE_ACCOUNT_NAME}...${NC}"
az storage account create \
  --name "${STORAGE_ACCOUNT_NAME}" \
  --resource-group "${RESOURCE_GROUP_NAME}" \
  --location "${LOCATION}" \
  --sku Standard_LRS \
  --encryption-services blob \
  --https-only true \
  --min-tls-version TLS1_2 \
  --allow-blob-public-access false \
  --tags ${TAGS}

# Enable versioning (recommended for state file protection)
echo -e "${GREEN}[3a/7] Enabling blob versioning...${NC}"
az storage account blob-service-properties update \
  --account-name "${STORAGE_ACCOUNT_NAME}" \
  --resource-group "${RESOURCE_GROUP_NAME}" \
  --enable-versioning true

# Create Blob Container
echo -e "${GREEN}[4/7] Creating blob container: ${CONTAINER_NAME}...${NC}"
az storage container create \
  --name "${CONTAINER_NAME}" \
  --account-name "${STORAGE_ACCOUNT_NAME}" \
  --auth-mode login

# Create Key Vault
echo -e "${GREEN}[5/7] Creating Key Vault: ${KEY_VAULT_NAME}...${NC}"
az keyvault create \
  --name "${KEY_VAULT_NAME}" \
  --resource-group "${RESOURCE_GROUP_NAME}" \
  --location "${LOCATION}" \
  --enabled-for-deployment true \
  --enabled-for-template-deployment true \
  --tags ${TAGS}

# Set Key Vault access policy for current user
if [ -n "${CURRENT_USER_OBJECT_ID}" ]; then
  echo -e "${GREEN}[5a/7] Setting Key Vault access policy for current user...${NC}"
  az keyvault set-policy \
    --name "${KEY_VAULT_NAME}" \
    --object-id "${CURRENT_USER_OBJECT_ID}" \
    --secret-permissions get list set delete \
    --key-permissions get list create delete \
    --certificate-permissions get list create delete
fi

# Get storage account key and store in Key Vault
echo -e "${GREEN}[6/7] Storing storage account key in Key Vault...${NC}"
STORAGE_ACCOUNT_KEY=$(az storage account keys list \
  --resource-group "${RESOURCE_GROUP_NAME}" \
  --account-name "${STORAGE_ACCOUNT_NAME}" \
  --query '[0].value' -o tsv)

az keyvault secret set \
  --vault-name "${KEY_VAULT_NAME}" \
  --name "terraform-backend-storage-key" \
  --value "${STORAGE_ACCOUNT_KEY}" \
  --description "Storage account key for Terraform backend"

# Store configuration in Key Vault
echo -e "${GREEN}[7/7] Storing configuration in Key Vault...${NC}"
az keyvault secret set \
  --vault-name "${KEY_VAULT_NAME}" \
  --name "terraform-backend-storage-account" \
  --value "${STORAGE_ACCOUNT_NAME}"

az keyvault secret set \
  --vault-name "${KEY_VAULT_NAME}" \
  --name "terraform-backend-container" \
  --value "${CONTAINER_NAME}"

az keyvault secret set \
  --vault-name "${KEY_VAULT_NAME}" \
  --name "terraform-backend-resource-group" \
  --value "${RESOURCE_GROUP_NAME}"

# Get subscription ID
SUBSCRIPTION_ID=$(az account show --query id -o tsv)

echo ""
echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}  Setup Complete!${NC}"
echo -e "${GREEN}================================================${NC}"
echo ""
echo -e "${YELLOW}Terraform Backend Configuration:${NC}"
echo ""
echo "Add this to your Terraform configuration (backend.tf):"
echo ""
echo "---------------------------------------------------"
cat << EOF
terraform {
  backend "azurerm" {
    resource_group_name  = "${RESOURCE_GROUP_NAME}"
    storage_account_name = "${STORAGE_ACCOUNT_NAME}"
    container_name       = "${CONTAINER_NAME}"
    key                  = "terraform.tfstate"
  }
}
EOF
echo "---------------------------------------------------"
echo ""
echo -e "${YELLOW}Environment Variables (optional):${NC}"
echo ""
echo "export ARM_SUBSCRIPTION_ID=\"${SUBSCRIPTION_ID}\""
echo "export ARM_ACCESS_KEY=\"\$(az keyvault secret show --vault-name ${KEY_VAULT_NAME} --name terraform-backend-storage-key --query value -o tsv)\""
echo ""
echo -e "${YELLOW}Or use Azure CLI authentication:${NC}"
echo ""
echo "az login"
echo "terraform init"
echo ""
echo -e "${YELLOW}Key Vault Information:${NC}"
echo "Key Vault Name: ${KEY_VAULT_NAME}"
echo "Secrets stored:"
echo "  - terraform-backend-storage-key"
echo "  - terraform-backend-storage-account"
echo "  - terraform-backend-container"
echo "  - terraform-backend-resource-group"
echo ""
echo -e "${GREEN}All resources have been created successfully!${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Initialize Terraform: terraform init"
echo "2. The state file will be stored in: ${STORAGE_ACCOUNT_NAME}/${CONTAINER_NAME}/terraform.tfstate"
echo "3. Access storage key from Key Vault if needed:"
echo "   az keyvault secret show --vault-name ${KEY_VAULT_NAME} --name terraform-backend-storage-key --query value -o tsv"
echo ""
