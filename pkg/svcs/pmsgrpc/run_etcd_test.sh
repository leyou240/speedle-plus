#!/bin/bash

shell_dir=$(dirname $0)
rm -rf ./speedle.etcd

set -ex
source ${GOPATH}/src/github.com/leyou240/speedle-plus/setTestEnv.sh

go clean -testcache

startPMS etcd --config-file ${shell_dir}/../pmsrest/config_etcd.json

go test ${TEST_OPTS} github.com/leyou240/speedle-plus/pkg/svcs/pmsgrpc -run=TestMats
rm -rf ./speedle.etcd
