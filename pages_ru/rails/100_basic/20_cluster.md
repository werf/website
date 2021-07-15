---
title: Подготовка окружения
permalink: rails/100_basic/20_cluster.html
chapter_initial_prepare_cluster: false
chapter_initial_prepare_repo: false
description: |
  В этой главе мы подготовим кластер Kubernetes, Container Registry и локальное окружение для развертывания приложений.
---

## Развертывание кластера

Выберите, каким образом будет развернут Kubernetes:

<div class="tabs">
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn active" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__minikube_linux')">Linux: minikube</a>
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__minikube_macos')">macOS: minikube</a>
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__minikube_windows')">Windows: minikube</a>
{% comment %} TODO(lesikov): раскомментить как переработаю 20_cluster_has_cluster.md.
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__ihave')">Свой кластер</a>
{% endcomment %}
</div>

<div id="tab__install__minikube_linux" class="tabs__content tabs__install__content active" markdown="1">
{% include_relative 20_cluster_minikube_linux.md %}
</div>
<div id="tab__install__minikube_macos" class="tabs__content tabs__install__content" markdown="1">
{% include_relative 20_cluster_minikube_macos.md %}
</div>
<div id="tab__install__minikube_windows" class="tabs__content tabs__install__content" markdown="1">
{% include_relative 20_cluster_minikube_windows.md %}
</div>

{% comment %} TODO(lesikov): раскомментить как переработаю 20_cluster_has_cluster.md.
<div id="tab__install__ihave" class="tabs__content tabs__install__content" markdown="1">
{% include_relative 20_cluster_has_cluster.md %}
</div>
{% endcomment %}

## Проверка

Для проверки работоспособности выполните:

```shell
curl http://werf-guide-app/ping
```

Если всё работает как надо, то NGINX Ingress Controller вернёт ошибку 404:

```html
<html>
<head><title>404 Not Found</title></head>
<body>
<center><h1>404 Not Found</h1></center>
<hr><center>nginx</center>
</body>
</html>
```

> Если у вас уже поднято приложение, на которое смотрит Ingress с доменом `werf-guide-app`, то вместо ошибки 404 вы получите ответ от своего приложения.

## Настройка Container Registry

Далее мы будем использовать Docker Hub Container Registry, но для этого руководства подойдет и любой другой Registry с TLS и аутентификацией ([GitHub Container Registry](https://github.com/features/packages), [GitLab Container Registry](https://docs.gitlab.com/ee/user/packages/container_registry/), ...).

### Авторизация в Docker Hub
Регистрируемся на [Docker Hub](https://hub.docker.com/signup), после чего [создаём приватный репозиторий](https://hub.docker.com/repository/create) с именем `werf-guide-app`, в котором будем хранить собираемые образы.

С помощью `docker login` получим доступ с текущего компьютера к новому репозиторию, введя логин и пароль от нашего пользователя на Docker Hub:
```shell
docker login
# Введем имя пользователя Docker Hub
Username: <имя пользователя Docker Hub>
# Введем пароль пользователя Docker Hub
Password: <пароль пользователя Docker Hub>
```

В случае успеха получим ответ:
```shell
Login Succeeded
```

### Создание namespace
Создадим namespace для нашего приложения:
```shell
kubectl create namespace werf-guide-app
```

В ответ отобразится следующее:
```shell
namespace/werf-guide-app created
```

### Создание Secret
Создадим Secret в кластере, который обеспечит доступ к нашему новому приватному репозиторию уже нашему приложению:
```shell
kubectl create secret docker-registry registrysecret \
  --docker-server='https://index.docker.io/v1/' \
  --docker-username='<имя пользователя Docker Hub>' \
  --docker-password='<пароль пользователя Docker Hub>'
```

В ответ отобразится следующее:
```shell
secret/registrysecret created
```

{% offtopic title="Если нужно пересоздать secret" %}
kubectl delete secret docker-registry
{% endofftopic %}


> Стоит обратить внимание на опцию `--docker-server`, параметр которой должен соответствовать адресу используемого
> registry. К примеру, для GitHub Container Registry необходимо иcпользовать `ghcr.io`, а для Docker Hub можно обойтись 
> без опции и использовать значение по умолчанию.

Теперь окружение для работы готово.

## Рабочее окружение работало, но перестало

{% include 100_basic/20_cluster_frequent_issues.md.liquid %}
