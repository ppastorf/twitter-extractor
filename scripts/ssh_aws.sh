#!/usr/bin/env bash

set -euo pipefail

aws_credentials="secret/aws_access"
terraform_dir="deploy/aws/terraform"

source $aws_credentials

export AWS_VM_01_PUBLIC_IP=$(
    terraform -chdir=$terraform_dir output \
    | awk '/vm_01_public_ip/{print $3}' \
    | tr -d \"
)

ssh -ti ${AWS_PRIVATE_KEY_FILE} ubuntu@${AWS_VM_01_PUBLIC_IP} "sudo su"
