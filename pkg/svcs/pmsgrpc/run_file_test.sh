#!/bin/bash

shell_dir=$(dirname $0)
temp_policy_file=/tmp/speedle-test-file-store.json

set -ex
source ${GOPATH}/src/github.com/leyou240/speedle-plus/setTestEnv.sh

startPMS file --config-file ${shell_dir}/../pmsrest/config_file.json

go test ${TEST_OPTS} github.com/leyou240/speedle-plus/pkg/svcs/pmsgrpc -run=TestMats
