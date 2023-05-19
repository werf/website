#!/bin/bash

demoscriptsdir="$(realpath $(dirname -- "$BASH_SOURCE"))"

source scripts/demo/env.sh
./scripts/demo/prepare_sample_app.sh

werf build --repo k3d-registry.sample-app.test:5000/sample-app --save-build-report --build-report-path=/tmp/images.json
kubectl delete ns --ignore-not-found sample-app

asciinema rec --title "successful deploy with werf" --command "cd ${WERF_DIR} ; WERF_REQUIRE_BUILT_IMAGES=1 INPUT_THROTTLING=11 ${demoscriptsdir}/_type_and_execute.sh ${demoscriptsdir}/raw_scenario/successful_deploy_with_werf.sh" assets/demo/successful_deploy_with_werf.cast --overwrite
