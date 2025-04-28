#!/usr/bin/env bash

set -euo pipefail

# Function to detect architecture and OS
get_arch_os() {
  ARCH=$(uname -m)
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  case "$ARCH" in
    x86_64|amd64)
      ARCH="x86_64"
      ;;
    aarch64|arm64)
      ARCH="aarch64"
      ;;
    armv7l)
      ARCH="arm"
      ;;
    i686)
      ARCH="i686"
      ;;
    *)
      echo "Unsupported architecture: $ARCH" >&2
      exit 1
      ;;
  esac

  case "$OS" in
    linux)
      OS="unknown-linux-gnu"
      ;;
    darwin)
      OS="apple-darwin"
      ;;
    *)
      echo "Unsupported OS: $OS" >&2
      exit 1
      ;;
  esac

  echo "$ARCH-$OS"
}

# Get latest version from GitHub API
get_latest_version() {
  curl -s https://api.github.com/repos/dandavison/delta/releases/latest | grep 'tag_name' | cut -d '"' -f 4
}

# Main install function
install_git_delta() {
  ARCH_OS=$(get_arch_os)
  VERSION=$(get_latest_version)

  FILENAME="delta-${VERSION#v}-$ARCH_OS.tar.gz"
  URL="https://github.com/dandavison/delta/releases/download/$VERSION/$FILENAME"

  TMPDIR=$(mktemp -d)
  trap 'rm -rf "$TMPDIR"' EXIT

  curl -sSL -o "$TMPDIR/$FILENAME" "$URL"
  tar -xzf "$TMPDIR/$FILENAME" -C "$TMPDIR"

  DELTA_PATH=$(find "$TMPDIR" -type f -name delta | head -n1)
  if [ -z "$DELTA_PATH" ]; then
    echo "Failed to find delta binary in archive." >&2
    exit 1
  fi

  install -m 0755 "$DELTA_PATH" /usr/local/bin/delta

  echo "git-delta installed successfully!"
}

install_git_delta
