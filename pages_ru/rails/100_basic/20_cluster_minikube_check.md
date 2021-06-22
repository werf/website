### Проверка

Для проверки работоспособности необходимо открыть страницу в браузере или использовать `сurl` в консоли:

```shell
curl example.com
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