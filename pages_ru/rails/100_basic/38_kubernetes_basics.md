---
title: Основы Kubernetes
permalink: rails/100_basic/38_kubernetes_basics.html
chapter_initial_prepare_cluster: false
chapter_initial_prepare_repo: false
description: |
  В этой главе мы рассмотрим основные ресурсы Kubernetes для развертывания приложений и обеспечения доступа к ним изнутри и снаружи кластера.
---

## Шаблоны, манифесты и ресурсы

werf при деплое использует YAML-манифесты, описывающие Kubernetes-ресурсы. Эти манифесты получаются из Helm-шаблонов, которые лежат в `.helm/templates` и `.helm/charts`.

_Helm-шаблон_, описывающий Pod (один из Kubernetes-ресурсов), может лежать в `.helm/templates/pod.yaml` и выглядеть так:
{% raw %}
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: standalone-pod
spec:
  containers:
  - name: main
    image: {{ $.Values.image }}  # Helm-шаблонизация, позволяет параметризовать имя образа для контейнера
    command: ["tail", "-f", "/dev/null"]
```
{% endraw %}

Перед деплоем этот Helm-шаблон с помощью werf преобразуется в _манифест_, который выглядит так:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: standalone-pod
spec:
  containers:
  - name: main
    image: alpine  # если пользователь в качестве имени образа указал "alpine"
    command: ["tail", "-f", "/dev/null"]
```

А уже во время деплоя этот манифест становится _ресурсом_ Pod в Kubernetes-кластере. Посмотреть, как этот ресурс выглядит в кластере, можно с помощью команды `kubectl get`:
```yaml
kubectl get pod standalone-pod --output yaml
```

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: standalone-pod
  namespace: default
spec:
  containers:
  - name: main
    image: alpine
    command: ["tail", "-f", "/dev/null"]
    ...
status:
  phase: Running
  podIP: 172.17.0.7
  startTime: "2021-06-02T13:17:47Z"
```

## Запуск приложений

Использование ресурса Pod — простейший способ запустить в Kubernetes один или несколько контейнеров. В манифесте Pod'а описываются контейнеры и их конфигурация. Но на практике сами по себе Pod'ы стараются не деплоить — вместо них деплоят Controller'ы, которые создают и управляют Pod'ами за вас. Один из таких контроллеров — **Deployment**. Создавая Pod'ы с помощью Deployment, вы сильно упрощаете управление Pod'ами.

Вот некоторые возможности, которые предоставляет Deployment (и которых нет у Pod'ов самих по себе):
* При автоматическом или ручном удалении Pod'а в Deployment'е этот Pod будет перезапущен;
* Большая часть конфигурации Pod'а не может быть обновлена. Чтобы обновить его конфигурацию, Pod нужно пересоздать. Конфигурацию же Pod'а в Deployment'е можно обновлять без пересоздания Deployment'а;
* Обновление конфигурации Pod'ов происходит без простоя: при обновлении часть Pod'ов со старой конфигурацией остаётся активной до тех пор, пока новые Pod'ы не начнут успешно запускаться;
* Одним Deployment'ом можно запускать нескольких Pod'ов сразу, в том числе на разных Node'ах.

Разные контроллеры имеют разные возможности, связанные с созданием и управлением Pod'ами. Вот основные контроллеры и их практическое применение:
* [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) — стандарт для деплоя stateless-приложений;
* [StatefulSet](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) — стандарт для деплоя stateful-приложений;
* [DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/) — для деплоя приложений, которые должны быть запущены только по одному экземпляру на каждом узле (агенты для логирования, мониторинга);
* [Job](https://kubernetes.io/docs/concepts/workloads/controllers/job/) — для запуска разовых задач в Pod'ах (например, миграции базы данных);
* [CronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) — для многоразового запуска задач 
  в Pod'ах по расписанию (например, регулярная очистка чего-либо).

Получить список ресурсов определенного типа в кластере можно с помощью всё той же команды `kubectl get` (в данном примере команды последовательно возвратят список всех Pod'ов, Deployment, Statefulset, Job, Cronjob из всех namespace):

```shell
kubectl get --all-namespaces pod
```

{% offtopic title="Посмотреть ответ" %}
```shell
NAMESPACE           NAME                                        READY   STATUS      RESTARTS   AGE
ingress-nginx       ingress-nginx-admission-create-8bgk7        0/1     Completed   0          11d
ingress-nginx       ingress-nginx-admission-patch-8fkgl         0/1     Completed   1          11d
ingress-nginx       ingress-nginx-controller-5d88495688-6lgx9   1/1     Running     1          11d
kube-system         coredns-74ff55c5b-hgzzx                     1/1     Running     1          13d
kube-system         etcd-minikube                               1/1     Running     1          13d
kube-system         kube-apiserver-minikube                     1/1     Running     1          13d
kube-system         kube-controller-manager-minikube            1/1     Running     1          13d
kube-system         kube-proxy-gtrcq                            1/1     Running     1          13d
kube-system         kube-scheduler-minikube                     1/1     Running     1          13d
kube-system         storage-provisioner                         1/1     Running     2          13d
```
{% endofftopic %}

```shell
kubectl get --all-namespaces deployment
```

{% offtopic title="Посмотреть ответ" %}
```shell
NAMESPACE           NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
ingress-nginx       ingress-nginx-controller   1/1     1            1           11d
kube-system         coredns                    1/1     1            1           13d
```
{% endofftopic %}

```shell
kubectl get --all-namespaces statefulset
kubectl get --all-namespaces job
kubectl get --all-namespaces cronjob
```

Также можно получить полную конфигурацию ресурса в yaml-формате, если добавить в команду `kubectl get` опцию `--output yaml`:

```shell
kubectl get deployment somedeployment --output yaml
```

В ответ отобразится следующее:
```yaml
...
kind: Deployment
metadata:
  name: somedeployment
