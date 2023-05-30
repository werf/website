#!/bin/bash

set -e

arg_site_lang="${1:?ERROR: Site language \'en\' or \'ru\' should be specified as the first argument.}"

if [ -n "$2" ]; then
  arg_target_page=$2
fi

script=$(cat <<EOF
cd /spelling/$arg_site_lang && \
  ./container_spell_check.sh $arg_site_lang $arg_target_page
EOF
)

werf run spell-checker --dev --docker-options="--entrypoint=bash" -- -c "$script"
