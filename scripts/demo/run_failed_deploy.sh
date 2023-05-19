#!/bin/bash

demoscriptsdir="$(realpath $(dirname -- "$BASH_SOURCE"))"

tmux new-session  -d "source ${demoscriptsdir}/env.sh ; cd ${WERF_DIR} ; INSTANT_TYPE=1 TYPE_DELAY=1 APP_IMAGE=$(cat /tmp/images.json | jq -r .Images.app.DockerImageName) source ${demoscriptsdir}/_type_and_execute.sh ${demoscriptsdir}/raw_scenario/failed_deploy_with_helm.sh ${APP_IMAGE} ; bash" \; \
     split-window -h -l 65% "source ${demoscriptsdir}/env.sh ; cd ${WERF_DIR} ; WERF_LOG_TERMINAL_WIDTH=90 WERF_REQUIRE_BUILT_IMAGES=1 INSTANT_TYPE=1 TYPE_DELAY=1 source ${demoscriptsdir}/_type_and_execute.sh ${demoscriptsdir}/raw_scenario/failed_deploy_with_werf.sh ; bash" \; \
     attach
