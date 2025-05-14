#!/bin/bash

PROJECT_ROOT_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")" && cd .. && pwd)"
CMD_DIR="${PROJECT_ROOT_PATH}/cmd"
CMD_MAIN_FILE="${CMD_DIR}/main.go"

export STORY_SHIFT_CONFIG_FILE="${PROJECT_ROOT_PATH}/config/test_config.yaml"

echo "CMD File: ${CMD_MAIN_FILE}"

start_cmd_main() {
    go run $CMD_MAIN_FILE
}

start_cmd_main