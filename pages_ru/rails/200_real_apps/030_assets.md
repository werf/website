---
title: Раздача ассетов
permalink: rails/200_real_apps/030_assets.html
examples_initial: examples/rails/020_logging
examples: examples/rails/030_assets
base_url: https://github.com/werf/werf-guides/blob/master/
description: |
  В этой главе мы реализуем в приложении работу со статическими файлами и покажем, как правильно отдавать их клиенту.
---

## Добавление страницы `/image` в приложение

Добавим нашему приложению новый endpoint `/image`, который будет отдавать страницу, использующую набор статических файлов. Для сборки JS-, CSS- и медиафайлов будем использовать Webpacker (без Sprockets).

Сейчас наше приложение представляет собой простой API без адекватных средств управления статичными файлами и frontend-кодом. Чтобы сделать из этого API простое веб-приложение, мы сгенерировали скелет нового Rails-приложения, но без ранее использованных флагов `--api` и `--skip-javascript`:
```shell
rails new --skip-keeps --skip-action-mailer --skip-action-mailbox \
  --skip-action-text --skip-active-record --skip-active-job \
  --skip-active-storage --skip-action-cable --skip-sprockets \
  --skip-spring --skip-listen --skip-turbolinks --skip-jbuilder \
  --skip-test --skip-system-test --skip-bootsnap .
```

Из этого скелета мы позаимствовали всё, что связано с управлением JS-кодом и статическими файлами. Основные изменения, сделанные в нашем приложении:
1. Добавление `webpacker` в [Gemfile]({{ page.base_url | append: page.examples | append: "/Gemfile" }}).
1. Создание стандартных конфигурационных файлов Webpacker:
  * [config/webpacker.yml]({{ page.base_url | append: page.examples | append: "/config/webpacker.yml" }})
  * [config/webpack/environment.js]({{ page.base_url | append: page.examples | append: "/config/webpack/environment.js" }})
  * [config/webpack/production.js]({{ page.base_url | append: page.examples | append: "/config/webpack/production.js" }})
1. Создание набора директорий для размещения исходных статических файлов:
  * [app/javascript/styles]({{ page.base_url | append: page.examples | append: "/app/javascript/styles" }})
  * [app/javascript/src]({{ page.base_url | append: page.examples | append: "/app/javascript/src" }})
  * [app/javascript/images]({{ page.base_url | append: page.examples | append: "/app/javascript/images" }})

Теперь добавим шаблон HTML-страницы, которая будет доступна по адресу `/image` и отображать кнопку _Get image_:
{% include snippetcut_example path="app/views/layouts/image.html.erb" syntax="erb" examples=page.examples %}

При нажатии на кнопку _Get image_ должен выполняться Ajax-запрос, с помощью которого мы получим и отобразим [SVG-картинку]({{ page.base_url | append: page.examples | append: "/app/javascript/images/werf-logo.svg" }}):
{% include snippetcut_example path="app/javascript/src/image.js" syntax="js" examples=page.examples %}

Также на странице будет использоваться CSS-файл [app/javascript/styles/site.css]({{ page.base_url | append: page.examples | append: "/app/javascript/styles/site.css" }}).

JS-файл, CSS-файл и SVG-картинка будут собраны с Webpack и доступны по адресу `/packs/...`:
{% include snippetcut_example path="app/javascript/packs/application.js" syntax="js" examples=page.examples %}

Обновим контроллер и маршруты, добавив в них новый endpoint `/image`, обрабатывающий созданный выше HTML-шаблон:
{% include snippetcut_example path="app/controllers/application_controller.rb" syntax="rb" examples=page.examples %}
{% include snippetcut_example path="config/routes.rb" syntax="rb" examples=page.examples %}

Приложение обновлено: теперь, в дополнение к уже знакомому по прошлым главам `/ping`, у приложения есть новый endpoint `/image`. На последнем отображается страница, использующая для работы разные типы статических файлов.

