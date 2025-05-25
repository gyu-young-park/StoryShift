#!/bin/bash

SCRIPT_PATH="$(realpath "$0")"
SCRIPT_DIR="$(dirname ${SCRIPT_PATH})"
PROJECT_PATH="$(dirname "${SCRIPT_DIR}")"
DOCKER_PATH=${PROJECT_PATH}
CACHE_SCRIPT="run_redis.sh"

IMAGE_NAME=story-shift
IMAGE_TAG=latest
NAMESPACE=""
REGISTRY=""
DOCKER_FILE="${DOCKER_PATH}/Dockerfile"
RUN_OPT=""
NETWORK="story-shift"

COMMAND=""
BUILD=false
RUN=false

function eval_command() {
    command="$1"
    echo "command: ${command}"

    if [ "${command}" = "build" ]; then
        BUILD=true
    elif [ "${command}" = "run" ]; then
        RUN=true
    else
        BUILD=true
        RUN=true
    fi
}

function get_parameter() {
    while getopts "i:t:n:f:r:c:R:" opt "$@"; do
        case "$opt" in
            i) IMAGE_NAME="$OPTARG" ;;
            t) IMAGE_TAG="$OPTARG" ;;
            n) NAMESPACE="$OPTARG" ;;
            f) DOCKER_FILE="$OPTARG" ;;
            r) REGISTRY="$OPTARG" ;;
            c) COMMAND="$OPTARG" ;;
            R) RUN_OPT="$OPTARG" ;;
            *) echo "Usage: "; exit 1 ;;
        esac
    done

    eval_command "${COMMAND}"
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

    echo "----------------------IMAGE BUILD--------------------"
    echo "CONTAINER_RUNTIME: ${container_runtime}"
    echo "PROJECT: ${PROJECT_PATH}"
    echo "DOCKER_FILE: ${DOCKER_FILE}"
    echo "IMAGE_NAME: ${image_name}"
    echo "BUILD_OPT: ${build_opt}"

    pushd ${PROJECT_PATH}
        ${container_runtime} build ${build_opt} -f ${DOCKER_FILE} -t ${image_name} ${PROJECT_PATH}
    popd
    echo "-----------------------------------------------------"
}

function docker_run() {
    container_runtime=$1
    run_opt=$2
    image_name=$3
    
    echo "---------------------------RUN------------------------"
    echo "CONTAINER_RUNTIME: ${container_runtime}"
    echo "RUN_OPT: ${run_opt}"
    echo "IMAGE_NAME: ${image_name}"
    echo "Network: ${NETWORK}"
    ${container_runtime} run ${run_opt} --network ${NETWORK} --rm ${image_name}
    echo "-----------------------------------------------------"
}

function create_docker_network() {
    network_id="$(docker network ls | grep ${NETWORK} | awk '{print $1}')"
    if [ "${network_id}" = "" ]; then
        ${container_runtime} network create ${NETWORK}
    else
        echo "Alreay network is running"
    fi
}

function docker_run_cache() {
    cache_script="${SCRIPT_DIR}/${CACHE_SCRIPT}"
    echo "Cache script: ${cache_script}"
    bash ${cache_script} -t redis -n "${NETWORK}"
}

function main() {
    get_parameter "$@"
    container_runtime=$(get_container_runtime)
    namespace=$(get_namespace $container_runtime)
    image_name=$(get_image_name)
    build_opt="${namespace}"

    if $BUILD; then
        docker_build "${container_runtime}" "${build_opt}" "${image_name}"
    fi

    if $RUN; then
        create_docker_network
        docker_run_cache
        docker_run "${container_runtime}" "${RUN_OPT}" "${image_name}"
    fi
}

main "$@"