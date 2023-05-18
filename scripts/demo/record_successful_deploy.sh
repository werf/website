#!/bin/bash

source scripts/demo/env.sh
./scripts/demo/prepare_sample_app.sh

werf build --repo k3d-registry.sample-app.test:5000/sample-app --save-build-report --build-report-path=/tmp/images.json
kubectl delete ns --ignore-not-found sample-app sample-app-2

asciinema rec --title "successful deploy" --command ./scripts/demo/run_successful_deploy.sh assets/demo/successful_deploy.cast --overwrite
