---
title: Сборка образа
permalink: java_springboot/100_basic/10_build.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- Dockerfile
- werf.yaml
{% endfilesused %}

Классическим способом сборки является использование `Dockerfile`. Возможно, в ваших приложениях уже реализована сборка с помощью этого механизма, поэтому мы начнём с него и подключим его в werf. В следующих главах мы ускорим сборку воспользовавшись дополнительными механизмами, но сейчас — сфокусируемся на быстром получении результата.

## Сборка и Dockerfile

Для работы нашего приложения нужно

- взять базовый образ с nodejs (`gradle:jdk8-openj9` подойдёт)
- пробросить туда код приложения
- запустить сборку gradle и переместить собранный jar в удобное место
- указать конфигурацию приложения с помощью переменных окружения

Реализуем это в виде `Dockerfile`:

{% snippetcut name="Dockerfile" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/010_build/Dockerfile" %}
{% raw %}
```Dockerfile
FROM gradle:jdk8-openj9
WORKDIR /app
ENV SQLITE_FILE="app.db"
COPY . .

RUN gradle build --no-daemon
RUN cp /app/build/libs/*.jar /app/demo.jar

EXPOSE 8080
CMD ['java','-jar','/app/demo.jar']
```
{% endraw %}
{% endsnippetcut %}

## Подключаем Dockerfile к werf

Подключим уже готовый Dockerfile к werf и осуществим сборку. Для этого сделаем в корне репозитория файл `werf.yaml`, описывающий сборку всего проекта.

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/011_build_werf/werf.yaml" %}
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

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/011_build_werf/werf.yaml" %}
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

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/011_build_werf/werf.yaml" %}
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

Для того, чтобы начать работу, нужно инициализировать git в папке с исходным кодом приложения:

{% raw %}
```bash
git init
```
{% endraw %}

werf реализует подход "What you git is what you get" и опирается на коммиты, игнорируя не закомиченные изменения. Поэтому перед каждым запуском сборки и/или деплоя не забывайте сделать коммит, что-то вроде:

{% raw %}
```bash
git add .
git commit -m "work in progress"
```
{% endraw %}

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
│ │ basicapp/dockerfile  Sending build context to Docker daemon  116.7kB
│ │ basicapp/dockerfile  Step 1/15 : FROM gradle:jdk8-openj9
│ │ basicapp/dockerfile   ---> 2fb781988fa5
│ │ basicapp/dockerfile  Step 2/15 : WORKDIR /app
│ │ basicapp/dockerfile   ---> Using cache
<..>
│ │ basicapp/dockerfile  Successfully built e0d6df14df8b
│ │ basicapp/dockerfile  Successfully tagged ee51ea7f-c498-45a5-a435-0fd830fbb576:latest
│ ├ Info
│ │       name: werf-guided-project:50558f3f54d2ebbbd817824c6d7194aabe725bff6d7beae4df9c5e29-1606128099580
│ │       size: 738.6 MiB
│ └ Building stage basicapp/dockerfile (86.12 seconds)
└ ⛵ image basicapp (86.32 seconds)

Running time 86.37 seconds
```

Запустим собранный образ с помощью [werf run]({{ site.docsurl }}/documentation/cli/main/run.html):

```bash
$ werf run --docker-options="-d -p 8080:8080 --restart=always" -- java -jar /app/demo.jar
```

Обратите внимание, что мы задаем [параметры docker](https://docs.docker.com/engine/reference/run/) и после двойного дефиса команду, с которой запустить образ.

_В подглаве "Организация локальной разработки" мы рассмотрим более корректные способы организовать локальную разработку, в том числе — автоматическую локальную пересборку и перезапуск контейнеров при коммите._

Теперь приложение доступно локально на порту 8080:

![](/applications_guide_ru/images/applications-guide/020-hello-world-in-browser.png)

<div id="go-forth-button">
    <go-forth url="210_cluster.html" label="Сборка" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
