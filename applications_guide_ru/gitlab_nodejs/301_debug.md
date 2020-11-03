---
title: Организация локальной разработки
permalink: gitlab_nodejs/301_debug.html
---




## Отладка кода инфраструктуры

Самая частая проблема: после деплоя приложения 

Если вы всё правильно сделали, уже должны корректно отрабатывать команды [`werf helm render`]({{ site.docsurl }}/documentation/cli/management/helm/render.html) и [`werf deploy`]({{ site.docsurl }}/documentation/cli/main/deploy.html). _Примечание: при локальном запуске эти команды могут жаловаться на нехватку данных, которые в ином случае были бы проброшены из CI. Например, на данные о теге собранного образа. Это нормально._

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

По шаблону `[[ project ]]-[[ env ]]`, где `[[ project ]]` — имя проекта, а `[[ env ]]` — имя окружения. Подробнее можно почитать [в документации]({{ site.docsurl }}/documentation/configuration/deploy_into_kubernetes.html#namespace-%D0%B2-kubernetes).

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



### Селекторы


Обратите внимание на поле `selector` у Service: он должен совпадать с аналогичным полем у Deployment. Ошибки в этой части — самая частая проблема с настройкой маршрута до приложения.

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/020-basic-1/.helm/templates/deployment.yaml" %}
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
