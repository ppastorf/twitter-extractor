#!/usr/bin/env bash

set -euo pipefail
source "secret/secrets"
export $(scripts/aws_vm_ip.sh)

ssh -ti ${AWS_PRIVATE_KEY_FILE} ubuntu@${AWS_VM_01_PUBLIC_IP} "sudo su"
