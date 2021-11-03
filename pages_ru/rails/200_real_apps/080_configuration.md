---
title: Конфигурация и секреты
permalink: rails/200_real_apps/080_config.html
examples_initial: examples/rails/050_s3
examples: examples/rails/080_configuration
base_url: https://github.com/werf/werf-guides/blob/master/
description: |
  В этой главе мы покажем, как правильно использовать и хранить секретную и несекретную конфигурацию приложения. 

  В предыдущих главах конфигурация добавлялась прямо в контейнеры при сборке или использовалась как есть в переменных окружения контейнеров при выкате.
  
  Теперь для безопасности и гибкости конфигурация будет сохраняться в ConfigMap и Secret. А в дополнение к параметрам Helm-чарта (Values) и секретам werf будут продемонстрированы подходы параметризации и переиспользования конфигурации, а также хранения конфиденциальных данных вместе с кодом в Git-репозитории проекта.
---

## ConfigMap и Secret

В Kubernetes есть специальные типы ресурсов ConfigMap и Secret, которые предназначены для отделения конфигурации, зависящей от среды и окружения, от образов контейнеров. 

Оба ресурса позволяют хранить данные в парах «ключ-значение» и впоследствии использовать их как переменные окружения, аргументы командной строки или как файлы, примонтированные в выбранный контейнер.

ConfigMap предназначен для хранения неконфиденциальных данных, в отличие от Secret, который используется для хранения различных типов секретов.

Подробнее про эти типы ресурсов можно почитать в официальной документации Kubernetes ([ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/), [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)). Далее мы разберём частые случаи их использования на примере нашего приложения.

## Хранение конфигурационных файлов приложения в ConfigMap

Сейчас конфигурационный файл [nginx.conf]({{ page.base_url | append: page.examples_initial | append: "/.werf/nginx.conf" }}) копируется в образ прямо во время сборки. Из-за этого при каждом его изменении будут происходить пересборка образа и перезапуск Pod'ов. Также сейчас нет простой возможности шаблонизировать `nginx.conf`.

Эти проблемы решаются, если перенести `.werf/nginx.conf` в отдельный ConfigMap — для того, чтобы монтировать `nginx.conf` во время деплоя, а не добавлять файл при сборке:
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

Использование Helm-чарта для выката приложения даёт ряд преимуществ — например, возможность шаблонизировать манифесты. Ключевой встроенный объект шаблонизации — это `Values`, который предоставляет доступ к значениям, передаваемым в чарт из различных источников.

При использовании werf все данные, передаваемые в chart, можно условно разделить на несколько категорий:
- [Обычные пользовательские данные]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/values.html#обычные-пользовательские-данные): параметры из файлов `values.yaml` и соответствующих опций. 
- [Пользовательские секреты]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/values.html#пользовательские-секреты): параметры из файлов `secret-values.yaml`.
- [Сервисные данные]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/values.html#сервисные-данные): информация о проекте, релизе, собранных образах и т.д.

## Переиспользование конфигурации с Values и ConfigMap

Для упрощения конфигурации манифестов часто используемые параметры можно выносить в `.helm/values.yaml`:
{% include snippetcut_example path=".helm/values.yaml" syntax="yaml" snippet="mysql" examples=page.examples %}
... и далее использовать в манифесте следующим образом:
{% include snippetcut_example path=".helm/templates/database.yaml" syntax="yaml" snippet="volume_claim_templates" examples=page.examples %}
Это также позволит иметь разные значения для разных окружений, о чем более подробно будет рассказано в следующих разделах самоучителя.

Но особенно полезным будет перенос в `.helm/values.yaml` повторяющейся конфигурации, например, переменных окружения для приложения, которые мы используем в нескольких местах сразу:
{% include snippetcut_example path=".helm/values.yaml" syntax="yaml" snippet="app" examples=page.examples %}
Теперь переменные окружения из `.Values.app.envs` можно либо подставлять в манифест в `env` контейнера, как было сделано ранее, либо вынести эти переменные окружения в ConfigMap, который потом подключить в контейнер через `envFrom`.

