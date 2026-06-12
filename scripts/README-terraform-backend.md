# Azure Terraform Backend Setup Scripts

These scripts automate the creation of all Azure resources needed to store Terraform state files securely.

## What Gets Created

1. **Resource Group** - Contains all backend resources
2. **Storage Account** - Stores the Terraform state file
   - Encryption enabled
   - HTTPS only
   - TLS 1.2 minimum
   - Public access disabled
   - Blob versioning enabled (for state file protection)
3. **Blob Container** - Named "tfstate" for state files
4. **Key Vault** - Securely stores credentials and configuration
   - Storage account key
   - Backend configuration details
   - Access policies configured

## Prerequisites

- Azure CLI installed
- Logged in: `az login`
- Appropriate permissions to create resources

## Usage

### PowerShell (Windows)

```powershell
# Basic usage (creates dev environment)
.\setup-terraform-backend-azure.ps1

# Specify environment
.\setup-terraform-backend-azure.ps1 -Environment "production"

# Specify environment and location
.\setup-terraform-backend-azure.ps1 -Environment "staging" -Location "westus2"

# Custom prefix
.\setup-terraform-backend-azure.ps1 -Environment "prod" -Prefix "mycompany"
```

### Bash (Linux/Mac)

```bash
# Make script executable
chmod +x setup-terraform-backend-azure.sh

# Basic usage (creates dev environment)
./setup-terraform-backend-azure.sh

# Specify environment
./setup-terraform-backend-azure.sh production

# With location override
AZURE_LOCATION=westus2 ./setup-terraform-backend-azure.sh staging
```

## Script Parameters

### PowerShell
- `-Environment` - Environment name (default: "dev")
- `-Location` - Azure region (default: "eastus")
- `-Prefix` - Resource name prefix (default: "terraform")

### Bash
- `$1` - Environment name (default: "dev")
- `AZURE_LOCATION` environment variable - Azure region (default: "eastus")

## Output

The script will create resources and output:

1. **Terraform backend configuration** to add to your `backend.tf`
2. **Environment variables** (optional authentication method)
3. **Key Vault name and secrets** for accessing configuration
4. **Next steps** for using the backend

## Using the Backend in Terraform

### Method 1: Backend Configuration File

Create `backend.tf`:

```hcl
terraform {
  backend "azurerm" {
    resource_group_name  = "terraform-prod-rg"
    storage_account_name = "terraformprodsa12345"
    container_name       = "tfstate"
    key                  = "terraform.tfstate"
  }
}
```

Then run:
```bash
az login
terraform init
```

### Method 2: Environment Variables

Set environment variables:

**PowerShell:**
```powershell
$env:ARM_SUBSCRIPTION_ID = "your-subscription-id"
$env:ARM_ACCESS_KEY = (az keyvault secret show --vault-name terraform-prod-kv-12345 --name terraform-backend-storage-key --query value -o tsv)
```

**Bash:**
```bash
export ARM_SUBSCRIPTION_ID="your-subscription-id"
export ARM_ACCESS_KEY=$(az keyvault secret show --vault-name terraform-prod-kv-12345 --name terraform-backend-storage-key --query value -o tsv)
```

Then initialize:
```bash
terraform init
```

### Method 3: Use with devgitsecops

```bash
# Install az cli and link it
devgitsecops install az --auto

# Run the setup script
devgitsecops az account set --subscription "your-subscription"
devgitsecops az login

# Run setup (using PowerShell on Windows)
.\scripts\setup-terraform-backend-azure.ps1 -Environment "production"

# Use terraform through devgitsecops
devgitsecops terraform init
devgitsecops terraform plan
devgitsecops terraform apply
```

## Resource Naming Convention

Resources are named using this pattern:

- Resource Group: `{prefix}-{environment}-rg`
- Storage Account: `{prefix}{environment}sa{random}`
- Key Vault: `{prefix}-{environment}-kv-{random}`
- Container: `tfstate`

Example for production environment:
- `terraform-prod-rg`
- `terraformprodsa12345`
- `terraform-prod-kv-12345`
- `tfstate`

## Security Features

✅ **Encryption at rest** - Storage account encryption enabled  
✅ **Encryption in transit** - HTTPS only, TLS 1.2 minimum  
✅ **No public access** - Blob public access disabled  
✅ **Version control** - Blob versioning enabled for state file recovery  
✅ **Secure storage** - Storage keys stored in Key Vault  
✅ **Access control** - Key Vault access policies configured  

## State File Versioning

Blob versioning is enabled, which means:
- Every state file change creates a new version
- Previous versions are retained
- You can recover from accidental deletions or corruption

To list versions:
```bash
az storage blob list \
  --account-name {storage-account-name} \
  --container-name tfstate \
  --include v \
  --auth-mode login
```

## Cleanup

To delete all resources:

```bash
# PowerShell
az group delete --name terraform-prod-rg --yes --no-wait

# Bash
az group delete --name terraform-prod-rg --yes --no-wait
```

## Troubleshooting

### "Storage account name already exists"
Storage account names must be globally unique. The script adds a random suffix, but if you get this error, re-run the script to generate a new name.

### "Key Vault name already exists"
Key Vault names must be globally unique. Re-run the script to generate a new random suffix.

### "Insufficient permissions"
Ensure you have:
- Contributor or Owner role on the subscription
- Permissions to create resource groups
- Permissions to create Key Vault and set access policies

### Cannot access Key Vault secrets
Run:
```bash
az keyvault set-policy \
  --name {key-vault-name} \
  --object-id $(az ad signed-in-user show --query id -o tsv) \
  --secret-permissions get list
```

## Best Practices

1. **Use separate backends per environment** - Run script once for dev, staging, and production
2. **Store Key Vault name** - Document the Key Vault name for your team
3. **Use service principals for CI/CD** - Don't use personal accounts in automation
4. **Enable soft delete** - Key Vault soft delete is enabled by default (90 days recovery)
5. **Monitor access** - Review Key Vault access logs regularly
6. **Backup state files** - While versioning helps, consider additional backups for critical environments

## Integration with CI/CD

### Azure DevOps

```yaml
- task: AzureCLI@2
  inputs:
    azureSubscription: 'your-service-connection'
    scriptType: 'bash'
    scriptLocation: 'inlineScript'
    inlineScript: |
      export ARM_ACCESS_KEY=$(az keyvault secret show \
        --vault-name terraform-prod-kv-12345 \
        --name terraform-backend-storage-key \
        --query value -o tsv)
      terraform init
      terraform plan
```

### GitHub Actions

```yaml
- name: Setup Terraform Backend
  run: |
    export ARM_ACCESS_KEY=$(az keyvault secret show \
      --vault-name terraform-prod-kv-12345 \
      --name terraform-backend-storage-key \
      --query value -o tsv)
    terraform init
  env:
    ARM_SUBSCRIPTION_ID: ${{ secrets.AZURE_SUBSCRIPTION_ID }}
```

## Additional Resources

- [Terraform Azure Backend Documentation](https://www.terraform.io/docs/language/settings/backends/azurerm.html)
- [Azure Storage Security](https://docs.microsoft.com/en-us/azure/storage/common/storage-security-guide)
- [Azure Key Vault Best Practices](https://docs.microsoft.com/en-us/azure/key-vault/general/best-practices)
