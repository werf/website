#!/bin/bash

scriptdir="$( dirname -- "$BASH_SOURCE"; )";
source $scriptdir/env.sh

export CACHE_VERSION=2
werf build --repo k3d-registry.sample-app.test:5000/sample-app

export CACHE_VERSION=3
werf build --repo k3d-registry.sample-app.test:5000/sample-app

asciinema rec --title "cleanup" --command "../type.sh '# Cleanup no-longer-relevant images from the container registry' ; ../type-and-execute.sh werf cleanup --repo k3d-registry.sample-app.test:5000/sample-app" $@
