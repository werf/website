---
title: Подготовка окружения
permalink: rails/100_basic/20_cluster.html
---

Выберите, каким образом будет развернут ваш Kubernetes-кластер:

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

### Настройка Container Registry

Далее мы будем использовать Docker Hub Container Registry, но для этого руководства подойдет любой Registry с TLS и аутентификацией (GitHub Container Registry, GitLab Container Registry, ...).

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
kubectl create secret docker-registry registrysecret ---docker-username=<имя пользователя Docker Hub> --docker-password=<пароль пользователя Docker Hub>
```

Теперь окружение для работы готово.
