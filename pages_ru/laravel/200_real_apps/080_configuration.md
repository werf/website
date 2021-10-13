---
title: Конфигурация и секреты
permalink: laravel/200_real_apps/080_config.html
examples_initial: examples/laravel/050_s3
examples: examples/laravel/080_configuration
base_url: https://github.com/werf/werf-guides/blob/master/
description: |
  В этой главе мы покажем, как правильно хранить и передавать приложениям секретную и несекретную конфигурацию. В предыдущих главах вся конфигурация находилась непосредственно в файлах конфигурации приложения или Kubernetes-манифестах. Для безопасности и гибкости часть этой конфигурации теперь будет передаваться приложениям с помощью values-файлов, ConfigMaps и Secrets.
---

## ConfigMaps для хранения файлов конфигурации

Сейчас конфигурационный файл [nginx.conf]({{ page.base_url | append: page.examples_initial | append: "/.werf/nginx.conf" }}) копируется в образ прямо во время сборки. Из-за этого при каждом его изменении будут происходить пересборка образа и перезапуск Pod'ов. Также сейчас нет простой возможности шаблонизировать `nginx.conf`.

Эти проблемы решаются, если перенести `.werf/nginx.conf` в отдельный ConfigMap, для того, чтобы монтировать `nginx.conf` уже во время деплоя, а не во время сборок:
{% include snippetcut_example path=".helm/templates/configmap-nginx.yaml" syntax="yaml" examples=page.examples %}
Теперь добавим этот ConfigMap как volume Deployment'а и примонтируем его как файл внутрь контейнера `frontend`:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="nginx_conf" examples=page.examples %}
Не забудем удалить более ненужный файл `.werf/nginx.conf`, а также команду копирования этого файла в образ во время сборки, после чего сборка образа `frontend` станет выглядеть так:
{% include snippetcut_example path="Dockerfile" syntax="dockerfile" snippet="frontend" examples=page.examples %}

## Обновление Deployment при изменении его ConfigMaps/Secrets

По умолчанию, изменения в ConfigMaps/Secrets, примонтированных к Deployment, StatefulSet или DaemonSet, не приведут к перезапуску Pod'ов с новой конфигурацией. Чтобы Pod'ы обновились, их нужно аннотировать хеш-суммами всех используемых Pod'ом ConfigMaps и Secrets. Тогда при изменении ConfigMap/Secret аннотация изменится и Pod пересоздастся. Так выглядит аннотация с хеш-суммой ConfigMap с `nginx.conf`:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="nginx_conf_checksum" examples=page.examples %}
Для каждого подключенного ConfigMap/Secret потребуется отдельная аннотация.
>_Авторы самоучителя предпочитают вместо аннотаций с хеш-суммами использовать операторы вроде [stakater/Reloader](https://github.com/stakater/Reloader), т. к. они проще, гибче и удобнее в работе._

## Конфигурация через values.yaml

Для упрощения конфигурации манифестов часто используемые параметры можно выносить в `.helm/values.yaml`:
{% include snippetcut_example path=".helm/values.yaml" syntax="yaml" snippet="mysql" examples=page.examples %}
... и подставлять их обратно в манифест так:
{% include snippetcut_example path=".helm/templates/database.yaml" syntax="yaml" snippet="volume_claim_templates" examples=page.examples %}
Это также позволит иметь разные значения для разных окружений, о чем более подробно будет рассказано в следующих разделах.

Но особенно полезным будет перенос в `.helm/values.yaml` повторяющейся конфигурации, например переменных окружения для приложения, которые мы используем в нескольких местах сразу:
{% include snippetcut_example path=".helm/values.yaml" syntax="yaml" snippet="app" examples=page.examples %}
Теперь переменные окружения из `.Values.app.envs` можно либо подставлять в манифест в `env` контейнера, как было ранее, либо вынести эти переменные окружения в ConfigMap, который потом подключить в контейнер через `envFrom`.

Первый вариант проще, но вариант с ConfigMap удобнее при большом количестве общих переменных окружения, т. к. позволит избежать их дублирования между контроллерами. ConfigMap может выглядеть так:
{% include snippetcut_example path=".helm/templates/configmap-app-envs.yaml" syntax="yaml" examples=page.examples %}

Подключаем этот ConfigMap через `envFrom` в Deployment как набор переменных окружения:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="backend_conf" examples=page.examples %}
Не забываем про аннотации с хеш-суммами, которые приведут к пересозданию Pod'ов при изменении ConfigMaps:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="configmap_app_envs_checksum" examples=page.examples %}
Аналогично этот ConfigMap подключается и в остальные ресурсы, в которых требуются эти переменные окружения.

