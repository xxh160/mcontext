#!/bin/bash

WORK=$(pwd)
TARGET_DIR=mcontext

if [[ ! "$WORK" == *"/$TARGET_DIR"* ]]; then
    echo "Must be executed in $TARGET_DIR"
    exit 1
fi

rm "${WORK}"/go.mod
go mod init "${TARGET_DIR}" 
go mod tidy

