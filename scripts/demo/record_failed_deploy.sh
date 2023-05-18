#!/bin/bash

source scripts/demo/env.sh
./scripts/demo/prepare_sample_app.sh

kubectl delete ns --ignore-not-found sample-app sample-app-2

werf converge --repo k3d-registry.sample-app.test:5000/sample-app --save-build-report --build-report-path=/tmp/images.json

APP_IMAGE=$(cat /tmp/images.json | jq -r .Images.app.DockerImageName)
helm upgrade --wait --install --create-namespace --namespace sample-app-2 sample-app .helm/ --set werf.image.app=${APP_IMAGE} --set host=sample-app-2.test --timeout 60s

./scripts/demo/raw_scenario/introduce_error.sh

werf build --repo k3d-registry.sample-app.test:5000/sample-app --save-build-report --build-report-path=/tmp/images.json

asciinema rec --title "failed deploy" --command ./scripts/demo/run_failed_deploy.sh assets/demo/failed_deploy.cast --overwrite
