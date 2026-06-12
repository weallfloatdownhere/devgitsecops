# devgitsecops

A DevOps toolkit manager with automated setup commands for common infrastructure tasks. Easily install, manage, and automate your DevOps tools.

## Features

📥 **Auto-Download** - Automatically downloads and installs DevOps tools from official sources  
🔧 **Tool Management** - Install, check status, and manage multiple DevOps tools  
📦 **Easy Distribution** - Single binary that handles all tool installations  
⚡ **Quick Setup** - Just run `devgitsecops install --all` and you're ready to go!  
☁️ **Automation Commands** - Built-in commands for common tasks (Terraform backend setup, etc.)  
🔒 **Secure** - Downloads from official sources, stores credentials in Key Vault

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

Once installed, use the tools directly:

```bash
# Tools are installed and ready to use
kubectl get pods
helm list
terraform init
az login
aws s3 ls
kustomize build ./overlays/production

# Or use automation commands
devgitsecops terraform setup-backend-azure --environment production
```

## Usage Examples

### Tool Management

```bash
# Check what tools are installed
devgitsecops status

# Install specific tools
devgitsecops install kubectl
devgitsecops install terraform

# Install all supported tools
devgitsecops install --all

# Check version
devgitsecops version
```

### Automation Commands

#### Terraform Backend Setup

Quickly set up Azure infrastructure for Terraform remote state:

```bash
# Setup backend for development environment (handles Azure login automatically!)
devgitsecops terraform setup-backend-azure --environment dev

# Setup for production with custom location
devgitsecops terraform setup-backend-azure --environment production --location westus2

# Auto-approve without confirmation
devgitsecops terraform setup-backend-azure --environment staging --auto-approve
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
├── main.go                     # Entry point
├── cmd/                        # CLI commands
│   ├── root.go                # Root command
│   ├── install.go             # Tool installation
│   ├── status.go              # Status checking
│   ├── version.go             # Version info
│   ├── terraform.go           # Terraform parent command
│   └── setup_terraform_backend.go  # Azure backend setup subcommand
├── internal/
│   ├── executor/              # Tool execution logic
│   │   └── executor.go
│   └── downloader/            # Tool download logic
│       └── downloader.go
├── scripts/                    # Helper scripts (legacy)
├── docs/                       # Documentation
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

### Adding New Automation Commands

To add new automation commands like `setup-backend-azure`:

1. Create a parent command file if needed (e.g., `cmd/kubernetes.go`)
2. Create the subcommand file (e.g., `cmd/setup_cluster.go`)
3. Register the subcommand with the parent in its `init()` function
4. Follow the existing patterns for Azure login, confirmation, etc.

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