Первый вариант проще, но вариант с ConfigMap удобнее при большом количестве общих переменных окружения, т.к. позволит избежать их дублирования между контроллерами. ConfigMap может выглядеть так:
{% include snippetcut_example path=".helm/templates/configmap-app-envs.yaml" syntax="yaml" examples=page.examples %}

Подключаем этот ConfigMap через `envFrom` в Deployment как набор переменных окружения:

{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="backend_conf" examples=page.examples %}

Не забываем про аннотации с хеш-суммами, которые приведут к пересозданию Pod'ов при изменении ConfigMap:

{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="configmap_app_envs_checksum" examples=page.examples %}

Аналогично этот ConfigMap подключается и в остальные ресурсы, в которых требуются те же переменные окружения.

## Использование конфиденциальных данных с Values и Secret

Чтобы начать работу с секретами, сначала требуется сгенерировать симметричный ключ шифрования, что можно сделать командой `werf helm secret generate-secret-key`. Однако, поскольку мы уже заранее подготовили и зашифровали секреты, то и ключ шифрования мы сгенерировали за вас. Он находится в репозитории в файле `.werf_secret_key` и начнёт использоваться автоматически.
>_При работе с реальными приложениями ключ НЕЛЬЗЯ хранить в репозитории. Вместо этого рекомендуется передавать его через переменную окружения `WERF_SECRET_KEY`, храня его в безопасном месте. Подробнее про работу с ключами шифрования можно прочитать в [документации werf]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/secrets.html#%D0%BA%D0%BB%D1%8E%D1%87-%D1%88%D0%B8%D1%84%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D1%8F)._

Конфигурация приложения содержит данные, которые не должны храниться в незашифрованном виде в репозитории — например, логин и пароль от базы данных. Вместо того, чтобы хранить логин и пароль в открытом виде в файле конфигурации, как мы делаем сейчас:
{% include snippetcut_example path="config/database.yml" syntax="yaml" examples=page.examples_initial %}
... мы будем хранить их зашифрованными в `.helm/secret-values.yaml` вместе с другими секретными параметрами:
{% include snippetcut_example path=".helm/secret-values.yaml" syntax="yaml" examples=page.examples %}

Чтобы увидеть расшифрованными секреты, хранящиеся в `.helm/secret-values.yaml`, можно выполнить следующую команду:
```bash
werf helm secret values decrypt .helm/secret-values.yaml
```

Ожидаемый результат:
```yaml
app:
  secretEnvs:
    S3_USERNAME: minioadmin
    S3_PASSWORD: minioadmin
    DB_USERNAME: root
    DB_PASSWORD: password
    SECRET_KEY_BASE: something
mysql:
  secretEnvs:
    MYSQL_ROOT_PASSWORD: password
minio:
  secretEnvs:
    MINIO_ROOT_USER: minioadmin
    MINIO_ROOT_PASSWORD: minioadmin
```

Теперь нужно передать секреты из `.helm/secret-values.yaml` обратно в конфигурационный файл приложения. Для этого сначала передадим их в ресурс Secret:
{% include snippetcut_example path=".helm/templates/secret-app-envs.yaml" syntax="yaml" examples=page.examples %}

... а затем примонтируем Secret в контейнеры как набор переменных окружения:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="backend_secret" examples=page.examples %}

Изменения в остальных файлах `.helm/templates` для вынесения секретной конфигурации в Secret аналогичны, поэтому приводить здесь листинги не будем. (На их содержимое можно посмотреть в [репозитории]({{ page.base_url | append: page.examples | append: "/.helm/templates/" }}).)

После передачи переменных окружения в контейнер остаётся только подставить их в конфигурационных файлах приложения:
{% include snippetcut_example path="config/database.yml" syntax="yaml" examples=page.examples %}
{% include snippetcut_example path="config/storage.yml" syntax="yaml" examples=page.examples %}
{% include snippetcut_example path="config/secrets.yml" syntax="yaml" examples=page.examples %}

