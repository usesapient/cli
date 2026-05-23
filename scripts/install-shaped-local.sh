#!/usr/bin/env bash

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

resolve_install_path() {
    if [[ -n "${SAPIENT_INSTALL_PATH:-}" ]]; then
        printf '%s\n' "$SAPIENT_INSTALL_PATH"
        return
    fi

    if command -v sapient >/dev/null 2>&1; then
        command -v sapient
        return
    fi

    local gobin
    gobin="$(go env GOBIN)"
    if [[ -z "$gobin" ]]; then
        gobin="$(go env GOPATH)/bin"
    fi
    printf '%s\n' "$gobin/sapient"
}

install_path="$(resolve_install_path)"
if [[ -e "$install_path" ]] && command -v realpath >/dev/null 2>&1; then
    install_path="$(realpath "$install_path")"
fi

version="${SAPIENT_VERSION:-$(cd "$ROOT" && go run ./cmd/sapient-shaped version | awk 'NR == 1 { print $2 }')}"
build_time="${SAPIENT_BUILD_TIME:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"

mkdir -p "$(dirname "$install_path")"

if [[ -e "$install_path" ]]; then
    backup_path="${install_path}.before-shaped-cli.$(date -u +%Y%m%d%H%M%S)"
    cp -p "$install_path" "$backup_path"
    printf 'Backed up existing sapient binary to %s\n' "$backup_path"
fi

cd "$ROOT"
go build \
    -ldflags "-X main.version=$version -X main.buildTime=$build_time" \
    -o "$install_path" \
    ./cmd/sapient-shaped

printf 'Installed shaped sapient %s to %s\n' "$version" "$install_path"
