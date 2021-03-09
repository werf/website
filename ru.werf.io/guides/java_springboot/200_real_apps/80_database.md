---
title: Подключаем Managed PostgreSQL
permalink: java_springboot/200_real_apps/80_database.html
layout: "development"
---


{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/templates/migrations.yaml
- .helm/requirements.yaml
- .helm/values.yaml
- .helm/secret-values.yaml
- queries.js
- postgres-pv.yaml
- .gitlab-ci.yml
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении продвинутую работу с базой данных, включающую в себя вопросы выполнения миграций. В качестве базы данных возьмём PostgreSQL.

Мы не будем вносить изменения в сборку — будем использовать образы с DockerHub, конфигурировать их и инфраструктуру.

<a name="kubeconfig" />

## Сконфигурировать PostgreSQL в Kubernetes

Подключение базы данных само по себе не является сложной задачей, особенно если она находится вне Kubernetes. Особенность подключения заключается только лишь в указании логина, пароля, хоста, порта и названия самой базы данных внутри инстанса с БД в переменных окружения к вашему приложению.

Мы же рассмотрим вариант с базой внутри Kubernetes-кластера. Удобство такого способа в том, что можно быстро развернуть полноценное окружение для нашего приложения. Однако предстоит решить непростую задачу конфигурирования базы данных, в том числе — вопрос конфигурирования места хранения данных базы.

{% offtopic title="А что с местом хранения?" %}
Kubernetes автоматически разворачивает приложение и переносит его с одного сервера на другой. А значит, «по умолчанию» приложение может сначала сохранить данные на одном сервере, а потом быть запущено на другом и оказаться без данных. Хранилища нужно конфигурировать и связывать с приложением; есть несколько способов это сделать.
{% endofftopic %}

Есть два способа подключить нашу БД: прописать Helm-чарт самостоятельно или подключить внешний чарт. Мы рассмотрим второй вариант: подключим PostgreSQL как внешний subchart.

Более того, мы рассмотрим вопрос того, каким образом можно обеспечить хранение данных для такой непостоянной сущности, как Pod в Kubernetes.

Для этого нужно:
1. Указать Postgres как зависимый сабчарт в `requirements.yaml`;
2. Сконфигурировать в werf работу с зависимостями;
3. Убедиться, что кластер настроен на работу с персистентным хранилищем;
3. Сконфигурировать подключённый сабчарт;
5. Убедиться, что создаётся Pod с PostgreSQL.

### Конфигурирование чарта

Пропишем Helm-зависимости:

{% snippetcut name=".helm/requirements.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/.helm/requirements.yaml" %}
{% raw %}
```yaml
dependencies:
- name: postgresql
  version: "8.9.2"
  repository: https://charts.bitnami.com/bitnami
  condition: postgresql.enabled
```
{% endraw %}
{% endsnippetcut %}

Для того, чтобы werf при деплое загрузила необходимые нам сабчарты, нужно прописать в `.gitlab-ci.yml` работу с зависимостями:

{% snippetcut name=".gitlab-ci.yml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/.gitlab-ci.yml" %}
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

Для того, чтобы подключённые сабчарты заработали, нужно указать настройки в `values.yaml`:

Изучив [документацию сабчарта](https://github.com/bitnami/charts/tree/master/bitnami/postgresql#parameters), можно увидеть, что название основной базы данных, пользователя, хоста и пароля и даже версии PostgreSQL задаются через следующие переменные:

{% snippetcut name=".helm/values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/.helm/values.yaml" %}
{% raw %}
```yaml
postgresql:
  postgresqlDatabase: guided-database
  postgresqlUsername: guide-username
  postgresqlHost: postgresql
  imageTag: "12"
  fullnameOverride: postgresql
  nameOverride: postgresql
```
{% endraw %}
{% endsnippetcut %}

Пароль от базы данных мы тоже конфигурируем, но храним его в секретных переменных. Для этого стоит использовать механизм секретных переменных. *Вопрос работы с секретными переменными рассматривался подробнее, когда мы [делали базовое приложение](020_basic/20_iac.html#secret-values-yaml).*

{% snippetcut name=".helm/secret-values.yaml (зашифрованный)" url="#" ignore-tests %}
{% raw %}
```yaml
postgresql:
  postgresqlPassword: 100087f3686b89a04c05a7e256ec51f8bbff3924a22ca0283f136a100896d23bf8b3
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="А где база будет хранить свои данные?" %}
Наверное, один из самых главных вопросов при разворачивании базы данных - это вопрос о месте её хранения. Каким образом мы можем сохранять данные в Kubernetes, если контейнер по определению является чем-то непостоянным?

В качестве хранилища данных в Kubernetes может быть настроен специальный [storage class](https://kubernetes.io/docs/concepts/storage/storage-classes/), который позвляет динамически создавать хранилища с помощью внешних ресурсов. Мы не будем подробно расписывать способы настройки этой сущности, т.к. её конфигурация может отличаться в зависимости от версии вашего кластера и места развертывания(на своих ресурсах или в облаке). Предлагаем лучше прокосультироваться по этому вопросу с людьми, поддерживающими ваш K8s-кластер.

Хорошим и простым примером `storage-class` является `local-storage`, который и будет использоваться в качестве типа хранилища для нашей базы данных.

`local-storage` позволяет хранить данные непосредственно на диске самих узлов кластера.
{% endofftopic %}

Далее мы указываем настройки `persistence`, с помощью которых настроим хранилище для БД:

{% snippetcut name=".helm/values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/.helm/values.yaml" %}
{% raw %}
```yaml
postgresql:
<...>
  persistence:
    storageClass: "local-storage"
    accessModes:
    - ReadWriteOnce
    size: 8Gi
    annotations: {}
```
{% endraw %}
{% endsnippetcut %}
На данный момент, возможно, осталась непонятной только одна настройка:
* `accessModes` - эта настройка определят тип доступа к нашему хранилищу. Значение `ReadWriteOnce` сообщает о том, что к хранилищу единовременно может обращаться только один Pod.

### Хранилище

Для того, чтобы описанный выше чарт смог выкатиться, описанное хранилище нужно создать. Для этого потребуется создать сущность [PersitentVolume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/). Конфигурацию объекта `PersistentVolume` не получится положить в репозиторий с исходным кодом приложения — его нужно создать вручную через `kubectl`.

Конфигурация объекта выглядит _примерно_ следующим образом:

{% snippetcut name="postgres-pv.yaml" url="#" ignore-tests %}
{% raw %}
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  finalizers:
  - kubernetes.io/pv-protection
  name: posgresql-data
spec:
  accessModes:
  - ReadWriteOnce
  capacity:
    storage: 8Gi
  local:
    path: /mnt/guided-postgresql-stage/posgresql-data-0
  nodeAffinity:
    required:
      nodeSelectorTerms:
      - matchExpressions:
        - key: kubernetes.io/hostname
          operator: In
          values:
            - kube-node-2
  persistentVolumeReclaimPolicy: Retain
  storageClassName: local-storage
  volumeMode: Filesystem
```
{% endraw %}
{% endsnippetcut %}

Мы не будем подробно останавливаться на каждом параметре этой сущности, т.к. сама её настройка достойна отдельной главы. Для глубокого понимания вы можете обратиться к [официальной документации Kubernetes](https://kubernetes.io/docs/concepts/storage/).

Для того, чтобы PersistentVolume корректно создал хранилище, нужно правильно указать, **на каком узле** и **в какой директории** будет наш localstorage.

В атрибуте `nodeAffinity` указываем нужный узел (в примере ниже этот узел называется `article-kube-node-2`, а посмотреть весь список узлов можно через `kubectl get nodes`):

{% snippetcut name="postgres-pv.yaml" url="#" ignore-tests %}
{% raw %}
```yaml
  nodeAffinity:
    required:
      nodeSelectorTerms:
      - matchExpressions:
        - key: kubernetes.io/hostname
          operator: In
          values:
          - article-kube-node-2
```
{% endraw %}
{% endsnippetcut %}

На указанном узле требуется **вручную создать директорию** и прописать её в конфиге PersistentVolume:

{% snippetcut name="postgres-pv.yaml" url="#" ignore-tests %}
{% raw %}
```yaml
  local:
    path: /mnt/guided-postgresql-stage/posgresql-data-0
```
{% endraw %}
{% endsnippetcut %}

После корректировки конфигурации — его нужно **применить к каждому namespace'у вручную**. Создайте файл `postgres-pv.yaml` и примените его к каждому окружению:

```shell
kubectl -n werf-guided-project-production apply -f postgres-pv.yaml
kubectl -n werf-guided-project-staging apply -f postgres-pv.yaml
```

{% offtopic title="Но каким образом Pod с нашей базой поймет, что ему нужно использовать именно этот PersistentVolume?" %}

В данном случае сабчарт, который мы указали, автоматически создаст еще одну сущность под названием [PersitentVolumeClaim](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims) (за его создание отвечают эти [cтроки](https://github.com/bitnami/charts/blob/master/bitnami/postgresql/templates/statefulset.yaml#L442-L460)).

Эта сущность является не чем иным, как прослойкой между хранилищем и Pod'ом с приложением.

Можно подключать её в качестве `volume` прямо в Pod:
{% snippetcut name="postgres-pv.yaml" url="#" ignore-tests %}
{% raw %}
```yaml
  volumes:
    - name: data
      persistentVolumeClaim:
        claimName: postgres-data
```
{% endraw %}
{% endsnippetcut %}
… и затем монитровать в нужное место в контейнере:
{% snippetcut name="postgres-pv.yaml" url="#" ignore-tests %}
{% raw %}
```yaml
      volumeMounts:
        - mountPath: "/var/lib/postgres"
          name: data
```
{% endraw %}
{% endsnippetcut %}
Хороший пример того, как это можно сделать, со всеми подробностями доступен в [документации Kubernetes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-persistent-volume-storage/).

Каким образом PersitentVolumeClaim подключается к PersitentVolume? В Kubernetes работает механизм [binding](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#binding), который позволяет подобрать нашему PVC любой PV, удовлетворяющий его по указанным параметрам (размер, тип доступа и т.д.). Указав одинаково эти параметры в PV и в `values.yaml`, мы гарантируем, что хранилище будет подключено к БД.

Есть и более очевидный способ соединить PV: добавить в него настройку, в которой указывается имя PVC и пространства имён, в котором он находится (имя namespace'а легко взять из переменной, которую генерирует `werf` при деплое):
{% snippetcut name="postgres-pv.yaml" url="#" ignore-tests %}
{% raw %}
```yaml
  claimRef:
    namespace: {{ .Values.global.namespace }}
    name: postgres-data
```
{% endraw %}
{% endsnippetcut %}
{% endofftopic %}

Обратите внимание: вы не сможете просто так удалить PersistentVolume из-за встроенной защиты. Если вы выполните команду:

```shell
kubectl -n werf-guided-project-staging delete pv posgresql-data
```

… PV как будто «зависнет» в процессе удаления. Это корректное поведение, вызванное `pv-protection`.

{% offtopic title="Почему так?" %}

Если вы посмотрите описание PV:

```shell
kubectl -n werf-guided-project-production get pv posgresql-data
```

… то увидите атрибут `finalizers`:

```yaml
  finalizers:
  - kubernetes.io/pv-protection
```

Они защищают данные от случайного удаления. Если же вы настаиваете — потребуется отредактировать PV, удалив из манифеста конфигурации строки с `pv-protection`, описанные выше:

```shell
kubectl -n werf-guided-project-production edit pv posgresql-data
```
{% endofftopic %}

<a name="appconnect" />

## Подключение приложения к базе PostgreSQL

Для подключения Node.js-приложения к PostgreSQL необходимо установить npm-пакет `pg` и сконфигурировать:

{% snippetcut name="queries.js" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/queries.js" %}
{% raw %}
```js
const pg = require("pg");
const pgconnectionString = "postgresql://" +process.env.POSTGRESQL_LOGIN + ":" +
                                            process.env.POSTGRESQL_PASSWORD + "@" +
                                            process.env.POSTGRESQL_HOST + ":" +
                                            process.env.POSTGRESQL_PORT + "/" +
                                            process.env.POSTGRESQL_DATABASE || "postgresql://127.0.0.1/postgres";
// Postgres connect
const pool = new pg.Pool({
  connectionString: pgconnectionString,
});
```
{% endraw %}
{% endsnippetcut %}

Для подключения к базе данных нам, очевидно, нужно знать: хост, порт, имя базы данных, логин, пароль. В коде приложения используется несколько переменных окружения: `POSTGRESQL_HOST`, `POSTGRESQL_PORT`, `POSTGRESQL_DATABASE`, `POSTGRESQL_LOGIN`, `POSTGRESQL_PASSWORD`.

Настраиваем эти переменные окружения по аналогии с тем, как [настраивали Redis](070_redis.html).

{% offtopic title="Какие значения прописываются в переменные окружения?" %}
Будем **конфигурировать хост** через `values.yaml`:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        - name: POSTGRESQL_HOST
          value: {{ pluck .Values.global.env .Values.app.postgresql.host | first | default .Values.app.postgresql.host._default | quote }}
```
{% endraw %}
{% endsnippetcut %}

**Конфигурируем логин и порт** через `values.yaml`, просто прописывая значения:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        - name: POSTGRESQL_LOGIN
          value: {{ pluck .Values.global.env .Values.app.postgresql.login | first | default .Values.app.postgresql.login._default | quote }}
<...>
        - name: POSTGRESQL_PORT
          value: {{ pluck .Values.global.env .Values.app.postgresql.port | first | default .Values.app.postgresql.port._default | quote }}
```
{% endraw %}
{% endsnippetcut %}

{% snippetcut name=".helm/values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/.helm/values.yaml" %}
{% raw %}
```yaml
  postgresql:
    login:
      _default: guide-username
<...>
    port:
      _default: 5432
```
{% endraw %}
{% endsnippetcut %}

**Конфигурируем пароль** через `values.yaml`, просто прописывая значения:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        - name: POSTGRESQL_PASSWORD
          value: {{ pluck .Values.global.env .Values.app.postgresql.password | first | default .Values.app.postgresql.password._default | quote }}
```
{% endraw %}
{% endsnippetcut %}

{% snippetcut name=".helm/secret-values.yaml" url="#" ignore-tests %}
{% raw %}
```yaml
postgresql:
  postgresqlPassword: 100087f3686b89a04c05a7e256ec51f8bbff3924a22ca0283f136a100896d23bf8b3
```
{% endraw %}
{% endsnippetcut %}

{% endofftopic %}


{% offtopic title="Безопасно ли хранить пароли в переменных к нашему приложению?" %}

Мы выбрали самый быстрый и простой способ того, как можно более-менее безопасно хранить данные в репозитории и передавать их в приложение. Но важно помнить, что любой человек с доступом к приложению в кластере (в особенности, к самой сущности Pod'а) сможет получить любой пароль, просто выполнив команду `env` внутри запущенного контейнера.

Можно избежать этого, не выдавая доступы на исполнение команд внутри контейнеров кому-либо. Кроме того, можно собирать свои собственные образы с нуля, убирая из них небезопасные команды, а переменные помещать в [сущность Secret](https://kubernetes.io/docs/concepts/configuration/secret/#using-secrets-as-environment-variables), после чего запрещать доступ к ним всем кроме доверенных лиц.

Но на данный момент мы настоятельно не рекомендуем использовать сущности самого Kubernetes для хранения любой секретной информации. Вместо этого попробуйте использовать инструменты, созданные специально для таких задач — например, [Vault](https://www.vaultproject.io/). Интегрировав клиент Vault прямо в ваше приложение, можно получать любые настройки непосредственно при запуске.

{% endofftopic %}
<a name="migrations" />

## Выполнение миграций

Работа реальных приложений почти немыслима без выполнения миграций. С точки зрения Kubernetes миграции выполняются созданием объекта Job, который разово запускает Pod с необходимыми контейнерами. Запуск миграций мы пропишем после каждого деплоя приложения.

Для выполнения миграций в БД будет использоваться пакет `node-pg-migrate`. На его примере и рассмотрим выполнение миграций.

{% offtopic title="Как настраиваем node-pg-migrate?" %}

Запуск миграции поместим в `package.json`, чтобы его можно было вызывать с помощью скрипта в `npm`:
```json
   "migrate": "node-pg-migrate"
```
Сама конфигурация миграций находится в отдельной директории `migrations`, которую мы создали на уровне исходного кода приложения:
```
node
├── migrations
│   ├── 1588019669425_001-users.js
│   └── 1588172704904_add-avatar-status.js
├── src
├── package.json
...
```
{% endofftopic %}

Далее необходимо добавить запуск миграций непосредственно в Kubernetes.

Как уже упоминалось, для этого будет создан ресурс Job в Kubernetes. Он обеспечивает единоразовый запуск Pod'а с необходимыми нам контейнерами; его предназначение — выполнение конечной функции, после чего он завершится. Повторно он может быть запущен только при последующих деплоях.

{% offtopic title="Как конфигурируем сам Job?" %}

Подробнее о конфигурировании объекта Job можно почитать в [документации Kubernetes](https://kubernetes.io/docs/concepts/workloads/controllers/job/).

Также мы воспользуемся аннотациями Helm [`helm.sh/hook` и `helm.sh/weight`](https://helm.sh/docs/topics/charts_hooks/), чтобы Job выполнялся после того, как применится новая конфигурация.

{% snippetcut name=".helm/templates/migrations.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/.helm/templates/migrations.yaml" %}
{% raw %}
```yaml
    "helm.sh/hook": post-install,post-upgrade
    "helm.sh/weight": "5"
```
{% endraw %}
{% endsnippetcut %}

Вопросы по настройке Job могут возникнуть уже на этапе, когда мы видим блок `annotations`. Он содержит в себе настройки для Helm, которые определяют, когда именно нужно запускать Job (подробнее про них можно узнать в [документации Helm](https://v2.helm.sh/docs/charts_hooks/)).

В данном случае первой аннотацией мы указываем, что Job нужно запускать только после того, как все объекты чарта будут загружены и запущены в Kubernetes.

А вторая аннотация отвечает за порядок запуска этого Job. Например, если мы имеем в нашем чарте несколько Job разного назначения и не хотим запускать их единовременно, а только по порядку, можно указать для них веса.

{% endofftopic %}

Так как состояние кластера постоянно меняется, нельзя быть уверенными, что на момент запуска миграций база работает и доступна. Поэтому в Job добавляется `initContainer`, который не даёт запуститься скрипту миграции, пока не станет доступна база данных:

{% snippetcut name=".helm/templates/migrations.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/.helm/templates/migrations.yaml" %}
{% raw %}
```yaml
      initContainers:
      - name: wait-postgres
        image: alpine:3.6
        command: ['/bin/sh', '-c', 'while ! getent ahostsv4 {{ pluck .Values.global.env .Values.app.postgresql.host | first | default .Values.app.postgresql.host._default }}; do sleep 1; done']
```
{% endraw %}
{% endsnippetcut %}

Остался непосредственный запуск миграции. В нём мы используем тот же самый образ, что и в ресурсе Deployment у приложения:

{% snippetcut name=".helm/templates/migrations.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/.helm/templates/migrations.yaml" %}
{% raw %}
```yaml
      containers:
      - name: init-tables
{{ tuple "basicapp" . | include "werf_container_image" | indent 8 }}
        command: ['node']
        args: ['node_modules/node-pg-migrate/bin/node-pg-migrate', 'up']
        workingDir: /app
        env:
        - name: DATABASE_URL
          value: "postgres://{{ pluck .Values.global.env .Values.app.postgresql.login | first | default .Values.app.postgresql.login._default }}:{{ pluck .Values.global.env .Values.app.postgresql.password | first | default .Values.app.postgresql.password._default }}@{{ pluck .Values.global.env .Values.app.postgresql.host | first | default .Values.app.postgresql.host._default }}:{{ pluck .Values.global.env .Values.app.postgresql.port | first | default .Values.app.postgresql.port._default }}/{{ pluck .Values.global.env .Values.app.postgresql.db | first | default .Values.app.postgresql.db._default }}"
```
{% endraw %}
{% endsnippetcut %}

<div id="go-forth-button">
    <go-forth url="90_stateful.html" label="Сборка образа" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