Для хранения и монтирования секретных конфигурационных файлов целиком также можно использовать Secret. Выглядит это аналогично использованию ConfigMap для монтирования несекретных файлов конфигурации, [описанному выше](#хранение-конфигурационных-файлов-приложения-в-configmap), за тем лишь исключением, что само содержание Secret должно храниться зашифрованным в файлах `.helm/secret/...` или в `.helm/secret-values.yaml`. Подробнее об этом можно прочитать в [документации werf]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/secrets.html#%D1%81%D0%B5%D0%BA%D1%80%D0%B5%D1%82%D0%BD%D1%8B%D0%B5-%D1%84%D0%B0%D0%B9%D0%BB%D1%8B).

Больше информации по работе с секретами также доступно в [документации werf]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/secrets.html).

## Проверка работоспособности

Убедимся, что изменения в конфигурации не повлияли на работоспособность приложения:
```shell
werf converge --repo <ИМЯ ПОЛЬЗОВАТЕЛЯ DOCKER HUB>/werf-guide-app
```

Ожидаемый результат:
```
...
┌ ⛵ image backend
│ Use cache image for backend/dockerfile
│      name: .../werf-guide-app:8df39a540816b5851314934d932f0152a8fa409369a58dbfeccce4be-1632852281567
│        id: 0825950bfd70
│   created: 2021-09-28 21:04:41 +0300 MSK
│      size: 366.3 MiB
└ ⛵ image backend (7.19 seconds)

┌ ⛵ image frontend
│ Use cache image for frontend/dockerfile
│      name: .../werf-guide-app:18b3a25a107798890200664b7de037c998bef8ed9d9bbe4d66efaf53-1632852281351
│        id: 7d029b84655e
│   created: 2021-09-28 21:04:41 +0300 MSK
│      size: 9.5 MiB
└ ⛵ image frontend (7.60 seconds)

┌ Waiting for release resources to become ready
...
│ ┌ Status progress
│ │ DEPLOYMENT      REPLICAS  AVAILABLE  UP-TO-DATE
│ │ werf-guide-app  1/1       1          1
│ │ │    POD                         READY  RESTARTS  STATUS
│ │ ├──  guide-app-85bd8c68f5-w9nff  2/2    0         Running
│ │ └──  guide-app-865b6d68bc-bkt9l  2/2    0         Terminating
│ │ STATEFULSET  REPLICAS  READY  UP-TO-DATE
│ │ minio        1/1       1      1
│ │ │    POD  READY  RESTARTS  STATUS
│ │ └──  0    1/1    0         Running
│ │ mysql        1/1       1      1
│ │ │    POD  READY  RESTARTS  STATUS
│ │ └──  0    1/1    0         Running
│ │ JOB                        ACTIVE  DURATION  SUCCEEDED/FAILED
│ │ setup-and-migrate-db-rev2  0       24s       1/0
│ │ │    POD                        READY  RESTARTS  STATUS
│ │ └──  and-migrate-db-rev2-lhnq7  0/1    0         Completed
│ │ setup-minio-rev2           0       33s       0->1/0
│ │ │    POD               READY  RESTARTS  STATUS
│ │ └──  minio-rev2-wr7tw  0/1    0         Running  ->
│ │ Completed
│ └ Status progress
└ Waiting for release resources to become ready (33.05 seconds)

Release "werf-guide-app" has been upgraded. Happy Helming!
NAME: werf-guide-app
LAST DEPLOYED: Tue Sep 28 21:26:48 2021
NAMESPACE: werf-guide-app
STATUS: deployed
REVISION: 2
TEST SUITE: None
Running time 43.88 seconds
```

Проверим доступность приложения:
```shell
curl http://werf-guide-app/ping
```

Ожидаемый результат:
```
pong
```

Теперь секретная конфигурация приложения хранится безопасно, дублирование конфигурации минимизировано, а файлы конфигурации приложений могут шаблонизироваться и больше не требуют пересборок при изменениях.
