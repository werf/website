#!/bin/bash

scriptdir="$( dirname -- "$BASH_SOURCE"; )";
source $scriptdir/env.sh

asciinema rec --title "deploy" --command "../type.sh '# Pack the defined state and push the OCI Helm Chart to the container registry' ; ../type-and-execute.sh werf bundle publish --repo k3d-registry.sample-app.test:5000/sample-app" $@
