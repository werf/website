#!/bin/bash

demoscriptsdir="$(realpath $(dirname -- "$BASH_SOURCE"))"

tmux new-session  -d "source ${demoscriptsdir}/env.sh ; cd ${WERF_DIR} ; INPUT_THROTTLING=18 APP_IMAGE=$(cat /tmp/images.json | jq -r .Images.app.DockerImageName) source ${demoscriptsdir}/_type_and_execute.sh ${demoscriptsdir}/raw_scenario/failed_deploy_with_helm.sh ${APP_IMAGE} ; bash" \; \
     split-window -h -l 65% "source ${demoscriptsdir}/env.sh ; cd ${WERF_DIR} ; WERF_LOG_TERMINAL_WIDTH=90 WERF_REQUIRE_BUILT_IMAGES=1 INPUT_THROTTLING=11 source ${demoscriptsdir}/_type_and_execute.sh ${demoscriptsdir}/raw_scenario/failed_deploy_with_werf.sh ; bash" \; \
     attach
