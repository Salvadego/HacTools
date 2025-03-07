#!/usr/bin/env bash
# hactools installer
#
# Usage:
#   ./install.sh [options]
#
# Options:
#   --help       Show this help message
#   --bin=PATH   Install binaries to PATH (e.g. --bin=$HOME/bin)
#   --no-update  Skip updating existing installation
#
# Environment variables:
#   HACTOOLS_INSTALL_DIR  Installation directory (default: $HOME/.hactools)
#   HACTOOLS_BIN_DIR      Binary directory (default: $HACTOOLS_INSTALL_DIR/bin)

set -e

HACTOOLS_REPO="github.com/SalvadegoDev/HacTools"
HACTOOLS_VERSION="latest"
HACTOOLS_INSTALL_DIR="${HACTOOLS_INSTALL_DIR:-$HOME/.config/hactools}"
HACTOOLS_BIN_DIR="${HACTOOLS_INSTALL_DIR/bin}"
HACTOOLS_CMDS=("xf" "xg" "ii")
NO_UPDATE=0

function print_help() {
  cat <<EOF
hactools installer

Usage:
  ./install.sh [options]

Options:
  --help       Show this help message
  --bin=PATH   Install binaries to PATH (e.g. --bin=\$HOME/bin)
  --no-update  Skip updating existing installation

Environment variables:
  HACTOOLS_INSTALL_DIR  Installation directory (default: \$HOME/.hactools)
  HACTOOLS_BIN_DIR      Binary directory (default: \$HACTOOLS_INSTALL_DIR/bin)
EOF
}

# Parse arguments
for arg in "$@"; do
  case $arg in
    --help)
      print_help
      exit 0
      ;;
    --bin=*)
      HACTOOLS_BIN_DIR="${arg#*=}"
      ;;
    --no-update)
      NO_UPDATE=1
      ;;
    *)
      echo "Unknown option: $arg"
      print_help
      exit 1
      ;;
  esac
done

if [ -z "$HACTOOLS_BIN_DIR" ]; then
  HACTOOLS_BIN_DIR="$HACTOOLS_INSTALL_DIR/bin"
fi

# Check for Go installation
if ! command -v go &> /dev/null; then
  echo "Error: Go is not installed. Please install Go first."
  echo "Visit https://golang.org/doc/install for installation instructions."
  exit 1
fi

# Create installation directories
mkdir -p "$HACTOOLS_INSTALL_DIR" "$HACTOOLS_BIN_DIR"

# Backup current binaries if they exist
backup_dir="$HACTOOLS_INSTALL_DIR/backup-$(date +%Y%m%d%H%M%S)"
need_backup=0

for cmd in "${HACTOOLS_CMDS[@]}"; do
  if [ -f "$HACTOOLS_BIN_DIR/$cmd" ]; then
    need_backup=1
    break
  fi
done

if [ $need_backup -eq 1 ] && [ $NO_UPDATE -eq 0 ]; then
  echo "Backing up existing installation..."
  mkdir -p "$backup_dir"
  for cmd in "${HACTOOLS_CMDS[@]}"; do
    if [ -f "$HACTOOLS_BIN_DIR/$cmd" ]; then
      cp "$HACTOOLS_BIN_DIR/$cmd" "$backup_dir/" 2>/dev/null || true
    fi
  done
fi

echo "Installing hactools to $HACTOOLS_INSTALL_DIR..."

# Download and install the tools
temp_dir=$(mktemp -d)
trap 'rm -rf $temp_dir' EXIT

cd "$temp_dir"

echo "Downloading hactools..."
if ! go install "$HACTOOLS_REPO/cmd/flex@$HACTOOLS_VERSION" "$HACTOOLS_REPO/cmd/groovy@$HACTOOLS_VERSION" "$HACTOOLS_REPO/cmd/impex@$HACTOOLS_VERSION"; then
  echo "Error: Failed to download and install hactools."
  exit 1
fi

# Copy the binaries to the installation directory
for cmd in "${HACTOOLS_CMDS[@]}"; do
  src_cmd=""
  case $cmd in
    xf) src_cmd="flex" ;;
    xg) src_cmd="groovy" ;;
    ii) src_cmd="impex" ;;
  esac
  
  if [ -f "$GOPATH/bin/$src_cmd" ]; then
    cp "$GOPATH/bin/$src_cmd" "$HACTOOLS_BIN_DIR/$cmd"
  elif [ -f "$HOME/go/bin/$src_cmd" ]; then
    cp "$HOME/go/bin/$src_cmd" "$HACTOOLS_BIN_DIR/$cmd"
  else
    echo "Warning: Could not find $src_cmd binary."
  fi
done

# Make binaries executable
chmod +x "$HACTOOLS_BIN_DIR"/* 2>/dev/null || true

# Generate shell configuration
shell_config=""
case $SHELL in
  */zsh)
    shell_config="${HOME}/.zshrc"
    ;;
  */bash)
    shell_config="${HOME}/.bashrc"
    if [[ "$OSTYPE" == "darwin"* ]]; then
      shell_config="${HOME}/.bash_profile"
    fi
    ;;
  */fish)
    shell_config="${HOME}/.config/fish/config.fish"
    ;;
  *)
    shell_config="${HOME}/.bashrc"
    ;;
esac

# Check if the bin directory is in PATH
if [[ ":$PATH:" != *":$HACTOOLS_BIN_DIR:"* ]]; then
  echo
  echo "To enable hactools, add the following line to $shell_config:"
  echo
  echo "  export PATH=\"\$PATH:$HACTOOLS_BIN_DIR\""
  echo
  echo "Then restart your shell or run:"
  echo
  echo "  export PATH=\"\$PATH:$HACTOOLS_BIN_DIR\""
  echo
fi

# Add example usage
echo "Installation completed!"

# Quick help
echo "For help, run:"
echo "  xf --help"
echo "  xg --help"
echo "  ii --help"
