#!/bin/bash

CONTAINER_RUNTIME=docker
TAG=9
PASSWORD=hello
PORT=3306
COMMAND=start # restart, delete, start
CONTAINER_NAME=mysql

function get_parameter() {
    while getopts "i:p:c:" opt "$@"; do
        case "$opt" in
            v) TAG="$OPTARG" ;;
            p) PASSWORD="$OPTARG" ;;
            c) COMMAND="$OPTARG" ;;
            *) echo "Usage: "; exit 1 ;;
        esac
    done
}

function check_is_image_exist() {
    ret=$($CONTAINER_RUNTIME images | grep mysql | awk '{print $1":"$2}')
    if [ -z "${ret}" ]; then
        return 1
    fi
    return 0
}

function pull_mysql_image() {
    image=$1
    echo "---------------pull docker image: ${image}---------------"
    $CONTAINER_RUNTIME pull ${image}
    echo "------------------------Done-----------------------------"
}

function check_is_container_run() {
    image=$1
    ret=$($CONTAINER_RUNTIME ps -a | grep ${image})
    if [ -z "${ret}" ]; then
        return 1
    fi
    return 0
}

function run_mysql_container() {
    image=$1
    password=$2
    port=$3
    echo "CMD: $CONTAINER_RUNTIME run --name ${CONTAINER_NAME} -e MYSQL_ROOT_PASSWORD=${password} -d -p ${port}:${port} ${image}"
    $CONTAINER_RUNTIME run --name ${CONTAINER_NAME} -e MYSQL_ROOT_PASSWORD=${password} -d -p ${port}:${port} ${image}
}

function run_command() {
    image=$1
    if [[ "${COMMAND}" == "restart" ]]; then
        echo "---------restart-------------"
        run_restart "${image}"
    elif [[ "${COMMAND}" == "delete" ]]; then
        echo "---------delete-------------"
        run_delete "${image}"
    else
        echo "---------start-------------"
        run_start "${image}"
    fi
}

function run_restart() {
    image=$1
    run_delete "${image}" 
    run_start "${image}"
}

function run_delete() {
    image=$1
    if $(check_is_container_run ${image}); then
        $CONTAINER_RUNTIME rm -f ${CONTAINER_NAME}
    else
        echo "mysql container is not running"
    fi
}

function run_start() {
    image=$1
    if $(check_is_image_exist); then
        echo "Already image is exists in local repo"
    else
        pull_mysql_image "${image}"
    fi

    if $(check_is_container_run ${image}); then
        echo "Already container is running"
    else
        run_mysql_container "${image}" "${PASSWORD}" "${PORT}"
    fi
}

function main() {
    get_parameter "$@"
    image="mysql:${TAG}"
    run_command ${image}
}

main "$@"