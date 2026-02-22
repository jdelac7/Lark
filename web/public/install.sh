#!/bin/sh
set -e

# Lark CLI installer
# Usage: curl -fsSL https://lark.black/install.sh | sh

GITHUB_REPO="jdelac7/Lark"
INSTALL_DIR="${LARK_INSTALL_DIR:-$HOME/.local/bin}"

# Colors (disabled if not a terminal)
if [ -t 1 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    CYAN='\033[0;36m'
    BOLD='\033[1m'
    RESET='\033[0m'
else
    RED=''
    GREEN=''
    YELLOW=''
    CYAN=''
    BOLD=''
    RESET=''
fi

info() {
    printf "${CYAN}%s${RESET}\n" "$1"
}

success() {
    printf "${GREEN}%s${RESET}\n" "$1"
}

warn() {
    printf "${YELLOW}%s${RESET}\n" "$1"
}

error() {
    printf "${RED}error: %s${RESET}\n" "$1" >&2
    exit 1
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        *)       error "Unsupported operating system: $(uname -s). Lark supports Linux and macOS." ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        *)             error "Unsupported architecture: $(uname -m). Lark supports amd64 and arm64." ;;
    esac
}

# Download a URL to a file (supports curl and wget)
download() {
    url="$1"
    output="$2"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$url" -o "$output"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$output" "$url"
    else
        error "Neither curl nor wget found. Please install one and try again."
    fi
}

# Download URL contents to stdout
download_stdout() {
    url="$1"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$url"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "$url"
    else
        error "Neither curl nor wget found. Please install one and try again."
    fi
}

# Fetch the latest release tag from GitHub
get_latest_version() {
    download_stdout "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" \
        | grep '"tag_name"' \
        | sed -E 's/.*"tag_name":\s*"([^"]+)".*/\1/'
}

main() {
    printf "\n${BOLD}Lark CLI Installer${RESET}\n\n"

    OS=$(detect_os)
    ARCH=$(detect_arch)
    info "Detected platform: ${OS}/${ARCH}"

    info "Fetching latest release..."
    VERSION=$(get_latest_version)
    if [ -z "$VERSION" ]; then
        error "Could not determine latest version. Check https://github.com/${GITHUB_REPO}/releases"
    fi
    info "Latest version: ${VERSION}"

    ARCHIVE="lark_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${ARCHIVE}"
    CHECKSUMS_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/checksums.txt"

    # Create temp directory with cleanup trap
    TMPDIR=$(mktemp -d)
    trap 'rm -rf "$TMPDIR"' EXIT

    info "Downloading ${ARCHIVE}..."
    download "$DOWNLOAD_URL" "$TMPDIR/$ARCHIVE"

    info "Verifying checksum..."
    download "$CHECKSUMS_URL" "$TMPDIR/checksums.txt"

    EXPECTED=$(grep "$ARCHIVE" "$TMPDIR/checksums.txt" | awk '{print $1}')
    if [ -z "$EXPECTED" ]; then
        error "Could not find checksum for ${ARCHIVE}"
    fi

    if command -v sha256sum >/dev/null 2>&1; then
        ACTUAL=$(sha256sum "$TMPDIR/$ARCHIVE" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
        ACTUAL=$(shasum -a 256 "$TMPDIR/$ARCHIVE" | awk '{print $1}')
    else
        warn "Warning: No sha256sum or shasum found, skipping checksum verification"
        ACTUAL="$EXPECTED"
    fi

    if [ "$ACTUAL" != "$EXPECTED" ]; then
        error "Checksum verification failed.\n  Expected: ${EXPECTED}\n  Got:      ${ACTUAL}"
    fi
    success "Checksum verified."

    info "Installing to ${INSTALL_DIR}..."
    mkdir -p "$INSTALL_DIR"
    tar -xzf "$TMPDIR/$ARCHIVE" -C "$TMPDIR"
    install -m 755 "$TMPDIR/lark" "$INSTALL_DIR/lark"

    success "Lark ${VERSION} installed to ${INSTALL_DIR}/lark"

    # Check if install dir is in PATH
    case ":$PATH:" in
        *":${INSTALL_DIR}:"*) ;;
        *)
            printf "\n"
            warn "Warning: ${INSTALL_DIR} is not in your PATH."
            printf "\n  Add it by appending this to your shell profile:\n\n"
            printf "    ${CYAN}export PATH=\"%s:\$PATH\"${RESET}\n\n" "$INSTALL_DIR"
            printf "  Then restart your shell or run: ${CYAN}source ~/.bashrc${RESET}\n"
            ;;
    esac

    printf "\n${BOLD}Get started:${RESET}\n\n"
    printf "  ${CYAN}1.${RESET} Subscribe at ${CYAN}https://lark.black${RESET} (\$2.99/month)\n"
    printf "  ${CYAN}2.${RESET} Activate:  ${GREEN}lark activate YOUR-LICENSE-KEY${RESET}\n"
    printf "  ${CYAN}3.${RESET} Play:      ${GREEN}lark${RESET}\n\n"
}

main
