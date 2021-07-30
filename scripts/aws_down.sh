#!/usr/bin/env bash

set -euxo pipefail

source "secret/secrets"

terraform_dir="deploy/aws/terraform"
terraform -chdir=$terraform_dir destroy -auto-approve

export $(scripts/aws_vm_ip.sh)
