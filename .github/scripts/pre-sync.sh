#!/bin/bash
set -e

mkdir -p ${SYNC_HOME}
cd ${SYNC_HOME}

cp synctl ${SYNC_HOME}
cp image_set.json ${SYNC_HOME}
ls -lh

docker run --rm -tid --name syncdb shipengqi/google_containers_sync_db top
docker cp syncdb:/sync.bolt.db ${SYNC_HOME}

docker kill syncdb
