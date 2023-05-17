#!/bin/bash

set -e

arg_site_lang="${1:?ERROR: Site language \'en\' or \'ru\' should be specified as the first argument.}"

script=$(cat <<EOF
cd /spelling/$arg_site_lang && \
  hunspell
EOF
)

werf run doc_$arg_site_lang --docker-options="--entrypoint=bash" -- -c "$script"
