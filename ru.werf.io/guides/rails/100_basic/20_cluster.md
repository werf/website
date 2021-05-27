---
title: Подготовка окружения
permalink: rails/100_basic/20_cluster.html
---

Выберите, с помощью чего вы хотите развернуть тестовый Kubernetes-кластер:

<div class="tabs">
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__docker')">Docker Desktop (Windows, macOS)</a>
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__minikube')">Minikube (Linux, Windows, macOS)</a>
{% comment %} TODO(lesikov): раскомментить как переработаю 20_cluster_has_cluster.md.
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__ihave')">Свой кластер</a>
{% endcomment %}
</div>

<div id="tab__install__docker" class="tabs__content tabs__install__content" markdown="1">
{% include_relative 20_cluster_docker_desktop.md %}
</div>

<div id="tab__install__minikube" class="tabs__content tabs__install__content" markdown="1">
{% include_relative 20_cluster_minikube.md %}
</div>


{% comment %} TODO(lesikov): раскомментить как переработаю 20_cluster_has_cluster.md.
<div id="tab__install__ihave" class="tabs__content tabs__install__content" markdown="1">
{% include_relative 20_cluster_has_cluster.md %}
</div>
{% endcomment %}

## Проверяем кластер

Если hosts-файл правильно настроен, то эти команды должны показать IP, соответствующие тем, которые мы указали в hosts-файле:
```shell
nslookup example.com
nslookup registry.example.com
```

Проверим, загружаются ли образы в Container Registry:
```shell
docker pull busybox
docker tag busybox registry.example.com:80/busybox
docker push registry.example.com:80/busybox
```

А теперь проверим, работает ли NGINX Ingress Controller:
```shell
curl example.com
```
… должно выдать ошибку 404 от Nginx:
```html
<html>
<head><title>404 Not Found</title></head>
<body>
<center><h1>404 Not Found</h1></center>
<hr><center>nginx</center>
</body>
</html>
```
Это значит, что NGINX Ingress Controller успешно запущен.

_После того, как мы задеплоим в кластер наше приложение, вместо ошибки 404 на `example.com` нам будет отвечать наше приложение._

<div id="go-forth-button">
    <go-forth url="30_deploy.html" label="Деплой приложения" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
