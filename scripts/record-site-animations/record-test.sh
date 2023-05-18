#!/bin/bash

scriptdir="$( dirname -- "$BASH_SOURCE"; )";
source $scriptdir/env.sh

asciinema rec --title "test" --command "../type.sh '# Run tests in Kubernetes' ; ../type-and-execute.sh werf kube-run --repo k3d-registry.sample-app.test:5000/sample-app -- npm test" $@
