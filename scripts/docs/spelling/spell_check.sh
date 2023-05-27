#!/bin/bash

set -e

arg_site_lang="${1:?ERROR: Site language \'en\' or \'ru\' should be specified as the first argument.}"

script=$(cat <<EOF
cd /spelling/$arg_site_lang && \
  ./container_spell_check.sh $arg_site_lang
EOF
)

werf run spell-checker --dev --docker-options="--entrypoint=bash" -- -c "$script"