## Секреты

Чтобы начать работу с секретами, сначала требуется сгенерировать симметричный ключ шифрования командой `werf helm secret generate-secret-key`. Но так как мы заранее подготовили и зашифровали секреты, то ключ шифрования мы сгенерировали за вас. Он находится в репозитории в файле `.werf_secret_key` и начнёт использоваться автоматически.
>_При работе с реальными приложениями ключ НЕЛЬЗЯ хранить в репозитории. Вместо этого рекомендуется передавать его через переменную окружения `WERF_SECRET_KEY`, храня его в безопасном месте. Подробнее про работу с ключами шифрования [здесь](https://ru.werf.io/documentation/v1.2/advanced/helm/configuration/secrets.html#%D0%BA%D0%BB%D1%8E%D1%87-%D1%88%D0%B8%D1%84%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D1%8F)._

В Kubernetes-манифестах у нас имеется конфигурация, которая не должна храниться в незашифрованном виде в репозитории, например, логин и пароль от базы данных. Вместо того, чтобы хранить логин и пароль в открытом виде в манифесте контроллера, мы будем хранить их зашифрованными в `.helm/secret-values.yaml` вместе с другими секретными параметрами:
{% include snippetcut_example path=".helm/secret-values.yaml" syntax="yaml" examples=page.examples %}

Увидеть расшифрованными секреты хранящиеся в `.helm/secret-values.yaml` можно следующей командой:
```bash
werf helm secret values decrypt .helm/secret-values.yaml
```

Ожидаемый результат:
```yaml
app:
  secretEnvs:
    APP_KEY: base64:GcPVmSxMZwsOJtNOJ9eVNNeU6B5buHuln93+w0TSvfE=
    DB_USERNAME: root
    DB_PASSWORD: password
    AWS_ACCESS_KEY_ID: minioadmin
    AWS_SECRET_ACCESS_KEY: minioadmin
mysql:
  secretEnvs:
    MYSQL_ROOT_PASSWORD: password
minio:
  secretEnvs:
    MINIO_ROOT_USER: minioadmin
    MINIO_ROOT_PASSWORD: minioadmin
```

Теперь нужно передать секреты из `.helm/secret-values.yaml` обратно в конфигурационный файл приложения. Для этого сначала передадим их в Secret-ресурс:
{% include snippetcut_example path=".helm/templates/secret-app-envs.yaml" syntax="yaml" examples=page.examples %}

... а затем примонтируем Secret в контейнеры как набор переменных окружения:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="backend_secret" examples=page.examples %}

Изменения в остальных файлах [.helm/templates]({{ page.base_url | append: page.examples | append: "/.helm/templates/" }}) для вынесения секретной конфигурации в Secret аналогичны, поэтому приводить их здесь не будем.

Также, для хранения и монтирования секретных конфигурационных файлов целиком можно использовать Secret. Выглядит это также, как использование ConfigMap для монтирования несекретных файлов конфигурации, [описанное выше](#configmaps-для-хранения-файлов-конфигурации), только само содержание Secret должно храниться зашифрованным в файлах `.helm/secret/...` или в `.helm/secret-values.yaml`, подробнее [здесь](https://ru.werf.io/documentation/v1.2/advanced/helm/configuration/secrets.html#%D1%81%D0%B5%D0%BA%D1%80%D0%B5%D1%82%D0%BD%D1%8B%D0%B5-%D1%84%D0%B0%D0%B9%D0%BB%D1%8B).

Больше информации по работе с секретами в [документации](https://ru.werf.io/documentation/v1.2/advanced/helm/configuration/secrets.html).

## Проверка работоспособности

Убедимся, что изменения в конфигурации не повлияли на работоспособность приложения:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guide-app
```

Ожидаемый результат:
{% include snippetcut_example path="werf-converge.log" syntax="shell" examples=page.examples %}

Проверим доступность приложения:
```shell
curl http://werf-guide-app/ping
```

Ожидаемый результат:
```
pong
```

Теперь секретная конфигурация приложения хранится безопасно, дублирование конфигурации минимизировано, а файлы конфигурации приложений могут шаблонизироваться и больше не требуют пересборок при изменениях.