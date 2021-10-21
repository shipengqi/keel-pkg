#!/bin/bash
set -e

SYNC_HOME=${HOME}/sync
mkdir -p ${SYNC_HOME}
cp synctl ${SYNC_HOME}
cp image_set.json ${SYNC_HOME}

docker run --rm -tid --name syncdb shipengqi/google_containers_sync_db top
docker cp syncdb:/sync.bolt.db ${SYNC_HOME}
docker kill syncdb

ls -lh  ${SYNC_HOME}
