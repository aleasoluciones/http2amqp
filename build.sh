#!/bin/bash
set -e

SRC_DIR=/go/src/github.com/aleasoluciones/http2amqp
GO_VERSION=1.9.2

# We need a temp dir to mount to the build process 
# because when using this script from docker, the 
# real dir mounted is from the host
WORKDIR=/tmp/$$/$(hostname)
rm -rf $WORKDIR
mkdir -p ${WORKDIR}
cp -a . ${WORKDIR}
docker run --rm -v ${WORKDIR}:${SRC_DIR} -e CGO_ENABLED=0 -e GOOS=linux golang:${GO_VERSION} bash -c "cd ${SRC_DIR};make jenkins"
cp -v ${WORKDIR}/http2amqp .
cp -v ${WORKDIR}/echoservice .
docker build --no-cache -t aleasoluciones/http2amqp .

rm -rf ${WORKDIR}
