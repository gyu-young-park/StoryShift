CUR_PATH="$(realpath "$0")"
CUR_DIR="$(dirname "$CUR_PATH")"
TEST_DIR="${CUR_DIR}/$(date '+%Y%m%d%H%M%S')_TEST"
TEST_LOG_FILE="test.log"
QUERT_SCEANARIO_FILE="${CUR_DIR}/query_scenario.txt"
# Get port dynamically
PORT="9596"

START_API_SERVER_SCRIPT_PATH=${CUR_DIR}/start.sh

function get_parameter() {
    while getopts "d:o:s:p" opt "$@"; do
        case "$opt" in
            d) TEST_DIR="$OPTARG" ;;
            s) QUERT_SCEANARIO_FILE="$OPTARG" ;;
            p) PORT=="$OPTARG" ;;
            *) echo "Usage: $0 -d <directory name of test log>"; exit 1 ;;
        esac
    done
    
    echo "Test Log Directory: ${TEST_LOG_FILE}"
}

function get_pid_with_port() {
    port=":$1"
    pid=$(lsof -i "$port" | awk 'NR==2 {print $2}')
    echo "${pid}" 
}

function get_port() {
    pid=$(get_pid_with_port "${PORT}")
    if [ ! -z "$pid" ]; then
        echo "Port is already used" >&2
        kill -9 $$
    fi 
    echo "${PORT}"
}

function get_query_scenario_file() {
    if [ ! -f "${QUERT_SCEANARIO_FILE}" ]; then
        echo "QUERT_SCEANARIO_FILE[ ${QUERT_SCEANARIO_FILE} ] is not found, exit 1" >&2
        kill -9 $$
    fi
    echo "${QUERT_SCEANARIO_FILE}"
}

function get_test_dir() {
    mkdir -p "${TEST_DIR}"
    echo "${TEST_DIR}"
}

function get_test_log_file {
    test_dir="$1"
    test_log_file="$TEST_LOG_FILE"

    test_log_file_path="${test_dir}/${test_log_file}"
    touch "${test_log_file_path}"
    echo "${test_log_file_path}"
}

function start_script() {
    log_file="$1"
    bash "${START_API_SERVER_SCRIPT_PATH}" &> "${log_file}" &
    # TODO: readiness check and retry
}

function run_test_query() {
    query_sleep_term=5
    scenario_file="$1"
    test_dir="$2"

    echo "SCENARIO FILE: ${scenario_file}"

    while IFS= read -r line; do
        no=$(echo "$line" | cut -d'|' -f1)
        name=$(echo "$line" | cut -d'|' -f2)
        cmd=$(echo "$line" | cut -d'|' -f3- | sed "s|--output \([^ ]*\)|--output ${test_dir}/\1|")
        echo "[$no] $name: $cmd"
        eval "$cmd"
        sleep ${query_sleep_term}
    done < "${scenario_file}"
}

function unzip_results() {
    # TODO unzip all results to see all of them at once
    echo "unzip result"
}

function kill_server() {
    port="$1"
    pid=$(get_pid_with_port ${port})
    if [ ! -z "$pid" ]; then
        kill -9 $pid || echo "Kill failed"
    else
        echo "Server is Not Running Port: ${port}"
    fi 
}

function main() {
    echo "start query test"
    get_parameter "$@"

    # Get All Data from global variable
    port=$(get_port)
    query_scenario=$(get_query_scenario_file)
    test_dir=$(get_test_dir)
    test_log_file=$(get_test_log_file "$test_dir")
    echo ${test_log_file}
    
    # Start server
    sleep 10
    start_script "${test_log_file}"

    # run scenario all test scripts
    run_test_query "${query_scenario}" "${test_dir}"

    # TearDown
    kill_server_wait_seconds=10
    echo "Sleep ${kill_server_wait_seconds} seconds..."
    sleep $kill_server_wait_seconds

    kill_server "${port}"
    echo "query test done!"
}

main "$@"
