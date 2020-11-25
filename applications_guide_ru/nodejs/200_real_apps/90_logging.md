---
title: Логгирование
permalink: nodejs/200_real_apps/90_logging.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- app.js
{% endfilesused %}

Некоторые приложения сконфигурированы так, чтобы писать логи в файл. Это некорректно для приложений, запускаемых в kubernetes: docker-образ в такой ситуации будет бесконтрольно разрастаться, а при перевыкате приложения — логи исчезнут. 

Правильный путь — отправлять логи в [stdout и stderr](https://habr.com/ru/post/55136/). Kubernetes собирает такие логи и есть различные решения, позволяющие централизованно с ними работать.

Для Node.JS способ правильно отправлять логи — это функции `console.log()` и `console.error()`, к примеру:

```js
// Connection to SQLite
global_errors = [];
let sqlite_file = process.env.SQLITE_FILE;
if (! sqlite_file) {
  console.log('Environment variable SQLITE_FILE is not set! I will use in-memory database.');
  sqlite_file = ':memory:';
}
```

За отправляемыми логами можно следить через `kubectl`, например:

```bash
kubectl -n werf-guided-project logs <имя пода> -f
```

<div id="go-forth-button">
    <go-forth url="201_build.html" label="Сборка образа" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
