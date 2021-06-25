---
title: Подготовка окружения
permalink: rails/100_basic/20_cluster.html
---

В этой главе мы подготовим кластер Kubernetes, Container Registry и локальное окружение для развертывания приложений.

## Развертывание кластера

Выберите, каким образом будет развернут Kubernetes:

<div class="tabs">
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__minikube_windows')">Minikube (Windows)</a>
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__minikube_macos')">Minikube (macOS)</a>
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__minikube_linux')">Minikube (Linux)</a>
{% comment %} TODO(lesikov): раскомментить как переработаю 20_cluster_has_cluster.md.
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__ihave')">Свой кластер</a>
{% endcomment %}
</div>

<div id="tab__install__minikube_windows" class="tabs__content tabs__install__content" markdown="1">
{% include_relative 20_cluster_minikube_windows.md %}
</div>
<div id="tab__install__minikube_macos" class="tabs__content tabs__install__content" markdown="1">
{% include_relative 20_cluster_minikube_macos.md %}
</div>
<div id="tab__install__minikube_linux" class="tabs__content tabs__install__content" markdown="1">
{% include_relative 20_cluster_minikube_linux.md %}
</div>

{% comment %} TODO(lesikov): раскомментить как переработаю 20_cluster_has_cluster.md.
<div id="tab__install__ihave" class="tabs__content tabs__install__content" markdown="1">
{% include_relative 20_cluster_has_cluster.md %}
</div>
{% endcomment %}

## Проверка

Для проверки работоспособности необходимо открыть страницу в браузере или использовать `сurl` в консоли:

```shell
curl http://example.com
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

> Если у вас уже поднято приложение, на которое смотрит Ingress с доменом `example.com`, то вместо ошибки 404 вы получите ответ от своего приложения.

## Настройка Container Registry

Далее мы будем использовать Docker Hub Container Registry, но для этого руководства подойдет и любой другой Registry с TLS и аутентификацией ([GitHub Container Registry](https://github.com/features/packages), [GitLab Container Registry](https://docs.gitlab.com/ee/user/packages/container_registry/), ...).

Регистрируемся на [Docker Hub](https://hub.docker.com/signup), после чего [создаём приватный репозиторий](https://hub.docker.com/repository/create) с именем `werf-guided-rails`, в котором будем хранить собираемые образы.

С помощью `docker login` получаем доступ с текущего компьютера к новому репозиторию, вводя логин и пароль от нашего пользователя на Docker Hub:
```bash
$ docker login
Username: <имя пользователя Docker Hub>
Password: <пароль пользователя Docker Hub>
Login Succeeded
```

Создаём Secret в кластере, который поможет получить доступ к новому репозиторию уже нашим будущим приложениям:
```bash
kubectl create namespace werf-guided-rails  # namespace для Secret'а ещё не существует, создадим его
kubectl create secret docker-registry registrysecret \
  --docker-server='https://index.docker.io/v1/' \
  --docker-username='<имя пользователя Docker Hub>' \
  --docker-password='<пароль пользователя Docker Hub>'
```

> Стоит обратить внимание на опцию `--docker-server`, параметр которой должен соответствовать адресу используемого 
> регистри. К примеру, для GitHub Container Registry необходимо иcпользовать `ghcr.io`, а для Docker Hub можно обойтись 
> без опции, использовать значение по умолчание.

Теперь окружение для работы готово.

## Рабочее окружение работало, но перестало

<div class="tabs">
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__stopped_working_minikube_windows')">Windows</a>
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__stopped_working_minikube_macos')">macOS</a>
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__stopped_working_minikube_linux')">Linux</a>
</div>

<div id="tab__install__stopped_working_minikube_windows" class="tabs__content tabs__install__content" markdown="1">
{% include 100_basic/20_cluster_stopped_working_minikube_windows.md.liquid %}
</div>
<div id="tab__install__stopped_working_minikube_macos" class="tabs__content tabs__install__content" markdown="1">
{% include 100_basic/20_cluster_stopped_working_minikube_macos.md.liquid %}
</div>
<div id="tab__install__stopped_working_minikube_linux" class="tabs__content tabs__install__content" markdown="1">
{% include 100_basic/20_cluster_stopped_working_minikube_linux.md.liquid %}
</div>
