#!/bin/bash

source scripts/demo/env.sh
./scripts/demo/prepare_sample_app.sh

./scripts/demo/reset.sh

buildah images --format '{{.Name}} {{.ID}}' | grep sample-app | cut -d' ' -f4 | xargs buildah rmi -f

buildah pull node:18-bullseye
buildah pull ghcr.io/werf/test-registry:client-id-8bc33c32-2726-4f65-9045-27fd8d1a10e6-334965913652
buildah tag ghcr.io/werf/test-registry:client-id-8bc33c32-2726-4f65-9045-27fd8d1a10e6-334965913652 k3d-registry.sample-app.test:5000/sample-app:client-id-8bc33c32-2726-4f65-9045-27fd8d1a10e6-334965913652
buildah push --tls-verify=false k3d-registry.sample-app.test:5000/sample-app:client-id-8bc33c32-2726-4f65-9045-27fd8d1a10e6-334965913652

asciinema rec --title "build" --command "INSTANT_TYPE=1 ./scripts/demo/_type_and_execute.sh ./scripts/demo/raw_scenario/build.sh" assets/demo/build.cast --overwrite
