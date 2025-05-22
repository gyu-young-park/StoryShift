#!/bin/bash

function get_parameter() {
    while getopts "d:o;s;p" opt "$@"; do
        case "$opt" in
            d) TEST_DIR="$OPTARG" ;;
            s) QUERT_SCEANARIO_FILE="$OPTARG" ;;
            p) PORT=="$OPTARG" ;;
            *) echo "Usage: $0 -d <directory name of test log>"; exit 1 ;;
        esac
    done
    
    echo "Test Log Directory: ${TEST_LOG_FILE}"
}

docker build -f ./Dockerfile -t StoryShift:0.0.1 .
docker run --name storyshift -e STORY_SHIFT_CONFIG_FILE=./config/test_config.yaml docker.io/library/story-shift:0.0.1