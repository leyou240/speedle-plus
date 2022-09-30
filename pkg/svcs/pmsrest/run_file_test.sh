#!/bin/bash

shell_dir=$(dirname $0)

set -ex
source ${GOPATH}/src/github.com/leyou240/speedle-plus/setTestEnv.sh

startPMS file --config-file ${shell_dir}/config_file.json

go test ${TEST_OPTS} github.com/leyou240/speedle-plus/pkg/svcs/pmsrest $*
