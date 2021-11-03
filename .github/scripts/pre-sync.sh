#!/bin/bash
set -e

SYNC_HOME=${HOME}/sync
mkdir -p ${SYNC_HOME}
cp image_set.json ${SYNC_HOME}
mv synctl ${SYNC_HOME}

docker run --rm -tid --name syncdb ${DOCKERHUB_USER}/${SYNC_DB_REGISTRY} top
docker cp syncdb:/sync.bolt.db ${SYNC_HOME}
docker kill syncdb

ls -lh  ${SYNC_HOME}
cat ${SYNC_HOME}/image_set.json
