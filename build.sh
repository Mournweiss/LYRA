#!/usr/bin/env bash

set -euo pipefail

# ANSI color codes
COLOR_INFO="\033[0m"       # White
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
FOREGROUND_MODE=false

ARTIFACT_PATHS=(
    "services/api-gateway/proto-context"
    "services/whisper-service/proto-context"
    "services/telegram-bot/proto-context"
)

show_help() {
    cat << EOF
LYRA Build Script

Usage: $0 [OPTIONS]

Options:
    -p, --podman                Use podman-compose as orchestrator
    -d, --docker                Use docker-compose as orchestrator
    -t, --telegram-token TOKEN  Set Telegram bot token in .env file
    -f, --foreground            Run containers in foreground mode (not detached)
    -h, --help                  Show this help message

EOF
}

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
            --foreground|-f)
                FOREGROUND_MODE=true
                shift
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                warn "Unknown argument: $1"
                shift
                ;;
        esac
    done
}

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

build_project() {
    local compose_cmd="$1"
    local compose_args="-f compose.yml up --build"

    if [ "$FOREGROUND_MODE" = false ]; then
        compose_args="$compose_args -d"
        info "Building and starting project in daemon mode using $compose_cmd..."
    else
        info "Building and starting project in foreground mode using $compose_cmd..."
    fi

    $compose_cmd $compose_args
    success "Build and startup completed"
}

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