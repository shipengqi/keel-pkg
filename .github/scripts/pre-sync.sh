#!/bin/bash
set -e

mkdir -p /var/run/keel/sync/
cd /var/run/keel/sync/

cp synctl /var/run/keel/sync/
ls -lh

docker run --rm -tid --name syncdb shipengqi/google_containers_sync_db top
docker cp syncdb:/sync.bolt.db /var/run/keel/sync/

docker kill syncdb
