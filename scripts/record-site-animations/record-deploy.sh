#!/bin/bash

scriptdir="$( dirname -- "$BASH_SOURCE"; )";
source $scriptdir/env.sh

asciinema rec --title "deploy" --command "../type.sh '# Sync Kubernetes to the defined state' ; ../type-and-execute.sh werf converge --repo k3d-registry.sample-app.test:5000/sample-app ; ../type-and-execute.sh curl sample-app.test/ping" $@
