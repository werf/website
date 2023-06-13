#!/bin/bash

set -e

arg_site_lang="${1:?ERROR: Site language \'en\' or \'ru\' should be specified as the first argument.}"

if [ -n "$2" ]; then
  arg_target_page=$2
fi

if [ -n "$3" ]; then
  arg_get_plain_text=$3
fi

script=$(cat <<EOF
cd /spelling/$arg_site_lang && \
  ./container_spell_check.sh $arg_site_lang $arg_target_page $arg_get_plain_text
EOF
)

werf run spell-checker --env='test' --dev --docker-options="--entrypoint=bash" -- -c "$script"
