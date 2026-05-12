#!/bin/bash
set -e

BINARY="agent-chat"
OUTDIR="./exe"
LDFLAGS='-s -w -extldflags "-static"'
PLATFORM="${1:-both}"

mkdir -p "${OUTDIR}"

build() {
    local os=$1
    local arch=$2
    local output="${OUTDIR}/${BINARY}-${os}-${arch}"
    echo "Building ${output}..."
    CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -ldflags "${LDFLAGS}" -o "${output}"
    echo "Done: ${output}"
}

case "${PLATFORM}" in
    arm64)
        build linux arm64
        ;;
    amd64)
        build linux amd64
        ;;
    darwin-arm64)
        build darwin arm64
        ;;
    darwin-amd64)
        build darwin amd64
        ;;
    all|both)
        build linux arm64
        build linux amd64
        build darwin arm64
        build darwin amd64
        ;;
    *)
        echo "Usage: $0 {arm64|amd64|darwin-arm64|darwin-amd64|all}"
        exit 1
        ;;
esac
