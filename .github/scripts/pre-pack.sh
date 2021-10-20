#!/bin/bash
set -e

mkdir -p ${PACK_HOME}
cd ${PACK_HOME}

cp packer ${PACK_HOME}
cp image_set.json ${PACK_HOME}
ls -lh



docker kill syncdb
