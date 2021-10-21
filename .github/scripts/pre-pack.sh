#!/bin/bash
set -e

PACK_HOME=${HOME}/pack
mkdir -p ${PACK_HOME}
mkdir -p ${PACK_HOME}/src
cp packer ${PACK_HOME}
cp versions.json ${PACK_HOME}
cp -rf package/* ${PACK_HOME}/src

cd ${PACK_HOME}
ls -lh

# workaround https://github.com/actions/virtual-environments/issues/2875
sudo rm -rf /usr/share/dotnet
sudo rm -rf /opt/ghc
sudo rm -rf "/usr/local/share/boost"
sudo rm -rf "$AGENT_TOOLSDIRECTORY"
