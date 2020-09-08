#!/bin/sh

set -e

PROJECT_DIR=$(cd "$(dirname $0)";pwd)

INSTRUMENT_SRC=git_instrument.go
INSTRUMENT_SRC_PRIVATE=${PROJECT_DIR}/../${INSTRUMENT_SRC}
INSTRUMENT_SRC_PUB=${PROJECT_DIR}/${INSTRUMENT_SRC}
INSTRUMENT_SRC_PUB_BACKUP=${INSTRUMENT_SRC_PUB}.backup

TARGET_DIR=${PROJECT_DIR}/target

rm -rf ${TARGET_DIR}
cd ${PROJECT_DIR}

mv ${INSTRUMENT_SRC_PUB} ${INSTRUMENT_SRC_PUB_BACKUP}
if [ -f ${INSTRUMENT_SRC_PRIVATE} ];then
  cp ${INSTRUMENT_SRC_PRIVATE} ${INSTRUMENT_SRC_PUB}
fi

go_build() {
    local _OS=$1
    local _PREFIX=$2
    local _OS_TARGET_DIR=${TARGET_DIR}/${_OS}

    mkdir -p ${_OS_TARGET_DIR}
    GOOS=${_OS} GOARCH=amd64 go build -o ${_OS_TARGET_DIR}/fgit${_PREFIX}
}

go_build linux .linux
go_build darwin .darwin
go_build windows .exe

rm ${INSTRUMENT_SRC_PUB}
mv ${INSTRUMENT_SRC_PUB_BACKUP} ${INSTRUMENT_SRC_PUB}
