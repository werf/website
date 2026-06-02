{{- define "resources" }}
resources:
  requests:
    memory: {{ pluck .Values.werf.env .Values.resources.requests.memory | first | default .Values.resources.requests.memory._default }}
  limits:
    memory: {{ pluck .Values.werf.env .Values.resources.requests.memory | first | default .Values.resources.requests.memory._default }}
{{- end }}

{{- define "imagePullSecrets" }}
{{- with .Values.imagePullSecrets }}
imagePullSecrets:
{{- range . }}
  - name: {{ . | quote }}
{{- end }}
{{- end }}
{{- end }}

{{- define "targetCluster" }}
{{- if eq .Values.werf.env "production" }}
{{- $targetCluster := .Values.global.targetCluster | default "" -}}
{{- if or (eq $targetCluster "eu") (eq $targetCluster "ru") -}}
{{- $targetCluster -}}
{{- else -}}
{{- fail "For production, set global.targetCluster to either 'eu' or 'ru'." -}}
{{- end -}}
{{- else -}}
eu
{{- end }}
{{- end }}

{{- define "ingressClassName" }}
{{- pluck .Values.werf.env .Values.ingressClassName | first | default .Values.ingressClassName._default -}}
{{- end }}

{{- define "certificateIssuerName" }}
{{- pluck .Values.werf.env .Values.certificateIssuerName | first | default .Values.certificateIssuerName._default -}}
{{- end }}

{{- define "clusterPlacement" }}
{{- $targetCluster := include "targetCluster" . -}}
{{- $clusterConfig := get (.Values.clusters | default dict) $targetCluster | default dict -}}
{{- $placement := get $clusterConfig "placement" | default dict -}}
{{- with (get $placement "nodeSelector") }}
nodeSelector:
{{ toYaml . | indent 2 }}
{{- end }}
{{- with (get $placement "tolerations") }}
tolerations:
{{ toYaml . | indent 2 }}
{{- end }}
{{- with (get $placement "affinity") }}
affinity:
{{ toYaml . | indent 2 }}
{{- end }}
{{- end }}

{{- define "docsCurrentMajor" -}}
{{- $current := pluck .Values.werf.env .Values.docsRouting.currentMajor | first | default .Values.docsRouting.currentMajor._default -}}
{{- $current = trim $current -}}
{{- $current = trimPrefix "/docs/" $current -}}
{{- $current = trimAll "/" $current -}}
{{- $current = lower $current -}}
{{- if hasPrefix "v" $current -}}
{{- $current -}}
{{- else -}}
{{- printf "v%s" $current -}}
{{- end -}}
{{- end }}

{{- define "docsSupportedRootsCSV" -}}
{{- $roots := pluck .Values.werf.env .Values.docsRouting.supportedRoots | first | default .Values.docsRouting.supportedRoots._default -}}
{{- $items := list (include "docsCurrentMajor" .) -}}
{{- range $root := $roots }}
  {{- $normalized := trim $root -}}
  {{- $normalized = trimPrefix "/docs/" $normalized -}}
  {{- $normalized = trimAll "/" $normalized -}}
  {{- $normalized = lower $normalized -}}
  {{- if and $normalized (regexMatch "^v[0-9]+$" $normalized) (not (has $normalized $items)) -}}
    {{- $items = append $items $normalized -}}
  {{- end -}}
{{- end -}}
{{- join "," $items -}}
{{- end }}

{{- define "docsSupportedRootsPattern" -}}
{{- join "|" (splitList "," (include "docsSupportedRootsCSV" .)) -}}
{{- end }}

{{- define "docsLatestAliasEnabled" -}}
{{- $configured := pluck .Values.werf.env .Values.docsRouting.latestAliasEnabled | first | default .Values.docsRouting.latestAliasEnabled._default | default "auto" -}}
{{- $configured = lower (trim (printf "%v" $configured)) -}}
{{- if or (eq $configured "true") (eq $configured "1") (eq $configured "yes") (eq $configured "on") -}}
true
{{- else if or (eq $configured "false") (eq $configured "0") (eq $configured "no") (eq $configured "off") -}}
false
{{- else if eq (include "docsSupportedRootsCSV" .) "v1" -}}
false
{{- else -}}
true
{{- end -}}
{{- end }}

{{- define "docsModernStructuredRootsPattern" -}}
{{- $roots := without (splitList "," (include "docsSupportedRootsCSV" .)) "v1" -}}
{{- if gt (len $roots) 0 -}}
{{- join "|" $roots -}}
{{- end -}}
{{- end }}

