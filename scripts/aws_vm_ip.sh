#!/usr/bin/env bash

terraform_dir="deploy/aws/terraform"

echo AWS_VM_01_PUBLIC_IP=$(
    terraform -chdir=$terraform_dir output \
        | awk '/vm_01_public_ip/{print $3}' | tr -d \"
)
