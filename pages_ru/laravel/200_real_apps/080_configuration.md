---
title: Конфигурация и секреты
permalink: laravel/200_real_apps/080_config.html
examples_initial: examples/laravel/050_s3
examples: examples/laravel/080_configuration
base_url: https://github.com/werf/werf-guides/blob/master/
description: |
  В этой главе мы покажем, как правильно использовать и хранить секретную и несекретную конфигурацию приложения. 

  В предыдущих главах конфигурация добавлялась прямо в контейнеры при сборке или использовалась как есть в переменных окружения контейнеров при выкате.
  
  Теперь для безопасности и гибкости конфигурация будет сохраняться в ConfigMap и Secret. А в дополнение к параметрам helm-чарта (Values) и секретам werf будут продемонстрированы подходы параметризации и переиспользования конфигурации, а также хранения конфиденциальных данных вместе с кодом в git-репозитории проекта.
---

## ConfigMap и Secret

В Kubernetes есть специальные типы ресурсов ConfigMap и Secret, которые предназначены для отделения конфигурации, зависящей от среды и окружения, от образов контейнеров. 

Оба ресурса позволяют хранить данные в парах "ключ-значение" и впоследствии использовать их как переменные окружения, аргументы командной строки или как файлы примонтированные в выбранный контейнер.

ConfigMap предназначен для хранения неконфиденциальных данных, в отличие от Secret, который используется для хранения различных типов секретов.

Подробнее про эти типы ресурсов можно почитать в документации Kubernetes ([ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/), [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)). Далее мы разберём частые случаи их использования на примере нашего приложения.

## Хранение конфигурационных файлов приложения в ConfigMap

Сейчас конфигурационный файл [nginx.conf]({{ page.base_url | append: page.examples_initial | append: "/.werf/nginx.conf" }}) копируется в образ прямо во время сборки. Из-за этого при каждом его изменении будут происходить пересборка образа и перезапуск Pod'ов. Также сейчас нет простой возможности шаблонизировать `nginx.conf`.

Эти проблемы решаются, если перенести `.werf/nginx.conf` в отдельный ConfigMap, для того, чтобы монтировать `nginx.conf` во время деплоя, а не добавлять файл при сборке:
{% include snippetcut_example path=".helm/templates/configmap-nginx.yaml" syntax="yaml" examples=page.examples %}
Теперь добавим этот ConfigMap в Deployment, примонтировав его как файл внутрь контейнера `frontend`:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="nginx_conf" examples=page.examples %}
Не забудем удалить более ненужный файл `.werf/nginx.conf`, а также команду копирования этого файла в образ во время сборки, после чего сборка образа `frontend` станет выглядеть так:
{% include snippetcut_example path="Dockerfile" syntax="dockerfile" snippet="frontend" examples=page.examples %}

## Перевыкат Deployment при изменении ConfigMap и Secret

