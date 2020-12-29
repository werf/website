---
title: Деплой приложения
permalink: nodejs/100_basic/30_deploy.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/templates/secret.yaml
- .helm/templates/ingress.yaml
- .helm/templates/service.yaml
{% endfilesused %}

В предыдущей главе мы описали IaC для сборки, теперь нужно описать её для запуска в Kubernetes. В нашем случае потребуются следующие объекты Kubernetes: Deployment, Service и Ingress.

{% offtopic title="Как быть, если я не знаю Kubernetes?" %}
В самоучителе будет приведён исходный код инфраструктуры и вы сможете интуитивно догадаться, что там написано. Чтобы научиться писать подобный код самостоятельно — вы можете воспользоваться обучающими материалами/документацией, например, в [официальной документации Kubernetes](https://kubernetes.io/docs/tutorials/kubernetes-basics/). Учебников, рассказывающих об объектах Kubernetes и возможных их настройках, на данный момент существует достаточно.
{% endofftopic %}

werf поддерживает весь функционал шаблонизатора helm, а также [предоставляет дополнительные функции и значения]({{ site.docsurl }}/documentation/advanced/helm/basics.html). Разберём самые необходимые из них. Подробнее о шаблонизации и правилах написания Kubernetes-объектов мы разберёмся позже, в главе "Конфигурирование инфраструктуры в виде кода", пока что — добьёмся, чтобы приложение заработало в реальном кластере.

## Deployment

Объект Deployment позволяет создать объект Pod, который содержит в себе и управляет контейнерами с приложениями. У создаваемого нами Pod будет один контейнер — `basicapp`.

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/015_deploy_app/.helm/templates/deployment.yaml" %}
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
        command: ["node","/app/app.js"]
        image: {{ tuple "basicapp" . | werf_image }}
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

Обратите внимание на конструкцию {% raw %}`image: {{ tuple "basicapp" . | werf_image }}`{% endraw %} — с помощью неё подставляется актуальное имя образа (`REPO:TAG`).

werf пересобирает контейнеры только при необходимости, если есть изменения в связанных файлах в git или в конфигурации werf.yaml. Аналогично, werf отслеживает изменение образов, и делает так, чтобы перевыкат Pod-а так же только при необходимости.

## Registry Secret

Kubernetes-кластер для запуска приложения использует образы из registry. Для этого ему нужно авторизоваться с помощью логина и пароля (вы сталкивались с ними в главе "Подготовка окружения"). Эти логин и пароль мы сообщаем кластеру с помощью объекта типа Secret с именем `registrysecret` (мы упомянули его ранее, в Deployment-е: `imagePullSecrets` - `registrysecret`). Опишем, как вам сформировать свой файл `secret.yaml`.

Допустим, ваш логин `admin`, и пароль тоже `admin`. Закодируем их с помощью base64:

```bash
echo -n "admin:admin" | base64
```

В результате получаем строку `YWRtaW46YWRtaW4=`. Её мы используем, чтобы сформировать подобную JSON:

```json
{"auths":{"localhost":{"username":"admin","password":"admin","email":"admin","auth":"YWRtaW46YWRtaW4="}}}
```

И закодируем его в base64:

```bash
echo -n '{"auths":{"localhost":{"username":"admin","password":"admin","email":"admin","auth":"YWRtaW46YWRtaW4="}}}' | base64
```

Получаем строку 

```
eyJhdXRocyI6eyJsb2NhbGhvc3QiOnsidXNlcm5hbWUiOiJhZG1pbiIsInBhc3N3b3JkIjoiYWRtaW4iLCJlbWFpbCI6ImFkbWluIiwiYXV0aCI6IllXUnRhVzQ2WVdSdGFXND0ifX19
```

Которую мы сможем использовать в своём `secret.yaml`:

{% snippetcut name=".helm/templates/secret.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/015_deploy_app/.helm/templates/secret.yaml" %}
{% raw %}
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: registrysecret
type: kubernetes.io/dockerconfigjson
data:
  .dockerconfigjson: eyJhdXRocyI6eyJsb2NhbGhvc3QiOnsidXNlcm5hbWUiOiJhZG1pbiIsInBhc3N3b3JkIjoiYWRtaW4iLCJlbWFpbCI6ImFkbWluIiwiYXV0aCI6IllXUnRhVzQ2WVdSdGFXND0ifX19
```
{% endraw %}
{% endsnippetcut %}

_Примечание: в приведённом примере ключи доступа хранятся в не зашифрованном (а только закодированном) виде. Это небезопасно. Вопросы безопасного хранения ключей мы рассмотрим в главе "Организация не локальной разработки"_

## Service

Объект Service позволяет приложениям в кластере обнаруживать друг друга. Пропишем его:

{% snippetcut name=".helm/templates/service.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/015_deploy_app/.helm/templates/service.yaml" %}
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

{% offtopic title="Что за объект Ingress и как он связан с балансировщиком?" %}
Возможна коллизия терминов:

* Есть Ingress в смысле [NGINX Ingress Controller](https://github.com/kubernetes/ingress-nginx), который работает в кластере и принимает входящие извне запросы.
* А ещё есть [объект Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/), который фактически описывает настройки для NGINX Ingress Controller.

В статьях и бытовой речи оба этих термина зачастую называют «Ingress», так что догадываться нужно по контексту.
{% endofftopic %}

{% snippetcut name=".helm/templates/ingress.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/015_deploy_app/.helm/templates/ingress.yaml" %}
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

## Выкат в кластер

Воспользуемся [командой `converge`]({{ site.docsurl }}/documentation/reference/cli/werf_converge.html) для того, чтобы собрать образ, загрузить собранный образ в registry и задеплоить приложение в Kubernetes. Единственной опцией указывается репозиторий для хранения образов `--repo registry.example.com/werf-guided-nodejs`.

{% offtopic title="Как будут храниться образы в репозитории?" %}
При организации деплоя без использования werf зачастую приходится формализовать принципы, по которым именуются образы в Registry. В нашем случае об этом нет необходимости думать — werf берёт организацию тегирования на себя.

werf реализует content-based тегирование: образы меняются и выкатываются автоматически, если меняется состояние контента в git.

Если вы хотите узнать подробности — читайте [в документации]({{ site.docsurl }}/documentation/internals/stages_and_storage.html#%D0%B8%D0%BC%D0%B5%D0%BD%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D0%B5-%D1%81%D1%82%D0%B0%D0%B4%D0%B8%D0%B9) и [статье на medium](https://medium.com/flant-com/content-based-tagging-in-werf-eb96d22ac509).
{% endofftopic %}

Сделайте коммит изменений в репозитории с кодом и затем выполните:

```bash
werf converge --repo registry.example.com/werf-guided-nodejs
```

В результате вы должны увидеть логи примерно такого вида:

```
│ │ basicapp/dockerfile  Successfully built 7e38465ee6de
│ │ basicapp/dockerfile  Successfully tagged cbb1cef2-a03a-432f-b13d-b95f0f0cb4e9:latest
│ ├ Info
│ │       name: localhost:5005/werf-guided-nodejs:017ce9df8dbd7d3505546c95557f1c1f39ce1e6666aaae29e8c12608-1605619646009
│ │       size: 375.8 MiB
│ └ Building stage basicapp/dockerfile (209.48 seconds)
└ ⛵ image basicapp (213.60 seconds)

Release "werf-guided-nodejs" does not exist. Installing it now.

┌ Waiting for release resources to become ready
│ ┌ Status progress
│ │ DEPLOYMENT                                                                                                                                                      REPLICAS                      AVAILABLE                        UP-TO-DATE
│ │ basicapp                                                                                                                                                        1/1                           1                                1
│ │ │   POD                                                           READY                  RESTARTS                       STATUS
│ │ └── 6cf5b444bc-6rh4j                                              1/1                    0                              Running
│ └ Status progress
└ Waiting for release resources to become ready (4.02 seconds)

NAME: werf-guided-nodejs
LAST DEPLOYED: Tue Nov 17 16:29:16 2020
NAMESPACE: werf-guided-nodejs
STATUS: deployed
REVISION: 1
TEST SUITE: None
Running time 222.54 seconds
```

После этого приложение должно быть доступно в браузере:

![](/applications_guide_ru/images/template/100_30_app_in_browser.png)

<div id="go-forth-button">
    <go-forth url="40_optimize.html" label="Ускорение сборки" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
