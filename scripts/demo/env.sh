#!/bin/bash

demoscriptsdir="$(realpath $(dirname -- "$BASH_SOURCE"))"

for var in $(env | grep LC_ | grep -v LC_LANG | cut -d'=' -f1) ; do
	eval "export $var=en_US.UTF-8"
done
export LC_LANG=en
export TZ=

export CACHE_VERSION=""
export WERF_BUILDAH_MODE=native-chroot
export WERF_INSECURE_REGISTRY=1
export WERF_FORCE_STAGED_DOCKERFILE=1
export WERF_SKIP_TLS_VERIFY_REGISTRY=1
export WERF_DISABLE_AUTO_HOST_CLEANUP=1
export WERF_BUNDLE_INSECURE_REGISTRY=1
export WERF_BUNDLE_SKIP_TLS_VERIFY_REGISTRY=1
#export WERF_LOG_TERMINAL_WIDTH=130
export WERF_DEV=1
export WERF_DIR=${demoscriptsdir}/sample-app
