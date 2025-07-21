#!/bin/bash

CONTAINER_RUNTIME=docker
VERSION=9
PASSWORD=hello
PORT=3306

function get_parameter() {
    while getopts "i:p:" opt "$@"; do
        case "$opt" in
            v) VERSION="$OPTARG" ;;
            p) PASSWORD="$OPTARG" ;;
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
    $CONTAINER_RUNTIME run --name mysql -e MYSQL_ROOT_PASSWORD=${password} -d -p ${port}:${port} ${image}
}

function main() {
    get_parameter "$@"

    image="mysql:${VERSION}"
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

main "$@"