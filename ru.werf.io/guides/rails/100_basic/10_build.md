---
title: Сборка образа
permalink: rails/100_basic/10_build.html
---

В этой главе мы соберём образ с демо-приложением, используя werf и [Dockerfile](https://docs.docker.com/engine/reference/builder/), а потом запустим приложение в контейнере на основе собранного образа.

## Создадим новый репозиторий с демо-приложением

[Установите werf]({{ site.docsurl }}/installation.html), после чего [установите Docker](https://docs.docker.com/get-docker/).

В отдельной директории на своём компьютере выполните:
```shell
git clone https://github.com/werf/werf-guides
cp -r werf-guides/examples/rails/000_app rails-app
cd rails-app
git init
git add .
git commit -m "initial"
```

## Создадим Dockerfile

Реализуем логику сборки нашего приложения в [Dockerfile](https://docs.docker.com/engine/reference/builder/):

{% snippetcut name="Dockerfile" url="https://github.com/werf/werf-guides/blob/master/examples/rails/010_build/Dockerfile" %}
{% raw %}
```Dockerfile
FROM ruby:2.7.1

# Добавляем в образ всё содержимое нашего репозитория, включая код приложения:
WORKDIR /app
COPY . .

# Устанавливаем системные зависимости:
RUN apt-get update -qq && apt-get install -y build-essential libpq-dev libxml2-dev libxslt1-dev curl

# Устанавливаем зависимости приложения:
RUN bundle install
```
{% endraw %}
{% endsnippetcut %}

## Интеграция werf и Dockerfile

Создадим в корне репозитория файл `werf.yaml`, в котором укажем, какой Dockerfile должен будет использоваться при сборке с werf:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-rails  # Имя проекта, по умолчанию подставляется в имя Helm-релиза и имя Namespace при деплое.
configVersion: 1

---
image: basicapp  # Имя образа, может использоваться в Helm-шаблонах и для вызова некоторых команд werf.
dockerfile: Dockerfile  # Путь к Dockerfile, в котором описана сборка для этого образа.
```
{% endraw %}
{% endsnippetcut %}

В `werf.yaml` может описываться сборка сразу несколько образов. Также для сборки образа существует ряд дополнительных настроек, с которыми можно ознакомиться [по ссылке]({{ site.docsurl }}/documentation/reference/werf_yaml.html#dockerfile-image-section-image).

## Сборка с werf

Перед выполнением сборки необходимо сделать коммит с нашими изменениями:
```shell
git add werf.yaml Dockerfile
git commit -m "Add build configuration"
```

> Почему изменения должны быть добавлены в коммит и как обойтись без этого при локальной разработке мы разберём далее в главе «Необходимо знать».

Сборка выполняется командой [`werf build`]({{ site.docsurl }}/documentation/reference/cli/werf_build.html):

{% raw %}
```shell
$ werf build
    ...
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

## Запуск контейнера с приложением локально

Запустить контейнер локально на основе собранного образа можно командой [werf run]({{ site.docsurl }}/documentation/cli/main/run.html):
```shell
werf run --docker-options="--rm -p 3000:3000" basicapp -- sh -ec "bundle exec rails db:migrate RAILS_ENV=development && bundle exec puma"
```

[Параметры Docker](https://docs.docker.com/engine/reference/run/) здесь задаются опцией `--docker-options`, а команда для выполнения в контейнере указывается в конце, после двух дефисов.

Теперь приложение доступно на [http://127.0.0.1:3000/](http://127.0.0.1:3000/):

![](/guides/images/rails/100_10_app_in_browser.png)

<div id="go-forth-button">
    <go-forth url="20_cluster.html" label="Подготовка кластера" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
