#!/bin/bash

source scripts/demo/env.sh
./scripts/demo/prepare_sample_app.sh

kubectl delete ns --ignore-not-found sample-app sample-app
APP_IMAGE=$(cat /tmp/images.json | jq -r .Images.app.DockerImageName)
helm upgrade --wait --install --create-namespace --namespace sample-app sample-app .helm/ --set werf.image.app=${APP_IMAGE} --set host=sample-app.test --timeout 60s

./scripts/demo/raw_scenario/introduce_error.sh
werf build --repo k3d-registry.sample-app.test:5000/sample-app --save-build-report --build-report-path=/tmp/images.json

asciinema rec --title "failed deploy with helm" --command "cd ${WERF_DIR} ; WERF_REQUIRE_BUILT_IMAGES=1 INPUT_THROTTLING=11 APP_IMAGE=$(cat /tmp/images.json | jq -r .Images.app.DockerImageName) ${demoscriptsdir}/_type_and_execute.sh ${demoscriptsdir}/raw_scenario/failed_deploy_with_helm-2.sh" assets/demo/failed_deploy_with_helm.cast --overwrite
