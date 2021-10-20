#!/bin/bash
set -e

docker login -u ${DOCKERHUB_USER} -p ${DOCKERHUB_PASS}

cd $HOME
mkdir -p ${SYNC_HOME}/build
cd ${SYNC_HOME}/build

cp ${SYNC_HOME}/sync.bolt.db .
cp ${SYNC_HOME}/image_set.json .
ls -lh

cat>Dockerfile<<EOF
FROM alpine:3.14
COPY sync.bolt.db /
COPY image_set.json /
EOF

docker build -t shipengqi/google_containers_sync_db .
docker push shipengqi/google_containers_sync_db
