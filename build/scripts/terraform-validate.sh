#!/bin/bash
set -e

fail()
{
  >&2 echo "$*"
  exit 1
}

[ -z "${TFPATH}"  ] && fail "TFPATH environment variable is not set."

cd "$TFPATH"

terraform init -backend=false
terraform validate
