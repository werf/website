Если вы используете внешний registry, он наверняка закрыт логином и паролем. Вы должны были сталкиваться с ними в главе «Подготовка окружения».

Эти логин и пароль мы сообщаем кластеру с помощью объекта типа Secret с именем `registrysecret`. Допустим, ваш логин — `admin`, пароль — `admin`, а registry находится по адресу `registry.example.com`. Сформируем файл `registry-secret.yaml`:

{% snippetcut name=".helm/templates/registry-secret.yaml" url="#" %}
{% raw %}
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: registrysecret
  annotations:
    "helm.sh/hook": pre-install
type: kubernetes.io/dockerconfigjson
data:
  {{- $registry := "registry.example.com" }}
  {{- $login := "admin" }}
  {{- $password := "admin" }}
  .dockerconfigjson: {{ printf "{\"auths\": {\"%s\": {\"auth\": \"%s\"}}}" $registry (printf "%s:%s" $login $password | b64enc) | b64enc }}
```
{% endraw %}
{% endsnippetcut %}

_Примечание: в приведённом примере ключи доступа хранятся в незашифрованном виде. Это небезопасно. Вопросы безопасного хранения ключей мы рассмотрим в главе «Организация не локальной разработки»._

Теперь нужно указать этот Secret в каждом объекте Deployment в атрибуте `imagePullSecrets`, чтобы кластер мог выгрузить образ:

{% snippetcut name=".helm/templates/deployment.yaml" url="#" %}
{% raw %}
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: basicapp
spec:
  selector:
    matchLabels:
      app: basicapp
  revisionHistoryLimit: 3
  strategy:
    type: RollingUpdate
  replicas: 1
  template:
    metadata:
      labels:
        app: basicapp
    spec:
      imagePullSecrets:                           # Added line
      - name: "registrysecret"                    # Added line
      containers:
      - name: basicapp
        command: ["java","-jar","/app/demo.jar"]
        image: {{ .Values.werf.image.basicapp }}
        workingDir: /app
        ports:
        - containerPort: 8080
          protocol: TCP
        env:
        - name: "SQLITE_FILE"
          value: "app.db"
```
{% endraw %}
{% endsnippetcut %}
