#!/bin/bash
PORT=6379
IMAGE=redis
TAG=7.2.5
SELECT=""
TYPE="redis"
NETWORK="bridge"

function get_parameter() {
    while getopts "t:n:" opt "$@"; do
        case "$opt" in
            t) TYPE="$OPTARG" ;;
            n) NETWORK="$OPTARG" ;;
            *) echo "Usage: "; exit 1 ;;
        esac
    done
}

function get_image() {
    image="${IMAGE}"
    tag="${TAG}"
    if [ "${TYPE}" = "valkey" ] || [ "${TYPE}" = "v" ]; then
        image="valkey/valkey"
        tag="8-alpine3.21"
    fi
    echo "${IMAGE}:${TAG}"
}

function docker_run() {
    image="$1"
    port="$2"
    container_id=$(docker ps | grep cache | awk '{print $1}')
    
    if [ "${container_id}" = "" ]; then
        docker run -d --name cache -p ${port}:${port} --network ${NETWORK} ${image}
    else
        echo "Alreay container is running: cache"
        docker ps -a | grep "cache"
    fi
}

function main() {
    get_parameter "$@"
    image=$(get_image)
    docker_run "${image}" "${PORT}"
    return 0
}

main "$@"