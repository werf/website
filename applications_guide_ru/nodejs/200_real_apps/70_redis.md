---
title: Подключаем Managed Redis
permalink: nodejs/200_real_apps/70_redis.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- package.json
- app.js
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении работу с простейшей базой данных типа in-memory — [Redis](https://redis.io/) (другим популярным примером является [memcached](https://memcached.org/)). Это означает, что база данных будет stateless.

Мы предполагаем, что база данных уже где-то реализована: вы либо подняли её на отдельном сервере, либо воспользовались Managed сервисом у своего cloud provider-а. Вопросы, связанные с самостоятельной установкой БД в кластер мы разберём в главе "Работа с инфраструктурой" 

## Сконфигурировать Redis в Kubernetes

TODO: создать Endpoint и объяснить, нахуя оно так.



{% snippetcut name=".helm/templates/service-mysql.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/270_managed_redis/.helm/templates/service-mysql.yaml" %}
{% raw %}
```yaml

.....

```
{% endraw %}
{% endsnippetcut %}



## Подключение Node.js-приложения к базе Redis

В нашем приложении мы будем подключаться к узлу Redis.

Мы не будем приводить полный код приложения — его можно [посмотреть в github](https://github.com/werf/werf-guides/blob/master/examples/nodejs/270_managed_redis/app.js). Приложение реализует то же API, что и описанное в главе "Быстрый старт разработчика", но вместо SQLite использует Redis.

Мы установим npm-пакет `redis` и реализуем подключение, например, так:

{% snippetcut name="app.js" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/270_managed_redis/app.js" %}
{% raw %}
```js
const redis = require("redis");

// Connection to Redis
const host = process.env.REDIS_HOST;
const port = process.env.REDIS_PORT;

const client = redis.createClient(port, host);
client.on("error", function(error) {
  console.error(error);
  process.exit(1)
});
```
{% endraw %}
{% endsnippetcut %}

Для подключения к базе данных нам, очевидно, нужно знать хост и порт. В коде приложения мы используем несколько переменных окружения: `REDIS_HOST`, `REDIS_PORT`.

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/270_managed_redis/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml

    ......

```
{% endraw %}
{% endsnippetcut %}

Тыры пыры конечно надо утащить это в values.yaml как мы писали тыры пыры

<div id="go-forth-button">
    <go-forth url="201_build.html" label="Сборка образа" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
