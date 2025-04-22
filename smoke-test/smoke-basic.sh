#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

set -e

. ${SCRIPT_DIR}/smoke.common.sh
trap cleanup EXIT

deleteCluster
createCluster

eval $(ssh-agent) 
ssh-add id_rsa_k0s

echo "* Starting test"
go test -v -cover ./...
echo "* Test OK"
echo "* Done"
