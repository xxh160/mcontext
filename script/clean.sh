#!/bin/bash
WORK=$(pwd)
TARGET_DIR=mcontext
EXE="${WORK}"/mcontext-main

if [[ ! "$WORK" == *"/$TARGET_DIR"* ]]; then
    echo "Error: must be executed in $TARGET_DIR"
    exit 1
fi

rm -f "${EXE}"

