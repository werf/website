### Установка и запуск minikube

Установим minikube, [следуя инструкциям здесь](https://minikube.sigs.k8s.io/docs/start/).

Разрешаем доступ в **Container Registry без TLS для docker**:
{% offtopic title="Windows" %}
В файл, находящийся по умолчанию в `%programdata%\docker\config\daemon.json`, добавим новый ключ:
```json
{
  "insecure-registries": ["registry.example.com:80"]
}
```
Перезапустим Docker Desktop через меню, открывающееся правым кликом по иконке Docker Desktop в трее.
{% endofftopic %}
{% offtopic title="macOS" %}
В файл, по умолчанию находящийся в `/etc/docker/daemon.json`, добавим новый ключ:
```json
{
  "insecure-registries": ["registry.example.com:80"]
}
```
Перезапустим Docker через меню Docker Desktop.
{% endofftopic %}
{% offtopic title="Linux" %}
В файл, по умолчанию находящийся в `/etc/docker/daemon.json`, добавим новый ключ:
```json
{
  "insecure-registries": ["registry.example.com:80"]
}
```
Перезапустим Docker:
```shell
sudo systemctl restart docker
```
{% endofftopic %}

Чтобы разрешить доступ в **Container Registry без TLS для werf** потребуется либо указывать флаг `--insecure-registry` для каждого вызова команд werf, взаимодействующих с Container Registry (`werf converge/cleanup/...`), либо глобально указать переменную окружения `WERF_INSECURE_REGISTRY=1` в сессии вашего терминала.

{% offtopic title="Windows" %}
TODO
{% endofftopic %}
{% offtopic title="macOS/Linux" %}
Добавим переменную окружения в `~/.bashrc`:

```shell
echo export WERF_INSECURE_REGISTRY=1 | tee -a ~/.bashrc
```
{% endofftopic %}

Запускаем minikube и создаём в нём новый Kubernetes-кластер:
```shell
minikube delete  # Удалим существующий minikube-кластер (если он есть).
minikube start --driver=docker --insecure-registry registry.example.com:80
```

[Устанавливаем утилиту kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/), после чего проверяем работоспособность нового кластера Kubernetes:
```shell
kubectl get --all-namespaces pod  # Должно показать список всех запущенных в кластере Pod'ов.
```

Если все Pod'ы из полученного списка находятся в состояниях `Running` или `Completed` (4-й столбец), а в 3-м столбце в выражениях вроде `1\1` цифра слева от `\` равна цифре справа (т.е. контейнеры Pod'а успешно запустились) — значит кластер Kubernetes запущен и работает. Если не все Pod'ы успешно запустились, то подождите и снова выполните команду выше для получения статуса всех Pod'ов.

### Установка NGINX Ingress Controller

Устанавливаем NGINX Ingress Controller:
```shell
minikube addons enable ingress
```

Немного подождём, после чего убедимся, что Ingress Controller успешно запустился:
```shell
kubectl -n ingress-nginx get pod
```

### Установка Container Registry для хранения образов

Установим и запустим Container Registry:
```shell
minikube addons enable registry
```

Создадим Ingress для доступа к Container Registry:
```yaml
kubectl apply -f - << EOF
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: registry
  namespace: kube-system
  annotations:
    nginx.ingress.kubernetes.io/proxy-body-size: "0"
spec:
  rules:
  - host: registry.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: registry
            port:
              number: 80
EOF
```

### Обновление файла hosts

Мы будем использовать домен `example.com` для доступа к приложению и домен `registry.example.com` для доступа к Container Registry.

Обновим файл hosts:
{% offtopic title="Windows" %}
Сначала получите IP-адрес minikube:
```shell
minikube ip
```

Используя полученный выше IP-адрес minikube, добавьте в конец файла `C:\Windows\System32\drivers\etc\hosts` следующую строку:
```
<IP-адрес minikube>    example.com registry.example.com
```
Должно получиться примерно так:
```
192.168.99.99          example.com registry.example.com
```
{% endofftopic %}
{% offtopic title="macOS/Linux" %}
Выполните команду в терминале:
```shell
echo "$(minikube ip) example.com registry.example.com" | sudo tee -a /etc/hosts
```
{% endofftopic %}

А теперь обновим файл hosts на мастер-ноде нашего кластера:
```shell
minikube ssh -- "echo $(minikube ip) example.com registry.example.com | sudo tee -a /etc/hosts"
```