>_[В начале главы](#подготовка-репозитория) описаны команды, с помощью которых можно увидеть полный, исчерпывающий список изменений, сделанных с приложением в этой главе._

## Организация раздачи статических файлов

По умолчанию в production-окружениях Rails даже не пытается раздавать статические файлы самостоятельно. Вместо этого разработчики Rails предлагают использовать для их раздачи reverse proxy вроде NGINX. Рекомендуется это потому, что производительность и эффективность reverse proxy для раздачи статических файлов на порядки выше, чем у Rails и Puma.

На практике обходиться без reverse proxy перед приложением допускается только во время разработки. Кроме более эффективной отдачи статических файлов reverse proxy хорошо дополняет application server (Puma), давая много дополнительных возможностей, которые обычно не реализуют на уровне application-серверов, а также предоставляя быструю и гибкую маршрутизацию запросов.

Есть несколько способов, как закрыть application server за reverse proxy в Kubernetes. Мы будем использовать распространенный и простой способ, который, тем не менее, хорошо масштабируется. В нем перед каждым Rails/Puma-контейнером поднимается NGINX-контейнер (в том же Pod), через который в Rails/Puma проксируются все запросы кроме запросов на статические файлы. Статические файлы отдаёт сам NGINX-контейнер, доставая их из своей файловой системы.

Приступим к непосредственной реализации.

## Обновление сборки и деплоя

Начнём с реорганизации сборки приложения. Теперь нам требуется собрать не только образ с Rails и Puma, но и образ с NGINX, где должны находиться статические файлы приложения, которые NGINX сможет сразу отдать клиенту:
{% include snippetcut_example path="Dockerfile" syntax="dockerfile" examples=page.examples %}

Во время сборки в образ с NGINX будет добавлен его конфигурационный файл:
{% include snippetcut_example path=".werf/nginx.conf" syntax="nginx" examples=page.examples %}

Обновим конфигурацию `werf.yaml`, чтобы werf собрал и сохранил два образа (backend, frontend) вместо одного:
{% include snippetcut_example path="werf.yaml" syntax="yaml" examples=page.examples %}

Добавим новый NGINX-контейнер в Deployment приложения:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" examples=page.examples %}

Теперь Service и Ingress должны обращаться на 80-й порт, а не 3000-й, чтобы все запросы шли через NGINX, а не напрямую на Rails/Puma:
{% include snippetcut_example path=".helm/templates/service.yaml" syntax="yaml" examples=page.examples %}
{% include snippetcut_example path=".helm/templates/ingress.yaml" syntax="yaml" examples=page.examples %}

## Проверка

Теперь попробуем переразвернуть приложение:
```shell
werf converge --repo <ИМЯ ПОЛЬЗОВАТЕЛЯ DOCKER HUB>/werf-guide-app
```

Ожидаемый результат:
```shell
...
┌ ⛵ image backend
│ ┌ Building stage backend/dockerfile
│ │ backend/dockerfile  Sending build context to Docker daemon    320kB
│ │ backend/dockerfile  Step 1/24 : FROM ruby:2.7 as base
│ │ backend/dockerfile   ---> 1faa5f2f8ca3
...
│ │ backend/dockerfile  Step 24/24 : LABEL werf-version=v1.2.12+fix2
│ │ backend/dockerfile   ---> Running in d6126083b530
│ │ backend/dockerfile  Removing intermediate container d6126083b530
│ │ backend/dockerfile   ---> 434a85e3df9d
│ │ backend/dockerfile  Successfully built 434a85e3df9d
│ │ backend/dockerfile  Successfully tagged 145074fa-9024-4bcb-9144-e245e1bbaf87:latest
│ │ ┌ Store stage into .../werf-guide-app
│ │ └ Store stage into .../werf-guide-app (14.27 seconds)
│ ├ Info
│ │      name: .../werf-guide-app:7d52fed109415751ad6f28dd259d7a201dc2ffe34863e70f0a404adf-1629129180941
│ │        id: 434a85e3df9d
│ │   created: 2021-08-16 18:53:00 +0300 MSK
│ │      size: 335.1 MiB
│ └ Building stage backend/dockerfile (23.70 seconds)
└ ⛵ image backend (30.13 seconds)

┌ ⛵ image frontend
│ ┌ Building stage frontend/dockerfile
│ │ frontend/dockerfile  Sending build context to Docker daemon    320kB
│ │ frontend/dockerfile  Step 1/28 : FROM ruby:2.7 as base
│ │ frontend/dockerfile   ---> 1faa5f2f8ca3
...
│ │ frontend/dockerfile  Step 28/28 : LABEL werf-version=v1.2.12+fix2
│ │ frontend/dockerfile   ---> Running in 29323e1512c7
│ │ frontend/dockerfile  Removing intermediate container 29323e1512c7
│ │ frontend/dockerfile   ---> 8f637e588672
│ │ frontend/dockerfile  Successfully built 8f637e588672
│ │ frontend/dockerfile  Successfully tagged 215f9544-2637-4ecf-9b81-8cbbcd6777fe:latest
│ │ ┌ Store stage into .../werf-guide-app
│ │ └ Store stage into .../werf-guide-app (9.83 seconds)
│ ├ Info
│ │      name: .../werf-guide-app:bdb809c26315846927553a069f8d6fe72c80e33d2a599df9c6ed2be0-1629129180706
│ │        id: 8f637e588672
│ │   created: 2021-08-16 18:53:00 +0300 MSK
│ │      size: 9.4 MiB
│ └ Building stage frontend/dockerfile (19.02 seconds)
└ ⛵ image frontend (29.28 seconds)

┌ Waiting for release resources to become ready
│ ┌ Status progress
│ │  DEPLOYMENT      REPLICAS                    AVAILABLE  UP-TO-DATE
│ │  werf-guide-app  2/1                         1          1
│ │  │               POD                         READY      RESTARTS    STATUS             ---
│ │  ├──             guide-app-566dbd7cdc-n8hxj  2/2        0           Running            Waiting  for:  replicas  2->1
│ │  └──             guide-app-c9f746755-2m98g   0/2        0           ContainerCreating
│ └ Status progress
│
│ ┌ Status progress
│ │  DEPLOYMENT      REPLICAS                    AVAILABLE  UP-TO-DATE
│ │  werf-guide-app  2->1/1                      1          1
│ │  │               POD                         READY      RESTARTS    STATUS
│ │  ├──             guide-app-566dbd7cdc-n8hxj  2/2        0           Running            ->  Terminating
│ │  └──             guide-app-c9f746755-2m98g   2/2        0           ContainerCreating  ->  Running
│ └ Status progress
└ Waiting for release resources to become ready (9.11 seconds)

Release "werf-guide-app" has been upgraded. Happy Helming!
NAME: werf-guide-app
LAST DEPLOYED: Mon Aug 16 18:53:19 2021
NAMESPACE: werf-guide-app
STATUS: deployed
REVISION: 23
TEST SUITE: None
Running time 41.55 seconds
```

Откроем в браузере [http://werf-guide-app/image](http://werf-guide-app/image) и нажмём на кнопку _Get image_. Ожидаемый результат:

{% asset guides/rails/030_assets_success.png %}

Также обратим внимание на то, какие ресурсы были запрошены и по каким ссылкам (последний ресурс здесь получен через Ajax-запрос):

{% asset guides/rails/030_assets_resources.png %}

Теперь наше приложение является не просто API, но веб-приложением, которое имеет средства для эффективного менеджмента статических файлов и JavaScript.

Также наше приложение готово выдерживать приличные нагрузки при большом количестве запросов к статическим файлам, и эти запросы не будут сказываться на работе приложения в целом. Масштабирование же Puma (отвечает за динамический контент) и NGINX (статический контент) происходит простым увеличением количества реплик (`replicas`) в Deployment'е приложения.
