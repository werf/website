#!/bin/sh
set -eu

normalize_major() {
  root=$(printf '%s' "$1" | tr '[:upper:]' '[:lower:]')
  root=${root#/docs/}
  root=${root#/}
  root=${root%/}

  case "$root" in
    v[0-9]*)
      suffix=${root#v}
      case "$suffix" in
        ''|*[!0-9]*) return 1 ;;
      esac
      printf 'v%s' "$suffix"
      return 0
      ;;
    [0-9]*)
      case "$root" in
        ''|*[!0-9]*) return 1 ;;
      esac
      printf 'v%s' "$root"
      return 0
      ;;
  esac

  return 1
}

normalize_bool() {
  value=$(printf '%s' "$1" | tr '[:upper:]' '[:lower:]' | tr -d '[:space:]')

  case "$value" in
    1|true|yes|on)
      printf 'true'
      return 0
      ;;
    0|false|no|off)
      printf 'false'
      return 0
      ;;
  esac

  return 1
}

append_root() {
  root="$1"
  case ",${supported_docs_roots_csv}," in
    *,${root},*) ;;
    *)
      if [ -n "$supported_docs_roots_csv" ]; then
        supported_docs_roots_csv="${supported_docs_roots_csv},${root}"
      else
        supported_docs_roots_csv="$root"
      fi
      ;;
  esac
}

current_docs_major="$(normalize_major "${CURRENT_DOCS_MAJOR:-}" || true)"
supported_docs_roots_csv=""

if [ -n "$current_docs_major" ]; then
  append_root "$current_docs_major"
fi

for raw_root in $(printf '%s' "${SUPPORTED_DOCS_MAJOR_VERSIONS:-}" | tr ',;[:space:]' '\n'); do
  normalized_root="$(normalize_major "$raw_root" || true)"
  if [ -n "$normalized_root" ]; then
    append_root "$normalized_root"
  fi
done

if [ -z "$current_docs_major" ] && [ -n "$supported_docs_roots_csv" ]; then
  current_docs_major=${supported_docs_roots_csv%%,*}
fi

if [ -z "$current_docs_major" ]; then
  current_docs_major="v1"
  append_root "$current_docs_major"
fi

docs_supported_roots_regex=$(printf '%s' "$supported_docs_roots_csv" | tr ',' '|')
docs_latest_alias_enabled="$(normalize_bool "${DOCS_LATEST_ALIAS_ENABLED:-}" || true)"

if [ -z "$docs_latest_alias_enabled" ]; then
  if [ "$current_docs_major" = "v1" ] && [ "$supported_docs_roots_csv" = "v1" ]; then
    docs_latest_alias_enabled="false"
  else
    docs_latest_alias_enabled="true"
  fi
fi

export CURRENT_DOCS_MAJOR="$current_docs_major"
export DOCS_SUPPORTED_ROOTS_REGEX="$docs_supported_roots_regex"
export DOCS_LATEST_ALIAS_ENABLED="$docs_latest_alias_enabled"
envsubst '$CURRENT_DOCS_MAJOR $DOCS_SUPPORTED_ROOTS_REGEX $DOCS_LATEST_ALIAS_ENABLED' \
  < /etc/nginx/templates/nginx.conf.template \
  > /etc/nginx/nginx.conf




