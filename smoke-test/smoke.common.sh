#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

BOOTLOOSE_TEMPLATE=${BOOTLOOSE_TEMPLATE:-"bootloose.yaml.tpl"}

export LINUX_IMAGE="${LINUX_IMAGE:-"quay.io/k0sproject/bootloose-ubuntu20.04"}"
export PRESERVE_CLUSTER="${PRESERVE_CLUSTER:-""}"
export K0S_VERSION

createCluster() {
  envsubst < "${SCRIPT_DIR}/${BOOTLOOSE_TEMPLATE}" > bootloose.yaml
  go tool github.com/k0sproject/bootloose create
}

deleteCluster() {
  # cleanup any existing cluster
  envsubst < "${SCRIPT_DIR}/${BOOTLOOSE_TEMPLATE}" > bootloose.yaml
  go tool github.com/k0sproject/bootloose delete && docker volume prune -f
}


cleanup() {
    echo "Cleaning up..."

    if [ -z "${PRESERVE_CLUSTER}" ]; then
      deleteCluster
    fi
}