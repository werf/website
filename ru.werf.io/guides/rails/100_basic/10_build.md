---
title: Сборка образа
permalink: rails/100_basic/10_build.html
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
cp -r werf-guides/examples/rails/000_app ./
cd 000_app
git init
git add .
git commit -m "initial commit"
```

_Так вы скопируете себе код [приложения на Ruby on Rails](https://github.com/werf/werf-guides/tree/master/examples/rails/000_app) и инициируете Git в каталоге с ним._

## Реализация сборки в Dockerfile

Конфигурация сборки нашего приложения будет состоять из следующих шагов:

- взять базовый образ с Ruby (`ruby:2.7.1` подойдёт);
- добавить в него код приложения;
- установить зависимости и gem'ы.

Реализуем это в `Dockerfile`:

{% snippetcut name="Dockerfile" url="https://github.com/werf/werf-guides/blob/master/examples/rails/010_build/Dockerfile" %}
{% raw %}
```Dockerfile
FROM ruby:2.7.1
WORKDIR /app
COPY . .

RUN apt-get update -qq && apt-get install -y \
 build-essential libpq-dev libxml2-dev libxslt1-dev curl
RUN bundle install

CMD ["bundle", "exec", "puma"]
```
{% endraw %}
{% endsnippetcut %}

## Интеграция Dockerfile с werf

Подключим готовый Dockerfile к werf. Для этого, создадим в корне репозитория файл `werf.yaml`, описывающий сборку всего проекта:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-rails
configVersion: 1
---
image: basicapp
dockerfile: Dockerfile
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="Что тут написано?" %}

`werf.yaml` начинается с обязательной **секции мета-информации**:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-rails
configVersion: 1
```
{% endraw %}
{% endsnippetcut %}

Здесь:
- **_project_** — поле, задающее уникальное имя проекта приложения. Имя проекта по умолчанию используется при генерации имени Helm-релиза и пространства имен, `namespace`, в Kubernetes. Изменение имени у активного проекта затруднительно и требует ряда действий, которые необходимо выполнить вручную (подробнее о возможных последствиях можно прочитать [здесь]({{ site.docsurl }}/documentation/reference/werf_yaml.html#последствия-смены-имени-проекта));
- **_configVersion_** определяет версию синтаксиса, используемую в `werf.yaml` (на данный момент поддерживается только версия `1`).

Следующая секция конфигурации, которая и будет основной для сборки: [**image config section**]({{ site.docsurl }}/documentation/reference/werf_yaml.html#секция-image).

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/011_build_werf/werf.yaml" %}
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
│ │ basicapp/dockerfile  Sending build context to Docker daemon  11.64MB
│ │ basicapp/dockerfile  Step 1/14 : FROM ruby:2.7.1
│ │ basicapp/dockerfile   ---> d8ca85855516
│ │ basicapp/dockerfile  Step 2/14 : WORKDIR /app
│ │ basicapp/dockerfile   ---> Using cache
<..>
│ │ basicapp/dockerfile  Successfully built 3a4ede4e9556
│ │ basicapp/dockerfile  Successfully tagged 0041b344-efe4-416d-baff-5e50fbb712b0:latest
│ ├ Info
│ │       name: werf-guided-rails:31e0e7436c3055fa816fc770ebda185bacb7e8ef53775b8e5488a83f-1611855308907
│ │       size: 929.2 MiB
│ └ Building stage basicapp/dockerfile (94.47 seconds)
└ ⛵ image basicapp (96.07 seconds)

Running time 96.38 seconds
```
{% endraw %}

## Запуск

Запуск контейнера выполняется командой [werf run]({{ site.docsurl }}/documentation/cli/main/run.html):

```shell
werf run --docker-options="--rm -p 3000:3000" -- bash -ec "bundle exec rails db:migrate RAILS_ENV=development && bundle exec puma"
```

Обратите внимание, что [параметры docker](https://docs.docker.com/engine/reference/run/) задаются опцией `--docker-options`, а команда запуска после двух дефисов.

В `werf.yaml` может описываться произвольное количество образов. Для запуска контейнера определённого `image` из werf.yaml необходимо использовать позиционный аргумент команды (`werf run basicapp ...`).

_Можно заметить, что вызов `werf run` осуществляет сборку недостающих образов, таким образом предварительная сборка командой `werf build` не обязательна._

Теперь приложение доступно локально на порту 3000:

![](/guides/images/rails/100_10_app_in_browser.png)

<div id="go-forth-button">
    <go-forth url="20_cluster.html" label="Подготовка кластера" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
