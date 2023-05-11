#!/bin/bash

set -e

arg_site_lang="${1:?ERROR: Site language \'en\' or \'ru\' should be specified as the first argument.}"

script=$(cat <<EOF
cd /app && \
  bundle exec htmlproofer \
    --allow-hash-href \
    --empty-alt-ignore \
    --check_html \
    --file_ignore "/____________/" \
    --url_ignore "/localhost/,/t.me/,/slack.com/,/cncf.io/,/example.com/,/github.com/werf/,/habr.com/,/apple-touch-icon.png/,/site.webmanifest/,/favicon/,/safari-pinned-tab/,/documentation/,/guides.html/" \
    --http-status-ignore "0,429,403,500" \
    /app/_site
EOF
)

werf run doc_$arg_site_lang --dev --docker-options="--entrypoint=bash" -- -c "$script"
