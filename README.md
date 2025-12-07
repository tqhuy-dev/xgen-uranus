# 🚀 Uranus CLI

A powerful CLI tool designed to help developers quickly scaffold and manage microservice applications and repositories.

## 📦 Installation

### Option 1: Install from GitHub (Recommended)

```bash
go install github.com/tqhuy-dev/xgen-uranus/cmd/uranus@latest
```

After installation, the `uranus` binary will be available in your `$GOPATH/bin` directory.

> **Note:** Make sure `$GOPATH/bin` is in your `PATH`. Add this to your shell config if needed:
> ```bash
> export PATH=$PATH:$(go env GOPATH)/bin
> ```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/tqhuy-dev/xgen-uranus.git
cd xgen-uranus

# Build the binary
make build

# Or install to /usr/local/bin
make install
```

### Option 3: Download Pre-built Binary

Download the latest release from the [Releases](https://github.com/tqhuy-dev/xgen-uranus/releases) page.

## 🎯 Usage

### Show Help

```bash
uranus --help
```

### Generate Application

Create a new microservice application with a complete project structure:

```bash
# Generate app in current directory
uranus generate app --name simple_app

# Generate app with short flag
uranus generate app -n my_service
```

This creates:
```
simple_app/
├── cmd/
│   └── main.go
├── internal/
│   ├── handlers/
│   ├── services/
│   ├── models/
│   └── repository/
├── pkg/
├── configs/
│   └── config.yaml
├── api/
│   └── proto/
├── go.mod
└── README.md
```

### Generate Repository

Create a new repository with interface and implementation:

```bash
# Generate repo in current directory
uranus generate repo --name simple_repo

# Generate repo in a specific path
uranus generate repo --name path/simple_repo

# Example: create user repository inside internal folder
uranus generate repo --name internal/repository/user_repo
```

This creates:
```
simple_repo/
├── repository.go        # Repository interface
├── repository_impl.go   # In-memory implementation
└── repository_test.go   # Unit tests
```

### List Repositories

List all repositories in the current directory or a specified path:

```bash
# List repos in current directory
uranus list repo

# List repos in a specific path
uranus list repo --path ./internal
```

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- Make

### Build Commands

```bash
# Build binary for current OS
make build

# Build for all platforms (Linux, macOS, Windows)
make build-all

# Build for specific platform
make build-linux
make build-darwin
make build-windows

# Install to /usr/local/bin
make install

# Uninstall
make uninstall

# Clean build artifacts
make clean

# Run without building
make run ARGS="--help"

# Run tests
make test

# Format code
make fmt
```

### Project Structure

```
xgen-uranus/
├── cmd/
│   └── uranus/         # CLI entry point
│       └── main.go
├── commands/           # CLI commands (Cobra)
│   ├── root.go         # Root command
│   ├── generate.go     # Generate parent command
│   ├── generate_app.go # Generate app subcommand
│   ├── generate_repo.go# Generate repo subcommand
│   ├── generate_tpl/   # Templates for code generation
│   ├── list.go         # List parent command
│   └── list_repo.go    # List repo subcommand
├── common/             # Common utilities
├── interceptors/       # gRPC/HTTP interceptors
├── transport/          # Transport layer (gRPC, HTTP)
├── grpc_third_party/   # Third-party proto files
├── Makefile            # Build automation
└── README.md           # This file
```

## 📋 Commands Reference

| Command | Description |
|---------|-------------|
| `uranus --help` | Show help information |
| `uranus generate app --name <name>` | Generate a new application |
| `uranus generate repo --name <name>` | Generate a new repository |
| `uranus generate repo --name <path/name>` | Generate repo in specific path |
| `uranus list repo` | List all repositories |
| `uranus list repo --path <path>` | List repos in specific path |

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [xgen](https://github.com/tqhuy-dev/xgen) - Core library

