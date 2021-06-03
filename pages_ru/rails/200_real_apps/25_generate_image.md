---
title: Backend для генерации изображений
permalink: rails/200_real_apps/25_generate_image.html
layout: development
---

В этой статье мы:
- создадим отдельный backend нашего приложения, взяв за основу приложение из главы basic-apps;
- добавим метод для генерации изображения по входному тексту;
- рассмотрим тему логгирования приложения и просмотр логов запущенного приложения.

- В предыдущей главе basic-apps мы закончили на приложении, которое позволяет создавать labels по url /api/labels.
- Для начала переименуем компонент приложения basicapp из главы basic-apps в backend
    - Зачем: этот компонент далее будет отвечать только за внутренний API нашего приложения.
    - Чтобы это сделать
        - Переименовываем файл Dockerfile в backend.Dockerfile
        - Переименовываем image basicapp в werf.yaml в backend
            - examples/rails/020_rename_to_backend/werf.yaml
            - Меняем .Values.werf.image.basicapp на .Values.werf.image.backend в deployment
                - examples/rails/020_rename_to_backend/.helm/templates/deployment.yaml
        - Переименовываем basicapp в backend в deployment
            - examples/rails/020_rename_to_backend/.helm/templates/deployment.yaml
        - Переименовываем basicapp в backend в service
            - examples/rails/020_rename_to_backend/.helm/templates/service.yaml
        - Переименовываем serviceName = basicapp в backend в ingress (само имя ingress не трогаем)
            - examples/rails/020_rename_to_backend/.helm/templates/ingress.yaml

- Добавим новый метод POST /api/generate-image
    - Заводим новый images controller с методом api: generate_image
        - На вход получаем текст из параметра text.
        - Обращаем внимание на то, что используется стандартный rails-логгер.
            - Как его правильно сконфигурировать и как потом смотреть логи будет расмотрено ниже.
        - examples/rails/021_generate_image/app/controllers/api/images_controller.rb
    - Создаём route на новый метод api `images#generate_image`:
        - examples/rails/021_generate_image/config/routes.rb
    - Обновляем Gemfile и Gemfile.lock: подключается новая зависимость приложения rmagick
        - examples/rails/021_generate_image/Gemfile
        - examples/rails/021_generate_image/Gemfile.lock
    - Установим также для корректной работы библиотеки rmagick пакет gsfonts в базовый образ:
        - examples/rails/021_generate_image/backend.Dockerfile

- Сконфигурируем логгирование и окружение нашего приложения.
    - По умолчанию rails-приложение предоставляет несколько окружений: development, test и production.
        - В окружении production при установке переменной окружения `RAILS_LOG_TO_STDOUT=1` стандартный логгер будет автоматом писать в stdout.
    - => Переделаем стандартный логгер рельсов для всех окружений на такой логгер, который будет писать в stdout и stderr — это наиболее правильный способ логирования приложений, запускаемых в кубах.
        - Дальнейший сбор и хранение логов (при необходимости) будет осуществляться отдельными решениями.
            - Принципиальный момент в том, что в приложении мы об этом не думаем и просто пишем в stdout/stderr.
        - Добавляем конфигурацию логгера приложения для всех окружений:
            - examples/rails/021_generate_image/config/application.rb
        - Удаляем стандартную конфигурацию логгера для production-окружения
            - examples/rails/021_generate_image/config/environments/production.rb

    - => Поменяем стандартный rails setup так, чтобы во всех окружениях была доступна переменная окружения `RAILS_LOG_TO_STDOUT=1` 
        - examples/rails/021_generate_image/

- Мы готовы выкатить новую версию приложения
    - добавляем все внесённые изменения в гит
        - делаем git add . ; git commit
    - werf converge ...

- Проверяем результат:
    
    ```
    curl -v -X POST "example.com/api/generate-image?text=do%20it" -o /tmp/out
    ```

    - В файле /tmp/out должно появится изображение с указанным тестом.
    - Если вдруг не получилось — не беда, сейчас будем смотреть логи приложения, также в случае ошибки в файле /tmp/out будет текст ошибки.

- Как посмотреть логи?

    - За api нашего приложения отвечает backend, который запущен в Deployment.
        - Задача deployment — создать указанное количество реплик Pod'ов.
        - Каждый Pod — это запущенный экземпляр нашего приложения на рельсах.
        - Соответственно чтобы увидеть логи backend-а, надо запросить эти логи для одного из созданных подов.

    - Получаем список подов:

    ```
    kubectl -n werf-guided-rails get pod
    ```

    =>

    ```
    NAME                      READY   STATUS    RESTARTS   AGE
    backend-7964b6b68-n28v5   1/1     Running   0          21m
    backend-7964b6b68-wzgw4   1/1     Running   0          21m
    mysql-666f76d7cb-967xb    1/1     Running   3          4d6h
    ```

    - Выбираем один из подов с префиксом `backend` и смотрим его логи следующей командой:

    ```
    kubectl -n werf-guided-rails logs -f backend-7964b6b68-n28v5
    ```

    =>

    ```
    I, [2021-06-07T16:47:28.032279 #1]  INFO -- : Started POST "/api/generate-image?text=do%20it" for 172.17.0.2 at 2021-06-07 16:47:28 +0000
    I, [2021-06-07T16:47:28.035225 #1]  INFO -- : Processing by Api::ImagesController#generate_image as JSON
    I, [2021-06-07T16:47:28.035426 #1]  INFO -- :   Parameters: {"text"=>"do it"}
    D, [2021-06-07T16:47:28.035787 #1] DEBUG -- : received generate_image request
    D, [2021-06-07T16:47:28.036043 #1] DEBUG -- : start generating image for text "do it"
    D, [2021-06-07T16:47:28.498286 #1] DEBUG -- : finish generation image for text "do it"
    I, [2021-06-07T16:47:28.498673 #1]  INFO -- :   Rendering text template
    I, [2021-06-07T16:47:28.498811 #1]  INFO -- :   Rendered text template (Duration: 0.0ms | Allocations: 1)
    I, [2021-06-07T16:47:28.498897 #1]  INFO -- : Sent data  (0.5ms)
    I, [2021-06-07T16:47:28.498972 #1]  INFO -- : Completed 200 OK in 463ms (Views: 0.4ms | ActiveRecord: 0.0ms | Allocations: 247)
    ```
