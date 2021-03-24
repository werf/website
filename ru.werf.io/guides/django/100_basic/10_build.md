---
title: Сборка образа
permalink: django/100_basic/10_build.html
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
cp -r werf-guides/examples/django/000_app ./
cd 000_app 
git init
git add .
git commit -m "initial commit"
```

_Так вы скопируете себе код [приложения на Django](https://github.com/werf/werf-guides/tree/master/examples/django/000_app) и инициируете Git в каталоге с ним._

## Реализация сборки в Dockerfile

Конфигурация сборки нашего приложения будет состоять из следующих шагов:

- взять базовый образ с Python3 (`python:3.9.1-alpine3.12` подойдёт);
- добавить в него код приложения;
- установить зависимости npm.

Реализуем это в `Dockerfile`:

{% snippetcut name="Dockerfile" url="https://github.com/werf/werf-guides/blob/master/examples/django/010_build/Dockerfile" %}
{% raw %}
```Dockerfile
FROM python:3.9.1-alpine3.12
WORKDIR /app
COPY project /app

RUN pip3 install -r requirements.txt

ENTRYPOINT ["python3", "manage.py"]
CMD ["runserver"]
```
{% endraw %}
{% endsnippetcut %}

## Интеграция Dockerfile с werf

Подключим готовый Dockerfile к werf. Для этого, создадим в корне репозитория файл `werf.yaml`, описывающий сборку всего проекта:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/django/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-django
configVersion: 1
---
image: basicapp
dockerfile: Dockerfile
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="Что тут написано?" %}

`werf.yaml` начинается с обязательной **секции мета-информации**:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/django/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-django
configVersion: 1
```
{% endraw %}
{% endsnippetcut %}

Здесь:
- **_project_** — поле, задающее уникальное имя проекта приложения. Имя проекта по умолчанию используется при генерации имени Helm-релиза и пространства имен, `namespace`, в Kubernetes. Изменение имени у активного проекта затруднительно и требует ряда действий, которые необходимо выполнить вручную (подробнее о возможных последствиях можно прочитать [здесь]({{ site.docsurl }}/documentation/reference/werf_yaml.html#последствия-смены-имени-проекта));
- **_configVersion_** определяет версию синтаксиса, используемую в `werf.yaml` (на данный момент поддерживается только версия `1`).

Следующая секция конфигурации, которая и будет основной для сборки: [**image config section**]({{ site.docsurl }}/documentation/reference/werf_yaml.html#секция-image).

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/django/011_build_werf/werf.yaml" %}
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
│ │ basicapp/dockerfile  Step 1/16 : FROM python:3.9.1-alpine3.12
│ │ basicapp/dockerfile   ---> b90fa0d7cbd1
│ │ basicapp/dockerfile  Step 2/16 : WORKDIR /app
│ │ basicapp/dockerfile   ---> Using cache
<..>
│ │ basicapp/dockerfile  Successfully built 02a0a425890a
│ │ basicapp/dockerfile  Successfully tagged a1cbf6dc-343f-4a77-b846-d0f12a700cb7:latest
│ ├ Info
│ │       name: werf-guided-django:a473b87e1ad65f102fa90f8c6647b03056e5ae95ff1ef3c5e7fd2c31-1605597979927
│ │       size: 953.1 MiB
│ └ Building stage basicapp/dockerfile (21.94 seconds)
└ ⛵ image basicapp (22.04 seconds)

Running time 22.07 seconds
```
{% endraw %}

## Запуск

Запуск контейнера выполняется командой [werf run]({{ site.docsurl }}/documentation/cli/main/run.html):

```shell
werf run --docker-options="-d -p 8000:8000 --restart=always" -- runserver 0.0.0.0:8000
```

Обратите внимание, что [параметры docker](https://docs.docker.com/engine/reference/run/) задаются опцией `--docker-options`, а команда запуска после двух дефисов.

В `werf.yaml` может описываться произвольное количество образов. Для запуска контейнера определённого `image` из werf.yaml необходимо использовать позиционный аргумент команды (`werf run basicapp ...`).

_Можно заметить, что вызов `werf run` осуществляет сборку, т.е. предварительная сборка не требуется._

Теперь приложение доступно локально на порту 8000:

![](/guides/images/django/100_10_app_in_browser.png)

<div id="go-forth-button">
    <go-forth url="20_cluster.html" label="Подготовка кластера" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
