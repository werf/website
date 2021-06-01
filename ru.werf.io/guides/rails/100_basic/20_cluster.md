---
title: Подготовка окружения
permalink: rails/100_basic/20_cluster.html
---

Выберите, каким образом будет развернут ваш Kubernetes-кластер:

<div class="tabs">
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__minikube')">Minikube (Windows, macOS, Linux)</a>
{% comment %} TODO(lesikov): раскомментить как переработаю 20_cluster_has_cluster.md.
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__ihave')">Свой кластер</a>
{% endcomment %}
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

Проверим, работает ли NGINX Ingress Controller и Registry:
```shell
curl http://registry.example.com/v2/
```

Проверим, можем ли мы загрузить в Registry образы:
```shell
docker pull busybox
docker tag busybox registry.example.com:80/busybox
docker push registry.example.com:80/busybox
```

Если всё в порядке, то Kubernetes и окружение для работы готовы.

<div id="go-forth-button">
    <go-forth url="30_deploy.html" label="Деплой приложения" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
