#!/bin/bash

scriptdir="$( dirname -- "$BASH_SOURCE"; )";
source $scriptdir/env.sh

asciinema rec --title "distribute" --command "INSTANT_TYPE=1 WERF_REQUIRE_BUILT_IMAGES=1 ./scripts/demo/_type_and_execute.sh ./scripts/demo/raw_scenario/distribute.sh" assets/demo/distribute.cast --overwrite
