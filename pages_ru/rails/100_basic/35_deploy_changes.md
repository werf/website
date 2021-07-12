---
title: Внесение изменений
permalink: rails/100_basic/35_deploy_changes.html
examples_initial: examples/basic/002_deploy
examples_scale: examples/basic/003_scale
examples_modify_app: examples/basic/004_modify_app
chapter_initial_prepare_cluster: true
chapter_initial_prepare_repo: true
description: |
    Внесём изменения в уже задеплоенное приложение и его инфраструктуру. Продемонстрируем, как работает подход infrastructure-as-code (IaC).
---

## Работаем с инфраструктурой

> Сейчас попробуем некоторые команды по взаимодействию с кластером, а в следующей главе рассмотрим их более подробно.

### Просмотр состояния

Чтобы посмотреть состояние запущенного приложения в Kubernetes, существует несколько команд.

Получить список запущенных Deployment'ов и Pod'ов:
```shell
kubectl get deployment
```

{% offtopic title="Посмотреть ответ" %}
```shell
NAME             READY   UP-TO-DATE   AVAILABLE   AGE
werf-guide-app   1/1     1            1           5m41s
```
{% endofftopic %}

```shell
kubectl get pod
```

{% offtopic title="Посмотреть ответ" %}
```shell
NAME                              READY   STATUS    RESTARTS   AGE
werf-guide-app-57f854484f-8ktqp   1/1     Running   0          7m24s
```
{% endofftopic %}

Посмотреть summary-информацию по Deployment'у `werf-guide-app`:
```shell
kubectl describe deployment werf-guide-app
```

{% offtopic title="Посмотреть ответ" %}
```shell
Name:                   werf-guide-app
Namespace:              werf-guide-app
CreationTimestamp:      Mon, 12 Jul 2021 10:30:14 +0300
Labels:                 app.kubernetes.io/managed-by=Helm
Annotations:            deployment.kubernetes.io/revision: 1
                        meta.helm.sh/release-name: werf-guide-app
                        meta.helm.sh/release-namespace: werf-guide-app
                        project.werf.io/env: 
                        project.werf.io/name: werf-guide-app
                        werf.io/version: v1.2.11+fix10
Selector:               app=werf-guide-app
Replicas:               1 desired | 1 updated | 1 total | 1 available | 0 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  25% max unavailable, 25% max surge
Pod Template:
  Labels:  app=werf-guide-app
  Containers:
   app:
    Image:      <dockerhubuser>/werf-guide-app:e5c6ebcd2718ccfe74d01069a0d758e03d5a2554155ccdc01be0daff-1625753768312
    Port:       8000/TCP
    Host Port:  0/TCP
    Command:
      /app/start.sh
    Environment:  <none>
    Mounts:       <none>
  Volumes:        <none>
Conditions:
  Type           Status  Reason
  ----           ------  ------
  Available      True    MinimumReplicasAvailable
  Progressing    True    NewReplicaSetAvailable
OldReplicaSets:  <none>
NewReplicaSet:   werf-guide-app-57f854484f (1/1 replicas created)
Events:
  Type    Reason             Age    From                   Message
  ----    ------             ----   ----                   -------
  Normal  ScalingReplicaSet  8m14s  deployment-controller  Scaled up replica set werf-guide-app-57f854484f to 1
```
{% endofftopic %}

Deployment запускает Pod'ы. Логи пишутся в Pod'ах. Получить логи одного из запущенных Pod'ов можно командой:
```shell
kubectl logs app-<podId>
```

{% offtopic title="Посмотреть ответ" %}
```shell
GET /ping HTTP/1.1
Host: werf-guide-app
X-Request-ID: 3cdfedee7ff8f4fb27566446d587ccbd
X-Real-IP: 192.168.49.1
X-Forwarded-For: 192.168.49.1
X-Forwarded-Host: werf-guide-app
X-Forwarded-Port: 80
X-Forwarded-Proto: http
X-Scheme: http
User-Agent: curl/7.68.0
Accept: */*
```
{% endofftopic %}

