#!/bin/bash

demoscriptsdir="$(realpath $(dirname -- "$BASH_SOURCE"))"

source scripts/demo/env.sh
./scripts/demo/prepare_sample_app.sh

kubectl delete ns --ignore-not-found sample-app
werf converge --repo k3d-registry.sample-app.test:5000/sample-app --save-build-report --build-report-path=/tmp/images.json

./scripts/demo/raw_scenario/introduce_error.sh
werf build --repo k3d-registry.sample-app.test:5000/sample-app --save-build-report --build-report-path=/tmp/images.json

asciinema rec --title "failed deploy with werf" --command "cd ${WERF_DIR} ; WERF_REQUIRE_BUILT_IMAGES=1 INPUT_THROTTLING=11 ${demoscriptsdir}/_type_and_execute.sh ${demoscriptsdir}/raw_scenario/failed_deploy_with_werf.sh" assets/demo/failed_deploy_with_werf.cast --overwrite
