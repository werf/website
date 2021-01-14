---
title: Подключаем внешний PostgreSQL
permalink: nodejs/200_real_apps/80_database.html
layout: "wip"
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/templates/job-migrations.yaml
- migrations/1588019669425_001-labels.js
- app.js
- package.json
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении продвинутую работу с базой данных, включающую в себя вопросы выполнения миграций. 

В качестве базы данных возьмём PostgreSQL. Мы предполагаем, что база данных уже где-то реализована: вы либо подняли её на отдельном сервере, либо воспользовались Managed сервисом у своего cloud provider-а. Вопросы, связанные с самостоятельной установкой БД в кластер мы разберём в главе "Работа с инфраструктурой"

{% offtopic title="А что если БД в кластере?" %}
Если база данных задеплоена в кластер, то всё, что написано ниже, остаётся в силе, изменятся только то, откуда берутся реквизиты доступа. Вопрос написания Helm-чартов с базой данных в данный момент не описан в самоучителе, однако описан во множестве других учебников по Kubernetes.

Размещение базы данных в кластере вручную — комплексный вопрос, требующий организации хранилища для данных, реквизитов доступа и т.п. Для того, чтобы всё сделать корректно требуются продвинутые навыки.
{% endofftopic %}

## Подключение приложения к базе PostgreSQL

Наше javascript-приложение мы перепишем, оставив то же API, но заменив место хранения данных: вместо sqlite теперь будет использоваться PostgreSQL. Соответственно, изменятся и переменные окружения, которыми будет конфигурироваться приложение.

Мы не будем приводить полный код приложения — его можно [посмотреть в github](https://github.com/werf/werf-guides/blob/master/examples/nodejs/280_database/app.js).

Для подключения Node.js-приложения к PostgreSQL необходимо установить npm-пакет `pg` и сконфигурировать:

{% snippetcut name="app.js" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/280_database/app.js" %}
{% raw %}
```js
// Connection to Redis
const Pool = require('pg').Pool
let client;
try {
  client = new Pool({
    user: process.env.POSTGRESQL_USER,
    host: process.env.POSTGRESQL_HOST,
    database: process.env.POSTGRESQL_DB,
    password: process.env.POSTGRESQL_PASSWORD,
    port: process.env.POSTGRESQL_PORT,
  })
} catch (err) {
  console.error('connection to PSQL failed')
  console.error(err);
  process.exit(1);
}
```
{% endraw %}
{% endsnippetcut %}

Далее, внесём изменения в объект Deployment: пропишем новые переменные окружения и реализуем `initContainer`, который не даст запустить приложение, пока база данных не запустится.

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/.helm/templates/deployment.yaml" %}
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
      imagePullSecrets:
      - name: "registrysecret"
      initContainers:
      - name: wait-postgres
        image: "foobic/pg_isready:latest"
        command: ['/scripts/pg_isready.sh']
        env:
        - name: HOST
          value: "host.docker.internal"
        - name: PORT
          value: "5432"
        - name: DBNAME
          value: "postgres"
        - name: RETRIES
          value: "-1"
      containers:
      - name: basicapp
        command: ["node","/app/app.js"]
        image: {{ tuple "basicapp" . | werf_image }}
        workingDir: /app
        ports:
        - containerPort: 3000
          protocol: TCP
        env:
          - name: "POSTGRESQL_HOST"
            value: "host.docker.internal"
          - name: "POSTGRESQL_PORT"
            value: "5432"
          - name: "POSTGRESQL_DB"
            value: "postgres"
          - name: "POSTGRESQL_USER"
            value: "postgres"
          - name: "POSTGRESQL_PASSWORD"
            value: "mysecretpassword"
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="Что здесь написано?" %}
[Init Container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/)-ы запускаются перед основными контейнерами. Пока они не завершатся успешно — основные контейнеры не будут запущены.

Мы используем [образ `foobic/pg_isready`](https://hub.docker.com/r/foobic/pg_isready/dockerfile), который представляет из себя `alpine` с консольной утилитой `pg_isready`, позволяющей проверять готовность базы Postgre. Он будет пытаться подсоединиться к базе раз за разом.

В нашем примере в качестве `HOST` используется значение `host.docker.internal`. Это значение необходимо, если вы ведёте разработку в локальном kubernetes (например, minikube), а базу данных развернули прямо на хост-машине. См. подробнее документацию Docker [о том, как из контейнера добраться до хост-машины](https://docs.docker.com/docker-for-mac/networking/#use-cases-and-workarounds).

Соответственно, аналогично прописываем переменные окружения у контейнера `basicapp`.
{% endofftopic %}

Оставлять ключи подключения в helm-чарте — плохая практика. Чтобы иметь возможность конфигурировать разные значения для разных стендов (тестового, продакшн и т.п.) и обезопасить себя — воспользуйтесь подходами, описанными в главе "Конфигурирование инфраструктуры в виде кода".

<a name="migrations" />

## Выполнение миграций

Работа реальных приложений почти немыслима без выполнения миграций. С точки зрения Kubernetes миграции выполняются созданием объекта Job, который разово запускает Pod с необходимыми контейнерами. Запуск миграций мы пропишем после каждого деплоя приложения.

Для выполнения миграций в БД нами будет использоваться пакет `node-pg-migrate`. Установим его: `npm install node-pg-migrate`.

{% offtopic title="Как настраиваем node-pg-migrate?" %}

Запуск миграции поместим в `package.json`, чтобы его можно было вызывать его через `npm run migrate`:

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

{% snippetcut name=".helm/templates/job-migrations.yaml" url="#" %}
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

По аналогии с Deployment мы будем использовать `initContainer`.

{% snippetcut name=".helm/templates/migrations.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/080-database/.helm/templates/migrations.yaml" %}
{% raw %}
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: basicapp-migrations
  annotations:
    "helm.sh/hook": post-install, post-upgrade
    "helm.sh/weight": "5"
spec:
  activeDeadlineSeconds: 600
  template:
    metadata:
      name: basicapp-migrations
    spec:
      imagePullSecrets:
        - name: "registrysecret"
      restartPolicy: OnFailure
      initContainers:
        - name: wait-postgres
          image: "foobic/pg_isready:latest"
          command: ['/scripts/pg_isready.sh']
          env:
            - name: HOST
              value: "host.docker.internal"
            - name: PORT
              value: "5432"
            - name: DBNAME
              value: "postgres"
            - name: RETRIES
              value: "-1"
      containers:
        - name: init-tables
          image: {{ tuple "basicapp" . | werf_image }}
          command: ['node']
          args: ['node_modules/node-pg-migrate/bin/node-pg-migrate', 'up']
          workingDir: /app
          env:
            - name: "DATABASE_URL"
              value: "postgres://postgres:mysecretpassword@host.docker.internal:5432/postgres"
```
{% endraw %}
{% endsnippetcut %}

## Деплой и тестирование

После того, как изменения в код внесены — нужно закоммитить изменения в git, задеплоить в кластер с помощью `werf converge` и протестировать методы `/api/labels`. Если миграции не отработают — ни один из методов не будет работать, т.к. не будет хватать таблицы.

<div id="go-forth-button">
    <go-forth url="90_stateful.html" label="Сборка образа" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