### Масштабирование

Веб-сервер запущен в Deployment'е `werf-guide-app`. Помотрим, сколько реплик запущено:

```shell
kubectl get pod
```

{% offtopic title="Посмотреть ответ" %}
```shell
NAME                              READY   STATUS    RESTARTS   AGE
werf-guide-app-57f854484f-8ktqp   1/1     Running   0          7m24s
```
{% endofftopic %}

Поменяем вручную на 4 реплики:
```shell
kubectl edit deployment werf-guide-app
```

В открывшемся редакторе выставляем `spec.replicas=4`, закрываем редактор.
Снова посмотрим, сколько реплик запущено:
```shell
kubectl get pod
```

{% offtopic title="Посмотреть ответ" %}
```shell
NAME                              READY   STATUS    RESTARTS   AGE
werf-guide-app-57f854484f-2xbzl   1/1     Running   0          119m
werf-guide-app-57f854484f-8ktqp   1/1     Running   0          6h59m
werf-guide-app-57f854484f-kr9sl   1/1     Running   0          119m
werf-guide-app-57f854484f-r99j4   1/1     Running   0          118m
```
{% endofftopic %}

Мы произвели масштабирование вручную. Теперь снова запустим `werf converge`:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guide-app
```

И снова посмотрим, сколько реплик запущено:
```shell
kubectl get pod
```

{% offtopic title="Посмотреть ответ" %}
```shell
NAME                              READY   STATUS    RESTARTS   AGE
werf-guide-app-57f854484f-2xbzl   1/1     Running   0          119m
```
{% endofftopic %}

Количество реплик снова соответствует таковому в Git-репозитории. Дело в том, что werf привела состояние кластера к состоянию, описанному в текущем Git-коммите. Этот принцип называется [**гитерминизмом**](https://ru.werf.io/documentation/v1.2/advanced/giterminism.html) (giterminism).

Как же соблюсти **гитерминизм** и сделать всё правильно? Выставим тот же `spec.replicas=4`, но уже через состояние приложения, описанное в Git:

{% include chapter_prepare_repo_commands.md.liquid examples_to=page.examples_scale from_scratch=false %}

Через файл `.helm/templates/deployment.yaml`:

{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" examples=page.examples_scale %}

Запустим деплой:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guide-app
```

Проверим сколько реплик запущено теперь:

```shell
kubectl get pod
```

{% offtopic title="Посмотреть ответ" %}
```shell
NAME                              READY   STATUS    RESTARTS   AGE
werf-guide-app-57f854484f-2xbzl   1/1     Running   0          119m
werf-guide-app-57f854484f-8ktqp   1/1     Running   0          6h59m
werf-guide-app-57f854484f-kr9sl   1/1     Running   0          119m
werf-guide-app-57f854484f-r99j4   1/1     Running   0          118m
```
{% endofftopic %}

Вернём состояние рабочей директории к предыдущему, с одной репликой:

{% include chapter_prepare_repo_commands.md.liquid examples_to=page.examples_initial from_scratch=false %}

Проверим содержимое `.helm/templates/deployment.yaml`:

{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" examples=page.examples_initial %}

Запустим деплой:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guide-app
```

## Меняем приложение

Наше приложение — это простейший эхо сервер, который при запросе при помощи curl:
```shell
curl "http://werf-guide-app/ping"
```

Ответит словом `pong`.

Внесём изменение в код приложения:

{% include chapter_prepare_repo_commands.md.liquid examples_to=page.examples_modify_app from_scratch=false %}

В коде приложения ответ сервера изменится на `hello world`:

{% include snippetcut_example path="start.sh" syntax="bash" examples=page.examples_modify_app %}

Запустим деплой:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guide-app
```

Проверим результат:
```shell
curl "http://werf-guide-app/ping"
```

В ответе увидим `hello world`.
Поздравляем, у нас всё получилось!
