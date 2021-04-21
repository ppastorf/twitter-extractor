#!/usr/bin/env bash

set -euxo pipefail
source "secret/aws_access"

terraform_dir="deploy/aws/terraform"

terraform -chdir=$terraform_dir destroy -auto-approve
