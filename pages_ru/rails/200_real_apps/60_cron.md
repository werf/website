---
title: Периодическое выполнение заданий
permalink: rails/200_real_apps/60_cron.html
layout: wip
---

В этой главе мы:
- Покажем как запускать периодические фоновые задания на примере rails rake tasks.
- Расскажем что такое CronJob и чем он отличается от Job.
- Расскажем для чего ещё можно использовать CronJob.
- Ответим на часто возникающие вопросы связанные с CronJob.

<cut>

<!-- TODO: Надо сделать шаг подготовка сворачиваемым и по умолчанию свёрнутым -->

## Подготовка

Возьмём за основу web-приложение из раздела "первые шаги". Состояние директории `rails-app` должно соответствовать шагу `examples/rails/018_fixup_consistency`:

```
git clone https://github.com/werf/werf-guides
cp -r werf-guides/examples/rails/018_fixup_consistency rails-app
cd rails-app
git init
git add .
git commit -m "initial"
```
</cut>

## Запускаем фоновые задания

[За основу взято наше web-приложение из раздела "первые шаги"](#подготовка). Данное приложение состоит из одного http api сервера.

Добавляем новые исходники в наше существующее приложение для выполнения периодических фоновых задач по обслуживанию нашего проекта:

```shell
cp ../werf-guides/examples/rails/800_cron/app/controllers/api/labels_controller.rb app/controllers/api/labels_controller.rb
cp ../werf-guides/examples/rails/800_cron/app/mailers/notifications_mailer.rb app/mailers/notifications_mailer.rb
cp ../werf-guides/examples/rails/800_cron/app/views/notifications_mailer/labels_count_report_email.html.erb app/views/notifications_mailer/labels_count_report_email.html.erb
cp ../werf-guides/examples/rails/800_cron/config/application.rb config/application.rb
cp ../werf-guides/examples/rails/800_cron/config/routes.rb config/routes.rb
cp ../werf-guides/examples/rails/800_cron/.helm/templates/cleanup-labels.yaml .helm/templates/cleanup-labels.yaml
cp ../werf-guides/examples/rails/800_cron/.helm/templates/debug-mails.yaml .helm/templates/debug-mails.yaml
cp ../werf-guides/examples/rails/800_cron/.helm/templates/ingress.yaml .helm/templates/ingress.yaml
cp ../werf-guides/examples/rails/800_cron/.helm/templates/send-report.yaml .helm/templates/send-report.yaml
cp ../werf-guides/examples/rails/800_cron/lib/tasks/crons.rake lib/tasks/crons.rake
```

### 1. Очистка устаревших записей в таблице labels

- Очищает те записи в таблице labels, у которых истёк time to live в 3 минуты.
- Данная задача может работать полностью автономно и не требует доступа к запущенному http-серверу.
- Особенность этой задачи в том, что она может долго выполняться в фоне, поэтому её выполнение внутри запущенного http-сервера не желательно.
- За реализацию очистки отвечает rake task `crons:cleanup_labels`.
    - Связанные исходники:
        - `lib/tasks/crons.rake`;
    - Вызывается из директории проекта командой `rake crons:cleanup_labels`. 
- Запуск данной очистки реализован через CronJob (о том, что это такое напишем далее) с именем `cleanup-labels` с запуском каждую минуту.
    - Связанные исходники:
        - `.helm/templates/cleanup-labels.yaml`.

### 2. Периодическая отсылка отчётов по почте

- Отправляет письмо с текущим количеством labels администратору системы.
- Данная задача реализована внутри http-сервера и инициируется специальным запросом по api: `/api/send-report`.
- Инициировать периодический запуск данного задания лучше вне http-сервера, а не внутренними механизмами сервера, т.к. он может быть запущен в нескольких экземплярах (подробности про реплики см. в "первых шагах") и в противном случае задание будет избыточно выполняться сразу на нескольких репликах.
- За реализацию отправки отвечает метод контроллера labels `send_report`:
    - Связанные исходники:
        - `app/controllers/api/labels_controller.rb`;
        - `app/mailers/notifications_mailer.rb`;
        - `app/views/notifications_mailer/labels_count_report_email.html.erb`.
    - Отправка происходит из самого http-сервера;
    - Отправка сделана через специальный дебаг-smtp-сервер (mailhog), который поднимается вместе с нашим приложением по адресу `debug-mails.example.com` и позволяет просматривать перехваченные письма — далее мы этим воспользуемся и посмотрим как оно работает.
        - Связанные исходники:
            - `.helm/templates/debug-mails.yaml`;
            - `.helm/templates/ingress.yaml`.
- Чтобы инициировать эту отправку периодически был создан CronJob `send-report` с запуском каждую минуту, который просто делает post-запрос на бекенд basicapp.
    - Связанные исходники:
        - `.helm/templates/send-report.yaml`.

### Выкат и проверка

Коммитим изменения:

```shell
git add .
git commit -m "Add cronjobs"
```

Запускаем выкат:

```shell
werf converge --repo REPO
```

#### Проверяем очистку устаревших labels

- Наша фоновая задача чистит те labels, время жизни которых больше 3х минут.
    - Чтобы было что очищать для начала создадим несколько штук:

        ```shell
        curl -XPOST "http://example.com/api/labels?label=aaa"
        curl -XPOST "http://example.com/api/labels?label=bbb"
        curl -XPOST "http://example.com/api/labels?label=ccc"
        curl -XPOST "http://example.com/api/labels?label=ddd"
        curl -XPOST "http://example.com/api/labels?label=eee"
        curl -XPOST "http://example.com/api/labels?label=fff"
        curl -XPOST "http://example.com/api/labels?label=ggg"
        ```

    - Проверим список вновь созданных labels:
    
        ```shell
        curl "http://example.com/api/labels"
        ```

- CronJob `cleanup-labels` будет создавать Job-ы по описанному jobTemplate-у. Затем эти Job-ы в свою очередь будут создавать Pod-ы. Каждый запускаемый Pod выполняет в отдельном контейнере команду `rake crons:cleanup_labels`.
- Для проверки работы Cronjob, после выката надо периодически запускать следующую команду до тех пор пока не увидим новые поды с с префиксом в имени `cleanup-labels-`:

    ```shell
    kubectl -n werf-guided-rails get cj,job,pod
    ```

- Затем скопируем имя одного из вновь созданных подов и посмотрим его логи:

    ```shell
    kubectl -n werf-guided-rails logs -f pod/cleanup-labels-1623946800-8bfhs
    ```

- Спустя 3 минуты снова проверим список существующих labels:

    ```shell
    curl "http://example.com/api/labels"
    ```

    - Список должен быть пустым.

#### Проверяем отправку репортов

- Наша фоновая задача выполняет отправку репортов о количестве записей в таблице labels каждую минуту.
- CronJob с именем `send-report` будет создавать Job-ы по описанному jobTemplate-у. Затем эти Job-ы в свою очередь будут создавать Pod-ы. Каждый запускаемый Pod выполняет в отдельном контейнере запрос на http-сервер `/api/send-report`. По запросу `/api/send-report` сервер уже выполняет отправку письма на почту администратору.
- Проверять отосланные письма мы будем через специальный дебаг smtp-сервер, который деплоится вместе с нашим приложением.
    - Для дебага наш http-сервер написан так, что пересылает письма через почтовый дебаг-smtp-сервер.
- Переходим по адресу [http://debug-mails.example.com](http://debug-mails.example.com).
    - Чтобы адрес заработал необходимо добавить `debug-mails.example.com` в хост-файл системы (см. [подготовка](#подготовка)).
    - Каждую минуту должно приходить новое письмо с репортом о текущем количестве записей в таблице Labels.

## Что такое CronJob

- Есть ресурс **Job**.
    - Один экземпляр ресурса представляет собой одно фоновое выполнение какой-то задачи.
    - Запускаемая команда будет работать в отдельном Pod-е согласно описанию шаблона `spec.template`.
    - После того как задача отработала, ресурс Job переходит в финальное состояние.
        - Экземпляр ресурса Job, однажды созданный в кластере, не может быть обновлён или изменён после создания.
            - Таким образом ресурс Job — одноразовый.
            - Чтобы инициировать повторное выполнение задачи необходимо создать новый экземпляр Job уже с другим именем.

- Чтобы организовать периодическое выполнение одной фоновой задачи есть ресурс **CronJob**.
    - Как видим из листингов ресурсов, мы как раз используем CronJob.
    - Также как из Job, CronJob ­— это про фоновое выполнение задач.
    - Но в отличие от Job, CronJob позволяет создавать многоразовые задачи
        - CronJob — это надстройка над Job, которая задаёт шаблон для генерации одноразовых Job-ов.
    - `spec.schedule` позволяет задать периодичность в привычном формате crontab.

- Рекомендуется также ознакомится с материалами про CronJob на сайте kubernetes:
    - [https://kubernetes.io/docs/tasks/job/automated-tasks-with-cron-jobs/](https://kubernetes.io/docs/tasks/job/automated-tasks-with-cron-jobs/);
    - [https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/).

## Для чего ещё может быть использован CronJob

Тут просто примеры применения.

- Периодическая индексация БД.
- Обработчик очереди заданий.
    - Задания поступают в некоторую очередь.
    - CronJob инициирует работу обработчика следующего задания из этой очереди.
- Инвалидация истекших пользовательских доступов в БД.
- Очистка временных данных и устаревших пользовательских сессий в БД.
- И т.п.

## Другие вопросы

- А что если я хочу запускать CronJob в том же поде и контейнере, где работает наш http-сервер?
    - Так делать не принято и через CronJob в чистом виде не получится.
        - Таким способом реализована очистка устаревших записей в таблице labels [описанная выше](#1-очистка-устаревших-записей-в-таблице-labels).
    - Для реализации такой схемы фоновое задание встраивается в само приложение (http-сервер в нашем случае).
        - Инициировать выполнение этого задания либо средствами самого приложения, либо через CronJob (предпочтительный вариант, по причинам описанным ранее).
        - Именно таким способом реализована периодическая отсылка отчётов [описанная выше](#2-периодическая-отсылка-отчётов-по-почте).