...

kubectl get pod somepod -o yaml
...
kind: Pod
metadata:
  name: somepod
...
```

Чаще всего вам придется сталкиваться с ресурсом Deployment, поэтому рассмотрим его подробнее. Про остальные типы контроллеров можно прочитать в [документации к Kubernetes](https://kubernetes.io/docs/concepts/workloads/).

## Deployment

Создадим в любой директории файл `deployment.yaml` и опишем в нём Deployment для stateless-приложения:
```yaml
# используемая версия API kubernetes
apiVersion: apps/v1
kind: Deployment  # Тип ресурса.
metadata:
  name: werf-kube-basics  # Имя Deployment'а.
spec:
  replicas: 2  # можно развернуть несколько Pod'ов сразу
  selector:
    # label, по которому будет происходить выборка
    matchLabels:
      app: werf-kube-basics
  # секция, описывающая шаблон, по которому приложению будут назначаться значения параметров
  template:
    metadata:
      # label ресурса
      labels:
        app: werf-kube-basics
    spec:
      terminationGracePeriodSeconds: 60  # сколько секунд есть у процессов Pod'а на graceful-завершение после получения TERM-сигнала при остановке Pod'а
      # описание конфигурации контейнеров Pod'а:
      containers:
      - name: main  # имя первого контейнера
        image: alpine  # имя и тег образа контейнера
        command:
          # основная команда контейнера, начнёт выполняться при его запуске:
          - sh
          - -ec
          - |
            while true; do
              echo -e "HTTP/1.1 200 OK\n\nAlive.\nOur \$MY_ENV_VAR value is \"$MY_ENV_VAR\"." | nc -l -s 0.0.0.0 -p 80
            done
        ports:
        - containerPort: 80  # открытый порт, на котором будет слушать приложение
        env:
        - name: "MY_ENV_VAR"  # имя дополнительной переменной окружения
          value: "myEnvVarValue"  # значение дополнительной переменной окружения
        resources:
          # требования к ресурсам
          requests:
            cpu: 50m  # не запускаться на Node'ах, которые не могут выделить 0.05 CPU
            memory: 50Mi  # не запускаться на Node'ах, которые не могут выделить 50 МБ RAM
          limits:
            cpu: 50m  # hard-лимит на CPU, не позволит контейнеру расходовать больше
            memory: 50Mi  # hard-лимит на RAM, при превышении ООМ-киллер остановит процессы контейнера
        lifecycle:
          preStop:
            exec:
              command: ["/bin/trigger-graceful-shutdown-for-my-app"]  # запустится перед завершением контейнера
        startupProbe:  # проверка готовности контейнера, при нескольких неудачах контейнер перезапустится
          httpGet:
            path: /startup  # будет выполняться GET-запрос на http://<PodIP>/startup.
            port: 80
        readinessProbe:  # проверка готовности контейнера, запускается после startupProbe. При нескольких неудачах трафик перестаёт идти на Pod через Service
          httpGet:
            path: /readiness
            port: 80
        livenessProbe:  # проверка работоспособности контейнера, запускается после startupProbe. При нескольких неудачах контейнер перезапускается
          httpGet:
            path: /liveness
            port: 80
      # InitContainers используются для разовых задач, выполняющихся перед запуском основных контейнеров
      # основные контейнеры Pod'а не запустятся, пока не выполнятся initContainers
      initContainers:
        - name: wait-postgres
          image: postgres
          command: ["echo", "pg_isready", "-h", "postgres.example.com"]
