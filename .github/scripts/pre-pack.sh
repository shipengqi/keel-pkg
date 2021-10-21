#!/bin/bash
set -e

mkdir -p ${PACK_HOME}
mkdir -p ${PACK_HOME}/src
cp packer ${PACK_HOME}
cp versions.json ${PACK_HOME}
cp -rf package/* ${PACK_HOME}/src

cd ${PACK_HOME}
ls -lh
