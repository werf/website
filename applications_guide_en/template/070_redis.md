---
title: Подключаем redis
sidebar: applications_guide
guide_code: ____________
permalink: ____________/070_redis.html
toc: false
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/requirements.yaml
- .helm/values.yaml
- .helm/secret-values.yaml
____________
- .gitlab-ci.yml
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении работу с простейшей базой данных типа in-memory — [Redis](https://redis.io/) (другим популярным примером является [memcached](https://memcached.org/)). Это означает, что база данных будет stateless.

{% offtopic title="Как быть, если база данных должна сохранять данные на диске?" %}
Этот вопрос мы разберём в следующей главе на примере [PostgreSQL](080_database.html). В рамках текущей главы разберёмся с общими вопросами: как базу данных в принципе завести в кластер, сконфигурировать и подключиться к ней из приложения.
{% endofftopic %}

В простейшем случае нет необходимости вносить изменения в сборку: уже собранные образы есть на DockerHub. Надо просто выбрать правильный образ, корректно сконфигурировать его в своей инфраструктуре, а потом подключиться к базе данных из ____________-приложения.

## Сконфигурировать Redis в Kubernetes

Для того, чтобы сконфигурировать Redis в кластере, необходимо определить объекты с помощью Helm. Мы можем сделать это самостоятельно, но рассмотрим вариант с подключением внешнего чарта. В любом случае потребуется указать: имя сервиса, порт, логин и пароль — и разобраться, как эти параметры пробросить в подключенный внешний чарт.

Итак, нам понадобится:

1. Указать Redis как зависимый сабчарт в `requirements.yaml`;
2. Сконфигурировать в werf работу с зависимостями;
3. Сконфигурировать подключённый сабчарт;
4. Убедиться, что создаётся кластер Redis в конфигурации master-slave.

Пропишем сабчарт с Redis:

{% snippetcut name=".helm/requirements.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/.helm/requirements.yaml" %}
{% raw %}
```yaml
dependencies:
<...>
- name: redis
  version: 9.3.2
  repository: https://kubernetes-charts.storage.googleapis.com/
  condition: redis.enabled
```
{% endraw %}
{% endsnippetcut %}

Для того, чтобы werf при деплое загрузил необходимые нам сабчарты, нужно прописать в `.gitlab-ci.yml` работу с зависимостями:

{% snippetcut name=".gitlab-ci.yml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/.gitlab-ci.yml" %}
{% raw %}
```yaml
.base_deploy: &base_deploy
  stage: deploy
  script:
    - werf helm repo init
    - werf helm dependency update
    - werf deploy --set "global.ci_url=$(cut -d / -f 3 <<< $CI_ENVIRONMENT_URL)"
```
{% endraw %}
{% endsnippetcut %}

… а также сконфигурировать имя сервиса, порт, логин и пароль, согласно [документации сабчарта](https://github.com/bitnami/charts/tree/master/bitnami/redis/#parameters):

{% snippetcut name=".helm/values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/.helm/values.yaml" %}
{% raw %}
```yaml
redis:
  fullnameOverride: redis
  nameOverride: redis
```
{% endraw %}
{% endsnippetcut %}


{% offtopic title="Откуда такой ключ Redis?" %}
Этот ключ должен совпадать с именем сабчарта-зависимости в файле `requirements.yaml` — тогда настройки будут пробрасываться в сабчарт.
{% endofftopic %}
{% snippetcut name="secret-values.yaml (расшифрованный)" url="#" ignore-tests %}
{% raw %}
```yaml
redis:
  password: "LYcj6c09D9M4htgGh64vXLxn95P4Wt"
```
{% endraw %}
{% endsnippetcut %}

Сконфигурировать логин и порт для подключения у этого сабчарта невозможно, но если изучить исходный код — можно найти значения, использующиеся в сабчарте. Пропишем нужные значения с понятными нам ключами — они понадобятся позже, когда будем конфигурировать приложение.

{% snippetcut name=".helm/values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/.helm/values.yaml" %}
{% raw %}
```yaml
redis:
<...>
   _port:
      _default: 6379
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="Почему мы пишем эти ключи с префиксом _ и вообще легально ли это?" %}
Когда мы пишем дополнительные ключи по соседству с ключами, пробрасывающимися в сабчарт, рискуем случайно «зацепить» лишнее. Поэтому нужно быть внимательным, сверяться с [документацией сабчарта](https://github.com/bitnami/charts/tree/master/bitnami/redis/#parameters) и не использовать пересекающиеся ключи.

Для надёжности — введём соглашение на использование знака подчёркивания `_` в начале таких ключей.
{% endofftopic %}


{% offtopic title="Как быть, если найти параметры не получается?" %}

Некоторые сервисы вообще не требуют аутентификации. В частности, Redis зачастую используется без неё.

{% endofftopic %}

Если посмотреть на рендер (`werf helm render`) нашего приложения с включенным сабчартом для Redis, то можем увидеть, какие будут созданы объекты типа Service:

```yaml
# Source: example_4/charts/redis/templates/redis-master-svc.yaml
apiVersion: v1
kind: Service
metadata:
  name: guided-redis-master

# Source: example_4/charts/redis/templates/redis-slave-svc.yaml
apiVersion: v1
kind: Service
metadata:
  name: guided-redis-slave
```

Знание этих Services требуется, чтобы потом к ним подключаться.

## Подключение ____________ приложения к базе Redis

В нашем приложении мы будем подключаться к master-узлу Redis. Необходимо, чтобы приложение при выкате в любое окружение подключалось к правильному Redis.

В нашем случае Redis будет использоваться как хранилище сессий.

{% snippetcut name="____________" url="#" %}
{% raw %}
```____________
____________
```
{% endraw %}
{% endsnippetcut %}

Для подключения к базе данных нам, очевидно, нужно знать: хост, порт, логин, пароль. В коде приложения мы используем несколько переменных окружения: `REDIS_HOST`, `REDIS_PORT`, `REDIS_LOGIN`, `REDIS_PASSWORD`____________. Часть значений уже определена в `values.yaml` для подключаемого сабчарта. Можно воспользоваться теми же значениями и дополнить их.

Будем **конфигурировать хост** через `values.yaml`:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        - name: REDIS_HOST
          value: {{ pluck .Values.global.env .Values.redis.host | first | default .Values.redis.host._default | quote }}
```
{% endraw %}
{% endsnippetcut %}

{% snippetcut name=".helm/values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/.helm/values.yaml" %}
{% raw %}
```yaml
redis:
<...>
  host:
    _default: redis-master
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="Зачем такие сложности? Может, просто прописать значения в шаблоне?" %}

Казалось бы, можно написать примерно так:

{% snippetcut name=".helm/templates/deployment.yaml" url="#" ignore-tests %}
{% raw %}
```yaml
        - name: REDIS_HOST
          value: {{ pluck .Values.global.env .Values.redis.host | first | default .Values.redis.host._default | quote }}
```
{% endraw %}
{% endsnippetcut %}

Но на практике иногда возникает необходимость переехать в другую базу данных или кастомизировать что-то: в этих случаях в разы удобнее работать через `values.yaml`. Причём значений для разных окружений мы не прописываем, а ограничиваемся значением по умолчанию:

{% snippetcut name="values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/.helm/values.yaml" %}
{% raw %}
```yaml
redis:
<...>
  host:
    _default: redis-master
```
{% endraw %}
{% endsnippetcut %}

И под конкретные окружения прописываем значения только в случаях, когда это действительно нужно.
{% endofftopic %}

**Конфигурируем логин и порт** через `values.yaml`, просто прописывая значения:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        - name: REDIS_PORT
          value: {{ pluck .Values.global.env .Values.redis._port | first | default .Values.redis._port._default | quote }}
```
{% endraw %}
{% endsnippetcut %}

Мы уже **сконфигурировали пароль** — используем прописанное ранее значение:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        - name: REDIS_PASSWORD
          value: {{ .Values.redis.password | quote }}
```
{% endraw %}
{% endsnippetcut %}

Также нам нужно **сконфигурировать переменные, необходимые приложению** для работы с Redis:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        - name: SESSION_TTL
          value: {{ pluck .Values.global.env .Values.app.redis.session_ttl | first | default .Values.app.redis.session_ttl._default | quote }}
```
{% endraw %}
{% endsnippetcut %}


{% snippetcut name="values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/.helm/values.yaml" %}
{% raw %}
```yaml
  redis:
    session_ttl:
      _default: "3600"
    cookie_secret:
      _default: "supersecret"
```
{% endraw %}
{% endsnippetcut %}

<div id="go-forth-button">
    <go-forth url="080_database.html" label="Подключение базы данных" base-url="applications_guide_ru"></go-forth>
</div>
