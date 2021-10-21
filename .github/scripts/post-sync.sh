#!/bin/bash
set -e

docker login -u ${DOCKERHUB_USER} -p ${DOCKERHUB_PASS}

SYNC_HOME=${HOME}/sync
mkdir -p ${SYNC_HOME}/build
cp sync.bolt.db ${SYNC_HOME}/build/
cp image_set.json ${SYNC_HOME}/build/

cd ${SYNC_HOME}/build
cat>Dockerfile<<EOF
FROM alpine:3.14
COPY sync.bolt.db /
COPY image_set.json /
EOF
ls -lh

docker build -t shipengqi/google_containers_sync_db .
docker push shipengqi/google_containers_sync_db
