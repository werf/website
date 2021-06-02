### Установка и запуск minikube

Установим minikube, [следуя инструкциям здесь](https://minikube.sigs.k8s.io/docs/start/).

Порты 80, 3000 должны быть свободны.

Чтобы разрешить Docker'у доступ в Registry по HTTP, в настройках Docker Desktop в секции Docker Engine добавим в ключ `insecure-registries` адрес нашего Registry:
```json
{
  "insecure-registries": ["registry.example.com:80"]
}
```
После чего перезапустим Docker Desktop через его меню.

А чтобы сам werf работал с Registry по HTTP для werf есть опция `--insecure-registry`. Чтобы не указывать её каждый раз рекомендуется выставить переменную окружения `WERF_INSECURE_REGISTRY=1` таким образом (PowerShell):
```powershell
[Environment]::SetEnvironmentVariable("WERF_INSECURE_REGISTRY", "1", "User")
```

Создаём новый Kubernetes-кластер с minikube:
```shell
minikube delete  # Удалим существующий minikube-кластер (если он есть).
minikube start --driver=docker --insecure-registry registry.example.com:80 --namespace werf-guided-rails
```

Если утилита kubectl всё ещё не установлена, то установим её, [следуя инструкциям](https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/).

Теперь проверим работоспособность нового кластера Kubernetes:
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

Установим и запустим Registry:
```shell
minikube addons enable registry
```

Создадим Ingress для доступа к Registry (PowerShell):
```powershell
kubectl -n kube-system expose service registry --type LoadBalancer --name registry-lb --port 80 --target-port 5000
@'
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
            name: registry-lb
            port:
              number: 80
'@ | kubectl apply -f -
```

### Обновление файла hosts

Мы будем использовать домен `example.com` для доступа к приложению и домен `registry.example.com` для доступа к Registry.

Обновим файлы hosts (PowerShell от администратора):
```powershell
Add-Content "C:\Windows\System32\drivers\etc\hosts" "`n127.0.0.1 example.com registry.example.com"
```

### Настройка Docker-окружения

При каждом открытии нового PowerShell-терминала нам потребуется настраивать окружение для работы с Docker, запущенным в minikube. Это делается командой
```
& minikube -p minikube docker-env | Invoke-Expression
```

Чтобы не вводить эту команду каждый раз, выполните в PowerShell от администратора команды (учтите, что это разрешит выполнение неподписанных скриптов, что может сказаться на безопасности):
```powershell
Set-ExecutionPolicy unrestricted
if (-Not(Test-Path -Path "$profile")) {
  New-Item "$profile" -type file
}
Add-Content -Path "$profile" -Value "if (minikube status) { & minikube -p minikube docker-env | Invoke-Expression }"
& minikube -p minikube docker-env | Invoke-Expression
```

### Проброс портов

Также нам всегда надо держать запущенным в отдельном окне `minikube tunnel` для доступа к ресурсам в кластере:
```powershell
minikube tunnel --cleanup=true
```
