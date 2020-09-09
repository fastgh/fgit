#!/bin/sh

set -e
#set -x

PROJECT_DIR=$(cd "$(dirname $0)";pwd)

BUILD_DIR=${PROJECT_DIR}/build

rm -rf ${BUILD_DIR}
cd ${PROJECT_DIR}

go_build() {
    local _OS=$1
    local _EXT=$_OS
    if [ $_OS == 'windows' ]; then
      _EXT='exe'
    fi

    mkdir -p ${BUILD_DIR}
    GOOS=${_OS} GOARCH=amd64 go build -o ${BUILD_DIR}/fgit.${_EXT}

    if [ $_OS == $(echo `uname -s` | tr A-Z a-z) ]; then
      sudo cp ${BUILD_DIR}/fgit.${_EXT} /usr/local/bin/fgit
      sudo chmod +x /usr/local/bin/fgit
    fi
}

go_build linux
go_build darwin
go_build windows

