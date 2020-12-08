---
title: Логгирование
permalink: java_springboot/200_real_apps/25_logging.html
---

Некоторые приложения сконфигурированы так, чтобы писать логи в файл. Это некорректно для приложений, запускаемых в kubernetes: docker-образ в такой ситуации будет бесконтрольно разрастаться, а при перевыкате приложения — логи исчезнут. 

Правильный путь — отправлять логи в [stdout и stderr](https://habr.com/ru/post/55136/). Kubernetes собирает такие логи и есть различные решения, позволяющие централизованно с ними работать.

Когда вы запускаете jar-файл — java пишет логи об этом в нужные потоки по умолчанию. Однако, если вы хотите отправлять дополнительную информацию, то необходимо [сконфигурировать и использовать логгер](https://www.baeldung.com/spring-boot-logging) в своём коде.

За отправляемыми логами можно следить через `kubectl`, например:

```bash
kubectl -n werf-guided-project logs <имя пода> -f
```

<div id="go-forth-button">
    <go-forth url="30_assets.html" label="Генерируем и раздаём ассеты" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