По умолчанию, изменения в ConfigMap и Secret, примонтированных к Deployment, StatefulSet или DaemonSet, не приведут к перезапуску Pod'ов с новой конфигурацией. Чтобы Pod'ы обновились, их нужно аннотировать хеш-суммами всех используемых Pod'ом ConfigMap и Secret. Тогда при изменении ConfigMap и Secret аннотация изменится и Pod пересоздастся. Так выглядит аннотация с хеш-суммой ConfigMap с `nginx.conf`:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="nginx_conf_checksum" examples=page.examples %}
Для каждого подключенного ConfigMap и Secret потребуется отдельная аннотация.
>_Авторы самоучителя предпочитают вместо аннотаций с хеш-суммами использовать операторы вроде [stakater/Reloader](https://github.com/stakater/Reloader), т.к. они проще, гибче и удобнее в работе._

## Values

Использование helm-чарта для выката приложения даёт ряд преимуществ, например, возможность шаблонизировать манифесты. Ключевой встроенный объект шаблонизации — это Values, который предоставляет доступ к значениям, передаваемым в чарт из различных источников.

При использовании werf все данные передаваемые в chart можно условно разделить на:
- [Обычные пользовательские данные]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/values.html#обычные-пользовательские-данные): параметры из `values.yaml`-файлов и соответствующих опций. 
- [Пользовательские секреты]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/values.html#пользовательские-секреты): параметры из `secret-values.yaml`-файлов.
- [Сервисные данные]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/values.html#сервисные-данные): информация о проекте, релизе, собранных образах и т.д.

## Переиспользование конфигурации с Values и ConfigMap

Для упрощения конфигурации манифестов часто используемые параметры можно выносить в `.helm/values.yaml`:
{% include snippetcut_example path=".helm/values.yaml" syntax="yaml" snippet="mysql" examples=page.examples %}
... и далее использовать в манифесте следующим образом:
{% include snippetcut_example path=".helm/templates/database.yaml" syntax="yaml" snippet="volume_claim_templates" examples=page.examples %}
Это также позволит иметь разные значения для разных окружений, о чем более подробно будет рассказано в следующих разделах самоучителя.

Но особенно полезным будет перенос в `.helm/values.yaml` повторяющейся конфигурации, например переменных окружения для приложения, которые мы используем в нескольких местах сразу:
{% include snippetcut_example path=".helm/values.yaml" syntax="yaml" snippet="app" examples=page.examples %}
Теперь переменные окружения из `.Values.app.envs` можно либо подставлять в манифест в `env` контейнера, как было ранее, либо вынести эти переменные окружения в ConfigMap, который потом подключить в контейнер через `envFrom`.

Первый вариант проще, но вариант с ConfigMap удобнее при большом количестве общих переменных окружения, т. к. позволит избежать их дублирования между контроллерами. ConfigMap может выглядеть так:
{% include snippetcut_example path=".helm/templates/configmap-app-envs.yaml" syntax="yaml" examples=page.examples %}

Подключаем этот ConfigMap через `envFrom` в Deployment как набор переменных окружения:

{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="backend_conf" examples=page.examples %}

Не забываем про аннотации с хеш-суммами, которые приведут к пересозданию Pod'ов при изменении ConfigMap:

{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="configmap_app_envs_checksum" examples=page.examples %}

Аналогично этот ConfigMap подключается и в остальные ресурсы, в которых требуются эти переменные окружения.

## Использование конфиденциальных данных с Values и Secret

Чтобы начать работу с секретами, сначала требуется сгенерировать симметричный ключ шифрования командой `werf helm secret generate-secret-key`. Но так как мы заранее подготовили и зашифровали секреты, то ключ шифрования мы сгенерировали за вас. Он находится в репозитории в файле `.werf_secret_key` и начнёт использоваться автоматически.
>_При работе с реальными приложениями ключ НЕЛЬЗЯ хранить в репозитории. Вместо этого рекомендуется передавать его через переменную окружения `WERF_SECRET_KEY`, храня его в безопасном месте. Подробнее про работу с ключами шифрования [здесь]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/secrets.html#%D0%BA%D0%BB%D1%8E%D1%87-%D1%88%D0%B8%D1%84%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D1%8F)._

Конфигурация приложения содержит данные, которые не должны храниться в незашифрованном виде в репозитории, например, логин и пароль от базы данных. Вместо того чтобы хранить логин и пароль в открытом виде в манифесте контроллера, мы будем хранить их зашифрованными в `.helm/secret-values.yaml` вместе с другими секретными параметрами:
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

Также, для хранения и монтирования секретных конфигурационных файлов целиком можно использовать Secret. Выглядит это так же, как использование ConfigMap для монтирования несекретных файлов конфигурации, [описанное выше](#хранение-конфигурационных-файлов-приложения-в-configmap), только само содержание Secret должно храниться зашифрованным в файлах `.helm/secret/...` или в `.helm/secret-values.yaml`, подробнее [здесь]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/secrets.html#%D1%81%D0%B5%D0%BA%D1%80%D0%B5%D1%82%D0%BD%D1%8B%D0%B5-%D1%84%D0%B0%D0%B9%D0%BB%D1%8B).

Больше информации по работе с секретами в [документации]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/secrets.html).

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