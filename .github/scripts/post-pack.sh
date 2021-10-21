#!/bin/bash
set -e

PACK_HOME=${HOME}/pack
cd ${PACK_HOME}
ls -lh ${PACK_HOME}

echo "Coping kubectl kubelet"
sudo chmod 755 kubectl kubelet
mkdir -p ${PACK_HOME}/src/kubernetes/bin/
mv kubectl ${PACK_HOME}/src/kubernetes/bin/
mv kubelet ${PACK_HOME}/src/kubernetes/bin/
ls -lh ${PACK_HOME}/src/kubernetes/bin/

echo "Coping runtime"
mkdir -p ${PACK_HOME}/src/runtime/containerd/bin/
tar -xf containerd-*${PACK_ARCH}.tar.gz -C ${PACK_HOME}/src/runtime/containerd/
tar -xf crictl-*${PACK_ARCH}.tar.gz -C ${PACK_HOME}/src/runtime/containerd/bin/
rm -rf containerd-*${PACK_ARCH}.tar.gz
rm -rf crictl-*${PACK_ARCH}.tar.gz
sudo chmod 755 runc.${PACK_ARCH}
mv runc.${PACK_ARCH} ${PACK_HOME}/src/runtime/containerd/bin/runc
ls -lh ${PACK_HOME}/src/runtime/containerd/bin/

echo "Coping images"
mkdir -p ${PACK_HOME}/src/images
cp ./images/* ${PACK_HOME}/src/images
rm -rf ./images
ls -lh ${PACK_HOME}/src/images

TAR_NAME=kube-${KUBERNETES_VERSION}-${PACK_ARCH}.tar.gz

if [ -n "${BETA_VERSION}" ];then
    TAR_NAME=kube-${KUBERNETES_VERSION}-${PACK_ARCH}-${BETA_VERSION}.tar.gz
fi

echo "Packing ${TAR_NAME}"
cd ${PACK_HOME}/src
tar -czvf ${PACK_HOME}/${TAR_NAME} .
cd ${PACK_HOME}
rm -rf ./src
echo "Pack ${TAR_NAME} done!"

echo "Pushing ${TAR_NAME} ..."
./packer push -k ${QINIU_ACCESS_KEY} -s ${QINIU_SECRET_KEY} --pkg-uri ${PACK_HOME}/${TAR_NAME}
echo "Push ${TAR_NAME} done!"
