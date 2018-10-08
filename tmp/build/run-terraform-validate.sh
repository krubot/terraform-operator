#!/usr/bin/env bash

set -eo pipefail

mkdir -p ${PWD}/.terraform

# output terraform version
terraform version

terraform init -upgrade=true -backend=false
terraform validate
