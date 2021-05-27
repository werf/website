### Установка Docker Desktop и запуск Kubernetes

Установим Docker Desktop на [Windows](https://docs.docker.com/docker-for-windows/install/) или [macOS](https://docs.docker.com/docker-for-mac/install/).

Разрешаем доступ в Container Registry без TLS:
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

[Включаем Kubernetes в Docker Desktop](https://docs.docker.com/desktop/kubernetes/). Выделяем достаточно ресурсов для запуска Kubernetes и приложений в нём, следуя инструкциям для [Windows](https://docs.docker.com/docker-for-windows/#resources) или [macOS](https://docs.docker.com/docker-for-mac/#resources). Если ресурсов будет недостаточно, то часть приложений или компонентов кластера запуститься не смогут. Дальнейшие инструкции были проверены при выделенных для кластера 6 CPU, 6 ГБ RAM и выделенном диске в 24 ГБ.
{% comment %} TODO(lesikov): чет слишком много ресурсов хотим, надо протестить на винде, сколько реально надо. Мне кажется для кластера + приложения достаточно будет 1 CPU и 1 ГБ RAM. {% endcomment %}

Утилита `kubectl` должна взаимодействовать именно с новым кластером Kubernetes:
```shell
kubectl config use-context docker-desktop
```

Проверим работоспособность Kubernetes:
```shell
kubectl --all-namespaces get pod # Должно показать список всех запущенных в кластере Pod'ов.
```

Если все Pod'ы из полученного списка находятся в состояниях `Running` или `Completed` (4-й столбец), а в 3-м столбце в выражениях вроде `1\1` цифра слева от `\` равна цифре справа (т.е. контейнеры Pod'а успешно запустились) — значит кластер Kubernetes запущен и работает. Если не все Pod'ы успешно запустились, то подождите и снова выполните команду выше для получения статуса Pod'ов.

### Установка NGINX Ingress Controller

{% comment %} TODO(lesikov): надо инструкции, как проверить, свободен ли порт для мака и винды. {% endcomment %}
80-й порт вашей машины должен быть свободен перед установкой. Устанавливаем NGINX Ingress Controller, [следуя инструкциям](https://kubernetes.github.io/ingress-nginx/deploy/#docker-desktop) (установка может занять до минуты).

Убедимся, что Ingress Controller успешно запустился:
```shell
kubectl -n ingress-nginx get pod
```

### Установка Container Registry для хранения образов

Запустим Container Registry в кластере:
```yaml
kubectl apply -f - << EOF
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: registry
  namespace: kube-system
spec:
  selector:
    matchLabels:
      app: registry
  template:
    metadata:
      labels:
        app: registry
    spec:
      containers:
      - name: registry
        image: "registry:2"
        ports:
        - containerPort: 5000
        volumeMounts:
        - name: registry-data
          mountPath: /var/lib/registry
      volumes:
      - name: registry-data
        emptyDir: {}

---
apiVersion: v1
kind: Service
metadata:
  name: registry
  namespace: kube-system
spec:
  selector:
    app: registry
  ports:
    - name: registry
      port: 80
      targetPort: 5000

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

Мы будем использовать домен `example.com` для доступа к приложению и `registry.example.com` для доступа к Container Registry. Обновим файл hosts:
{% offtopic title="Windows" %}
Добавьте в конец файла `C:\Windows\System32\drivers\etc\hosts` следующие строки:
```
127.0.0.1           example.com
127.0.0.1           registry.example.com
```
{% endofftopic %}
{% offtopic title="macOS" %}
Выполните команду в терминале:
```shell
echo "127.0.0.1 example.com registry.example.com" | sudo tee -a /etc/hosts
```
{% endofftopic %}
