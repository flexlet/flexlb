#!/bin/bash

BASE_DIR="$(readlink -f $(dirname $0)/..)"
BUILD_DIR="${BASE_DIR}/build"
RELEASE_DIR="${BASE_DIR}/release"

# inject build options
source ${BUILD_DIR}/profile

# build target
TARGET_DIR="${BUILD_DIR}/target"
PKG_DIR="${TARGET_DIR}/flexlb-${ARCH}-${VERSION}"

# build binary
mkdir -p ${PKG_DIR}/bin
cd ${BASE_DIR}
CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build -o ${PKG_DIR}/bin/flexlb cmd/flexlb-server/main.go

# copy config, service, install script
cp -r ${RELEASE_DIR}/* ${PKG_DIR}/

# create tarball
cd ${PKG_DIR}/
tar -zcf ${TARGET_DIR}/flexlb-${ARCH}-${VERSION}.tar.gz *
echo "${TARGET_DIR}/flexlb-${ARCH}-${VERSION}.tar.gz"
