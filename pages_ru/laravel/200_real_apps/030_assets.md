---
title: Раздача ассетов
permalink: laravel/200_real_apps/030_assets.html
examples_initial: examples/laravel/020_logging
examples: examples/laravel/030_assets
base_url: https://github.com/werf/werf-guides/blob/master/
description: |
  В этой главе мы реализуем в приложении работу со статическими файлами и покажем, как правильно отдавать их клиенту.
---

## Добавление страницы `/image` в приложение

Добавим нашему приложению новый endpoint `/image`, который будет отдавать страницу, использующую набор статических файлов. Для сборки JS-, CSS- и медиафайлов будем использовать штатный [laravel-mix](https://laravel.com/docs/8.x/mix), являющийся по сути надстройкой над webpack.

Сейчас наше приложение представляет собой простой API без адекватных средств управления статичными файлами и frontend-кодом. 

Чтобы сделать из этого API простое веб-приложение, мы позаимствовали всё, что связано с управлением JS-кодом и статическими файлами из [официального репозитория Laravel](https://github.com/laravel/laravel). Основные изменения, сделанные в нашем приложении:
1. Добавление файла [package.json]({{ page.base_url | append: page.examples | append: "/package.json" }}).
1. Создание конфигурационного файла laravel-mix [webpack.mix.js]({{ page.base_url | append: page.examples | append: "/webpack.mix.js" }}).
1. Создание набора директорий для размещения исходных статических файлов:
* [resources/css]({{ page.base_url | append: page.examples | append: "/resources/css" }})
* [resources/js]({{ page.base_url | append: page.examples | append: "/resources/js" }})
* [resources/images]({{ page.base_url | append: page.examples | append: "/resources/images" }})

Теперь добавим шаблон HTML-страницы, которая будет доступна по адресу `/image` и отображать кнопку _Get image_:
{% include snippetcut_example path="resources/views/image.blade.php" syntax="php" examples=page.examples %}

При нажатии на кнопку _Get image_ должен выполняться Ajax-запрос, с помощью которого мы получим и отобразим [SVG-картинку]({{ page.base_url | append: page.examples | append: "/resources/images/werf-logo.svg" }}):
{% include snippetcut_example path="resources/js/image.js" syntax="js" examples=page.examples %}

Также на странице будет использоваться CSS-файл [resources/css/site.css]({{ page.base_url | append: page.examples | append: "/resources/css/site.css" }}).

JS-файл, CSS-файл и SVG-картинка будут собраны с Webpack и будут доступны в поддиректориях `public/...`:
```shell
$ tree public/
public/
├── css
│   └── site.css
├── favicon.ico
├── images
│   └── werf-logo.svg
├── index.php
├── js
│   └── image.js
├── mix-manifest.json
└── robots.txt
```

Обновим маршруты, добавив в них новый endpoint `/image`, обрабатывающий созданный выше HTML-шаблон:
{% include snippetcut_example path="routes/web.php" syntax="php" examples=page.examples %}

Приложение обновлено: теперь, в дополнение к уже знакомому по прошлым главам `/ping`, у приложения есть новый endpoint `/image`. На последнем отображается страница, использующая для работы разные типы статических файлов.

>_[В начале главы](#подготовка-репозитория) описаны команды, с помощью которых можно увидеть полный, исчерпывающий список изменений, сделанных с приложением в этой главе._

## Организация раздачи статических файлов

Статические файлы отдаёт сам NGINX-контейнер, доставая их из своей файловой системы.

Приступим к непосредственной реализации.

## Обновление сборки и деплоя

Начнём с реорганизации сборки приложения. Добавляем новую стадию для сборки статических файлы приложения:
{% include snippetcut_example path="Dockerfile" syntax="dockerfile" snippet="assets" examples=page.examples %}
Полученный результат подкладываем в контейнер с NGINX:
{% include snippetcut_example path="Dockerfile" syntax="dockerfile" snippet="frontend" examples=page.examples %}

## Проверка

Теперь попробуем переразвернуть приложение:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guide-app
```

Ожидаемый результат:
{% include snippetcut_example path="werf-converge.log" syntax="shell" examples=page.examples %}

Откроем в браузере [http://werf-guide-app/image](http://werf-guide-app/image) и нажмём на кнопку _Get image_. Ожидаемый результат:

{% asset guides/laravel/030_assets_success.png %}

Также обратим внимание на то, какие ресурсы были запрошены и по каким ссылкам (последний ресурс здесь получен через Ajax-запрос):

{% asset guides/laravel/030_assets_resources.png %}

Теперь наше приложение готово выдерживать приличные нагрузки при большом количестве запросов к статическим файлам, и эти запросы не будут сказываться на работе приложения в целом. Масштабирование же php-fpm (отвечает за динамический контент) и NGINX (статический контент) происходит простым увеличением количества реплик (`replicas`) в Deployment'е приложения.
