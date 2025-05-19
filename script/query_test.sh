CUR_PATH="$(realpath "$0")"
CUR_DIR="$(dirname "$CUR_PATH")"
TEST_LOG_DIR="${CUR_DIR}/test_log"
TEST_LOG_FILE="$(date '+%Y%m%d%H%M%S')_QUERY_TEST.log"

START_API_SERVER_SCRIPT_PATH=${CUR_DIR}/start.sh

function get_parameter() {
    while getopts "d:" opt "$@"; do
        case "$opt" in
            d) TEST_LOG_DIR="$OPTARG" ;;
            *) echo "Usage: $0 -d <directory name of test log>"; exit 1 ;;
        esac
    done

    TEST_LOG_DIR_NAME="$(basename ${TEST_LOG_DIR})"
    TEST_LOG_DIR="${CUR_DIR}/${TEST_LOG_DIR_NAME}"

    echo "Test Log Directory: ${TEST_LOG_DIR}"
}

function make_test_log_file {
    dir="$1"
    if [ ! -d "$dir" ]; then
        mkdir "${dir}"
    fi

    touch "${TEST_LOG_DIR}/${TEST_LOG_FILE}"
    echo "${TEST_LOG_DIR}/${TEST_LOG_FILE}"
}

function start_script() {
    log_file="$1"
    bash "${START_API_SERVER_SCRIPT_PATH}" &> "${log_file}" &
    pid=$!
    echo "$pid"
}

function run_test_query() {
    echo "run1"
}

function kill_server() {
    pid="$1"
    kill -9 $pid || echo "Kill failed"
}

function main() {
    echo "start query test"
    get_parameter "$@"

    TEST_LOG_FILE_PATH=$(make_test_log_file "$TEST_LOG_DIR")
    echo ${TEST_LOG_FILE_PATH}
    
    SERVER_PID=$(start_script "${TEST_LOG_FILE_PATH}")
    echo "Server pid: ${SERVER_PID}"
    
    run_test_query

    kill_server_wait_seconds=10
    echo "Sleep ${kill_server_wait_seconds} seconds..."
    sleep $kill_server_wait_seconds
    kill_server $SERVER_PID

    echo "query test done!"
}

main "$@"
