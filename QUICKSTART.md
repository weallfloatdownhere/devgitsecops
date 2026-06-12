# devgitsecops - Quick Start Guide

## What is devgitsecops?

A unified command-line interface that **automatically downloads and embeds** multiple DevOps tools. Instead of managing 6+ different CLI tools, you can use one tool that handles everything!

## Installation Status

✅ **Go 1.26.4** - Installed at `C:\Go`  
✅ **devgitsecops** - Built at `bin\devgitsecops.exe`
✅ **Auto-Download** - Automatically downloads and installs tools on demand!

## Quick Start

### 1. Check Tool Status

```powershell
.\bin\devgitsecops.exe status
```

This shows which tools are installed and ready to use.

### 2. Install Tools

**The easiest way - automatic download:**

```powershell
# Download all supported tools automatically
.\bin\devgitsecops.exe install --all

# Or install specific tools
.\bin\devgitsecops.exe install kubectl     # Downloads kubectl automatically!
.\bin\devgitsecops.exe install helm        # Downloads helm automatically!
.\bin\devgitsecops.exe install terraform   # Downloads terraform automatically!
.\bin\devgitsecops.exe install kustomize   # Downloads kustomize automatically!
```

**Supported auto-download tools:**
- ✅ kubectl - Downloaded from official Kubernetes releases
- ✅ kustomize - Downloaded from GitHub releases  
- ✅ helm - Downloaded from official Helm releases
- ✅ terraform - Downloaded from HashiCorp releases
- ⚠️ az cli - Must be installed manually (large Python-based tool)
- ⚠️ aws cli - Must be installed manually (large Python-based tool)

**Alternative: Link existing installations**

If you already have tools installed:

```powershell
# Auto-detect and link from system PATH
.\bin\devgitsecops.exe install kubectl --auto

# Or manually specify path
.\bin\devgitsecops.exe install kubectl --from "C:\path\to\kubectl.exe"
```

### 3. Use the Tools

Once tools are installed, use them through devgitsecops:

```powershell
# Kubernetes operations
.\bin\devgitsecops.exe kubectl get pods
.\bin\devgitsecops.exe kubectl apply -f deployment.yaml

# Helm operations
.\bin\devgitsecops.exe helm list
.\bin\devgitsecops.exe helm install myapp ./chart

# Terraform operations
.\bin\devgitsecops.exe terraform init
.\bin\devgitsecops.exe terraform plan

# Azure CLI
.\bin\devgitsecops.exe az login
.\bin\devgitsecops.exe az vm list

# AWS CLI
.\bin\devgitsecops.exe aws s3 ls
.\bin\devgitsecops.exe aws ec2 describe-instances

# Kustomize
.\bin\devgitsecops.exe kustomize build ./overlays/production
```

## Add to PATH (Optional)

To use `devgitsecops` from anywhere:

### Temporary (Current Session Only)

```powershell
$env:Path = "$PWD\bin;" + $env:Path
devgitsecops status
```

### Permanent (System-Wide)

Run PowerShell as Administrator:

```powershell
$currentPath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
$newPath = "$PWD\bin;$currentPath"
[System.Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
```

Then restart your terminal and use:

```powershell
devgitsecops kubectl get pods
devgitsecops helm list
```

## Real-World Example Workflow

```powershell
# 1. Setup Azure Terraform backend (NEW! - includes auto-login)
devgitsecops setup-terraform-backend --environment production --location eastus

# 2. Check cluster connection
devgitsecops kubectl cluster-info

# 3. Build kustomize configuration
devgitsecops kustomize build ./environments/production

# 4. Apply to cluster
devgitsecops kubectl apply -k ./environments/production

# 5. Deploy with Helm
devgitsecops helm upgrade --install myapp ./charts/myapp

# 6. Provision infrastructure with Terraform
devgitsecops terraform init
devgitsecops terraform apply -var-file=production.tfvars

# 7. Manage Azure resources
devgitsecops az aks get-credentials --resource-group myRG --name myCluster

# 8. Manage AWS resources
devgitsecops aws eks update-kubeconfig --name my-cluster
```

## Tool Storage Location

All linked tools are stored in: `C:\Users\h4ckb\.devgitsecops\bin\`

This keeps your system clean and organized!

## Getting Help

```powershell
# General help
.\bin\devgitsecops.exe --help

# Command-specific help
.\bin\devgitsecops.exe install --help
.\bin\devgitsecops.exe kubectl --help

# Version information
.\bin\devgitsecops.exe version
```

## Next Steps

1. Install the DevOps tools you use (kubectl, helm, terraform, etc.)
2. Run `devgitsecops install --all --auto` to link them
3. Start using all your tools through one unified interface!

## Benefits

✅ **Single Binary** - One tool instead of 6+  
✅ **Consistent Interface** - Same command structure for all   
✅ **Cloud Automation** - Built-in commands for common cloud tasks  
✅ **Terraform Backend Setup** - One command to setup Azure backend infrastructuretools  
✅ **Easy Distribution** - Share one binary with your team  
✅ **Organized** - All tools in one managed location  
✅ **Version Tracking** - Check versions of all tools at once

## Troubleshooting

### Tool Not Found After Installation

```powershell
# Refresh environment and try again
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
devgitsecops status
```

### Manual Tool Installation

If a tool isn't in PATH, find it and link manually:

```powershell
# Find kubectl
where.exe kubectl

# Link it
devgitsecops install kubectl --from "C:\path\shown\above\kubectl.exe"
```

## Building from Source

Already done! But for reference:

```powershell
# Download dependencies
go mod tidy

# Build
go build -v -o bin\devgitsecops.exe main.go

# Or use the build script
.\build.bat
```

---

**Happy DevOps-ing! 🚀**
