#!/bin/bash

shell_dir=$(dirname $0)

set -ex
#source ${shell_dir}/start_etcd.sh
source ${GOPATH}/src/github.com/leyou240/speedle-plus/setTestEnv.sh

go clean -testcache

startPMS etcd --config-file ${shell_dir}/config_etcd.json

go test ${TEST_OPTS} github.com/leyou240/speedle-plus/pkg/svcs/pmsrest $*
