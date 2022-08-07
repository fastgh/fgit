#!/bin/bash

set -e
#set -x

PROJECT_DIR=$(cd "$(dirname $0)";pwd)

BUILD_DIR=${PROJECT_DIR}/build

rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

cd ${PROJECT_DIR}

go_build() {
    echo $1 $2
    local _OS=$1
    local _ARCH=$2

    local _GOARCH=$_ARCH
    if [ $_ARCH == 'x86_64' ]; then
      _GOARCH='amd64'
    fi

    local _EXT=$_OS
    if [ $_OS == 'windows' ]; then
      _EXT='exe'
    else
      _EXT=${_OS}.${_ARCH}
    fi

    GOOS=${_OS} GOARCH=${_GOARCH} go build -o ${BUILD_DIR}/fgit.${_EXT}

    #if [ $_OS == $(echo `uname -s` | tr A-Z a-z) ]; then
    #  sudo cp ${BUILD_DIR}/fgit.${_EXT} /usr/local/bin/fgit
    #  sudo chmod +x /usr/local/bin/fgit
    #fi
}

go_build linux x86_64
go_build darwin x86_64
go_build darwin arm64
go_build windows x86_64

