# HacTools

CLI tools for easy interaction with SAP Hybris HAC (Hybris Administration Console).

## Features

- **FlexSearch (xf)**: Execute flexible search queries directly from command line
- **Groovy (xg)**: Run Groovy scripts against your Hybris instance
- **Impex (ii)**: Import Impex data from files or direct input

## Installation

### Using Go Install

```bash
go install github.com/SalvdegoDev/HacTools/cmd/flex@latest
go install github.com/SalvdegoDev/HacTools/cmd/groovy@latest
go install github.com/SalvdegoDev/HacTools/cmd/impex@latest
```

### Using Installation Script

For convenience, you can use this installation script:

```bash
# Install or update to the latest version
curl -sSL https://raw.githubusercontent.com/SalvdegoDev/HacTools/main/install.sh | bash

# Install specific version
curl -sSL https://raw.githubusercontent.com/SalvdegoDev/HacTools/main/install.sh | bash -s v1.0.0
```

## Usage

### Configuration

Set your Hybris HAC credentials through environment variables:

```bash
export HYBRIS_HAC_URL="https://your-hybris-instance:9002/hac"
export HYBRIS_USER="admin"
export HYBRIS_PASSWORD="your-password"
```

Or pass them directly as command line arguments:

```bash
xf --address="https://your-hybris-instance:9002/hac" --user="admin" --password="your-password" "SELECT * FROM {Product}"
```

### FlexSearch (xf)

Execute FlexibleSearch queries against Hybris:

```bash
# Run a query directly
xf "SELECT COUNT(*) FROM {Product}"

# Run a query from a file
xf ./my-query.sql

# Limit results
xf --max-count=100 "SELECT * FROM {Product}"

# Turn on debugging
xf --log-level=debug "SELECT * FROM {Product}"
```

### Groovy (xg)

Execute Groovy scripts against Hybris:

```bash
# Run a Groovy script directly
xg "return modelService.getAllItems(ProductModel.class).size()"

# Run a script from a file
xg ./my-script.groovy

# Run with commit
xg --commit my-script.groovy

# Run JavaScript instead of Groovy
xg --type=javascript "print('Hello from JavaScript');"

# Run BeanShell
xg --type=beanshell "print('Hello from BeanShell');"
```

### Impex (ii)

Import Impex data:

```bash
# Import from a file
ii ./my-impex.impex

# Import directly
ii "INSERT_UPDATE Product;code[unique=true];name[lang=en];$catalogVersion
;myProduct;My Test Product;Default:Staged"

# Enable code execution in Impex
ii --exec ./my-impex.impex

# Enable legacy mode
ii --legacy ./my-impex.impex
```

## Options

All commands share these common options:

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--address` | `-s` | HAC URL | `$HYBRIS_HAC_URL` or `https://localhost:9002/hac` |
| `--user` | `-u` | Username | `$HYBRIS_USER` or `admin` |
| `--password` | `-p` | Password | `$HYBRIS_PASSWORD` or `nimda` |
| `--log-level` | `-l` | Log level (debug, info, error, none) | `error` |

### FlexSearch (xf) Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--max-count` | `-m` | Maximum number of results | `10` |
| `--no-analyze` | `-A` | Do not analyze PK | `false` |
| `--no-blacklist` | `-B` | Ignore column blacklist | `false` |

### Groovy (xg) Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--commit` | `-c` | Execute with commit | `false` |
| `--type` | `-t` | Script type (groovy, javascript, beanshell) | `groovy` |

### Impex (ii) Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--legacy` | `-L` | Enable legacyMode | `false` |
| `--exec` | `-c` | Enable code execution | `false` |
| `--distributed` | `-d` | Enable distributed mode | `false` |
| `--sld` | ` ` | Enable SLD | `false` |

## Building from Source

```bash
# Clone the repository
git clone https://github.com/SalvdegoDev/HacTools.git
cd HacTools

# Build all tools
go build -o bin/xf ./cmd/flex
go build -o bin/xg ./cmd/groovy
go build -o bin/ii ./cmd/impex
```

## Contributing

Pull requests are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the submission process.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
