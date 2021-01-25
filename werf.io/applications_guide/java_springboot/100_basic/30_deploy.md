---
title: Деплой приложения
permalink: java_springboot/100_basic/30_deploy.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/templates/ingress.yaml
- .helm/templates/service.yaml
{% endfilesused %}

Для того, чтобы приложение заработало в Kubernetes, необходимо описать инфраструктуру приложения как код ([IaC](https://ru.wikipedia.org/wiki/%D0%98%D0%BD%D1%84%D1%80%D0%B0%D1%81%D1%82%D1%80%D1%83%D0%BA%D1%82%D1%83%D1%80%D0%B0_%D0%BA%D0%B0%D0%BA_%D0%BA%D0%BE%D0%B4)). В нашем случае потребуются следующие объекты Kubernetes: Deployment, Service и Ingress.

{% offtopic title="Как быть, если я не знаю Kubernetes?" %}
В самоучителе будет приведён исходный код инфраструктуры и вы сможете интуитивно догадаться, что там написано. Чтобы научиться писать подобный код самостоятельно — вы можете воспользоваться обучающими материалами/документацией, например, [от вендора](https://kubernetes.io/docs/tutorials/kubernetes-basics/). Учебников, рассказывающих об объектах Kubernetes и возможных их настройках, на данный момент существует достаточно.

А ещё лучше: следуя сути DevOps, общайтесь с инженерами поддержки. Конфигурация объектов — это хороший язык для коммуникации на стыке разработки и поддержки. 
{% endofftopic %}

werf поддерживает весь функционал шаблонизатора helm, а также [предоставляет дополнительные функции и значения]({{ site.docsurl }}/documentation/advanced/helm/basics.html). Разберём самые необходимые из них. Подробнее о шаблонизации и правилах написания Kubernetes-объектов мы разберёмся позже, в главе "Конфигурирование инфраструктуры в виде кода", пока что — добьёмся, чтобы приложение заработало в реальном кластере.

## Deployment

Объект Deployment позволяет создать объект Pod, который содержит в себе и управляет контейнерами с приложениями. У создаваемого нами Pod будет один контейнер — `basicapp`.

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/015_deploy_app/.helm/templates/deployment.yaml" %}
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
      containers:
      - name: basicapp
        command: ["java","-jar","/app/demo.jar"]
        image: {{ tuple "basicapp" . | werf_image }}
        workingDir: /app
        ports:
        - containerPort: 8080
          protocol: TCP
        env:
        - name: "SQLITE_FILE"
          value: "app.db"
```
{% endraw %}
{% endsnippetcut %}

Обратите внимание на конструкцию {% raw %}`image: {{ tuple "basicapp" . | werf_image }}`{% endraw %} — с помощью неё подставляется актуальное имя образа (`REPO:TAG`).

werf пересобирает контейнеры только при необходимости, если есть изменения в связанных файлах в git или в конфигурации werf.yaml. Аналогично, werf отслеживает изменение образов, и делает так, чтобы перевыкат Pod-а так же только при необходимости.

## Service

Объект Service позволяет приложениям в кластере взаимодействовать друг с другом. Пропишем его:

{% snippetcut name=".helm/templates/service.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/015_deploy_app/.helm/templates/service.yaml" %}
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
    port: 8080
    protocol: TCP
```
{% endraw %}
{% endsnippetcut %}

## Ingress

Объект Ingress позволяет организовать маршрутизацию трафика на созданный Service для нужного домена (в нашем примере — `mydomain.io`).

{% offtopic title="Что за объект Ingress и как он связан с балансировщиком?" %}
Возможна коллизия терминов:

* Есть Ingress в смысле [NGINX Ingress Controller](https://github.com/kubernetes/ingress-nginx), который работает в кластере и принимает входящие извне запросы.
* А ещё есть [объект Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/), который фактически описывает настройки для NGINX Ingress Controller.

В статьях и бытовой речи оба этих термина зачастую называют «Ingress», так что догадываться нужно по контексту.
{% endofftopic %}

{% snippetcut name=".helm/templates/ingress.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/015_deploy_app/.helm/templates/ingress.yaml" %}
{% raw %}
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: nginx
  name: basicapp
spec:
  rules:
  - host: mydomain.io
    http:
      paths:
      - path: /
        backend:
          serviceName: basicapp
          servicePort: 8080
```
{% endraw %}
{% endsnippetcut %}

## Выкат в кластер

Воспользуемся [командой `converge`]({{ site.docsurl }}/documentation/reference/cli/werf_converge.html) для того, чтобы собрать образ, загрузить собранный образ в registry и задеплоить приложение в Kubernetes. Единственной опцией указывается репозиторий для хранения образов `--repo registry.mydomain.io/werf-guided-project`.

Сделайте коммит изменений в репозитории с кодом и затем выполните:

```bash
werf converge --repo registry.mydomain.io/werf-guided-project
```

В результате вы должны увидеть логи примерно такого вида:

```
│ │ basicapp/dockerfile  Successfully built 7e38465ee6de
│ │ basicapp/dockerfile  Successfully tagged cbb1cef2-a03a-432f-b13d-b95f0f0cb4e9:latest
│ ├ Info
│ │       name: localhost:5005/werf-guided-project:017ce9df8dbd7d3505546c95557f1c1f39ce1e6666aaae29e8c12608-1605619646009
│ │       size: 375.8 MiB
│ └ Building stage basicapp/dockerfile (209.48 seconds)
└ ⛵ image basicapp (213.60 seconds)

Release "werf-guided-project" does not exist. Installing it now.
W1117 16:29:16.809420   42045 warnings.go:67] networking.k8s.io/v1beta1 Ingress is deprecated in v1.19+, unavailable in v1.22+; use networking.k8s.io/v1 Ingress
W1117 16:29:16.936024   42045 warnings.go:67] networking.k8s.io/v1beta1 Ingress is deprecated in v1.19+, unavailable in v1.22+; use networking.k8s.io/v1 Ingress

┌ Waiting for release resources to become ready
│ ┌ Status progress
│ │ DEPLOYMENT                                                                                                                                                      REPLICAS                      AVAILABLE                        UP-TO-DATE
│ │ basicapp                                                                                                                                                        1/1                           1                                1
│ │ │   POD                                                           READY                  RESTARTS                       STATUS
│ │ └── 6cf5b444bc-6rh4j                                              1/1                    0                              Running
│ └ Status progress
└ Waiting for release resources to become ready (4.02 seconds)

NAME: werf-guided-project
LAST DEPLOYED: Tue Nov 17 16:29:16 2020
NAMESPACE: werf-guided-project
STATUS: deployed
REVISION: 1
TEST SUITE: None
Running time 222.54 seconds
```

После этого приложение должно быть доступно в браузере:

![](/applications_guide/images/template/100_30_app_in_browser.png)

<div id="go-forth-button">
    <go-forth url="40_optimize.html" label="Ускорение сборки" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
