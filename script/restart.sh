#!/bin/bash
WORK=$(pwd)
TARGET_DIR=mcontext
LOG_DIR="${WORK}/"data/logs
EXE="${WORK}"/mcontext-main

if [[ ! "$WORK" == *"/$TARGET_DIR"* ]]; then
    echo "Error: must be executed in $TARGET_DIR"
    exit 1
fi

PID=$(pgrep -f "${EXE}")

if [ -n "${PID}" ]; then
    kill "${PID}"
fi

mkdir -p "${LOG_DIR}"

if [[ -f "${EXE}" ]]; then
    rm "${EXE}"
fi

go build -o "${EXE}" cmd/main.go

nohup "${EXE}" > "${LOG_DIR}/run.log.$(date +%s)" 2>&1 &
