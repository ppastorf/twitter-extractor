#!/usr/bin/env bash

set -euxo pipefail

environment="$1"
playbook="$2"
source "secret/secrets"

if [[ $environment == "aws" ]]; then
    export $(scripts/aws_vm_ip.sh)
fi

ansible-playbook \
    -i "deploy/${environment}/hosts.yaml" \
    "deploy/playbooks/${playbook}.yaml"
