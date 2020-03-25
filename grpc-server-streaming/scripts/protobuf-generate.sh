#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

PROJECT_ROOT_DIR="$(dirname "$SCRIPT_DIR")"

mkdir -p $PROJECT_ROOT_DIR/pkg/generated

protoc -I=$PROJECT_ROOT_DIR --go_out=plugins=grpc:$PROJECT_ROOT_DIR/pkg/generated $PROJECT_ROOT_DIR/api/calculator/*.proto