---
title: Логирование
permalink: rails/200_real_apps/020_logging.html
examples_initial: examples/rails/010_basic_app
examples: examples/rails/020_logging
description: |
  В этой главе мы настроим и разберём особенности логирования приложения в Kubernetes, а также сделаем структурированный формат логов для последущеющего парсинга системами сбора и анализа логов.
---

## Вывод логов в stdout

У задеплоенных в Kubernetes приложений логи всегда должны отправляться в [stdout/stderr](https://ru.wikipedia.org/wiki/%D0%A1%D1%82%D0%B0%D0%BD%D0%B4%D0%B0%D1%80%D1%82%D0%BD%D1%8B%D0%B5_%D0%BF%D0%BE%D1%82%D0%BE%D0%BA%D0%B8#%D0%A1%D1%82%D0%B0%D0%BD%D0%B4%D0%B0%D1%80%D1%82%D0%BD%D1%8B%D0%B9_%D0%B2%D1%8B%D0%B2%D0%BE%D0%B4). Отправка логов в стандартные потоки делает логи доступными для Kubernetes и систем сбора логов, помогает избегать потерь логов при пересоздании контейнеров, а также переполнения логами дисков Kubernetes-узлов при сохранении их в файлах внутри контейнеров.

По умолчанию Rails не пишет логи ни в stdout, ни в stderr, а сохраняет их в файл. При использовании стандартной конфигурации необходимо использовать переменную окружения `RAILS_LOG_TO_STDOUT=1`, чтобы включить вывод всех логов (в т.ч. ошибок) в stdout.

Так как нам не понадобится писать логи куда-либо кроме stdout, то, независимо от переменных окружения, логи будем всегда направлять в stdout.

{% include snippetcut_example path="config/environments/production.rb" syntax="ruby" snippet="logger" examples=page.examples %}

## Формат логов

Логи Rails-приложений по умолчанию представляет собой обычный текст:
```shell
=> Booting Puma
=> Rails 6.1.4 application starting in production
=> Run "bin/rails server --help" for more startup options
Puma starting in single mode...
* Puma version: 5.3.2 (ruby 2.7.4-p191) ("Sweetnighter")
*  Min threads: 5
*  Max threads: 5
*  Environment: production
*          PID: 1
* Listening on http://0.0.0.0:3000
Use Ctrl-C to stop
I, [2021-06-07T16:47:28.498897 #1]  INFO -- : Started GET \"/ping\" for 192.168.49.1 at 2021-07-23 15:08:16 +0000
I, [2021-06-07T16:47:28.498972 #1]  INFO -- : Processing by ApplicationController#ping as */*
I, [2021-06-07T16:47:28.498897 #1]  INFO -- : Completed 200 OK in 0ms (Views: 0.1ms | Allocations: 166)
```

Обратите внимание, как логи приложения оказываются перемешаны с логами Rails и логами Puma, и все они имеют разный формат. Такой текст парсить системами сбора и анализа логов будет очень непросто.

Эту проблему можно решить, если отдавать логи в структурированном формате вроде JSON: система сбора логов обычно очень просто распознает и разбирает JSON-логи, а также корректно обрабатывает логи/сообщения в неожиданных, неструктурированных форматах, которые могут вклиниваться между логами в JSON-формате.

Реализуем свой `ActiveSupport::Logger::SimpleFormatter`, который будет отдавать логи в JSON вместо plain text:
{% include snippetcut_example path="lib/json_simple_formatter.rb" syntax="ruby" examples=page.examples %}

> _Добавление поддержки тегирования логов выходит за рамки этого руководства, но может быть реализовано при необходимости._

А теперь используем наш новый `JSONSimpleFormatter` по умолчанию, чтобы все логи отдавались в JSON:
{% include snippetcut_example path="config/environments/production.rb" syntax="ruby" snippet="log_formatter" examples=page.examples %}

## Управление уровнем логирования

По умолчанию уровень логирования приложения для production-окружения задан как `info`. Но иногда возникает необходимость изменить это.

К примеру, при диагностике проблем в production может помочь переключение логирования с `info` на `debug` для получения дополнительной отладочной информации. Но если приложение работает в большом количестве реплик, то переключать все реплики приложения на `debug` может быть не лучшей идеей, т.к. это может сказаться на безопасности, а также сильно увеличить нагрузку на компоненты, ответственные за сбор, хранение и анализ логов.

Решить эту проблему поможет возможность выставлять уровень логирования через переменную окружения. Это позволит, например, запустить рядом с уже существующим Deployment'ом приложения, у которого уровень логирования `info`, точно такой же Deployment, но в одной реплике и с уровнем логирования `debug`. Также для этого нового Deployment'а мы можем отключить централизованный сбор логов, если таковой имеется. Всё это в совокупности позволит нам не перегружать системы сбора логов и не сохранять в них потенциально небезопасные отладочные логи.

Возможность указывать уровень логирования через переменную окружения `RAILS_LOG_LEVEL` реализуется так:
{% include snippetcut_example path="config/environments/production.rb" syntax="ruby" snippet="log_level" examples=page.examples %}

Если переменная окружения не указана, будет использован стандартный уровень `info`.

## Фильтрация логов

Добавим следующую конфигурацию для скрытия секретных параметров в логах:
{% include snippetcut_example path="config/initializers/filter_parameter_logging.rb" syntax="ruby" examples=page.examples %}

> _Не забывайте обновлять этот список при добавлении новых секретов._

## Отображение логов при деплое с werf

По умолчанию werf при деплое показывает логи всех контейнеров приложения и делает это до тех пор, пока они не перейдут в состояние `Ready`. 
С помощью специальных werf-аннотаций можно организовать фильтрацию логов и выводить только те строки, которые удовлетворяют заданным шаблонам.
В дополнение к этому можно настроить вывод логов только для определённых контейнеров.

К примеру, следующим образом можно отключить вывод контейнера с именем `container_name` при выкате:
```yaml
annotations:
  werf.io/skip-logs-for-containers: container_name
```

А в следующем примере показывается, как можно выводить только те строки, которые удовлетворяют заданному регулярному выражению:
```yaml
annotations:
  werf.io/log-regex: ".*ERROR.*"
```

Список всех доступных аннотаций можно посмотреть [здесь](https://ru.werf.io/documentation/v1.2/reference/deploy_annotations.html).

> _Эти аннотации влияют только на то, как логи отображаются при деплое с werf. Они не оказывают никакого влияния на развертываемое приложение или его конфигурацию. Логи по-прежнему доступны в stdout/stderr контейнера в первоначальном виде._

## Проверка работоспособности

Теперь попробуем развернуть приложение:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guide-app
```

Ожидаемый результат:
```shell
...
┌ ⛵ image app
│ ┌ Building stage app/dockerfile
│ │ app/dockerfile  Sending build context to Docker daemon  30.72kB
│ │ app/dockerfile  Step 1/13 : FROM ruby:2.7
│ │ app/dockerfile   ---> 1faa5f2f8ca3
...
│ │ app/dockerfile  Step 13/13 : LABEL werf-version=v1.2.11+fix10
│ │ app/dockerfile   ---> Running in db6e76c3f427
│ │ app/dockerfile  Removing intermediate container db6e76c3f427
│ │ app/dockerfile   ---> d7f69fbfdedb
│ │ app/dockerfile  Successfully built d7f69fbfdedb
│ │ app/dockerfile  Successfully tagged eb50cb50-1191-4d0b-8bf2-a4d5bba18ecf:latest
│ ├ Info
│ │      name: .../werf-guide-app:...
│ │        id: d7f69fbfdedb
│ │   created: 2021-07-23 18:00:19 +0300 MSK
│ │      size: 327.3 MiB
│ └ Building stage app/dockerfile (12.72 seconds)
└ ⛵ image app (17.81 seconds)

┌ Waiting for release resources to become ready
│ ┌ Status progress
│ │ DEPLOYMENT      REPLICAS  AVAILABLE  UP-TO-DATE
│ │ werf-guide-app  1/1       1          1
│ │ │   POD                         READY  RESTARTS  STATUS
│ │ ├── guide-app-5f97776488-vwqfg  1/1    0         Terminating
│ │ └── guide-app-fcf7c4ff5-wvb62   1/1    0         Running
│ └ Status progress
└ Waiting for release resources to become ready (4.89 seconds)

Release "werf-guide-app" has been upgraded. Happy Helming!
NAME: werf-guide-app
LAST DEPLOYED: Fri Jul 23 18:00:34 2021
NAMESPACE: werf-guide-app
STATUS: deployed
REVISION: 14
TEST SUITE: None
Running time 26.67 seconds
```

Выполним пару запросов для генерации логов:
```shell
curl http://werf-guide-app/ping       # вернёт "pong" и код возврата 200
curl http://werf-guide-app/not_found  # ничего не вернёт, только код возврата 404
```

Проверим логи:
```shell
kubectl logs deploy/werf-guide-app
```

Ожидаемый результат:
```shell
=> Booting Puma
=> Rails 6.1.4 application starting in production
=> Run `bin/rails server --help` for more startup options
Puma starting in single mode...
* Puma version: 5.3.2 (ruby 2.7.4-p191) ("Sweetnighter")
*  Min threads: 5
*  Max threads: 5
*  Environment: production
*          PID: 1
* Listening on http://0.0.0.0:3000
Use Ctrl-C to stop
{"type":"INFO","time":"2021-07-23T15:11:36+00:00","message":"Started GET \"/ping\" for 192.168.49.1 at 2021-07-23 15:11:36 +0000"}
{"type":"INFO","time":"2021-07-23T15:11:36+00:00","message":"Processing by ApplicationController#ping as */*"}
{"type":"INFO","time":"2021-07-23T15:11:36+00:00","message":"Completed 200 OK in 0ms (Views: 0.1ms | Allocations: 135)"}
{"type":"INFO","time":"2021-07-23T15:11:30+00:00","message":"Started GET \"/not_found\" for 192.168.49.1 at 2021-07-23 15:11:30 +0000"}
{"type":"FATAL","time":"2021-07-23T15:11:30+00:00","message":"  \nActionController::RoutingError (No route matches [GET] \"/not_found\"):\n  "}
```

Как видим, логи приложения теперь отдаются в JSON и легко могут быть распарсены. Логи Rails и Puma по-прежнему отдаются обычным текстом, но главное, что системы сбора логов теперь не будут пытаться распарсить логи приложения и логи Rails/Puma так, как будто они имеют один формат: JSON-логи сохранятся отдельно и для них будет возможен поиск/фильтрация только по нужным полям.
