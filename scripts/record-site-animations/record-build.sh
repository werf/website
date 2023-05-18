#!/bin/bash

scriptdir="$( dirname -- "$BASH_SOURCE"; )";
source $scriptdir/env.sh

docker pull ghcr.io/werf/test-registry:client-id-8bc33c32-2726-4f65-9045-27fd8d1a10e6-334965913652
docker tag ghcr.io/werf/test-registry:client-id-8bc33c32-2726-4f65-9045-27fd8d1a10e6-334965913652 k3d-registry.sample-app.test:5000/sample-app:client-id-8bc33c32-2726-4f65-9045-27fd8d1a10e6-334965913652
docker push k3d-registry.sample-app.test:5000/sample-app:client-id-8bc33c32-2726-4f65-9045-27fd8d1a10e6-334965913652

asciinema rec --title "build" --command "../type.sh '# Sync the registry to the defined state' ; ../type-and-execute.sh werf build --repo k3d-registry.sample-app.test:5000/sample-app" $@
