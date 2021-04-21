#!/usr/bin/env bash

set -euo pipefail

image_version=${1:-"latest"}
image_tag="worker-extractor-${image_version}"
registry="registry.gitlab.com"
image_path="icmc-ssc0158-2021/2021/gcloud13"

source "secret/gitlab_registry"

docker login registry.gitlab.com \
    -u $GITLAB_DEPLOY_USER \
    -p $GITLAB_DEPLOY_TOKEN

docker build "app/worker-extractor" -t "${registry}/${image_path}:${image_tag}"

docker push "${registry}/${image_path}:${image_tag}"
