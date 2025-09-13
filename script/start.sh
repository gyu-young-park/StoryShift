#!/bin/bash

PROJECT_ROOT_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")" && cd .. && pwd)"
CMD_DIR="${PROJECT_ROOT_PATH}/cmd"
CMD_WEB_MAIN_FILE="${CMD_DIR}/web/main.go"
CMD_CLI_MAIN_FILE="${CMD_DIR}/cli/main.go"
CMD=""

export STORY_SHIFT_CONFIG_FILE="${PROJECT_ROOT_PATH}/config/test_config.yaml"

function parse_opt() {
    while [[ $# -gt 0 ]]; do
    case "$1" in
      -c)
        MODE="$2"
        shift 2
        ;;
      *)
        echo "Usage: $0 -c [web|cli]"
        exit 1
        ;;
    esac
  done

  case "$MODE" in
    web) CMD=$CMD_WEB_MAIN_FILE ;;
    cli) CMD=$CMD_CLI_MAIN_FILE ;;
    *)   echo "Invalid or missing mode"; exit 1 ;;
  esac
}

function start_cmd_main() {
    parse_opt "$@"

    echo "CMD File: ${CMD}"
    go run $CMD
}

start_cmd_main "$@"