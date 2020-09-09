{{ define "app_vars" }}
{{ $redis_pass := pluck .Values.global.env .Values.app.redis.password | first | default .Values.app.redis.password._default }}
{{ $redis_host := pluck .Values.global.env .Values.app.redis.host.master | first | default .Values.app.redis.host.master._default }}
redis_uri: {{ printf "redis://root:%s@%s:6379" $redis_pass $redis_host }}
{{ $db_user := pluck .Values.global.env .Values.app.postgresql.user | first | default .Values.app.postgresql.user._default }}
{{ $db_pass := pluck .Values.global.env .Values.app.postgresql.password | first | default .Values.app.postgresql.password._default }}
{{ $db_host := pluck .Values.global.env .Values.app.postgresql.host | first | default .Values.app.postgresql.host._default }}
{{ $db_port := pluck .Values.global.env .Values.app.postgresql.port | first | default .Values.app.postgresql.port._default }}
{{ $db_name := pluck .Values.global.env .Values.app.postgresql.db | first | default .Values.app.postgresql.db._default }}
db_uri: {{ printf "postgresql://%s:%s@%s:%s/%s" $db_user $db_pass $db_host ($db_port|toString) $db_name }}
{{ $rmq_user := pluck .Values.global.env .Values.app.rabbitmq.user | first | default .Values.app.rabbitmq.user._default }}
{{ $rmq_pass := pluck .Values.global.env .Values.app.rabbitmq.password | first | default .Values.app.rabbitmq.password._default }}
{{ $rmq_host := pluck .Values.global.env .Values.app.rabbitmq.host | first | default .Values.app.rabbitmq.host._default }}
{{ $rmq_port := pluck .Values.global.env .Values.app.rabbitmq.port | first | default .Values.app.rabbitmq.port._default }}
{{ $rmq_vhost := pluck .Values.global.env .Values.app.rabbitmq.vhost | first | default .Values.app.rabbitmq.vhost._default }}
rmq_uri: {{ printf "amqp://%s:%s@%s:%s/%s" $rmq_user $rmq_pass $rmq_host ($rmq_port|toString) $rmq_vhost }}
{{- end }}