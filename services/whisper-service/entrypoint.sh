#!/bin/bash

set -euo pipefail

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly SCRIPT_NAME="$(basename "$0")"

if [[ -f "${SCRIPT_DIR}/build-functions.sh" ]]; then
    source "${SCRIPT_DIR}/build-functions.sh"
else
    echo "ERROR: build-functions.sh not found in ${SCRIPT_DIR}" >&2
    exit 1
fi

# ANSI color codes
COLOR_INFO="\033[0m"       # White
COLOR_WARN="\033[1;33m"    # Yellow
COLOR_ERROR="\033[1;31m"   # Red
COLOR_SUCCESS="\033[1;32m" # Green
COLOR_RESET="\033[0m"

info()    { echo -e "${COLOR_INFO}$*${COLOR_RESET}" >&2; }
warn()    { echo -e "${COLOR_WARN}$*${COLOR_RESET}" >&2; }
error()   { echo -e "${COLOR_ERROR}$*${COLOR_RESET}" >&2; exit 1; }
success() { echo -e "${COLOR_SUCCESS}$*${COLOR_RESET}" >&2; }

show_help() {
    cat << EOF
LYRA Whisper Service Build Script

USAGE:
    $SCRIPT_NAME [COMMAND]

COMMANDS:
    install_deps     Install system dependencies
    build_libs       Build custom libraries (inih, curlpp, pugixml)
    create_configs   Create CMake configuration files
    build_grpc       Build and install gRPC
    build_minio      Build and install minio-cpp
    build_app        Build the application
    full_build       Run complete build process
    help             Show this help message

EXAMPLES:
    $SCRIPT_NAME install_deps
    $SCRIPT_NAME full_build

EOF
}

error_exit() {
    error "$*"
}

validate_environment() {
    if [[ ! -d "/usr/local" ]]; then
        error_exit "Directory /usr/local not found"
    fi

    command -v apt-get >/dev/null 2>&1 || error_exit "apt-get not found"
    command -v cmake >/dev/null 2>&1 || error_exit "cmake not found"
    command -v git >/dev/null 2>&1 || error_exit "git not found"
}

main() {
    info "Starting LYRA Whisper Service build process"
    validate_environment

    case "${1:-build}" in
        "install_deps")
            info "Installing system dependencies..."
            install_system_dependencies
            ;;
        "build_libs")
            info "Building custom libraries..."
            build_inih_library
            build_curlpp_library
            build_pugixml_static_library
            ;;
        "create_configs")
            info "Creating CMake configuration files..."
            create_cmake_configs
            ;;
        "build_grpc")
            info "Building gRPC..."
            build_grpc
            ;;
        "build_minio")
            info "Building minio-cpp..."
            build_minio_cpp
            ;;
        "build_app")
            info "Building application..."
            build_application
            ;;
        "full_build")
            info "Running full build process..."
            install_system_dependencies
            build_inih_library
            build_curlpp_library
            build_pugixml_static_library
            create_cmake_configs
            build_grpc
            build_minio_cpp
            build_application
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            error_exit "Unknown command: $1. Use '$SCRIPT_NAME help' for usage information."
            ;;
    esac

    success "Build process completed successfully"
}

main "$@"
