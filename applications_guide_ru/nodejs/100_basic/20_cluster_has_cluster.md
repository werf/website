## У меня уже есть кластер

Если у вас уже есть кластер, вероятно, вы уже сталкивались с его конфигурированием. Проверьте свой кластер по чек-листу вначале этой статьи.

Вне зависимости от того, как установлен кластер, необходимо в добыть ключи доступа к Kubernetes-кластеру в виде файла `.kube/config`

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


### Проверка работоспособности и доступа к кластеру

Проверьте свой кластер по чек-листу вначале этой статьи.

Итогом всех манипуляций должна стать возможность получить доступ к кластеру с помощью утилиты `kubectl` (возможно эту утилиту придётся установить отдельно), например:

```bash
kubectl get ns
```

Покажет вам список всех namespace-ов в кластере, а не сообщение об ошибке.

### Ingress

В вашем кластере должен быть установлен nginx ingress. Если нет, то на сайте решения есть [подробные инструкции для ключевых cloud provider-ов](https://kubernetes.github.io/ingress-nginx/deploy/).

### Registry

werf поддерживает [множество имплементаций]({{ site.docsurl }}/documentation/advanced/supported_registry_implementations.html) registry.

Убедитесь, что есть сетевая связность кластера и registry.

### Hosts

В самоучителе предполагается, что кластер (вернее, его nginx ingress) доступен по адресу `example.com`, а registry — по адресу `registry.example.com`. Мы предполагаем, что вы самостоятельно выберете и настроите нужные ns-записи, а при прохождении самоучителя будете подставлять в код ваши домены.

### Авторизация в Registry

Для того, чтобы werf смог загрузить собранный образ в registry — нужно на локальной машине авторизоваться с помощью `docker login` примерно так:

```bash
docker login <registry_domain> -u <account_login> -p <account_password>
```
