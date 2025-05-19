START_API_SERVER_SCRIPT_PATH=./start.sh
TEST_LOG_DIR="./test_log"
TEST_LOG_FILE="$(date '+%Y%m%d%H%M%S')_QUERY_TEST.log"

function get_parameter() {
    while getopts "d:" opt "$@"; do
        case "$opt" in
            d) TEST_LOG_DIR="$OPTARG" ;;
            *) echo "Usage: $0 -d <directory of test log>"; exit 1 ;;
        esac
    done

    echo "Test Log Directory: ${TEST_LOG_DIR}"
}

function make_log_file {
    dir="$1"
    if [ -d "$dir" ]; then
        return 0
    fi
    
    mkdir "${dir}"
    

    echo "${TEST_LOG_DIR}/${TEST_LOG_FILE}"
}

get_parameter "$@"
TEST_LOG_FILE=$(make_log_file "$TEST_LOG_DIR")
echo $TEST_LOG_FILE

# $START_API_SERVER_SCRIPT_PATH 2>&1 > ./test_log.txt
