#!/bin/bash

set -e

cd "$(dirname "$0")/../.."

current_root="v3"

chart_dir="$(mktemp -d)"
trap 'rm -rf "$chart_dir"' EXIT
cp -R .helm/. "$chart_dir/"
cat >"$chart_dir/Chart.yaml" <<EOF
apiVersion: v2
name: werfio
version: 0.0.0
EOF

fail() {
  echo "FAIL: $1" >&2
  exit 1
}

for cluster in eu ru; do
  rendered="$(helm template werfio-production "$chart_dir" \
    --set werf.env=production \
    --set werf.namespace=werfio-production \
    --set "werf.image.web-backend=stub" \
    --set "global.targetCluster=$cluster")"

  echo "$rendered" | grep -q -- "- name: CURRENT_DOCS_MAJOR" ||
    fail "$cluster: CURRENT_DOCS_MAJOR env is missing"
  echo "$rendered" | grep -A1 -- "- name: CURRENT_DOCS_MAJOR" | grep -q "value: \"$current_root\"" ||
    fail "$cluster: CURRENT_DOCS_MAJOR is not $current_root"
  echo "$rendered" | grep -A1 -- "- name: SUPPORTED_DOCS_MAJOR_VERSIONS" | grep -q 'value: "v3,v2,v1.2"' ||
    fail "$cluster: supported documentation versions changed"
  echo "$rendered" | grep -q "rewrite \^/docs/?\$ */docs/$current_root/ " ||
    fail "$cluster: /docs/ does not redirect to /docs/$current_root/"
  echo "$rendered" | grep -q "rewrite \^/docs/(?<ver>latest|.*)/?\$ */docs/\$ver/usage/project_configuration/overview.html" ||
    fail "$cluster: latest alias is missing in docs shortcuts"
done

grep -q "latest|pr-\[^/\]+|$(echo "$current_root" | sed 's/\./\\\\./')" .werf/nginx-dev.conf ||
  fail "nginx-dev.conf does not proxy /docs/$current_root/"

for lang in en ru; do
  grep -q "url: /docs/$current_root/" "_data/$lang/topnav.yml" ||
    fail "$lang topnav does not link to /docs/$current_root/"
done

grep -q "CURRENT_DOCS_MAJOR: \"$current_root\"" docker-compose.yml ||
  fail "local backend does not use $current_root"

for root in v3 v2; do
  grep -qx "Allow: /docs/$root" robots.txt || fail "robots.txt blocks $root"
  grep -q "/docs/$root/sitemap.xml" sitemap.xml || fail "sitemap.xml omits $root"
done

node scripts/docs/check_doc_group.js

echo "OK: docs roots are consistent with currentRoot $current_root"
