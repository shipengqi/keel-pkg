#!/bin/sh

Token=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
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
# images
https://hub.docker.com/r/redhat/ubi8

###Add more to make etcd image
tag=$1

etcd_version=3.4.16
if [ "X" = "X$tag" ]
then
  echo "please provide one tag for etcd image"
fi

mkdir -p docker/rootfs
rm -f etcd*
rm -f docker/etcd*

curl -o etcd.tar.gz -L https://svsartifactory.swinfra.net/artifactory/itom-mvn-cdf-group/com/microfocus/cdf/etcd/${etcd_version}/etcd-${etcd_version}-linux64.tar.gz
tar -xzf etcd.tar.gz
rm -f etcd.zip

chmod +x etcd*
cp etc* docker/

cd docker
docker build -t $tag .
