# devgitsecops

A unified command-line interface that embeds multiple DevOps tools into a single binary, making it easier to manage your infrastructure and deployments.

## Features

🚀 **Single Binary** - One tool to rule them all  
� **Auto-Download** - Automatically downloads and installs tools on first use  
🔧 **Multiple Tools** - Access all your favorite DevOps tools through one interface  
📦 **Easy Distribution** - Distribute one binary instead of managing multiple tool installations  
⚡ **Quick Setup** - Just run `devgitsecops install --all` and you're ready to go!  
☁️ **Cloud Automation** - Built-in commands for common cloud operations (Terraform backend setup, etc.)

## Embedded Tools

- **kubectl** - Kubernetes CLI
- **kustomize** - Kubernetes configuration management
- **helm** - Kubernetes package manager
- **terraform** - Infrastructure as Code
- **az** - Azure CLI
- **aws** - AWS CLI

## Installation

### Prerequisites

- Go 1.21 or higher
- The individual tools you want to use (kubectl, helm, terraform, etc.)

### Build from Source

```bash
# Clone the repository
cd projects/devgitsecops

# Download dependencies
go mod download

# Build the binary
make build

# Or install directly to $GOPATH/bin
make install
```

### Binary Releases

Download the latest release for your platform from the releases page.

## Quick Start

### 1. Build the Tool

```bash
make build
```

### 2. Check Tool Status

```bash
./bin/devgitsecops status
```

This will show which tools are installed and available.

### 3. Install/Link Tools

#### Automatic Download (Recommended)

The easiest way - let devgitsecops download the tools for you:

```bash
# Download all supported tools automatically
./bin/devgitsecops install --all

# Or download individual tools
./bin/devgitsecops install kubectl     # Auto-downloads kubectl
./bin/devgitsecops install helm        # Auto-downloads helm
./bin/devgitsecops install terraform   # Auto-downloads terraform
./bin/devgitsecops install kustomize   # Auto-downloads kustomize
```

**What gets downloaded:**
- ✅ **kubectl** - Latest stable from Kubernetes releases
- ✅ **kustomize** - Latest from GitHub releases
- ✅ **helm** - Latest from official Helm releases
- ✅ **terraform** - Latest from HashiCorp releases
- ⚠️ **az cli** - Must be installed manually (Python-based, very large)
- ⚠️ **aws cli** - Must be installed manually (Python-based, very large)

#### Link Existing Installations

If you prefer to use tools already installed on your system:

```bash
# Auto-detect from system PATH
./bin/devgitsecops install kubectl --auto

# Or specify exact path
./bin/devgitsecops install kubectl --from /usr/local/bin/kubectl
```

### 4. Use the Tools

Once installed, use the tools through devgitsecops:

```bash
# Kubernetes operations
./bin/devgitsecops kubectl get pods
./bin/devgitsecops kubectl apply -f deployment.yaml

# Helm operations
./bin/devgitsecops helm list
./bin/devgitsecops helm install myapp ./chart

# Terraform operations
./bin/devgitsecops terraform init
./bin/devgitsecops terraform plan
./bin/devgitsecops terraform apply

# Azure operations
./bin/devgitsecops az login
./bin/devgitsecops az vm list

# AWS operations
./bin/devgitsecops aws s3 ls
./bin/devgitsecops aws ec2 describe-instances

# Kustomize operations
./bin/devgitsecops kustomize build ./overlays/production
```

## Usage Examples

### Kubernetes Workflow

```bash
# Check cluster info
devgitsecops kubectl cluster-info

# Apply configurations
devgitsecops kubectl apply -k ./overlays/production

# Use helm
devgitsecops helm repo add stable https://charts.helm.sh/stable
devgitsecops helm install nginx stable/nginx-ingress
```

### Infrastructure as Code

```bash
# Initialize Terraform
devgitsecops terraform init

# Plan changes
devgitsecops terraform plan -out=tfplan

# Apply changes
devgitsecops terraform apply tfplan
```

### Cloud Operations

```bash
# Azure
devgitsecops az account list
devgitsecops az group create --name myResourceGroup --location eastus

# AWS
devgitsecops aws configure
devgitsecops aws s3 mb s3://my-bucket
```

### Terraform Backend Setup

Quickly set up Azure infrastructure for Terraform remote state:

```bash
# Setup backend for development environment (handles Azure login automatically!)
devgitsecops setup-terraform-backend --environment dev

# Setup for production with custom location
devgitsecops setup-terraform-backend --environment production --location westus2

# Auto-approve without confirmation
devgitsecops setup-terraform-backend --environment staging --auto-approve
```

This command automatically:
- ✅ Checks Azure login status and logs you in if needed
- ✅ Creates Resource Group
- ✅ Creates Storage Account (encrypted, versioned, HTTPS-only)
- ✅ Creates Blob Container for state files
- ✅ Creates Key Vault with stored credentials

After setup, you get a ready-to-use Terraform backend configuration!

## Configuration

devgitsecops stores tool binaries in `~/.devgitsecops/bin/`. You can customize this location by modifying the executor package.

### Config File

Create a config file at `~/.devgitsecops.yaml`:

```yaml
# Configuration options (future feature)
tools:
  kubectl:
    version: "1.28.0"
  helm:
    version: "3.12.0"
```

## Development

### Project Structure

```
devgitsecops/
├── main.go                 # Entry point
├── cmd/                    # CLI commands
│   ├── root.go            # Root command
│   ├── kubectl.go         # kubectl wrapper
│   ├── helm.go            # helm wrapper
│   ├── terraform.go       # terraform wrapper
│   ├── az.go              # az wrapper
│   ├── aws.go             # aws wrapper
│   ├── kustomize.go       # kustomize wrapper
│   ├── install.go         # Installation command
│   └── status.go          # Status command
├── internal/
│   └── executor/          # Tool execution logic
│       └── executor.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Format code
make fmt

# Run linter
make lint
```

### Adding New Tools

1. Create a new command file in `cmd/` directory (e.g., `cmd/newtool.go`)
2. Follow the pattern from existing tool commands
3. Add the tool to the status and install commands

## Troubleshooting

### Tool Not Found

```bash
# Check status
devgitsecops status

# Try auto-install
devgitsecops install <tool> --auto

# Or manually specify path
devgitsecops install <tool> --from /path/to/tool
```

### Permission Issues

On Unix systems, ensure binaries have execute permissions:

```bash
chmod +x ~/.devgitsecops/bin/*
```

### Path Issues

Add devgitsecops to your PATH:

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH="$PATH:/path/to/devgitsecops/bin"
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License

## Roadmap
x] Automatic binary downloads
- [x] Azure Terraform backend setup automation
- [ ] Version management for tools
- [ ] Update command to refresh tool versions
- [ ] AWS Terraform backend setup
- [ ] Kubernetes cluster setup helper
- [ ] Update command to refresh tool versions
- [ ] Plugin system for custom tools
- [ ] Shell completion scripts
- [ ] Interactive mode
- [ ] Tool configuration profiles
- [ ] Docker container support

## Credits

This tool is a wrapper around the following excellent tools:
- [kubectl](https://kubernetes.io/docs/reference/kubectl/)
- [kustomize](https://kustomize.io/)
- [helm](https://helm.sh/)
- [terraform](https://www.terraform.io/)
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/)
- [AWS CLI](https://aws.amazon.com/cli/)