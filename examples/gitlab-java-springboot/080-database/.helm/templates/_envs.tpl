{{- define "database_envs" }}
- name: POSTGRESQL_HOST
  value: {{ pluck .Values.postgresql.postgresqlHost | quote }}
- name: POSTGRESQL_LOGIN
  value: {{ pluck .Values.global.env .Values.postgresql.login | first | default .Values.postgresql.login_default | quote }}
- name: POSTGRESQL_PORT
  value: {{ pluck .Values.global.env .Values.postgresql.port | first | default .Values.postgresql.port_default | quote }}
- name: POSTGRESQL_PASSWORD
  value: {{ pluck .Values.global.env .Values.postgresql.password | first | default .Values.postgresql.password_default | quote }}
{{- end }}
