## Выбираем один из облачных провайдеров

Существует огромное количество услуг, позиционирующихся как [Managed Kubernetes](https://www.google.com/search?q=managed+kubernetes). Часть из них включает серверные мощности, часть — нет.

Проще всего взять услугу в AWS (EKS) или Google Cloud (GKE). При первой регистрации они представляют бонус, которого должно хватить на неделю-другую работы с кластером. Однако для регистрации понадобится ввести данные банковской карты.

Альтернативой может стать использование Yandex.Cloud: здесь также предоставляется Managed Kubernetes и тоже требуется ввести данные карты. Но в качестве карты можно воспользоваться, например, виртуальной картой от Яндекс.Денег, позволяющей проводить «подписочные» платежи внутри России.

Также можно попробовать развернуть Kubernetes самостоятельно [в Hetzner на основании их статьи](https://community.hetzner.com/tutorials/install-kubernetes-cluster) — это один из самых дешёвых облачных провайдеров. Однако надо понимать, что вам придётся самостоятельно разбираться с большим пластом работ по администрированию платформы. Для production найдите надёжного провайдера или команду поддержки.
Если вы ещё ни разу не устанавливали Kubernetes самостоятельно и/или не обладаете опытом системного администрирования — не стоит пытаться освоить такую объемную тему в рамках этого самоучителя.

Вне зависимости от выбранного облачного провайдера вам потребуется получить ключи доступа к Kubernetes-кластеру в виде файла `.kube/config`.

<!--  .kube/config -->
{% offtopic title="Что за .kube/config?" %}
Это файл, хранящий реквизиты доступа к кластеру. Если разместить его по пути `~/.kube/config`, утилиты, работающие с Kubernetes, смогут подключиться к кластеру.

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

Обратите внимание, что кроме ключей для подключения здесь содержится путь до API Kubernetes — например, `server: https://127.0.0.1:6445`. Важно, чтобы этот путь был корректным. Некоторые утилиты, генерирующие `.kube/config`, не подставляют туда публичный IP — тогда его нужно скорректировать вручную.

{% endofftopic %}
<!-- / .kube/config -->

### Проверка работоспособности и доступа к кластеру

Проверьте свой кластер по чек-листу в начале этой статьи.

Итогом всех манипуляций должна стать возможность получить доступ к кластеру с помощью утилиты `kubectl` (возможно, эту утилиту придётся установить отдельно). Например, вызов:

```shell
kubectl get ns
```

… покажет список всех namespace'ов в кластере, а не сообщение об ошибке.

### Ingress

Нужно установить в кластер Nginx Ingress. На сайте решения есть [подробные инструкции](https://kubernetes.github.io/ingress-nginx/deploy/) для ключевых провайдеров.

### Registry

Скорее всего облачный провайдер registry как одну из услуг, но не всегда.

{% offtopic title="Мой провайдер не предоставляет registry, как быть?" %}

Вы можете развернуть Docker Registry самостоятельно на отдельно стоящей виртуальной машине или воспользоваться каким-то облачным решением. werf поддерживает [множество имплементаций]({{ site.docsurl }}/documentation/advanced/supported_registry_implementations.html) реестров для контейнеров.

{% endofftopic %}

Обеспечьте сетевую связность кластера и registry: опять же, скорее всего облачный провайдер решает эту проблему самостоятельно, либо в нём есть специальные мануалы на этот счёт.

### Hosts

В самоучителе предполагается, что кластер (вернее, его Nginx Ingress) доступен по адресу `example.com`, а registry — по адресу `registry.example.com`. Мы ожидаем, что вы самостоятельно выберете и настроите нужные DNS-записи, а при прохождении самоучителя будете подставлять в код свои домены.

### Авторизация в registry

Для того, чтобы werf смог загрузить собранный образ в registry, нужно на локальной машине авторизоваться с помощью `docker login` примерно так:

```shell
docker login registry.example.com -u <account_login> -p <account_password>
```

где `<account_login>` — логин для доступа к registry, а `<account_password>` — пароль.
