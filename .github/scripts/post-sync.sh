#!/bin/bash
set -e

docker login -u shipengqi -p ${DOCKER_PASS}

cd $HOME
mkdir -p /var/run/keel/sync/build
cd /var/run/keel/sync/build

cp /var/run/keel/sync/sync.bolt.db .
ls -lh

cat>Dockerfile<<EOF
FROM alpine:3.14
COPY sync.bolt.db /
EOF

docker build -t shipengqi/google_containers_sync_db .
docker push shipengqi/google_containers_sync_db
