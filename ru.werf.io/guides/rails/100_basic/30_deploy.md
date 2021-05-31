---
title: Деплой приложения
permalink: rails/100_basic/30_deploy.html
---

В предыдущих главах мы научились собирать образ приложения, а в этой главе мы этот образ будем разворачивать в Kubernetes.

При деплое в Kubernetes используются Kubernetes-манифесты, которые описывают сущности, необходимые для работы приложений. Эти сущности включают в себя, к примеру, Deployment, отвечающий за запуск приложений в контейнерах, и Service/Ingress, отвечающие за доступ к запущенным приложениям изнутри и извне кластера.

Перед деплоем werf создаёт Kubernetes-манифесты из шаблонов. Шаблоны в werf — это Helm-шаблоны с некоторой [дополнительной функциональностью]({{ site.docsurl }}/documentation/v1.2/advanced/helm/overview.html).

## Deployment

Ресурс Deployment создаёт набор ресурсов, запускающих приложение. Создадим для него файл `.helm/templates/deployment.yaml` с содержанием:
{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/015_deploy_app/.helm/templates/deployment.yaml" %}
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
  template:
    metadata:
      labels:
        app: basicapp
    spec:
      containers:
      - name: basicapp
        command: ["/bin/bash", "-ec", "bundle exec rails db:migrate RAILS_ENV=development && bundle exec puma"]
        image: {{ .Values.werf.image.basicapp }}
        ports:
        - containerPort: 3000
        env:
        - name: "SQLITE_FILE"
          value: "app.db"
```
{% endraw %}
{% endsnippetcut %}

Конструкция {% raw %}`image: {{ .Values.werf.image.basicapp }}`{% endraw %} используется для того, чтобы werf автоматически подставлял генерируемый тег образа (и имя образа), так как тег становится известен только в процессе сборки.

Werf пересобирает образы только при изменениях в файлах, указанных в `werf.yaml`, а также при изменении самого `werf.yaml`. При пересборке изменится и тег образа, что приведёт к обновлению Deployment'а. Если же изменений в вышеупомянутых файлах нет, то образ не пересоберётся, а Deployment и создаваемые им ресурсы не перевыкатятся, т.к. в этом просто нет необходимости — в кластере уже самая свежая версия приложения.

{% comment %} TODO(lesikov): скроем, пока не сделаем инструкцию для деплоя из своего кластера.
## Registry Secret

Kubernetes-кластер для запуска приложения использует образы из registry. Поэтому важно, чтобы кластер мог авторизоваться в registry. Как правило, ситуация отличается для локального и внешнего registry.

<div class="tabs">
<a href="javascript:void(0)" class="tabs__btn tabs__secret__btn" onclick="openTab(event, 'tabs__secret__btn', 'tabs__secret__content', 'tab__secret__local')">Локальный registry</a>
<a href="javascript:void(0)" class="tabs__btn tabs__secret__btn" onclick="openTab(event, 'tabs__secret__btn', 'tabs__secret__content', 'tab__secret__remote')">Внешний registry</a>
</div>

<div id="tab__secret__local" class="tabs__content tabs__secret__content" markdown="1">
{% include_relative 30_deploy_registrysecret_local.md %}
</div>

<div id="tab__secret__remote" class="tabs__content tabs__secret__content" markdown="1">
{% include_relative 30_deploy_registrysecret_remote.md %}
</div>
{% endcomment %}

## Service

Ресурс Service позволит другим приложениям в кластере обращаться к нашему приложению. Создадим для него файл `.helm/templates/service.yaml` с содержанием:
{% offtopic title="Windows" %}
{% snippetcut name=".helm/templates/service.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/015_deploy_app/.helm/templates/service.yaml" %}
{% raw %}
```yaml
apiVersion: v1
kind: Service
metadata:
  name: basicapp
spec:
  type: LoadBalancer
  selector:
    app: basicapp
  ports:
  - name: http
    port: 3000
```
{% endraw %}
{% endsnippetcut %}
{% endofftopic %}
{% offtopic title="macOS/Linux" %}
{% snippetcut name=".helm/templates/service.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/015_deploy_app/.helm/templates/service.yaml" %}
{% raw %}
```yaml
apiVersion: v1
kind: Service
metadata:
  name: basicapp
spec:
  selector:
    app: basicapp
  ports:
  - name: http
    port: 3000
```
{% endraw %}
{% endsnippetcut %}
{% endofftopic %}

## Ingress

Ресурс Ingress позволяет открыть доступ к нашему приложению снаружи кластера, т.к. Service разрешает доступ только между приложениями внутри кластера. В Ingress'е мы указываем на какой Service должен пойти внешний трафик, который придет на домен `example.com`. Создадим для Ingress'а файл `.helm/templates/ingress.yaml` с содержимым:

{% snippetcut name=".helm/templates/ingress.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/015_deploy_app/.helm/templates/ingress.yaml" %}
{% raw %}
```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: nginx
  name: basicapp
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /
        backend:
          serviceName: basicapp
          servicePort: 3000
```
{% endraw %}
{% endsnippetcut %}

Также добавим наш домен в конфигурацию Rails (`config/application.rb`):
{% snippetcut name="config/application.rb" url="#" %}
{% raw %}
```ruby
module DemoApplication
  class Application < Rails::Application
    config.hosts << 'example.com'
    ...
  end
end
```
{% endraw %}
{% endsnippetcut %}

## Деплой в Kubernetes

Сохраним наши изменения перед сборкой/деплоем:
```shell
git add .helm config
git commit -m "Add deploy configuration"
```

Команда [werf converge]({{ site.docsurl }}/documentation/reference/cli/werf_converge.html) сделает сразу и сборку и развертывание приложения в Kubernetes:
```shell
werf converge --repo registry.example.com:80/werf-guided-rails
```

Результат выполнения в случае успешной сборки и деплоя:
```shell
...
│ basicapp/dockerfile  Successfully built 4c1054085159
│ │ basicapp/dockerfile  Successfully tagged 93c05bf8-c459-4768-b388-3cdbc80e2868:latest
│ ├ Info
│ │       name: localhost:5000/werf-guided-rails:f4caaa836701e5346c4a0514bb977362ba5fe4ae114d0176f6a6c8cc-1612277803607
│ │       size: 371.4 MiB
│ └ Building stage basicapp/dockerfile (40.31 seconds)
└ :boat: image basicapp (41.13 seconds)
...
┌ Waiting for release resources to become ready
│ ┌ Status progress
│ │ DEPLOYMENT                                                                                                                                                      REPLICAS                      AVAILABLE                        UP-TO-DATE
│ │ basicapp                                                                                                                                                        1/1                           1                                1
│ │ │   POD                                                           READY                  RESTARTS                       STATUS
│ │ └── 687f8cc569-n6gkw                                              1/1                    0                              Running
│ └ Status progress
└ Waiting for release resources to become ready (0.02 seconds)

NAME: werf-guided-rails
LAST DEPLOYED: Tue Feb  2 21:57:23 2021
NAMESPACE: werf-guided-rails
STATUS: deployed
REVISION: 1
TEST SUITE: None
Running time 62.66 seconds
```

Теперь приложение доступно по ссылке [http://example.com:3000/api/labels](http://example.com:3000/api/labels) для Windows и [http://example.com/api/labels](http://example.com/api/labels) для macOS/Linux.

![](/images/template/100_30_app_in_browser.png)

