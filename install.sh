#!/bin/bash
#
# UDA 一键安装脚本
# 使用方法: curl -sSf https://raw.githubusercontent.com/zhaoyilun/uda/main/install.sh | sh
# 或: sh install.sh
#

set -e

# 配置
UDA_VERSION="0.1.0"
UDA_REPO="https://github.com/zhaoyilun/uda"
INSTALL_DIR="${HOME}/.local/bin"
UDA_DIR="${HOME}/.uda"

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检测操作系统
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux";;
        Darwin*)    echo "macos";;
        CYGWIN*|MINGW*|MSYS*) echo "windows";;
        *)          echo "unknown";;
    esac
}

# 检测架构
detect_arch() {
    case "$(uname -m)" in
        x86_64)    echo "x86_64";;
        aarch64|arm64) echo "aarch64";;
        *)          echo "unknown";;
    esac
}

# 下载文件
download_file() {
    local url=$1
    local dest=$2

    if command -v curl &> /dev/null; then
        curl -fsSL "$url" -o "$dest"
    elif command -v wget &> /dev/null; then
        wget -q "$url" -O "$dest"
    else
        log_error "需要 curl 或 wget"
        exit 1
    fi
}

# 主安装流程
main() {
    log_info "UDA 安装脚本 v${UDA_VERSION}"
    echo ""

    # 检测平台
    OS=$(detect_os)
    ARCH=$(detect_arch)

    if [ "$OS" = "unknown" ]; then
        log_error "不支持的操作系统"
        exit 1
    fi

    log_info "检测到系统: ${OS} ${ARCH}"

    # 创建安装目录
    log_info "创建安装目录..."
    mkdir -p "${INSTALL_DIR}"
    mkdir -p "${UDA_DIR}"

    # 下载 uda 二进制
    log_info "下载 UDA..."

    local ext=""
    if [ "$OS" = "windows" ]; then
        ext=".exe"
    fi

    local shell_name="$(basename "${SHELL}")"
    if [ -z "$shell_name" ]; then
        shell_name="bash"
    fi

    local uda_url="${UDA_REPO}/releases/latest/download/uda${ext}"
    local uda_path="${INSTALL_DIR}/uda${ext}"

    # 尝试下载，如果失败则提示用户手动构建
    if download_file "$uda_url" "$uda_path" 2>/dev/null; then
        chmod +x "$uda_path"
    else
        log_warn "无法从 GitHub 下载预编译版本"
        log_info "请从源码构建: https://github.com/zhaoyilun/uda#构建"
        exit 1
    fi

    # 安装 uv
    log_info "安装 uv..."
    "${uda_path}" self install

    # 添加到 PATH 提示
    echo ""
    log_info "安装完成!"
    echo ""
    echo "请将以下内容添加到你的 shell 配置文件 (~/.bashrc, ~/.zshrc 等):"
    echo ""
    echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
    echo "  eval \"\$(uda init ${shell_name})\""
    echo ""
    echo "然后重新加载配置:"
    echo "  source ~/.bashrc  # 或 source ~/.zshrc"
    echo ""
    echo "开始使用:"
    echo "  uda create myenv --python 3.11"
    echo "  eval \"\$(uda init bash)\""
    echo "  uda activate myenv"
    echo ""
}

# 运行
main "$@"
