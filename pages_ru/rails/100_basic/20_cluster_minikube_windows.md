### Установка и запуск minikube

Установим/обновим minikube, [следуя инструкциям проекта](https://minikube.sigs.k8s.io/docs/start/). Убедимся, что у нас используется самая свежая версия minikube, доступная по приведённой ссылке.

Порт 80 должен быть свободен.

Создаём новый Kubernetes-кластер с minikube:
```bash
minikube delete  # удалим существующий minikube-кластер (если он есть)
minikube start --driver=docker --namespace werf-guided-rails
```

Если утилита kubectl всё ещё не установлена, то установим её, [следуя инструкциям](https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/).

Теперь проверим работоспособность нового кластера Kubernetes:
```bash
kubectl get --all-namespaces pod  # должно показать список всех Pod'ов, запущенных в кластере
```

Если все Pod'ы из полученного списка находятся в состояниях `Running` или `Completed` (4-й столбец), а в 3-м столбце (в выражениях вроде `1/1`) цифра слева от `/` равна цифре справа (т.е. контейнеры Pod'а успешно запустились) — кластер Kubernetes запущен и работает. Если не все Pod'ы успешно запустились, то подождите и снова выполните команду выше для получения статуса всех Pod'ов.

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

Теперь нам надо не забывать держать запущенным `minikube tunnel` в отдельном окне PowerShell. Это необходимо для доступа к ресурсам в кластере через Ingress:
```bash
minikube tunnel --cleanup=true
```

### Обновление файла hosts

Для доступа к приложению мы будем использовать домен `example.com`. Для этого обновим файл hosts (в PowerShell от администратора):
```powershell
Add-Content "C:\Windows\System32\drivers\etc\hosts" "`n127.0.0.1 example.com kubernetes-basics-app.example.com"
```

### Проверка

Для проверки работоспособности необходимо открыть страницу в браузере или выполнить запрос через PowerShell:

```shell
Invoke-WebRequest -Uri http://example.com
```

Если домен `example.com` не используется приложениями в Kubernetes, то NGINX должен вернуть страницу 404:

```html
<html>
<head><title>404 Not Found</title></head>
<body>
<center><h1>404 Not Found</h1></center>
<hr><center>nginx</center>
</body>
</html>
```
