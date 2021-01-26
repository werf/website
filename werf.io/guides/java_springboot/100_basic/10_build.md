---
title: Сборка образа
permalink: java_springboot/100_basic/10_build.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- Dockerfile
- werf.yaml
{% endfilesused %}

Классическим способом сборки является использование [Dockerfile](https://docs.docker.com/engine/reference/builder/). Возможно, в ваших приложениях уже реализована сборка с помощью этого механизма, поэтому мы начнём с него, а затем подключим его в werf. В следующих главах мы ускорим сборку, воспользовавшись альтернативным синтаксисом описания сборки, а сейчас — сфокусируемся на быстром получении результата.

## Реализация сборки в Dockerfile

Конфигурация сборки нашего приложения будет состоять из следующих шагов:

- взять базовый образ с nodejs (`gradle:jdk8-openj9` подойдёт)
- добавить в него код приложения
- собрать приложение с помощью gradle и переместить полученный jar в подходящее место
- указать конфигурацию приложения с помощью переменных окружения

Реализуем это в `Dockerfile`:

{% snippetcut name="Dockerfile" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/010_build/Dockerfile" %}
{% raw %}
```Dockerfile
FROM gradle:jdk8-openj9
WORKDIR /app
COPY . .

RUN gradle build --no-daemon
RUN cp /app/build/libs/*.jar /app/demo.jar

CMD ['java','-jar','/app/demo.jar']
```
{% endraw %}
{% endsnippetcut %}

## Интеграция Dockerfile с werf

Подключим готовый Dockerfile к werf. Для этого, создадим в корне репозитория файл `werf.yaml`, описывающий сборку всего проекта.

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

Начинается `werf.yaml` с обязательной **секции мета-информации**:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-project
configVersion: 1
```
{% endraw %}
{% endsnippetcut %}

Здесь:
- **_project_** — поле, задающее уникальное имя проекта приложения. Имя проекта по умолчанию используется при генерации имени Helm-релиза и пространства имен, `namespace`, в Kubernetes. Изменение имени у активного проекта требует должного уровня ответственности и ряда действий, которые необходимо выполнить вручную (подробнее о возможных последствиях можно прочитать [здесь]({{ site.docsurl }}/documentation/reference/werf_yaml.html#последствия-смены-имени-проекта)).
- **_configVersion_** — определяет версию синтаксиса, используемую в `werf.yaml` (на данный момент поддерживается только версия `1`).

Следующая секция конфигурации, которая и будет основной для сборки: [**image config section**]({{ site.docsurl }}/documentation/reference/werf_yaml.html#секция-image).

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
---
image: basicapp
dockerfile: Dockerfile
```
{% endraw %}
{% endsnippetcut %}

В строке `image: basicapp` определяется внутренний идентификатор образа, который впоследствии может использоваться при конфигурации выката, а также для вызова команд werf для определённого образа из werf.yaml (к примеру, `werf build basicapp`, `werf run basicapp` и т.д.). 

Строка `dockerfile: Dockerfile` указывает, что сборочная конфигурация описаны в существующем `Dockerfile`. 

Так же доступны и другие директивы, с которыми можно ознакомиться [по ссылке]({{ site.docsurl }}/documentation/reference/werf_yaml.html#dockerfile-image-section-image).

В одном werf.yaml может быть описано произвольное количество образов.

{% endofftopic %}

## Сборка

Для того, чтобы начать работу, нужно инициализировать git в папке с исходным кодом приложения:

{% raw %}
```bash
git init
```
{% endraw %}

werf реализует подход "What you git is what you get" и опирается на коммиты, игнорируя незакомиченные изменения. Поэтому при локальной разработке не забывайте фиксировать свои изменения коммитом. К примеру так:

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

_В подглаве "Ускорение сборки" мы переведём сборку с Dockerfile на альтернативный синтаксис werf под названием Stapel и получим расширенные возможности: инкрементальную пересборку с учетом истории git, возможность использовать Ansible, использование кэша между сборками, удобные инструменты диагностики и многое другое._

Но уже сейчас вы можете заметить, что werf делает расширенный вывод логов сборки:

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

## Запуск

Запустим собранный образ с помощью [werf run]({{ site.docsurl }}/documentation/cli/main/run.html):

```bash
$ werf run --docker-options="-d -p 8080:8080 --restart=always" -- java -jar /app/demo.jar
```

Обратите внимание, что мы задаем [параметры docker](https://docs.docker.com/engine/reference/run/) опцией `--docker-options`, а команду, указываем после двух дефисов.

_В подглаве "Организация локальной разработки" мы рассмотрим более корректные способы организации локальной разработки, в том числе — автоматическую локальную пересборку и перезапуск контейнеров при коммите._

Теперь приложение доступно локально на порту 8080:

![](/guides/images/springboot/100_10_app_in_browser.png)

<div id="go-forth-button">
    <go-forth url="20_cluster.html" label="Подготовка кластера" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
