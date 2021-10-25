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
sudo cp ./images/* ${PACK_HOME}/src/images
sudo rm -rf ./images
sudo chmod 755 ${PACK_HOME}/src/images/*
ls -lh ${PACK_HOME}/src/images

PACK_VERSION=${KUBERNETES_VERSION}-${PACK_ARCH}
if [ -n "${BETA_VERSION}" ];then
    PACK_VERSION=${KUBERNETES_VERSION}-${PACK_ARCH}-${BETA_VERSION}
fi
TAR_NAME=kube-${PACK_VERSION}.tar.gz

echo "Packing ${TAR_NAME}"
cd ${PACK_HOME}/src
sudo tar -czvf ${PACK_HOME}/${TAR_NAME} .
cd ${PACK_HOME}
sudo rm -rf ./src
sudo chmod 755 ${PACK_HOME}/${TAR_NAME}
echo "Pack ${TAR_NAME} done!"

echo "Pushing ${TAR_NAME} to ${PUSH_TO} ..."


if [ "${PUSH_TO}" = "dockerhub" ];then
    docker login -u ${DOCKERHUB_USER} -p ${DOCKERHUB_PASS}
    cat>Dockerfile<<EOF
FROM busybox:1.34.0
COPY ${TAR_NAME} /
COPY versions.json /
EOF
    cat Dockerfile
    docker build -t ${DOCKERHUB_USER}/${PACK_REGISTRY}:${PACK_VERSION} .
    docker push ${DOCKERHUB_USER}/${PACK_REGISTRY}:${PACK_VERSION}
else
    sudo ./packer push -k ${QINIU_ACCESS_KEY} -s ${QINIU_SECRET_KEY} -b ${QINIU_BUCKET} --pkg-uri ${PACK_HOME}/${TAR_NAME}
fi

echo "Push ${TAR_NAME} done!"
