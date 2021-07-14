### Установка и запуск minikube

Установим/обновим minikube, [следуя инструкциям проекта](https://minikube.sigs.k8s.io/docs/start/). Убедимся, что у нас используется самая свежая версия minikube, доступная по приведённой ссылке.

Создаём новый Kubernetes-кластер с minikube:
```bash
minikube delete  # удалим существующий minikube-кластер (если он есть)
minikube start --driver=docker
```

Выставим Namespace по умолчанию, чтобы не указывать его при каждом вызове `kubectl`:
```powershell
kubectl config set-context minikube --namespace=werf-guide-app
```

Если утилита kubectl всё ещё не установлена, то установим её, [следуя инструкциям](https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/).

Теперь проверим работоспособность нового кластера Kubernetes:
```bash
kubectl get --all-namespaces pod  # должно показать список всех Pod'ов, запущенных в кластере
```

Если все Pod'ы из полученного списка находятся в состояниях `Running` или `Completed` (4-й столбец), а в 3-м столбце (в выражениях вроде `1/1`) для `Running` цифра слева от `/` равна цифре справа (т.е. контейнеры Pod'а успешно запустились) — кластер Kubernetes запущен и работает. Если не все Pod'ы успешно запустились, то подождите и снова выполните команду выше для получения статуса всех Pod'ов.

### Установка NGINX Ingress Controller

Устанавливаем NGINX Ingress Controller:
```bash
minikube addons enable ingress
```

Немного подождём, после чего убедимся, что Ingress Controller успешно запустился:
```bash
kubectl -n ingress-nginx get pod
```

Сделаем NGINX Ingress Controller доступным на 80-м порту после запуска `minikube tunnel`:
```bash
kubectl expose service -n ingress-nginx ingress-nginx-controller --name ingress-nginx-controller-lb --type LoadBalancer --port 80 --target-port http
```

Важно проверить, что в нашей системе не занят 80-й порт. Следующая команда должна выдать пустой результат:
```bash
netstat -anb | grep :80
```

В случае если порт занят, то результат будет подобным:

```bash
TCP     0.0.0.0:80        0.0.0.0:0       LISTENING
TCP     [::]:80           [::]:0          LISTENING
```

— в этом случае необходимо найти и остановить запущенный сервер.

Теперь нам надо не забывать держать запущенным `minikube tunnel` в отдельном окне PowerShell. Это необходимо для доступа с нашего хоста к ресурсам в кластере через Ingress:
```bash
minikube tunnel --cleanup=true
```

### Обновление файла hosts

Для доступа к приложению мы будем использовать домен `werf-guide-app`. Для этого обновим файл hosts (в PowerShell от администратора):
```powershell
Add-Content "C:\Windows\System32\drivers\etc\hosts" "`n127.0.0.1 werf-guide-app"
```
