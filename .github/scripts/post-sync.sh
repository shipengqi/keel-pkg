#!/bin/bash
set -e

docker login -u ${DOCKERHUB_USER} -p ${DOCKERHUB_PASS}

SYNC_HOME=${HOME}/sync
mkdir -p ${SYNC_HOME}/build
cp ${SYNC_HOME}/sync.bolt.db ${SYNC_HOME}/build/
cp ${SYNC_HOME}/image_set.json ${SYNC_HOME}/build/

cd ${SYNC_HOME}/build
cat>Dockerfile<<EOF
FROM alpine:3.14
COPY sync.bolt.db /
COPY image_set.json /
EOF
ls -lh

docker build -t ${DOCKERHUB_USER}/${SYNC_DB_REGISTRY} .
docker push ${DOCKERHUB_USER}/${SYNC_DB_REGISTRY}
