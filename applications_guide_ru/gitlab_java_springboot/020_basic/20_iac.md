---
title: Конфигурирование IaC
sidebar: applications_guide
guide_code: gitlab_java_springboot
permalink: gitlab_java_springboot/020_basic/20_iac.html
toc: false
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/templates/ingress.yaml
- .helm/templates/service.yaml
- .helm/values.yaml
- .helm/secret-values.yaml
{% endfilesused %}

Для того, чтобы приложение заработало в Kubernetes, необходимо описать инфраструктуру приложения как код ([IaC](https://ru.wikipedia.org/wiki/%D0%98%D0%BD%D1%84%D1%80%D0%B0%D1%81%D1%82%D1%80%D1%83%D0%BA%D1%82%D1%83%D1%80%D0%B0_%D0%BA%D0%B0%D0%BA_%D0%BA%D0%BE%D0%B4)). В нашем случае потребуются следующие объекты Kubernetes: Pod, Service и Ingress.

Конфигурацию для Kubernetes нужно шаблонизировать. Один из популярных инструментов для такой шаблонизации — это [Helm](https://helm.sh/), его движок встроен в werf. Помимо этого, werf предоставляет возможности работы с секретными значениями, а также дополнительные Go-шаблоны для интеграции собранных образов.

В этой главе мы научимся описывать Helm-шаблоны, используя возможности werf, а также освоим встроенные инструменты отладки.


## Составление конфигов инфраструктуры

На сегодняшний день Helm — один из самых удобных (и самых распространённых) способов, которым можно описать свой деплой в Kubernetes. Он позволяет устанавливать готовые чарты с приложениями прямо из репозитория: введя одну команду, можно развернуть в своем кластере готовый Redis, PostgreSQL, RabbitMQ… Кроме того, Helm можно использовать для разработки собственных чартов, применяя удобный синтаксис для шаблонизации выката ваших приложений.

По этим причинам он был встроен в werf для решения соответствующих задач.

{% offtopic title="Что делать, если вы не работали с Helm?" %}

Не будем вдаваться в подробности [разработки YAML-манифестов с помощью Helm для Kubernetes](https://habr.com/ru/company/flant/blog/423239/). Если у вас есть вопросы о том, как именно описываются объекты Kubernetes, советуем посетить страницы документации Kubernetes о [концепциях](https://kubernetes.io/ru/docs/concepts/) и документацию Helm по [разработке шаблонов](https://helm.sh/docs/chart_template_guide/).

В первое время работа с Helm и конфигурацией для Kubernetes может быть очень сложной из-за нелепых мелочей вроде опечаток и пропущенных пробелов. Если вы только начали осваивать эти технологии — постарайтесь найти наставника, который поможет преодолеть эти сложности и посмотрит на ваши исходники сторонним взглядом.

В случае затруднений убедитесь, что вы:

- понимаете, как работает [indent](https://helm.sh/docs/chart_template_guide/function_list/#indent);
- понимаете, что такое [конструкция tuple](https://helm.sh/docs/chart_template_guide/control_structures/);
- понимаете, как Helm работает с хэш-массивами;
- очень внимательно следите за пробелами в YAML.

{% endofftopic %}

Итак, для работы рассматриваемого приложения в среде Kubernetes понадобится:

- описать сущности Deployment (он породит в кластере Pod) и Service;
- направить трафик на приложение, донастроив роутинг в кластере с помощью сущности Ingress;
- не забыть создать отдельную сущность Secret, которая позволит Kubernetes скачивать собранные образа из Registry.

### Создание Pod'а

Для того, чтобы в кластере появился Pod с нашим приложением, мы создадим объект Deployment. У создаваемого Pod будет один контейнер — `basicapp`. Укажем, **как этот контейнер будет запускаться**.

Здесь и далее будут показаны только фрагменты файлов. Если вам не знаком синтаксис Kubernetes-объектов и вы не можете дополнить приведённые сниппеты самостоятельно — обязательно сверяйтесь с файлами в [репозитории](https://github.com/werf/werf-guides/tree/master/examples/gitlab_java_springboot).

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/020-basic/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
      containers:
      - name: basicapp
        command: ["java"]
        args: ["-jar", "/app/demo.jar", "$JAVA_OPT"]
{{ tuple "basicapp" . | include "werf_container_image" | indent 8 }}
```
{% endraw %}
{% endsnippetcut %}

Обратите внимание на вызов [`werf_container_image`]({{ site.docsurl }}/v1.1-stable/documentation/reference/deploy_process/deploy_into_kubernetes.html#werf_container_image). Данная функция генерирует ключи `image` и `imagePullPolicy` со значениями, необходимыми для соответствующего контейнера Pod'а, что позволяет гарантировать перевыкат контейнера тогда, когда это нужно.

{% offtopic title="А в чём проблема?" %}
Kubernetes не знает ничего об изменении контейнера: он действует на основании описания объектов и сам выкачивает образы из Registry. Поэтому Kubernetes'у нужно в явном виде сообщать, что делать.

werf складывает собранные образы в Registry с разными именами — в зависимости от выбранной стратегии тегирования и деплоя (подробнее это разобрано в главе про CI). Как следствие, в описание контейнера нужно пробрасывать правильный путь до образа, а также дополнительные аннотации, связанные со стратегией деплоя.

Подробнее - можно посмотреть в [документации]({{ site.docsurl }}/v1.1-stable/documentation/reference/deploy_process/deploy_into_kubernetes.html#werf_container_image).
{% endofftopic %}

Для корректной работы приложения ему нужно узнать **переменные окружения**.

Например, для Java это `JAVA_OPT` - различные опции с которыми будет запускаться java. Аналогично в приложение может передаваться, например, пароль к базе данных. Добавим его также (хотя он и не используется в приложении) для иллюстрации механизма.

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/020-basic/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        env:
        - name: JAVA_OPT
          value: "--debug"
        - name: DBPASS
          value: "mysuperdbpassword"
<...>
{{ tuple "basicapp" . | include "werf_container_env" | indent 8 }}
```
{% endraw %}
{% endsnippetcut %}

Мы задали значение для `DBPASS` в явном виде — и это абсолютно не безопасный путь для хранения таких критичных данных. Мы разберём более правильный путь ниже, в главе "Разное поведение в разных окружениях".

Обратите также внимание на функцию [`werf_container_env`]({{ site.docsurl }}/v1.1-stable/documentation/reference/deploy_process/deploy_into_kubernetes.html#werf_container_env): с её помощью werf вставляет в описание объекта служебные переменые окружения.

<a name="helm-values-yaml" />

{% offtopic title="Как динамически подставлять в переменные окружения нужные значения?" %}

Helm — шаблонизатор, который поддерживает множество инструментов для подстановки значений. Один из центральных способов — подставлять значения из файла `values.yaml`. Наша конструкция могла бы иметь вид:

{% snippetcut name=".helm/templates/deployment.yaml" url="#" ignore-tests %}
{% raw %}
```yaml
      env:
<...>
      - name: "JAVA_OPT"
        value: "{{ .Values.app.java_opt}}"
```
{% endraw %}
{% endsnippetcut %}

… или даже более сложный — для того, чтобы значение основывалось на текущем окружении:

{% snippetcut name=".helm/templates/deployment.yaml" url="#" ignore-tests %}
{% raw %}
```yaml
      env:
<...>
      - name: "JAVA_OPT"
        value: "{{ pluck .Values.global.env .Values.app.java_opt | first | default .Values.app.java_opt._default }}"
```
{% endraw %}
{% endsnippetcut %}

{% snippetcut name=".helm/values.yaml" url="#" ignore-tests %}
{% raw %}
```yaml
app:
  java_opt:
    _default: "--debug"
    production: ""
    testing: "--debug"
```
{% endraw %}
{% endsnippetcut %}

{% endofftopic %}


При запуске приложения в Kubernetes **логи необходимо отправлять в stdout и stderr** — это нужно для простого сбора логов, например, через `filebeat`, а также для того, чтобы не разрастались Docker-образы запущенных приложений.

Spring-framework уже автоматически предоставляет логи в stdout. Однако мы можем переопределить уровень логирования в application.properties при необходимости. Подробнее - в [документации](https://docs.spring.io/spring-boot/docs/2.1.1.RELEASE/reference/html/boot-features-logging.html).

### Доступность Pod'а

Для того, чтобы запросы извне попали к нам в приложение, нужно открыть порт у Pod'а, создать объект Service и привязать его к Pod'у, а также создать объект Ingress.

{% offtopic title="Что за объект Ingress и как он связан с балансировщиком?" %}
Возможна коллизия терминов:

* Есть Ingress в смысле [NGINX Ingress Controller](https://github.com/kubernetes/ingress-nginx), который работает в кластере и принимает входящие извне запросы.
* А ещё есть [объект Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/), который фактически описывает настройки для NGINX Ingress Controller.

В статьях и бытовой речи оба этих термина зачастую называют «Ingress», так что догадываться нужно по контексту.
{% endofftopic %}

Наше приложение работает на стандартном порту `8080` — **откроем порт Pod'у**:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/020-basic/.helm/templates/deployment.yaml" %}
```yaml
        ports:
        - containerPort: 8080
          protocol: TCP
```
{% endsnippetcut %}

Затем **пропишем Service**, чтобы к Pod'у могли обращаться другие приложения кластера:

{% snippetcut name=".helm/templates/service.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/020-basic/.helm/templates/service.yaml" %}
{% raw %}
```yaml
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .Chart.Name }}
spec:
  selector:
    app: {{ .Chart.Name }}
  ports:
  - name: http
    port: 8080
    protocol: TCP
```
{% endraw %}
{% endsnippetcut %}

Обратите внимание на поле `selector` у Service: он должен совпадать с аналогичным полем у Deployment. Ошибки в этой части — самая частая проблема с настройкой маршрута до приложения.

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/020-basic/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Chart.Name }}
spec:
  selector:
    matchLabels:
      app: {{ .Chart.Name }}
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="Как убедиться, что выше всё сделано правильно?" %}

Мы будем деплоить приложение в Kubernetes позже, но если после деплоя у вас возникли проблемы в этом месте, то вернитесь сюда и проведите проверку, описанную ниже.

Попробуйте получить `endpoint` сервиса в нужном вам окружении. Если в нем будет фигурировать IP Pod'а — всё настроено правильно. А если нет, то проверьте еще раз, совпадают ли поля `selector` у Service и Deployment.

Название endpoint'а совпадает с названием сервиса, т.е. пример нужной команды:

{% raw %}
`kubectl -n <название окружения> get ep {{ .Chart.Name }}`
{% endraw %}

{% endofftopic %}

После этого можно настраивать **роутинг на Ingress**. Укажем, на какой домен, путь, сервис и порт направлять запросы:

{% snippetcut name=".helm/templates/ingress.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/020-basic/.helm/templates/ingress.yaml" %}
{% raw %}
```yaml
  rules:
  - host: mydomain.io
    http:
      paths:
      - path: /
        backend:
          serviceName: {{ .Chart.Name }}
          servicePort: 8080
```
{% endraw %}
{% endsnippetcut %}

### Разное поведение в разных окружениях

Некоторые настройки хочется видеть разными в разных окружениях. К примеру, домен, на котором будет открываться приложение, должен быть либо staging.mydomain.io, либо mydomain.io — в зависимости от того, куда мы задеплоились.

В werf для этого существует три механики:

1. Подстановка значений из `values.yaml` по аналогии с Helm.
2. Проброс значений через аргумент `--set` при работе в CLI-режиме, по аналогии с Helm.
3. Подстановка секретных значений из `secret-values.yaml`.

**Вариант с `values.yaml`** рассматривался ранее в главе ["Создание Pod'а"](#helm-values-yaml).

Второй вариант подразумевает **задание переменных через CLI**. Например, в случае выполнения команды `werf deploy --set "global.ci_url=mydomain.io"` в YAML'ах будет доступно значение {% raw %}`{{ .Values.global.ci_url }}`{% endraw %}.

Этот вариант удобен для проброса, например, имени домена для каждого окружения:

{% snippetcut name=".helm/templates/ingress.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/029-basic/.helm/templates/ingress.yaml" %}
{% raw %}
```yaml
  rules:
  - host: {{ .Values.global.ci_url }}
```
{% endraw %}
{% endsnippetcut %}

<a name="secret-values-yaml" />Отдельная проблема — **хранение и задание секретных переменных**, например, учётных данных аутентификации для сторонних сервисов, API-ключей и т.п.

Так как werf рассматривает Git как единственный источник правды, правильно хранить секретные переменные там же. Чтобы делать это корректно, мы [храним данные в шифрованном виде]({{ site.docsurl }}/v1.1-stable/documentation/reference/deploy_process/working_with_secrets.html). Подстановка значений из этого файла происходит при рендере шаблона, который также запускается при деплое.

Чтобы воспользоваться секретными переменными:

* [сгенерируйте ключ]({{ site.docsurl }}/v1.1-stable/documentation/cli/management/helm/secret/generate_secret_key.html) (`werf helm secret generate-secret-key`);
* определите ключ в переменных окружения для приложения, в текущей сессии консоли (например, `export WERF_SECRET_KEY=634f76ead513e5959d0e03a992372b8e`);
* пропишите полученный ключ в `Variables` для вашего репозитория в GitLab (раздел `Settings` → `CI/CD`), название переменной должно быть­`WERF_SECRET_KEY`:

![](/applications_guide_ru/images/applications-guide/020-werf-secret-key-in-gitlab.png)

После этого можно задать секретную переменную, например `DBPASS`. Зайдите в режим редактирования секретных значений:

```bash
$ werf helm secret values edit .helm/secret-values.yaml
```

Откроется консольный текстовый редактор с данными в расшифованном виде:

{% snippetcut name=".helm/secret-values.yaml в расшифрованном виде" url="#" ignore-tests %}
```yaml
app:
  password:
    _default: my-secret-password
    production: my-super-secret-password
```
{% endsnippetcut %}

После сохранения значения в файле зашифруются и примут примерно такой вид:

{% snippetcut name=".helm/secret-values.yaml в зашифрованном виде" url="#" ignore-tests %}
```yaml
app:
  password:
    _default: 10006755d101c5243fc400ababd7358689a921c19ee7e96a95f0ab82d46e4424573ab50ba666fcf5ce5e5dbd2b696c7706cf
    production: 1000bcd51061ebd1b2c2990041d30783be607b3a0aec8890c098f17bc96dc43e93765219651d743c7a37fb7361c10b703c7b
```
{% endsnippetcut %}

<a name="registryaccess" />

### Доступ кластера к Registry

Для того, чтобы кластер имел доступ к собранным образам, необходимо создать ключ доступа к Registry и прописать его в кластер. Этот ключ мы назовём `registrysecret`.

Сперва нужно создать API-ключ в GitLab: зайдите в настройки пользователя (`Settings`) и в разделе `Personal Access Tokens` создайте API-ключ с правами на `read_registry`. Корректнее всего создать отдельного служебного пользователя, чтобы не завязываться на персональный аккаунт.

Полученный ключ должен быть прописан в **каждом** пространстве имён в Kubernetes, куда осуществляется деплой, в виде объекта Secret. Сделать это можно, выполнив на master-узле кластера команду:

```bash
kubectl create secret docker-registry registrysecret -n <namespace> --docker-server=<registry_domain> --docker-username=<account_email> --docker-password=<account_password> --docker-email=<account_email>
```

Здесь:

- `<namespace>` — название пространства имён в Kubernetes (например, `werf-guided-project-production`);
- `<registry_domain>` — домен Registry (например, `registry.gitlab.com`);
- `<account_email>` — email вашей учётной записи в GitLab;
- `<account_password>` — созданный API-ключ.

В каждом Deployment'е также указывается имя секрета:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/029-basic/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
    spec:
      imagePullSecrets:
      - name: "registrysecret"
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="Это каждый раз копировать руками?" %}

В идеале проблема копирования секретов должна решаться на уровне платформы, но если это не сделано — можно прописать нужные команды в CI-процесс.

Вариант решения — завести секрет один раз в пространство имён `kube-system`, а затем в `.gitlab-ci.yaml` при деплое копировать этот секрет:

```bash
kubectl get ns ${CI_ENVIRONMENT_CUSTOM_SLUG:-${CI_ENVIRONMENT_SLUG}} || kubectl create namespace ${CI_ENVIRONMENT_CUSTOM_SLUG:-${CI_ENVIRONMENT_SLUG}}
kubectl get secret registrysecret -n kube-system -o json |
                      jq ".metadata.namespace = \"${CI_ENVIRONMENT_CUSTOM_SLUG:-${CI_ENVIRONMENT_SLUG}}\"|
                      del(.metadata.annotations,.metadata.creationTimestamp,.metadata.resourceVersion,.metadata.selfLink,.metadata.uid)" |
                      kubectl apply -f -
```
{% endofftopic %}


<a name="iac-debug-deploy" />

## Отладка конфигов инфраструктуры и деплой в Kubernetes

После того, как написана основная часть конфигов, хочется проверить корректность конфигов и задеплоить их в Kubernetes. Для того, чтобы отрендерить конфиги инфраструктуры, требуются сведения об окружении, на которое будет произведён деплой, ключ для расшифровки секретных значений и т.п.

Если мы запускаем werf вне Gitlab CI, необходимо сделать несколько операций вручную прежде, чем werf сможет рендерить конфиги и деплоить в Kubernetes. А именно:

* Вручную подключиться к GitLab Registry с помощью [`docker login`](https://docs.docker.com/engine/reference/commandline/login/) (если это не было сделано ранее);
* Установить переменную окружения `WERF_IMAGES_REPO` с путём до Registry (вида `registry.mydomain.io/myproject`);
* Установить переменную окружения `WERF_SECRET_KEY` со значением, [сгенерированным ранее в главе "Разное поведение в разных окружениях"](#secret-values-yaml);
* Установить переменную окружения `WERF_ENV` с названием окружения, в которое будет осуществляться деплой. Вопрос разных окружений будет затронут подробнее в процессе создания CI-процесса, а сейчас — просто установим значение `staging`. **Важно удалить эту переменную в финальном варианте деплоя**: иначе деплой всегда будет идти в один и тот же namespace.

Если вы всё правильно сделали, уже должны корректно отрабатывать команды [`werf helm render`]({{ site.docsurl }}/v1.1-stable/documentation/cli/management/helm/render.html) и [`werf deploy`]({{ site.docsurl }}/v1.1-stable/documentation/cli/main/deploy.html). _Примечание: при локальном запуске эти команды могут жаловаться на нехватку данных, которые в ином случае были бы проброшены из CI. Например, на данные о теге собранного образа. Это нормально._

{% offtopic title="Как вообще работает деплой?" %}

werf (по аналогии с Helm) берет YAML-шаблоны, которые описывают объекты Kubernetes, и генерирует из них общий манифест. Манифест отдается в API Kubernetes, который на его основе вносит все необходимые изменения в кластер.

werf отслеживает, как Kubernetes вносит изменения, и сигнализирует о результатах в реальном времени. Всё это возможно благодаря встроенной в werf библиотеке [kubedog](https://github.com/werf/kubedog). Уже сам Kubernetes занимается выкачиванием нужных образов из Registry и запуском их на нужных серверах с указанными настройками.

{% endofftopic %}

Запустите деплой и дождитесь успешного завершения:

```bash
werf deploy --stages-storage :local
```

А проверить, что приложение задеплоилось в кластер, можно с помощью kubectl. Должно получиться примерно следующее:

```bash
$ kubectl get namespace
NAME                                 STATUS               AGE
default                              Active               161d
werf-guided-project-production       Active               4m44s
werf-guided-project-staging          Active               3h2m
```

{% offtopic title="Как формируется имя namespace'а?" %}

По шаблону `[[ project ]]-[[ env ]]`, где `[[ project ]]` — имя проекта, а `[[ env ]]` — имя окружения. Подробнее можно почитать [в документации]({{ site.docsurl }}/v1.1-stable/documentation/configuration/deploy_into_kubernetes.html#namespace-%D0%B2-kubernetes).

При необходимости namespace можно переназначить.

{% endofftopic %}

```bash
$ kubectl -n example-1-staging get po
NAME                                 READY                STATUS   RESTARTS  AGE
werf-guided-project-9f6bd769f-rm8nz  1/1                  Running  0         6m12s
```

```bash
$ kubectl -n example-1-staging get ingress
NAME                                 HOSTS                ADDRESS  PORTS     AGE
werf-guided-project                  staging.mydomain.io           80        6m18s
```

А также вы должны увидеть сервис через браузер.

<div id="go-forth-button">
    <go-forth url="30_ci.html" label="Построение CI-процесса" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
