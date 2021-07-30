#!/usr/bin/env bash

set -euxo pipefail
source "secret/secrets"

terraform_dir="deploy/aws/terraform"
terraform fmt $terraform_dir
terraform -chdir=$terraform_dir init
terraform -chdir=$terraform_dir apply -auto-approve

export $(scripts/aws_vm_ip.sh)
