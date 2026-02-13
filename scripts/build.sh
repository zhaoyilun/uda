#!/bin/bash
#
# UDA 构建脚本
#

set -e

VERSION="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}"
OUTPUT="${OUTPUT:-uda}"
LDFLAGS="-s -w -X github.com/uda/uda/cmd.version=${VERSION}"

echo "Building UDA ${VERSION}..."

# 检测 Go
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

# 清理旧的构建
rm -f "${OUTPUT}"

# 构建
go build -ldflags="${LDFLAGS}" -o "${OUTPUT}" .

echo "Built: ${OUTPUT}"
ls -lh "${OUTPUT}"
