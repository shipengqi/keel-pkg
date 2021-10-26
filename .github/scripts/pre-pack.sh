#!/bin/bash
set -e

PACK_HOME=${HOME}/pack
mkdir -p ${PACK_HOME}
mkdir -p ${PACK_HOME}/src
sudo chmod 644 versions.json
cp versions.json ${PACK_HOME}
cp -rf package/* ${PACK_HOME}/src
mv packer ${PACK_HOME}

ls -lh ${PACK_HOME}

# workaround https://github.com/actions/virtual-environments/issues/2875
sudo rm -rf /usr/share/dotnet
# sudo rm -rf /opt/ghc
# sudo rm -rf "/usr/local/share/boost"
# sudo rm -rf "$AGENT_TOOLSDIRECTORY"

sudo swapoff -a
sudo rm -f /swapfile
sudo apt clean
df -h
