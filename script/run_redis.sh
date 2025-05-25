#!/bin/bash
PORT=6379
IMAGE=redis
TAG=7.2.5
SELECT=""
TYPE="redis"
CONTAINER_NAME="redis"

function get_parameter() {
    while getopts "t:" opt "$@"; do
        case "$opt" in
            t) TYPE="$OPTARG" ;;
            *) echo "Usage: "; exit 1 ;;
        esac
    done
}

function get_image() {
    if [ "${TYPE}" = "valkey" ] || [ "${TYPE}" = "v" ]; then
        IMAGE="valkey/valkey"
        TAG="8-alpine3.21"
        CONTAINER_NAME="valkey"
    fi
    echo "${IMAGE}:${TAG}"
}

function docker_run() {
    image="$1"
    port="$2"
    container_name="$3"
    container_id=$(docker ps | grep redis | awk '{print $1}')
    
    if [ "${container_id}" = "" ]; then
        docker run -d --name ${container_name} -p ${port}:${port} ${image}
    else
        echo "Alreay container is running: ${container_name}"
        docker ps -a | grep "${container_name}"
    fi
}

function main() {
    get_parameter "$@"
    image=$(get_image)
    docker_run "${image}" "${PORT}" "${CONTAINER_NAME}"
}

main "$@"