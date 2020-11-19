---
title: Сборка образа
permalink: nodejs/100_basic/10_build.html
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

Реализуем это в виде `Dockerfile`:

{% snippetcut name="Dockerfile" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/010_build/Dockerfile" %}
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

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/011_build_werf/werf.yaml" %}
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

Начинается `werf.yaml` с обязательной [**секции мета-информации**]({{ site.docsurl }}/documentation/configuration/introduction.html#секция-мета-информации):

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-project
configVersion: 1
```
{% endraw %}
{% endsnippetcut %}

Здесь:
- **_project_** — поле, задающее имя для проекта. Под проектом понимается набор из Docker-образов, которые мы будем собирать. Имя проекта по умолчанию используется в названии Helm-релиза и пространства имен в Kubernetes, куда будет выкатываться приложение. Данное имя не рекомендуется изменять (или подходить к таким изменениям с должным уровнем ответственности), так как после изменений уже имеющиеся ресурсы, которые выкачены в кластер, не будут автоматически переименованы.
- **_configVersion_** — в данном случае определяет версию синтаксиса, используемую в `werf.yaml` (на данный момент мы всегда используем `1`).

Следующая секция конфигурации, которая и будет основной для сборки: [**image config section**]({{ site.docsurl }}/documentation/configuration/introduction.html#%D1%81%D0%B5%D0%BA%D1%86%D0%B8%D1%8F-%D0%BE%D0%B1%D1%80%D0%B0%D0%B7%D0%B0).

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
---
image: basicapp
dockerfile: Dockerfile
```
{% endraw %}
{% endsnippetcut %}

В строке `image: basicapp` дано название для образа, который соберёт werf. Это имя будет впоследствии указываться при запуске контейнера. Строка `dockerfile: Dockerfile` указывает, что сборка осуществляться на основании `Dockerfile`.

Один werf.yaml может собирать несколько образов с разными настройками.

{% endofftopic %}

## Осуществляем сборку

Для того, чтобы запустить сборку — воспользуемся [командой `build`]({{ site.docsurl }}/documentation/reference/cli/werf_build.html):

{% raw %}
```bash
werf build
```
{% endraw %}

_В подглаве "Ускорение сборки" мы переведём сборку с Dockerfile на механизм werf под названием Stapel и получим расширенные возможности: оптимизацию, распределённую сборку, инструменты диагностики, content based теггирование и другие._

Но уже сейчас вы можете заметить, что werf делает расширенный вывод логов сборки, вроде:

```
┌ ⛵ image basicapp
│ ┌ Building stage basicapp/dockerfile
│ │ basicapp/dockerfile  Sending build context to Docker daemon  69.63kB
│ │ basicapp/dockerfile  Step 1/16 : FROM node:14-stretch
│ │ basicapp/dockerfile   ---> b90fa0d7cbd1
│ │ basicapp/dockerfile  Step 2/16 : WORKDIR /app
│ │ basicapp/dockerfile   ---> Using cache
│ │ basicapp/dockerfile   ---> 4bb99952fe98
<..>
│ │ basicapp/dockerfile  Successfully built 02a0a425890a
│ │ basicapp/dockerfile  Successfully tagged a1cbf6dc-343f-4a77-b846-d0f12a700cb7:latest
│ ├ Info
│ │       name: werf-guided-project:a473b87e1ad65f102fa90f8c6647b03056e5ae95ff1ef3c5e7fd2c31-1605597979927
│ │       size: 953.1 MiB
│ └ Building stage basicapp/dockerfile (21.94 seconds)
└ ⛵ image basicapp (22.04 seconds)

Running time 22.07 seconds
```

Запустим собранный образ с помощью [werf run]({{ site.docsurl }}/documentation/cli/main/run.html):

```bash
$ werf run --docker-options="-d -p 3000:3000 --restart=always" -- node /app/app.js
```

Обратите внимание, что мы задаем [параметры docker](https://docs.docker.com/engine/reference/run/) и после двойного дефиса команду, с которой запустить образ.

_В подглаве "Организация локальной разработки" мы рассмотрим более корректные способы организовать локальную разработку, в том числе — автоматическую локальную пересборку и перезапуск контейнеров при коммите._

Теперь приложение доступно локально на порту 3000:

![](/applications_guide_ru/images/applications-guide/020-hello-world-in-browser.png)

<div id="go-forth-button">
    <go-forth url="210_cluster.html" label="Сборка" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
