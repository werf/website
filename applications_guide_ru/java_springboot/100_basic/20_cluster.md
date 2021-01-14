---
title: Подготовка кластера
permalink: java_springboot/100_basic/20_cluster.html
---

Для дальнейшей работы нам понадобятся: Kubernetes кластер и хранилище образов (registry).

**Настройка окружения — это сложная задача**. Не рекомендуем закапываться в решение проблем своими руками: используйте существующие на рынке услуги, позовите на помощь коллег, которые имеют практику настройки инфраструктуры, или задавайте вопросы в [Telegram-сообществе werf](https://t.me/werf_ru).

В целях обучения мы рекомендуем воспользоваться одним из двух вариантов:

- Если у вас достаточно мощный компьютер — установить и настроить minikube
- В противном случае — воспользоваться предложением одного из облачных провайдеров

Также стоит отметить, что для серьёзной разработки понадобится "настоящий" кластер. Использовать minikube не получится. Но возможно к этому моменту вы найдёте, кто поднимет и настроит Kubernetes за вас. Подробнее к вопросу кластера и registry мы вернёмся в главе "Работа с инфраструктурой". 

Выберите, как будете реализовывать окружение:

<div style="display: flex; justify-content: space-between; margin: 0 10px 0 20px;">
<div class="button__blue button__blue_inline expand_columns_button" id="minikube_button"><a href="#">Установка minikube</a></div>
<div class="button__blue button__blue_inline expand_columns_button" id="cloud_provider_button"><a href="#">Использование cloud provider</a></div>
<div class="button__blue button__blue_inline expand_columns_button" id="has_cluster_button"><a href="#">У меня уже есть кластер</a></div>
</div>


{% expandonclick id="minikube_button__content" %}

## Установка Minikube

Minikube — это "облегчённая" версия Kubernetes, работающая только на одном узле.

- [Установите minikube](https://minikube.sigs.k8s.io/docs/start/) и запустите его:
```bash
minikube start
```
- Включите дополнение minikube registry: 
```bash
minikube addons enable registry
```
- Запустите registry с привязкой к порту 5000: 
```bash
kubectl -n kube-system expose rc/registry --type=ClusterIP --port=5000 --target-port=5000 --name=werf-registry --selector='actual-registry=true'
```
- В отдельном терминале пробросьте порт:
```bash
kubectl port-forward --namespace kube-system service/werf-registry 5000
```

В результате автоматически будет создан файл `.kube/config` с ключами доступа к локальному кластеру и werf сможет подключиться к локальному registry.
{% endexpandonclick %}

{% expandonclick id="cloud_provider_button__content" %}

## Выбираем один из cloud provider-ов

Существует огромное количество услуг, позиционирующихся как [Managed Kubernetes](https://www.google.com/search?q=managed+kubernetes). Часть из них включает серверные мощности, часть — нет.

Проще всего взять услугу в AWS (EKS) или Google Cloud (GKE). При первой регистрации они представляют бонус, которого должно хватить на неделю-другую работы с кластером. Однако для регистрации понадобится ввести данные банковской карты.

Альтернативой может стать использование Yandex.Cloud: здесь также предоставляется Managed Kubernetes и тоже требуется ввести данные карты. Но в качестве карты можно воспользоваться, например, виртуальной картой от Яндекс.Денег, позволяющей проводить «подписочные» платежи внутри России.

Также можно попробовать развернуть Kubernetes самостоятельно [в Hetzner на основании их статьи](https://community.hetzner.com/tutorials/install-kubernetes-cluster) — это один из самых дешёвых облачных провайдеров. Однако надо понимать, что вам придётся самостоятельно разбираться с большим пластом работ по администрированию платформы. Для production найдите надёжного провайдера или команду поддержки.
Если вы ещё ни разу не устанавливали Kubernetes самостоятельно и/или не обладаете опытом системного администрирования — не стоит пытаться освоить такую объемную тему в рамках этого самоучителя.

Вне зависимости от выбранного cloud provider-а, необходимо в итоге:

1 — Добыть ключи доступа к Kubernetes-кластеру в виде файла `.kube/config`
<!--  .kube/config -->
{% offtopic title="Что за .kube/config?" %}
Это файл, хранящий реквизиты доступа к кластеру. Если разместить его по пути `~/.kube/config` - утилиты, работающие с Kubernetes смогут подключиться к кластеру.

Файл выглядит примерно так:

{% snippetcut name=".kube/config" url="#" %}
{% raw %}
```yaml
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ALotOfNumbersAndLettersAndSoOnAveryVERYveryLongStringInBase64=
    server: https://127.0.0.1:6445
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: ManyNumbersAndLettersAVERYveryVerylongString=
    client-key-data: ManyLettersAndNumbersAveryVeryVERYlongString=
```
{% endraw %}
{% endsnippetcut %}

Обратите внимание, что кроме ключей для подключения, здесь содержится путь до API Kubernetes, например, `server: https://127.0.0.1:6445`. Важно, чтобы этот путь был корректным. Некоторые утилиты, генерирующие `.kube/config` не подставляют туда публичный IP и его нужно скорректировать руками.

{% endofftopic %}
<!-- / .kube/config -->
2 — Развернуть registry в облаке (скорее всего cloud provider предоставляет это как одну из услуг)

{% offtopic title="Мой провайдер не предоставляет registry, как быть?" %}

Вы можете развернуть docker registry самостоятельно или воспользоваться каким-то облачным решением. werf поддерживает [множество имплементаций]({{ site.docsurl }}/documentation/advanced/supported_registry_implementations.html) 

{% endofftopic %}

3 — Обеспечить доступ кластера к registry (скорее всего cloud provider решает эту проблему самостоятельно, либо в нём есть специальные мануалы на этот счёт)

{% offtopic title="Как сделать доступ кластера к registry?" %}

Нужно узнать ключи доступа к registry и прописать их в кластере. Полученные ключи должны быть прописаны в **каждом** пространстве имён в Kubernetes, куда осуществляется деплой, в виде объекта Secret.

Сперва создайте нужное пространство имён:

```bash
kubectl create namespace werf-guided-project
```

А затем создайте в нём Secret:

```bash
kubectl create secret docker-registry registrysecret -n <namespace> --docker-server=<registry_domain> --docker-username=<account_email> --docker-password=<account_password> --docker-email=<account_email>
```

Здесь:

- `<namespace>` — название пространства имён в Kubernetes (например, `werf-guided-project`);
- `<registry_domain>` — домен Registry;
- `<account_email>` — логин в Registry;
- `<account_password>` — пароль в Registry.

Позже, в главе "Деплой приложения" мы будем указывать имя созданного секрета (`registrysecret`) в конфигурации объектов в Kubernetes:

{% snippetcut name=".helm/templates/deployment.yaml" url="#" %}
{% raw %}
```yaml
    spec:
      imagePullSecrets:
      - name: "registrysecret"
```
{% endraw %}
{% endsnippetcut %}

На данный момент достаточно того, чтобы секрет был создан в пространстве имён `werf-guided-project`.

{% endofftopic %}

Итогом всех манипуляций должна стать возможность получить доступ к кластеру с помощью утилиты `kubectl` (возможно эту утилиту придётся установить отдельно), например:

```bash
kubectl get ns
```

Покажет вам список всех namespace-ов в кластере, а не сообщение об ошибке.
{% endexpandonclick %}

{% expandonclick id="has_cluster_button__content" %}

## У меня уже есть кластер

Если у вас уже есть кластер, вероятно, вы уже сталкивались с его конфигурированием. Убедитесь, что:

- У вас есть доступ с локального компьютера к кластеру
- У вас есть registry
- У кластера есть доступ к registry (в самоучителе предполагается, что в кластере есть соответствующий Secret)
{% offtopic title="Как оформить секрет с доступом?" %}

Нужно узнать ключи доступа к registry и прописать их в кластере. Полученные ключи должны быть прописаны в **каждом** пространстве имён в Kubernetes, куда осуществляется деплой, в виде объекта Secret. Сделать это можно, выполнив команду:

```bash
kubectl create secret docker-registry registrysecret -n <namespace> --docker-server=<registry_domain> --docker-username=<account_login> --docker-password=<account_password> --docker-email=<account_login>
```

Здесь:

- `<namespace>` — название пространства имён в Kubernetes (например, `werf-guided-project`);
- `<registry_domain>` — домен Registry;
- `<account_login>` — логин в Registry;
- `<account_password>` — пароль в Registry.

Позже, в главе "Деплой приложения" мы будем указывать имя созданного секрета (`registrysecret`) в конфигурации объектов в Kubernetes:

{% snippetcut name=".helm/templates/deployment.yaml" url="#" %}
{% raw %}
```yaml
    spec:
      imagePullSecrets:
      - name: "registrysecret"
```
{% endraw %}
{% endsnippetcut %}

На данный момент достаточно того, чтобы секрет был создан в пространстве имён `werf-guided-project`.

{% endofftopic %}

{% endexpandonclick %}

## Авторизация в Registry

Для того, чтобы werf смог загрузить собранный образ в registry — нужно авторизоваться с помощью `docker login` примерно так:

```bash
docker login <registry_domain> -u <account_login> -p <account_password>
```

## Домен

В самоучителе предполагается, что кластер доступен по адресу `example.com`. Именно этот домен и его поддомены указаны в дальнейшем в конфигах, в случае, если вы будете использовать другой — скорректируйте конфигурацию самостоятельно.

Если ваш кластер установлен локально — пропишите в локальном файле `/etc/hosts` строку вида

```
127.0.0.1           example.com
```

Если вы используете не локальный кластер — просьба сконфигурировать правильные DNS записи самостоятельно.

## Ingress

Не во всех инсталляциях кластера "из коробки" есть балансер, возможно его придётся дополнительно установить. К примеру:

- в Minikube нужно [включать addon](https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/#enable-the-ingress-controller)
- в Docker Desktop на MacOs нужно в явном виде [установить nginx ingress вручную](https://kubernetes.github.io/ingress-nginx/deploy/) 

Если этого не сделать, то когда вы в следующей главе будете пытаться открыть в браузере example.com — приложения там не увидите.

{% offtopic title="Как убедиться, что с балансером всё хорошо?" %}

По умолчанию можно считать, что всё хорошо, если не уверены, что тут написано — вернитесь к этому пункту позже.

- Балансер установлен
- Pod с балансером корректно поднялся (например, для nginx-ingress, это можно посмотреть так: `kubectl -n ingress-nginx get po`)
- На 80-ом порту (это можно посмотреть с помощью `lsof -n | grep LISTEN`) работает нужное приложение
{% endofftopic %}


<div id="go-forth-button">
    <go-forth url="30_deploy.html" label="Деплой приложения" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
