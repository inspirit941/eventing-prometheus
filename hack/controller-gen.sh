#!/usr/bin/env bash

controller-gen object:headerFile=./hack/boilerplate.go.txt paths=./pkg/apis/... +crd:crdVersions=v1 +output:crd:artifacts:config=manifests
# description 삭제 -> crd에 옵션 추가. i.e) +crd.crdVersions=v1,maxDesclen=0