---
title: Деплой приложения
permalink: rails/100_basic/30_deploy.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/templates/registry-secret.yaml
- .helm/templates/ingress.yaml
- .helm/templates/service.yaml
{% endfilesused %}

В предыдущей главе мы описали IaC для сборки, теперь нужно описать IaC для запуска в Kubernetes. В нашем случае потребуются следующие объекты Kubernetes: Deployment, Service и Ingress.

{% offtopic title="Как быть, если я не знаю Kubernetes?" %}
В самоучителе будет приведён исходный код инфраструктуры и вы сможете интуитивно догадаться, что там написано. Чтобы научиться писать подобный код самостоятельно, стоит воспользоваться обучающими материалами и/или, например, в [официальной документацией Kubernetes](https://kubernetes.io/docs/tutorials/kubernetes-basics/). Видеокурсов и учебников, рассказывающих об объектах Kubernetes и их возможных настройках, на данный момент существует достаточно.
{% endofftopic %}

werf поддерживает весь функционал шаблонизатора Helm (werf использует вкомпилированный в него Helm для деплоя), а также предоставляет [дополнительную функциональность]({{ site.docsurl }}/documentation/v1.2/advanced/helm/overview.html). Подробнее в шаблонизации и правилах написания Kubernetes-объектов мы разберёмся позже, в главе «Конфигурирование инфраструктуры в виде кода», а пока — добьёмся, чтобы приложение заработало в реальном кластере.

## Deployment

Объект Deployment позволяет создать объект Pod, который содержит в себе контейнеры с приложениями и управляет ими. У создаваемого нами Pod будет один контейнер — `basicapp`.

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
  revisionHistoryLimit: 3
  strategy:
    type: RollingUpdate
  replicas: 1
  template:
    metadata:
      labels:
        app: basicapp
    spec:
      containers:
      - name: basicapp
        command: ["bundle","exec", "puma"]
        image: {{ .Values.werf.image.basicapp }}
        workingDir: /app
        ports:
        - containerPort: 3000
          protocol: TCP
        env:
        - name: "SQLITE_FILE"
          value: "app.db"
```
{% endraw %}
{% endsnippetcut %}

Обратите внимание на конструкцию {% raw %}`image: {{ .Values.werf.image.basicapp }}`{% endraw %} — с помощью неё подставляется актуальное имя образа.

werf пересобирает контейнеры только при необходимости, если есть изменения в связанных файлах в Git или в конфигурации `werf.yaml`. При изменении образа меняется его тег, что приводит к перевыкату Pod'а с новым образом.

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

## Service

Объект Service позволяет приложениям в кластере обнаруживать друг друга. Пропишем его:

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
    protocol: TCP
```
{% endraw %}
{% endsnippetcut %}

## Ingress

Объект Ingress позволяет организовать маршрутизацию трафика на созданный Service для нужного домена (в нашем примере — `example.com`).

{% offtopic title="Что за объект Ingress?" %}
Возможна коллизия терминов:

* Есть [NGINX Ingress Controller](https://github.com/kubernetes/ingress-nginx), который работает в кластере и принимает входящие извне запросы.
* А ещё есть [объект Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/), который фактически описывает настройки для NGINX Ingress Controller.

В статьях и бытовой речи оба этих термина зачастую называют «Ingress», так что догадываться нужно по контексту.
{% endofftopic %}

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

Также нужно добавить хост в конфигурацию Rails:

{% snippetcut name="config/application.rb" url="#" %}
{% raw %}
```ruby
class Application < Rails::Application
  # ...
  config.hosts << 'example.com'
  # ...
end
```
{% endraw %}
{% endsnippetcut %}

После этого не забыть зафиксировать изменения!

## Выкат в кластер

Перед выполнением выката необходимо добавить изменения в git-репозиторий проекта:

```shell
git add .helm/templates
git commit -m "Add helm chart configuration"
```

> Почему изменения должны добавляться в git-репозиторий, что такое гитерминизм и режим разработчика, а также другие особенности работы с файлами проекта будут разобраны далее в главе «Необходимо знать»

Команда [werf converge]({{ site.docsurl }}/documentation/reference/cli/werf_converge.html) используется для сборки и загрузки образов в registry и последующего выката приложения в Kubernetes. Единственной обязательной опцией указывается репозиторий для хранения образов `--repo registry.example.com/werf-guided-springboot`.

```shell
werf converge --repo registry.example.com/werf-guided-rails
...
│ basicapp/dockerfile  Successfully built 4c1054085159
│ │ basicapp/dockerfile  Successfully tagged 93c05bf8-c459-4768-b388-3cdbc80e2868:latest
│ ├ Info
│ │       name: localhost:5000/werf-guided-rails:f4caaa836701e5346c4a0514bb977362ba5fe4ae114d0176f6a6c8cc-1612277803607
│ │       size: 371.4 MiB
│ └ Building stage basicapp/dockerfile (40.31 seconds)
└ :boat: image basicapp (41.13 seconds)

Release "werf-guided-rails" does not exist. Installing it now.

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

После этого приложение должно быть доступно в браузере:

![](/guides/images/template/100_30_app_in_browser.png)

<div id="go-forth-button">
    <go-forth url="40_optimize.html" label="Ускорение сборки" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
