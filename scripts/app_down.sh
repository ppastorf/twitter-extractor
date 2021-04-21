#!/usr/bin/env bash

set -euxo pipefail
source "secret/aws_access"

terraform_dir="deploy/aws/terraform"
inventory="deploy/hosts.yaml"
playbooks="deploy/playbooks"

export AWS_VM_01_PUBLIC_IP=$(
    terraform -chdir=$terraform_dir output \
    | awk '/vm_01_public_ip/{print $3}' \
    | tr -d \"
)

ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i "$inventory" "$playbooks/compose_down.yaml"
