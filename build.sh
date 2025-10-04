#!/usr/bin/env bash

set -euo pipefail

# ANSI color codes
COLOR_INFO="\033[0m"       # White (default)
COLOR_WARN="\033[1;33m"    # Yellow
COLOR_ERROR="\033[1;31m"   # Red
COLOR_SUCCESS="\033[1;32m" # Green
COLOR_RESET="\033[0m"

info()    { echo -e "${COLOR_INFO}$1${COLOR_RESET}"; }
warn()    { echo -e "${COLOR_WARN}$1${COLOR_RESET}"; }
error()   { echo -e "${COLOR_ERROR}$1${COLOR_RESET}" >&2; exit 1; }
success() { echo -e "${COLOR_SUCCESS}$1${COLOR_RESET}"; }

ORCHESTRATOR=""
TELEGRAM_TOKEN=""

ARTIFACT_PATHS=(
    "services/api-gateway/proto-context"
    "services/whisper-service/proto-context"
    "services/telegram-bot/proto-context"
)

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --podman|-p)
                ORCHESTRATOR="podman-compose"
                shift
                ;;
            --docker|-d)
                ORCHESTRATOR="docker-compose"
                shift
                ;;
            --telegram-token|-t)
                TELEGRAM_TOKEN="$2"
                shift 2
                ;;
            *)
                warn "Unknown argument: $1"
                shift
                ;;
        esac
    done
}

# Check if orchestrator is available
is_orchestrator_available() {
    case "$1" in
        podman-compose)
            command -v podman-compose &>/dev/null && return 0 || return 1
            ;;
        docker-compose)
            command -v docker-compose &>/dev/null && return 0 || return 1
            ;;
        "docker compose")
            command -v docker &>/dev/null && docker compose version &>/dev/null && return 0 || return 1
            ;;
        *)
            return 1
            ;;
    esac
}

# Select and validate orchestrator
select_orchestrator() {
    local candidates=("podman-compose" "docker-compose" "docker compose")
    if [ -n "$ORCHESTRATOR" ]; then
        if is_orchestrator_available "$ORCHESTRATOR"; then
            echo "$ORCHESTRATOR"
            return 0
        else
            error "$ORCHESTRATOR not found."
        fi
    else
        for orch in "${candidates[@]}"; do
            if is_orchestrator_available "$orch"; then
                echo "$orch"
                return 0
            fi
        done
        error "No supported container orchestrator found"
    fi
}

# Generate .env from .env.example
generate_env() {
    local env_file=".env"
    local example_file=".env.example"
    if [ ! -f "$env_file" ]; then
        if [ -f "$example_file" ]; then
            info "Generating .env from .env.example..."
            cp "$example_file" "$env_file"
            success ".env generated from .env.example"
        else
            error ".env.example not found, cannot generate .env"
        fi
    else
        info ".env already exists, skipping generation"
    fi
    if [ -n "$TELEGRAM_TOKEN" ]; then
        info "Setting TELEGRAM_BOT_TOKEN in .env..."
        if grep -q '^TELEGRAM_BOT_TOKEN=' "$env_file"; then
            sed -i "s|^TELEGRAM_BOT_TOKEN=.*|TELEGRAM_BOT_TOKEN=$TELEGRAM_TOKEN|" "$env_file"
        else
            echo "TELEGRAM_BOT_TOKEN=$TELEGRAM_TOKEN" >> "$env_file"
        fi
        success "TELEGRAM_BOT_TOKEN set in .env"
    fi
}

# Build project via compose.yml
build_project() {
    local compose_cmd="$1"
    info "Building and starting project using $compose_cmd..."
    $compose_cmd -f compose.yml up -d --build
    success "Build and startup completed"
}

# Remove all paths listed in ARTIFACT_PATHS.
clean_artifacts() {
    for path in "$@"; do
        if [ -e "$path" ]; then
            info "Removing $path..."
            rm -rf "$path" || error "Failed to remove $path"
        fi
    done
    success "Build artifact cleanup complete"
}

copy_proto_contexts() {
    local src_proto="proto"
    if [ ! -d "$src_proto" ]; then
        error "Source proto directory '$src_proto' does not exist"
    fi
    info "Copying proto/ to proto-context/ in all services..."
    cp -r "$src_proto" "services/api-gateway/proto-context"
    cp -r "$src_proto" "services/whisper-service/proto-context"
    cp -r "$src_proto" "services/telegram-bot/proto-context"
    success "Proto contexts copied successfully"
}

main() {
    parse_args "$@"
    generate_env
    clean_artifacts "${ARTIFACT_PATHS[@]}"
    copy_proto_contexts
    local orchestrator
    orchestrator=$(select_orchestrator)
    build_project "$orchestrator"
}

main "$@"