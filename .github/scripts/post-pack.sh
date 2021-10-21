#!/bin/bash
set -e

PACK_HOME=${HOME}/pack
cd ${PACK_HOME}
ls -lh ${PACK_HOME}

sudo chmod 755 kubectl kubelet
mkdir -p ${PACK_HOME}/src/kubernetes/bin/
mv kubectl ${PACK_HOME}/src/kubernetes/bin/
mv kubelet ${PACK_HOME}/src/kubernetes/bin/
ls -lh ${PACK_HOME}/src/runtime/containerd/bin/

mkdir -p ${PACK_HOME}/src/runtime/containerd/bin/
tar -xf containerd-*.tar.gz -C ${PACK_HOME}/src/runtime/containerd/
tar -xf crictl-*.tar.gz -C ${PACK_HOME}/src/runtime/containerd/bin/
rm -rf containerd-*.tar.gz
rm -rf crictl-*.tar.gz
sudo chmod 755 runc.${PACK_ARCH}
mv runc.${PACK_ARCH} ${PACK_HOME}/src/runtime/containerd/bin/
ls -lh ${PACK_HOME}/src/runtime/containerd/bin/

mkdir -p ${PACK_HOME}/src/images
cp ./images/* ${PACK_HOME}/src/images
rm -rf ./images
ls -lh ${PACK_HOME}/src/images

echo ${KUBERNETES_VERSION} ${PACK_ARCH}
cd ${PACK_HOME}/src

TAR_NAME=kube-${KUBERNETES_VERSION}-${PACK_ARCH}.tar.gz

if [ -n "${BETA_VERSION}" ];then
    TAR_NAME=kube-${KUBERNETES_VERSION}-${PACK_ARCH}-${BETA_VERSION}.tar.gz
fi

echo "Packing ${TAR_NAME}"
tar -czvf ${PACK_HOME}/${TAR_NAME} .
cd ${PACK_HOME}
rm -rf ./src
echo "Pack ${TAR_NAME} done!"

./packer push -k ${PACK_HOME} -s ${PACK_HOME} --pkg-uri ${PACK_HOME}/${TAR_NAME}

