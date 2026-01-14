#!/bin/bash

set -e

arg_site_lang="${1:?ERROR: Site language \'en\' or \'ru\' should be specified as the first argument.}"

script=$(cat <<EOF
cd /app && \
  bundle exec htmlproofer \
    --allow-hash-href \
    --empty-alt-ignore \
    --check_html \
    --file_ignore "/____________/,/changelog.html/" \
    --url_ignore "/localhost/,/cloudyuga.guru/,/t.me/,/slack.com/,/cncf.io/,/example.com/,/github.com/werf/,/habr.com/,/apple-touch-icon.png/,/site.webmanifest/,/favicon/,/safari-pinned-tab/,/docs/,/guides.html/" \
    --url-swap "github.com/werf/website/blob/main:github.com/werf/website/blob/$GITHUB_SHA,github.com/werf/website/tree/main:github.com/werf/website/tree/$GITHUB_SHA" \
    --http-status-ignore "0,429,403" \
    /app/_site
EOF
)

werf run doc_$arg_site_lang --docker-options="--entrypoint=bash" -- -c "$script"