```

> Этот пример не является примером того, как должен выглядеть production-ready Deployment. Лучшие практики по организации ресурсов для ваших приложений мы рассмотрим в следующих главах. Более подробное описание Deployment доступно в [официальной документации](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/).

Создадим Deployment-ресурс в кластере на основе созданного файла `deployment.yaml` с помощью `kubectl apply`:

```shell
kubectl apply -f deployment.yaml
```

> В этой главе для деплоя вместо `werf converge` мы будем использовать `kubectl apply`. Он удобен для быстрого и простого создания ресурсов в кластере, но для деплоя реального приложения необходимо использовать `werf converge`.

Убедимся, что наш Deployment создался:
```shell
kubectl get deployment werf-kube-basics
```

В ответ отобразится примерно следующее:
```shell
NAME                    READY   UP-TO-DATE   AVAILABLE   AGE
kubernetes-basics-app   0/2     2            0           25s
```

А теперь понаблюдаем за развертыванием Pod'ов, создаваемых нашим Deployment'ом:
```bash
kubectl get pod -l app=werf-kube-basics
```

В ответ отобразится примерно следующее:
```shell
NAME                                    READY   STATUS    RESTARTS   AGE
kubernetes-basics-app-6f59b7cf5-pjmrd   1/1     Running   0          109s
kubernetes-basics-app-6f59b7cf5-thzcc   1/1     Running   0          109s
```

## Service и Ingress

С помощью Deployment'а мы можем развернуть наше stateless-приложение, но если пользователям или другим приложениям потребуется связываться с этим приложением изнутри или снаружи кластера, то потребуются два дополнительных ресурса: Ingress и Service.

Создадим файл `ingress.yaml`:
```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: nginx
  name: werf-kube-basics
spec:
  rules:
  - host: werf-kube-basics  # домен, запросы на который будут обрабатываться в paths ниже
    http:
      paths:
      - path: /  # запросы с префиксом / (все запросы) перенаправятся на порт 80 нашего Service'а
        backend:
          serviceName: werf-kube-basics
          servicePort: 80
```

И создадим файл `service.yaml`:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: werf-kube-basics
spec:
  selector:
    app: werf-kube-basics  # этот Service перенаправляет трафик на Pod'ы с этим лейблом
  ports:
  - name: http
    port: 80  # перенаправить трафик с 80-го порта Service'а на 80-й порт Pod'а
```

А теперь создадим ресурсы в кластере на основе этих манифестов:
```shell
kubectl apply -f ingress.yaml -f service.yaml
```

Убедимся, что ресурс Service создался:
```shell
kubectl get service werf-kube-basics
```

В ответ получим следующее:
```shell
NAME                    TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
werf-kube-basics   ClusterIP   10.107.19.126   <none>        80/TCP    35s
```
Убедимся, что ресурс ingress создался:
```shell
kubectl get ingress werf-kube-basics
```
В ответ получим следующее:
```shell
NAME                    CLASS    HOSTS                               ADDRESS   PORTS   AGE
werf-kube-basics   <none>   werf-kube-basics.example.com             80      3m21s
```

Если несколько упростить, то эти два ресурса позволят HTTP-пакетам, приходящим на [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/how-it-works/), у которых есть заголовок `Host: werf-kube-basics`, быть перенаправленными на 80-й порт Service'а `werf-kube-basics`, а оттуда — на 80-й порт одного из Pod'ов нашего Deployment'а. В конфигурации по умолчанию Service будет перенаправлять запросы на все Pod'ы Deployment'а поровну.

Обратимся к нашему приложению через Ingress:

```shell
curl http://werf-kube-basics
```

В ответ отобразится следующее:
```shell
Alive.
Our $MY_ENV_VAR value is "myEnvVarValue".
```

При этом Service-ресурсы нужны не только для связи Ingress'ов и приложения. Service-ресурсы также дают возможность ресурсам внутри кластера общаться между собой. При создании Service'а создается доменное имя `<ServiceName>.<NamespaceName>.svc.cluster.local`, доступное изнутри кластера. Также Service доступен и по более коротким доменным именам:
* `<ServiceName>` — при обращении из того же Namespace'а,
* `<ServiceName>.<NamespaceName>` — при обращении из другого.

Создадим новый контейнер, не имеющий отношения к нашему приложению:
```shell
kubectl run werf-temporary-deployment --image=alpine --rm -it -- sh
```

В запустившемся контейнере обратимся к нашему приложению через Service:

```shell
apk add curl  # Установим curl внутри контейнера.
curl http://werf-kube-basics  # Обратимся к одному из Pod'ов нашего приложения через Service.
```

В ответ отобразится следующее:
```
Alive.
Our $MY_ENV_VAR value is "myEnvVarValue".
```

> Использование Ingress-ресурсов — не единственный способ получить доступ к приложению снаружи кластера. Service'ы типа `LoadBalancer` и `NodePort` позволяют предоставить доступ к приложению снаружи и без Ingress'ов. Почитать подробнее про Service'ы можно в [официальной документации](https://kubernetes.io/docs/concepts/services-networking/service/). А про Ingress'ы — [здесь](https://kubernetes.io/docs/concepts/services-networking/ingress/).

## Очистка

Удалим созданные нами ресурсы, так как они больше не понадобятся:
```shell
kubectl delete -f deployment.yaml -f service.yaml -f ingress.yaml
```

В ответ получим сообщение об успешном удалении ресурсов:
```shell
deployment.apps "kubernetes-basics-app" deleted
service "kubernetes-basics-app" deleted
ingress.networking.k8s.io "kubernetes-basics-app" deleted
```