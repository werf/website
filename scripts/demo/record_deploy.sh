#!/bin/bash

source scripts/demo/env.sh
./scripts/demo/prepare_sample_app.sh

werf build --repo k3d-registry.sample-app.test:5000/sample-app --save-build-report --build-report-path=/tmp/images.json
kubectl delete ns --ignore-not-found sample-app

asciinema rec --title "deploy" --command "INSTANT_TYPE=1 WERF_REQUIRE_BUILT_IMAGES=1 ./scripts/demo/_type_and_execute.sh ./scripts/demo/raw_scenario/deploy.sh" assets/demo/deploy.cast --overwrite
