---
title: Подключение базы данных
permalink: gitlab_java_springboot/080_database.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/requirements.yaml
- .helm/values.yaml
- .helm/secret-values.yaml
- .gitlab-ci.yml
- pom.xml
- postgres-pv.yaml
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

{% snippetcut name=".helm/requirements.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/080-database/.helm/requirements.yaml" %}
{% raw %}
```yaml
dependencies:
- name: postgresql
  version: 8.0.0
  repository: https://kubernetes-charts.storage.googleapis.com/
  condition: postgresql.enabled
```
{% endraw %}
{% endsnippetcut %}

Для того, чтобы werf при деплое загрузила необходимые нам сабчарты, нужно прописать в `.gitlab-ci.yml` работу с зависимостями:

{% snippetcut name=".gitlab-ci.yml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/080-database/.gitlab-ci.yml" %}
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

{% snippetcut name=".helm/values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/080-database/.helm/values.yaml" %}
{% raw %}
```yaml
postgresql:
  enabled: true
<...>
  _postgresqlHost:
    _default: "werf-guided-project-staging-postgresql"
    production: "werf-guided-project-production-postgresql"
  postgresqlDatabase: guided-database
  postgresqlUsername: guide-username
  servicePort: 5432
  imageTag: "12"
```
{% endraw %}
{% endsnippetcut %}

Пароль от базы данных мы тоже конфигурируем, но храним его в секретных переменных. Для этого стоит использовать механизм секретных переменных. *Вопрос работы с секретными переменными рассматривался подробнее, когда мы [делали базовое приложение](020_basic/20_iac.html#secret-values-yaml).*

{% snippetcut name=".helm/secret-values.yaml (зашифрованный)" url="#" ignore-tests %}
{% raw %}
```yaml
postgresql:
  postgresqlPassword: 1000b925471a9491456633bf605d7d3f74c3d5028f2b1e605b9cf39ba33962a4374c51f78637b20ce7f7cd27ccae2a3b5bcf
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

{% snippetcut name=".helm/values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/080-database/.helm/values.yaml" %}
{% raw %}
```yaml
postgresql:
<...>
  persistence:
    enabled: true
    storageClass: "local-storage"
    accessModes:
    - ReadWriteOnce
    size: 8Gi
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
          - article-kube-node-2
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

Для подключения Spring приложения к PostgreSQL необходимо установить postgresql-driver 'postgresql' и jpa (Java persistence api) и сконфигурировать:

{% snippetcut name="pom.xml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/080-database/pom.xml" %}
{% raw %}
```xml
		<dependency>
			<groupId>org.postgresql</groupId>
			<artifactId>postgresql</artifactId>
			<scope>runtime</scope>
		</dependency>
```
{% endraw %}
{% endsnippetcut %}

Для подключения к базе данных нам, очевидно, нужно знать: хост, порт, имя базы данных, логин, пароль. В коде приложения используется несколько переменных окружения: `POSTGRESQL_HOST`, `POSTGRESQL_PORT`, `POSTGRESQL_DATABASE`, `POSTGRESQL_LOGIN`, `POSTGRESQL_PASSWORD`.

Настраиваем эти переменные окружения по аналогии с тем, как [настраивали Redis](070_redis.html).

{% offtopic title="Какие значения прописываются в переменные окружения?" %}
Будем **конфигурировать хост** через `values.yaml`:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/080-database/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        - name: POSTGRESQL_HOST
          value: {{ pluck .Values.global.env .Values.postgresql._postgresqlHost | first | default .Values.postgresql._postgresqlHost._default  | quote }}
```
{% endraw %}
{% endsnippetcut %}

**Конфигурируем логин и порт** через `values.yaml`, просто прописывая значения:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/080-database/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        - name: POSTGRESQL_LOGIN
          value: {{ .Values.postgresql.postgresqlUsername | quote }}
        - name: POSTGRESQL_PORT
          value: {{ .Values.postgresql.servicePort | quote }}
```
{% endraw %}
{% endsnippetcut %}

{% snippetcut name=".helm/values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/080-database/.helm/values.yaml" %}
{% raw %}
```yaml
postgresql:
<...>
  postgresqlUsername: guide-username
  servicePort: 5432
```
{% endraw %}
{% endsnippetcut %}

**Конфигурируем пароль** через `values.yaml`, просто прописывая значения:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/080-database/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        - name: POSTGRESQL_PASSWORD
          value: {{ .Values.postgresql.postgresqlPassword | quote }}
```
{% endraw %}
{% endsnippetcut %}

{% snippetcut name=".helm/secret-values.yaml" url="#" ignore-tests %}
{% raw %}
```yaml
postgresql:
  password:
    _default: 100067e35229a23c5070ad5407b7406a7d58d4e54ecfa7b58a1072bc6c34cd5d443e
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

Обычной практикой при использовании kubernetes является выполнение миграций в БД с помощью объекта Job, единоразового выполнения команды в контейнере с определенными параметрами.

Если вы используете JPA — этот механизм самостоятельно выполняет миграции. Подробно о том как это делается и настраивается прописано в [документации](https://docs.spring.io/spring-boot/docs/2.1.1.RELEASE/reference/html/howto-database-initialization.html). В ней рассматривается несколько инструментов для выполнения миграций и работы с базой данных: JPA, Spring Batch Database, Flyway и Liquibase. С точки зрения helm-чартов или сборки ничего не изменится относительно описанных в статье настроек.

<a name="database-fixtures" />

## Накатка фикстур

Накатка фикстур производится фреймворком самостоятельно используя вышеназванные инструменты для миграций. Для накатывания фикстур в зависимости от стенда можно передавать переменную окружения в приложение. Например, для Flyway это будет `SPRING_FLYWAY_LOCATIONS`. Для Spring Batch Database нужно присвоить Java-переменной `spring.batch.initialize-schema` значение переменной из environment  (`always` либо `never`). Для JPA это `validate`, `update`, `create`, `create-drop`. Реализации могут отличаться, но подход остаётся один: обязательно нужно передавать переменную в контейнер.

<div id="go-forth-button">
    <go-forth url="090_unittesting.html" label="Юнит-тесты и Линтеры" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
