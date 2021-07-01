### Установка и запуск minikube

Установим/обновим minikube, [следуя инструкциям проекта](https://minikube.sigs.k8s.io/docs/start/). Убедимся, что у нас используется самая свежая версия minikube, доступная по приведённой ссылке.

Создаём новый Kubernetes-кластер с minikube:
```shell
minikube delete  # удалим существующий minikube-кластер (если он есть)
minikube start --driver=docker
```

Выставим Namespace по умолчанию, чтобы не указывать его при каждом вызове `kubectl`:
```shell
kubectl config set-context minikube --namespace=werf-guide-app
```

Если утилита kubectl всё ещё не установлена, то установим её, [следуя инструкциям](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/).

Теперь проверим работоспособность нового кластера Kubernetes:
```shell
kubectl get --all-namespaces pod  # должно показать список всех Pod'ов, запущенных в кластере
```

Если все Pod'ы из полученного списка находятся в состояниях `Running` или `Completed` (4-й столбец), а в 3-м столбце (в выражениях вроде `1/1`) цифра слева от `/` равна цифре справа (т.е. контейнеры Pod'а успешно запустились) — кластер Kubernetes запущен и работает. Если не все Pod'ы успешно запустились, то подождите и снова выполните команду выше для получения статуса всех Pod'ов.

### Установка NGINX Ingress Controller

Устанавливаем NGINX Ingress Controller:
```shell
minikube addons enable ingress
```

Немного подождём, после чего убедимся, что Ingress Controller успешно запустился:
```shell
kubectl -n ingress-nginx get pod
```

### Обновление файла hosts

Для доступа к приложению мы будем использовать домен `werf-guide-app`. Для этого обновим файл hosts:
```shell
echo "$(minikube ip) werf-guide-app werf-kube-basics" | sudo tee -a /etc/hosts
```
