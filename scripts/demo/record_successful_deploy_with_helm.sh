#!/bin/bash

demoscriptsdir="$(realpath $(dirname -- "$BASH_SOURCE"))"

source scripts/demo/env.sh
./scripts/demo/prepare_sample_app.sh

werf build --repo k3d-registry.sample-app.test:5000/sample-app --save-build-report --build-report-path=/tmp/images.json
kubectl delete ns --ignore-not-found sample-app

asciinema rec --title "successful deploy with helm" --command "cd ${WERF_DIR} ; WERF_REQUIRE_BUILT_IMAGES=1 INPUT_THROTTLING=11 APP_IMAGE=$(cat /tmp/images.json | jq -r .Images.app.DockerImageName) ${demoscriptsdir}/_type_and_execute.sh ${demoscriptsdir}/raw_scenario/successful_deploy_with_helm-2.sh" assets/demo/successful_deploy_with_helm.cast --overwrite
