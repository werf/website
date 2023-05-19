#!/bin/bash

source scripts/demo/env.sh
./scripts/demo/prepare_sample_app.sh

export CACHE_VERSION=2
werf build --repo k3d-registry.sample-app.test:5000/sample-app

export CACHE_VERSION=3
werf build --repo k3d-registry.sample-app.test:5000/sample-app

asciinema rec --title "cleanup" --command "INSTANT_TYPE=1 TYPE_DELAY=0.8 ./scripts/demo/_type_and_execute.sh ./scripts/demo/raw_scenario/cleanup.sh" assets/demo/cleanup.cast --overwrite
