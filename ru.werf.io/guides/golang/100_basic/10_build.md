---
title: Сборка образа
permalink: golang/100_basic/10_build.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- Dockerfile
- werf.yaml
{% endfilesused %}

В этой главе мы начнём работу с werf: контейнизируем тестовое приложение, используя werf и [Dockerfile](https://docs.docker.com/engine/reference/builder/), и запустим конечный образ локально в Docker.

## Подготовка рабочего места

Мы предполагаем, что вы уже [установили werf]({{ site.docsurl }}/installation.html) и Docker.

Создайте директорию на своём компьютере и выполните там следующие шаги:

```shell
git clone https://github.com/werf/werf-guides.git
cp -r werf-guides/examples/golang/000_app ./
cd 000_app 
git init
git add .
git commit -m "initial commit"
```

_Так вы скопируете себе код [приложения на Go](https://github.com/werf/werf-guides/tree/master/examples/golang/000_app) и инициируете Git в каталоге с ним._

## Реализация сборки в Dockerfile

Конфигурация сборки нашего приложения будет состоять из следующих шагов:

- взять официальный образ Go (`golang:1.15-alpine` подойдёт);
- добавить в него код приложения;
- установить модули-зависимости;
- скомпилировать бинарный файл приложения;
- взять базовый образ alpine (например, `alpine:3.13`);
- скопировать в него бинарный файл с приложением из сборочного образа.

Реализуем это в `Dockerfile`:

{% snippetcut name="Dockerfile" url="https://github.com/werf/werf-guides/blob/master/examples/golang/010_build/Dockerfile" %}
{% raw %}
```Dockerfile
FROM golang:1.15-alpine as builder
RUN apk add --update gcc musl-dev
WORKDIR /app
COPY go.* /app/
RUN go mod download
COPY /cmd /app/cmd
RUN go build ./cmd/demoapp

FROM alpine:3.13
COPY --from=builder /app/demoapp /app/demoapp
COPY /app.db /app/app.db

CMD ["/app/demoapp"]
```
{% endraw %}
{% endsnippetcut %}

## Интеграция Dockerfile с werf

Подключим готовый Dockerfile к werf. Для этого, создадим в корне репозитория файл `werf.yaml`, описывающий сборку всего проекта:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/golang/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-golang
configVersion: 1
---
image: basicapp
dockerfile: Dockerfile
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="Что тут написано?" %}

`werf.yaml` начинается с обязательной **секции мета-информации**:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/golang/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-golang
configVersion: 1
```
{% endraw %}
{% endsnippetcut %}

Здесь:
- **_project_** — поле, задающее уникальное имя проекта приложения. Имя проекта по умолчанию используется при генерации имени Helm-релиза и пространства имен, `namespace`, в Kubernetes. Изменение имени у активного проекта затруднительно и требует ряда действий, которые необходимо выполнить вручную (подробнее о возможных последствиях можно прочитать [здесь]({{ site.docsurl }}/documentation/reference/werf_yaml.html#последствия-смены-имени-проекта));
- **_configVersion_** определяет версию синтаксиса, используемую в `werf.yaml` (на данный момент поддерживается только версия `1`).

Следующая секция конфигурации, которая и будет основной для сборки: [**image config section**]({{ site.docsurl }}/documentation/reference/werf_yaml.html#секция-image).

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/golang/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
---
image: basicapp
dockerfile: Dockerfile
```
{% endraw %}
{% endsnippetcut %}

В строке `image: basicapp` определяется идентификатор образа, который впоследствии может использоваться при конфигурации выката, а также для вызова команд werf для определённого образа из `werf.yaml` (к примеру, `werf build basicapp`, `werf run basicapp` и т.д.). 

Строка `dockerfile: Dockerfile` указывает, что сборочная конфигурация описана в существующем файле, расположенном по пути `Dockerfile`. 

Также доступны и другие директивы, с которыми можно ознакомиться [по ссылке]({{ site.docsurl }}/documentation/reference/werf_yaml.html#dockerfile-image-section-image).

В одном `werf.yaml` может быть описано произвольное количество образов.

{% endofftopic %}

## Сборка

Перед выполнением сборки необходимо добавить изменения в git-репозиторий проекта:

```shell
git add werf.yaml Dockerfile
git commit -m "Add build configuration"
```

> Почему изменения должны добавляться в git-репозиторий, что такое гитерминизм и режим разработчика, а также другие особенности работы с файлами проекта будут разобраны далее в главе «Необходимо знать»

Сборка выполняется командой [`werf build`]({{ site.docsurl }}/documentation/reference/cli/werf_build.html):

{% raw %}
```shell
┌ ⛵ image basicapp
│ ┌ Building stage basicapp/dockerfile
│ │ basicapp/dockerfile  Sending build context to Docker daemon  69.63kB
│ │ basicapp/dockerfile  Step 1/16 : FROM golang:1.15-alpine as builder
│ │ basicapp/dockerfile   ---> 1de1afaeaa9a
│ │ basicapp/dockerfile  Step 2/16 : RUN apk add --update gcc musl-dev
│ │ basicapp/dockerfile   ---> Using cache
<..>
│ │ basicapp/dockerfile  Successfully built f2a49a68a666
│ │ basicapp/dockerfile  Successfully tagged c89c09f6-411f-4d4c-9c9e-43f4bdfad074:latest
│ ├ Info
│ │       name: werf-guided-golang:f276ddd4d73aafb69d657234505e718f78284bbdd816863f1540a912-1611832304608
│ │       size: 17.8 MiB
│ └ Building stage basicapp/dockerfile (55.42 seconds)
└ ⛵ image basicapp (55.75 seconds)


Running time 55.82 seconds
```
{% endraw %}

## Запуск

Запуск контейнера выполняется командой [werf run]({{ site.docsurl }}/documentation/cli/main/run.html):

```shell
werf run --docker-options="--rm -p 3000:3000" -- /app/demoapp
```

Обратите внимание, что [параметры docker](https://docs.docker.com/engine/reference/run/) задаются опцией `--docker-options`, а команда запуска после двух дефисов.

В `werf.yaml` может описываться произвольное количество образов. Для запуска контейнера определённого `image` из werf.yaml необходимо использовать позиционный аргумент команды (`werf run basicapp ...`).

_Можно заметить, что вызов `werf run` осуществляет сборку, т.е. предварительная сборка не требуется._

Теперь приложение доступно локально на порту 3000:

![](/images/golang/100_10_app_in_browser.png)

<div id="go-forth-button">
    <go-forth url="20_cluster.html" label="Подготовка кластера" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
