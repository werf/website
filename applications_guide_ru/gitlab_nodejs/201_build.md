---
title: Сборка образа
permalink: gitlab_nodejs/201_build.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- Dockerfile
- werf.yaml
{% endfilesused %}

Классическим способом сборки является использование `Dockerfile`. Возможно, в ваших приложениях уже реализована сборка с помощью этого механизма, поэтому мы начнём с него и подключим его в werf. В следующих главах мы ускорим сборку воспользовавшись дополнительными механизмами, но сейчас — сфокусируемся на быстром получении результата.

## Сборка и Dockerfile

Для работы нашего приложения нужно

- взять базовый образ с nodejs (`node:14-stretch` подойдёт)
- пробросить туда код приложения
- установить зависимости npm
- указать конфигурацию приложения с помощью переменных окружения

{% snippetcut name="Dockerfile" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab_nodejs/010_build/Dockerfile" %}
{% raw %}
```Dockerfile
FROM node:14-stretch
WORKDIR /app
ENV SQLITE_FILE="app.db"
COPY . .

RUN apt update
RUN apt install -y tzdata locales
RUN npm ci

EXPOSE 3000
CMD ['node','/app/app.js']
```
{% endraw %}
{% endsnippetcut %}

## Подключаем Dockerfile к werf

Подключим уже готовый Dockerfile к werf и осуществим сборку. Для этого сделаем в корне репозитория файл `werf.yaml`, описывающий сборку всего проекта.

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab_nodejs/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-project
configVersion: 1
---
image: basicapp
dockerfile: Dockerfile
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="Что тут написано?" %}

TODO: расписать, что а) имаджей может быть несколько, несколько докерфайлов б) подробнее о конфигурировании читай там-то там-то в доке.  https://ru.werf.io/documentation/reference/werf_yaml.html

{% endofftopic %}

TODO: запускаем сборку руками, верфь конвердж вот это всё.

TODO: делаем ремарку про отладку, оптимизацию сборки, распределённую сборку и тегирование — они будут позже, в следующих главах, когда мы займёмся оптимизацией. Кратко можно дать затравочку про идею слоёв.

TODO: но уже сейчас мы получаем прикольный вывод логов с доп. данными о процессе, которые помогают диагностировать проблемы + контент бейзд теггинг в базовом виде.

TODO: запускаем локально и смотрим глазами

<div id="go-forth-button">
    <go-forth url="210_cluster.html" label="Сборка" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
