#!/bin/bash

# DocRadar 打包脚本
# 用法: ./build.sh [选项]
# 选项:
#   all      - 构建所有平台 (默认)
#   mac      - 仅构建 macOS
#   win      - 仅构建 Windows
#   clean    - 清理构建目录

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目信息
APP_NAME="doc-radar"
VERSION="1.0.0"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BUILD_DIR="${SCRIPT_DIR}/build/bin"
DIST_DIR="${SCRIPT_DIR}/dist"

# 打印带颜色的消息
info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# 检查依赖
check_dependencies() {
    info "检查依赖..."
    
    if ! command -v wails &> /dev/null; then
        error "未找到 wails 命令，请先安装: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    fi
    
    if ! command -v go &> /dev/null; then
        error "未找到 go 命令，请先安装 Go"
    fi
    
    if ! command -v npm &> /dev/null; then
        error "未找到 npm 命令，请先安装 Node.js"
    fi
    
    success "依赖检查通过"
}

# 清理构建目录
clean() {
    info "清理构建目录..."
    rm -rf "$BUILD_DIR"
    rm -rf "$DIST_DIR"
    success "清理完成"
}

# 创建发布目录
prepare_dist() {
    info "准备发布目录..."
    mkdir -p "$DIST_DIR"
}

# 构建 macOS 版本
build_mac() {
    info "构建 macOS 版本..."
    
    # 构建 ARM64 (Apple Silicon)
    info "  构建 macOS ARM64..."
    wails build -platform darwin/arm64 -o "${APP_NAME}"
    
    # 构建 AMD64 (Intel)
    info "  构建 macOS AMD64..."
    wails build -platform darwin/amd64 -o "${APP_NAME}-intel"
    
    # 打包
    if [ -d "$BUILD_DIR/${APP_NAME}.app" ]; then
        info "  打包 macOS ARM64..."
        cd "$BUILD_DIR"
        zip -r -q "${DIST_DIR}/${APP_NAME}-${VERSION}-macos-arm64.zip" "${APP_NAME}.app"
        cd - > /dev/null
        success "  macOS ARM64 构建完成: dist/${APP_NAME}-${VERSION}-macos-arm64.zip"
    fi

    if [ -d "$BUILD_DIR/${APP_NAME}-intel.app" ]; then
        info "  打包 macOS AMD64..."
        cd "$BUILD_DIR"
        zip -r -q "${DIST_DIR}/${APP_NAME}-${VERSION}-macos-amd64.zip" "${APP_NAME}-intel.app"
        cd - > /dev/null
        success "  macOS AMD64 构建完成: dist/${APP_NAME}-${VERSION}-macos-amd64.zip"
    fi
}

# 构建 Windows 版本
build_win() {
    info "构建 Windows 版本..."
    
    # 构建 AMD64
    info "  构建 Windows AMD64..."
    wails build -platform windows/amd64 -o "${APP_NAME}.exe"
    
    # 构建 386 (32位)
    info "  构建 Windows 386..."
    wails build -platform windows/386 -o "${APP_NAME}-x86.exe"
    
    # 打包
    if [ -f "$BUILD_DIR/${APP_NAME}.exe" ]; then
        info "  打包 Windows AMD64..."
        cd "$BUILD_DIR"
        zip -q "${DIST_DIR}/${APP_NAME}-${VERSION}-windows-amd64.zip" "${APP_NAME}.exe"
        cd - > /dev/null
        success "  Windows AMD64 构建完成: dist/${APP_NAME}-${VERSION}-windows-amd64.zip"
    fi

    if [ -f "$BUILD_DIR/${APP_NAME}-x86.exe" ]; then
        info "  打包 Windows 386..."
        cd "$BUILD_DIR"
        zip -q "${DIST_DIR}/${APP_NAME}-${VERSION}-windows-386.zip" "${APP_NAME}-x86.exe"
        cd - > /dev/null
        success "  Windows 386 构建完成: dist/${APP_NAME}-${VERSION}-windows-386.zip"
    fi
}

# 构建所有平台
build_all() {
    check_dependencies
    clean
    prepare_dist
    build_mac
    build_win
    
    echo ""
    success "=========================================="
    success "所有平台构建完成!"
    success "=========================================="
    echo ""
    info "构建产物:"
    ls -lh "$DIST_DIR"
}

# 显示帮助
show_help() {
    echo "DocRadar 打包脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  all      构建所有平台 (默认)"
    echo "  mac      仅构建 macOS (ARM64 + AMD64)"
    echo "  win      仅构建 Windows (AMD64 + 386)"
    echo "  clean    清理构建目录"
    echo "  help     显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0           # 构建所有平台"
    echo "  $0 mac       # 仅构建 macOS"
    echo "  $0 win       # 仅构建 Windows"
    echo "  $0 clean     # 清理构建目录"
}

# 主函数
main() {
    case "${1:-all}" in
        all)
            build_all
            ;;
        mac)
            check_dependencies
            prepare_dist
            build_mac
            ;;
        win)
            check_dependencies
            prepare_dist
            build_win
            ;;
        clean)
            clean
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            error "未知选项: $1"
            show_help
            ;;
    esac
}

main "$@"

