# Azure Terraform Backend Setup Command

The `setup-terraform-backend` command is now built directly into devgitsecops! No need for separate scripts.

## Quick Start

```bash
# Make sure Azure CLI is installed
devgitsecops install az --auto

# Setup Terraform backend (one command - handles login automatically!)
devgitsecops setup-terraform-backend --environment production
```

**Note:** The command automatically checks if you're logged in to Azure and prompts for authentication if needed. No need to run `az login` separately!

## How It Works

The command follows this flow:

1. **Check Azure CLI** - Verifies `az` is installed
2. **Check Login Status** - Automatically detects if you're logged in
3. **Auto-Login** - If not logged in, opens browser for authentication
4. **Confirm Setup** - Shows what will be created (unless `--auto-approve`)
5. **Create Resources** - Sets up all backend infrastructure
6. **Store Credentials** - Saves keys securely in Key Vault
7. **Output Configuration** - Provides ready-to-use Terraform backend config

## Command Usage

```bash
devgitsecops setup-terraform-backend [flags]
```

### Flags

- `-e, --environment string` - Environment name (dev, staging, production) [default: "dev"]
- `-l, --location string` - Azure region [default: "eastus"]
- `-p, --prefix string` - Resource name prefix [default: "terraform"]
- `-y, --auto-approve` - Skip confirmation prompt

## Examples

### Basic Setup (Development)

```bash
devgitsecops setup-terraform-backend
```

Creates:
- Resource Group: `terraform-dev-rg`
- Storage Account: `terraformdevsa12345`
- Container: `tfstate`
- Key Vault: `terraform-dev-kv-12345`

### Production Environment

```bash
devgitsecops setup-terraform-backend --environment production --location westus2
```

### With Custom Prefix

```bash
devgitsecops setup-terraform-backend \
  --environment staging \
  --prefix mycompany \
  --location eastus
```

Creates:
- Resource Group: `mycompany-staging-rg`
- Storage Account: `mycompanystagingsa12345`
- Key Vault: `mycompany-staging-kv-12345`

### CI/CD (Non-Interactive)

```bash
devgitsecops setup-terraform-backend \
  --environment prod \
  --auto-approve
```

## What Gets Created

The command creates a complete Terraform backend infrastructure:

1. **Resource Group** - Container for all backend resources
2. **Storage Account**
   - Encryption enabled
   - HTTPS only
   - TLS 1.2 minimum
   - Blob versioning enabled (for state recovery)
   - Public access disabled
3. **Blob Container** - Named "tfstate"
4. **Key Vault** - Stores credentials securely
   - Storage account access key
   - Backend configuration details
   - Access policies configured for current user

## Output

After successful setup, the command outputs:

1. **Terraform backend configuration** (ready to copy into your `backend.tf`)
2. **Environment variables** for authentication
3. **Key Vault details** for accessing stored secrets
4. **Next steps** for using the backend

Example output:

```hcl
terraform {
  backend "azurerm" {
    resource_group_name  = "terraform-production-rg"
    storage_account_name = "terraformproductionsa12345"
    container_name       = "tfstate"
    key                  = "terraform.tfstate"
  }
}
```

## Using the Backend

### Step 1: Copy Configuration

Copy the output configuration into your `backend.tf` file.

### Step 2: Initialize Terraform

```bash
# Authenticate with Azure
devgitsecops az login

# Initialize Terraform
devgitsecops terraform init
```

### Step 3: Use Terraform

```bash
devgitsecops terraform plan
devgitsecops terraform apply
```

The state file is automatically stored in Azure Storage!

## Access Stored Credentials

Retrieve the storage account key from Key Vault:

```bash
devgitsecops az keyvault secret show \
  --vault-name terraform-production-kv-12345 \
  --name terraform-backend-storage-key \
  --query value -o tsv
```

## Complete Workflow Example

```bash
# 1. Install Azure CLI (if not already installed)
devgitsecops install az --auto

# 2. Create Terraform backend infrastructure (auto-login included!)
devgitsecops setup-terraform-backend --environment production

# 3. Create your Terraform project
mkdir my-terraform-project
cd my-terraform-project

# 4. Create backend.tf with the output from step 2
cat > backend.tf << 'EOF'
terraform {
  backend "azurerm" {
    resource_group_name  = "terraform-production-rg"
    storage_account_name = "terraformproductionsa12345"
    container_name       = "tfstate"
    key                  = "terraform.tfstate"
  }
}
EOF

# 5. Create your infrastructure code
cat > main.tf << 'EOF'
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "eastus"
}
EOF

# 6. Initialize and apply
devgitsecops terraform init  
✅ **Auto-login** - Handles Azure authentication automatically
devgitsecops terraform plan
devgitsecops terraform apply
```

## Security Features

✅ **Encryption at rest** - Storage account encryption enabled  
✅ **Encryption in transit** - HTTPS only, TLS 1.2 minimum  
✅ **No public access** - Blob public access disabled  
✅ **State versioning** - Blob versioning enabled for recovery  
✅ **Secure credentials** - Storage keys stored in Key Vault  
✅ **Access control** - RBAC with Key Vault access policies

## Troubleshooting

### Azure CLI not found

```bash
devgitsecops install az --auto
```

### Permission errors

Ensure you have:
- Contributor or Owner role on the subscription
- Permissions to create resource groups
- Permissions to create Key Vault and set access policies

### Storage account name conflicts

Storage account names must be globally unique. The command automatically adds a random suffix. If you still get conflicts, re-run the command to generate a new suffix.

### Cannot access Key Vault secrets

Set access policy manually:

```bash
devgitsecops az keyvault set-policy \
  --name terraform-production-kv-12345 \
  --object-id $(devgitsecops az ad signed-in-user show --query id -o tsv) \
  --secret-permissions get list
```

## Comparison: Script vs Built-in Command

### Before (PowerShell Script)

```powershell
# Login to Azure
az login

# Download script
# Make executable
# Run script with parameters
.\setup-terraform-backend-azure.ps1 -Environment "production"
```

### Now (Built-in Command)

```bash
# Just run it - auto-login included!
devgitsecops setup-terraform-backend --environment production
```

**Benefits:**
- ✅ No separate scripts to manage
- ✅ Automatic Azure login handling
- ✅ Same interface as all other commands
- ✅ Cross-platform (works on Windows, Linux, Mac)
- ✅ Built-in help: `devgitsecops setup-terraform-backend --help`
- ✅ Part of your unified DevOps tool

## Integration with CI/CD

### Azure DevOps

```yaml
- task: Bash@3
  inputs:
    targetType: 'inline'
    script: |
      # Ensure devgitsecops is available
      ./devgitsecops setup-terraform-backend \
        --environment $(Environment) \
        --auto-approve
      
      ./devgitsecops terraform init
      ./devgitsecops terraform plan
```

### GitHub Actions

```yaml
- name: Setup Terraform Backend
  run: |
    ./devgitsecops setup-terraform-backend \
      --environment production \
      --auto-approve
    
    ./devgitsecops terraform init
    ./devgitsecops terraform plan
  env:
    ARM_SUBSCRIPTION_ID: ${{ secrets.AZURE_SUBSCRIPTION_ID }}
```

## Multiple Environments

Create backends for different environments:

```bash
# Development
devgitsecops setup-terraform-backend --environment dev

# Staging
devgitsecops setup-terraform-backend --environment staging

# Production
devgitsecops setup-terraform-backend --environment production
```

Each environment gets its own isolated backend infrastructure!

## Cleanup

To remove all backend resources:

```bash
devgitsecops az group delete --name terraform-production-rg --yes --no-wait
```
