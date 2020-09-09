{{ define "app_vars" }}
{{ $rmq_user := pluck .Values.global.env .Values.app.rabbitmq.user | first | default .Values.app.rabbitmq.user._default }}
{{ $rmq_pass := pluck .Values.global.env .Values.app.rabbitmq.password | first | default .Values.app.rabbitmq.password._default }}
{{ $rmq_host := pluck .Values.global.env .Values.app.rabbitmq.host | first | default .Values.app.rabbitmq.host._default }}
{{ $rmq_port := pluck .Values.global.env .Values.app.rabbitmq.port | first | default .Values.app.rabbitmq.port._default }}
{{ $rmq_vhost := pluck .Values.global.env .Values.app.rabbitmq.vhost | first | default .Values.app.rabbitmq.vhost._default }}
rmq_uri: {{ printf "amqp://%s:%s@%s:%s/%s" $rmq_user $rmq_pass $rmq_host ($rmq_port|toString) $rmq_vhost }}
{{- end }}
