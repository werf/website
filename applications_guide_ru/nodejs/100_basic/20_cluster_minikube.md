## Minikube

Это "облегчённая" версия Kubernetes, работающая только на одном узле.

### Установка

[Установите minikube](https://minikube.sigs.k8s.io/docs/start/) и запустите его:

```bash
minikube start
```

В результате автоматически будет создан файл `.kube/config` с ключами доступа к локальному кластеру и werf сможет подключиться к локальному registry.

Итогом этих манипуляций должна стать возможность получить доступ к кластеру с помощью утилиты `kubectl` (возможно эту утилиту придётся установить отдельно), например:

```bash
kubectl get ns
```

Покажет вам список всех namespace-ов в кластере, а не сообщение об ошибке.

### Ingress

В minikube нужно [включить addon](https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/#enable-the-ingress-controller)

{% offtopic title="Как убедиться, что с Ingress всё хорошо?" %}

По умолчанию можно считать, что всё хорошо, если не уверены, что тут написано — вернитесь к этому пункту позже.

- Балансер установлен
- Pod с балансером корректно поднялся (для nginx-ingress, это можно посмотреть так: `kubectl -n ingress-nginx get po`)
- На 80-ом порту (это можно посмотреть с помощью `lsof -n | grep LISTEN`) работает нужное приложение
  {% endofftopic %}

### Registry

Включите addon minikube registry:

```bash
minikube addons enable registry
```

Запустите registry с привязкой к порту 5000: 
```bash
kubectl -n kube-system expose rc/registry --type=ClusterIP --port=5000 --target-port=5000 --name=werf-registry --selector='actual-registry=true'
```

В отдельном терминале пробросьте порт:
```bash
kubectl port-forward --namespace kube-system service/werf-registry 5000
```

### Hosts

В самоучителе предполагается, что кластер доступен по адресу `example.com`, а registry — по адресу `registry.example.com`. Именно этот домен и его поддомены указаны в дальнейшем в конфигах, в случае, если вы будете использовать другой — скорректируйте конфигурацию самостоятельно.

Пропишите в локальном файле `/etc/hosts` строки вида

```
127.0.0.1           example.com
127.0.0.1           registry.example.com
```

### Авторизация в Registry

Для того, чтобы werf смог загрузить собранный образ в registry — нужно на локальной машине авторизоваться с помощью `docker login` примерно так:

```bash
docker login <registry_domain> -u <account_login> -p <account_password>
```
