#!/bin/bash

asciinema rec --title "cleanup" --command "../type.sh '# Cleanup no-longer-relevant images from the container registry' ; ../type-and-execute.sh werf cleanup --repo k3d-registry.sample-app.test:5000/sample-app"
