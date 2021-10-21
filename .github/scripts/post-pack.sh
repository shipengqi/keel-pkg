#!/bin/bash
set -e

PACK_HOME=${HOME}/pack
cd ${PACK_HOME}

chmod 755 kubectl kubelet
mkdir -p ${PACK_HOME}/src/kubernetes/bin/
cp kubectl ${PACK_HOME}/src/kubernetes/bin/
cp kubelet ${PACK_HOME}/src/kubernetes/bin/
ls -lh ${PACK_HOME}/src/runtime/containerd/bin/

mkdir -p ${PACK_HOME}/src/runtime/containerd/bin/
tar -xf containerd-*.tar.gz -C ${PACK_HOME}/src/runtime/containerd/
tar -xf crictl-*.tar.gz -C ${PACK_HOME}/src/runtime/containerd/bin/
ls -lh ${PACK_HOME}/src/runtime/containerd/bin/

mkdir -p ${PACK_HOME}/src/images
cp ./images/* ${PACK_HOME}/src/images
ls -lh ${PACK_HOME}/src/images

echo ${KUBERNETES_VERSION} ${ARCH}
cd ${PACK_HOME}/src

TAR_NAME=kube-${KUBERNETES_VERSION}-${ARCH}.tar.gz

if [ "${BETA}" -eq "true" ];then
    TAR_NAME=kube-${KUBERNETES_VERSION}-${ARCH}-beta.tar.gz
fi

echo "Packing ${TAR_NAME}"
tar -czvf ${PACK_HOME}/kube-${KUBERNETES_VERSION}-${ARCH}.tar.gz .
echo "Pack ${TAR_NAME} done!"

./packer push -k ${PACK_HOME} -s ${PACK_HOME} --pkg-uri ${PACK_HOME}/kube-${KUBERNETES_VERSION}-${ARCH}.tar.gz

