#!/bin/bash

SCRIPT_PATH="$(realpath "$0")"
PROJECT_PATH="$(dirname "$(dirname ${SCRIPT_PATH})")"
DOCKER_PATH=${PROJECT_PATH}

IMAGE_NAME=story-shift
IMAGE_TAG=latest
NAMESPACE=""
REGISTRY=""
DOCKER_FILE="${DOCKER_PATH}/Dockerfile"

function get_parameter() {
    while getopts "i:t:n:f:r" opt "$@"; do
        case "$opt" in
            i) IMAGE_NAME="$OPTARG" ;;
            t) IMAGE_TAG="$OPTARG" ;;
            n) NAMESPACE="$OPTARG" ;;
            f) DOCKER_FILE="$OPTARG" ;;
            r) REGISTRY="$OPTARG" ;;
            *) echo "Usage: "; exit 1 ;;
        esac
    done
}

function get_container_runtime() {
    container_runtime="docker"
    if command -v nerdctl >/dev/null 2>&1; then
        container_runtime="nerdctl"
    fi
    echo ${container_runtime}
}

function get_namespace() {
    container_runtime=$1
    namespace=${NAMESPACE}
    if [ "$container_runtime" = "nerdctl" ]; then
        namespace=" -n ${namespace} "
    fi
    echo "$namespace"
}

function get_image_name() {
    echo "${IMAGE_NAME}:${IMAGE_TAG}"
}

function docker_build() {
    container_runtime=$1
    build_opt=$2
    image_name=$3

    echo "CONTAINER_RUNTIME: ${container_runtime}"
    echo "PROJECT: ${PROJECT_PATH}"
    echo "DOCKER_FILE: ${DOCKER_FILE}"
    echo "IMAGE_NAME: ${image_name}"
    echo "BUILD_OPT: ${build_opt}"

    pushd ${PROJECT_PATH}
        ${container_runtime} build ${build_opt} -f ${DOCKER_FILE} -t ${image_name} ${PROJECT_PATH}
    popd
}

function main() {
    get_parameter "$@"
    container_runtime=$(get_container_runtime)
    echo "Container Runtime: ${container_runtime}"

    namespace=$(get_namespace $container_runtime)
    build_opt="${namespace}"
    image_name=$(get_image_name)

    docker_build "${container_runtime}" "${build_opt}" "${image_name}"
}

main "$@"

# docker build -f ./Dockerfile -t StoryShift:0.0.1 .
# docker run --name storyshift -e STORY_SHIFT_CONFIG_FILE=./config/test_config.yaml docker.io/library/story-shift:0.0.1