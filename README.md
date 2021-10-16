# keel-pkg
keel offline package


cd src
tar -czvf ../kube-1.22.0.tar.gz ./

## Image Sync Tools
- https://github.com/mritd/imgsync
- https://github.com/zhangguanzhang/google_containers
- https://github.com/AliyunContainerService/image-syncer
- https://github.com/tkestack/image-transfer
- https://github.com/containers/skopeo
  - [skopeo-sync](https://github.com/containers/skopeo/blob/main/docs/skopeo-sync.1.md)
- https://github.com/containers/image

## Registry
- [Google GCR](https://console.cloud.google.com/gcr/images/google-containers)
- [AliCloud ACR](https://cr.console.aliyun.com/cn-hangzhou/instances/images)
  - [ACR doc](https://help.aliyun.com/document_detail/257112.html?spm=5176.166170.J_5253785160.5.286851646Ug5KU)
- [Huawei Swr](https://console-intl.huaweicloud.com/swr/?agencyId=1e02890d062a42f9be14b82feaa5b711&region=cn-east-3&locale=zh-cn#/app/swr/huaweiOfficialList)
  - [Swr doc](https://support.huaweicloud.com/intl/zh-cn/productdesc-swr/swr_03_0001.html)

## Reference
- https://blog.csdn.net/networken/article/details/84571373
- [github client](https://github.com/google/go-github)

## Build
```bash
CGO_ENABLED=0 go build -tags="containers_image_openpgp" -mod=mod -o synctl ./app/synchronizer/
```

## Know issues
- `go build -tags="containers_image_openpgp"` https://github.com/containers/image/issues/268
