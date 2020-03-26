SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

PROJECT_ROOT_DIR="$(dirname "$SCRIPT_DIR")"

go run $PROJECT_ROOT_DIR/cmd/client/client.go
