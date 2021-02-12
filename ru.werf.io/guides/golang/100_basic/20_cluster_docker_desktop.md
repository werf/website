## Docker Desktop

### Установка

Установите Docker Desktop на [Windows](https://docs.docker.com/docker-for-windows/install/) или [macOS](https://docs.docker.com/docker-for-mac/install/).

Включите Kubernetes и настройте выделенные на него ресурсы.

{% offtopic title="Почему ресурсы — это важно?" %}
Если у кластера будет сликшмо мало ресурсов, не сможет запуститься ваше приложение, Ingress или даже какой-то из необходимых самому оркестратору системных компонентов. Без опыта администрирования кластера разобраться в таких проблемах будет очень тяжело.

Составители самоучителя тестировали приложения на конфигурации 6 CPU, 6 ГБ памяти, swap на 1,5 ГБ, образ диска на 24 ГБ.
{% endofftopic %}

В результате автоматически будет создан файл `.kube/config` с ключами доступа к локальному кластеру и werf сможет подключиться к локальному registry.

Итогом этих манипуляций должна стать возможность получить доступ к кластеру с помощью утилиты `kubectl` (возможно, эту утилиту придётся установить отдельно). Например, вызов:

```shell
kubectl get ns
```

… покажет список всех namespace'ов в кластере, а не сообщение об ошибке.

### Ingress

Нужно в явном виде [установить Nginx Ingress вручную](https://kubernetes.github.io/ingress-nginx/deploy/).

В случае Docker Desktop иногда бывают сложности с доступом к Ingress: порты могут не проброситься на хост-машину. Чтобы быть уверенным, что Ingress корректно работает, проверьте:

- Pod с Ingress-контроллером корректно поднялся (для nginx-ingress это можно посмотреть так: `kubectl -n ingress-nginx get po`);
- наличие службы на 80-м порту (это можно посмотреть с помощью `lsof -n | grep LISTEN` или [аналогичным способом](https://www.google.com/search?q=check+used+ports&oq=check+used+ports)).

Если на 80-м порту (HTTP) нет Kubernetes — возможно, нужно пробросить порт вручную. Посмотрите имя Pod'а с ingress-controller:

```shell
kubectl -n ingress-nginx get po
```

… и затем запустите проброс порта:

```shell
kubectl port-forward --address 0.0.0.0 -n ingress-nginx ingress-nginx-controller-<тут_будут_буквы_цифры> 80:80
```

После описанных операций проверьте, что появилось на 80-м порту (например, `lsof -n | grep LISTEN`).

### Registry

Docker Desktop не содержит встроенного registry. Самый простой способ — это запустить его вручную как Docker-образ:

```shell
docker run -d -p 5000:5000 --restart=always --name registry registry:2
```

Обратите внимание, что в этом случае по умолчанию registry будет работать без SSL-шифрования. А значит, при работе с werf потребуется добавлять [параметр](https://werf.io/documentation/reference/cli/werf_managed_images_add.html#options) `--insecure-registry=true`.

### Hosts

В самоучителе предполагается, что кластер (вернее, его Nginx Ingress) доступен по адресу `example.com`, а registry — по адресу `registry.example.com`. Именно этот домен и его поддомены указаны в дальнейшем в конфигах. В случае, если вы будете использовать другие адреса, скорректируйте конфигурацию самостоятельно.

Пропишите в локальном файле `/etc/hosts` строки вида:

```
127.0.0.1           example.com
127.0.0.1           registry.example.com
```

### Авторизация в Registry

Если вы запустили registry локально, как приведено в примере выше, то ваш локальный registry работает без пароля и ничего делать не нужно.

{% offtopic title="Я использую внешний registry или установил логин/пароль" %}
Для того, чтобы werf смог загрузить собранный образ в registry, нужно на локальной машине авторизоваться с помощью `docker login` примерно так:

```shell
docker login <registry_domain> -u <account_login> -p <account_password>
```
{% endofftopic %}
