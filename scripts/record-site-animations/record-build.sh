#!/bin/bash

asciinema rec --title "build" --command "../type.sh '# Sync the registry to the defined state' ; ../type-and-execute.sh werf build --repo k3d-registry.sample-app.test:5000/sample-app"
