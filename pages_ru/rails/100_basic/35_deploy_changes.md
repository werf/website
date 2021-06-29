---
title: Внесение изменений
permalink: rails/100_basic/35_deploy_changes.html
examples_initial: examples/rails/015_deploy_app
examples: examples/rails/016_scale
examples_unscale: examples/rails/017_unscale
examples_add_fields: examples/rails/018_add_fields
examples_consistency: examples/rails/019_fixup_consistency
description: |
    Внесём изменения в уже выкаченное приложение и его инфраструктуру. Продемонстрируем, как работает подход infrastructure-as-code (IaC).
---

## Работаем с инфраструктурой
### Просмотр состояния

Чтобы посмотреть состояние запущенного приложения в Kubernetes, существует несколько команд.

Получить список запущенных Deployment'ов и Pod'ов:
```shell
kubectl get deployment
kubectl get pod
```

Посмотреть summary-информацию по Deployment'у `basicapp`:
```shell
kubectl describe deployment basicapp
```

Deployment запускает Pod'ы. Логи пишутся в Pod'ах. Такой командой можно получить логи одного из запущенных Pod'ов:
```shell
kubectl logs basicapp-<podId>
```

### Масштабирование

Веб-сервер запущен в Deployment'е `basicapp`. Помотрим, сколько реплик запущено:

```shell
kubectl get pod
```

Ответ будет примерно таким:
```shell
NAME                      READY   STATUS    RESTARTS   AGE
basicapp-57789b68-kxcb9   1/1     Running   0          72m
```

Поменяем вручную на 4 реплики:
```shell
kubectl edit deployment basicapp
```

В открывшемся редакторе выставляем `spec.replicas=4`, закрываем редактор.
Снова посмотрим, сколько реплик запущено:
```shell
kubectl get pod
```

Ответ станет примерно таким:
```shell
NAME                      READY   STATUS    RESTARTS   AGE
basicapp-57789b68-nbdjb   1/1     Running   0          6s
basicapp-57789b68-dkbgx   1/1     Running   0          6s
basicapp-57789b68-c4thw   1/1     Running   0          6s
basicapp-57789b68-kxcb9   1/1     Running   0          72m
```

Мы произвели масштабирование вручную. Теперь снова запустим `werf converge`:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guided-rails
```

И снова посмотрим, сколько реплик запущено:
```shell
kubectl get pod
```
```shell
NAME                      READY   STATUS    RESTARTS   AGE
basicapp-57789b68-kxcb9   1/1     Running   0          72m
```

Количество реплик снова соответствует таковому в Git-репозитории. Дело в том, что werf привела состояние кластера к состоянию, описанному в текущем Git-коммите. Этот принцип называется **гитерминизмом** (giterminism).

Как же соблюсти **гитерминизм** и сделать всё правильно? Выставим тот же `spec.replicas=4`, но уже через состояние приложения, описанное в Git — в файле `.helm/templates/deployment.yaml`:

{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" examples=page.examples %}

Запустим деплой:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guided-rails
```

Проверим сколько реплик запущено теперь:

```shell
kubectl get pod
```
```shell
NAME                      READY   STATUS    RESTARTS   AGE
basicapp-57789b68-fsxlw   1/1     Running   0          7s
basicapp-57789b68-pqs2n   1/1     Running   0          7s
basicapp-57789b68-vx88n   1/1     Running   0          7s
basicapp-57789b68-kxcb9   1/1     Running   0          72m
```

Вернём состояние рабочей директории к предыдущему, с одной репликой:

{% include chapter_prepare_repo_commands.md.liquid examples_to=page.examples_unscale from_scratch=false %}

Проверим содержимое `.helm/templates/deployment.yaml`:

{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" examples=page.examples_unscale %}

Запустим деплой:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guided-rails
```

## Меняем приложение

Наше приложение — это [CRUD](https://ru.wikipedia.org/wiki/CRUD) для создания labels. Ниже рассмотрим некоторые действия, которые можно выполнить с приложением.

### Работа с приложением

Продемонстрируем получение списка labels из консоли:
```shell
curl "http://URL:PORT/api/labels"
```

Создадим новые labels из консоли:
```shell
curl -X POST "http://URL:PORT/api/labels/?label=red"
curl -X POST "http://URL:PORT/api/labels/?label=hot"
curl -X POST "http://URL:PORT/api/labels/?label=blue"
curl -X POST "http://URL:PORT/api/labels/?label=cold"
```

Получим обновлённый список labels:
```shell
curl "http://URL:PORT/api/labels"
```

### Добавление полей

Добавим новые timestamp-поля для label: время его создания и обновления.

Внесём все необходимые изменения в нашу рабочую директорию:

{% include chapter_prepare_repo_commands.md.liquid examples_to=page.examples_add_fields from_scratch=false %}

Внесённые изменения включают в себя:
- новый файл миграций `20210526202700_add_timestamps_to_labels.rb`:

    {% include snippetcut_example path="db/migrate/20210526202700_add_timestamps_to_labels.rb" syntax="ruby" examples=page.examples_add_fields %}

- обновлённую схему БД `schema.rb`:

    {% include snippetcut_example path="db/schema.rb" syntax="ruby" examples=page.examples_add_fields %}

- правки в `_label.json.jbuilder`, чтобы API выдавал не только `id` и имя `label`, но и поля `created_at` и `updated_at`:

    {% include snippetcut_example path="app/views/api/labels/_label.json.jbuilder" syntax="ruby" examples=page.examples_add_fields %}

Запустим деплой:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guided-rails
```

Проверим результат:
```shell
curl -X POST "http://URL:PORT/api/labels/?label=black"
curl -X POST "http://URL:PORT/api/labels/?label=white"
curl "http://URL:PORT/api/labels"
```

В ответе увидим новые поля. Поздравляем, у нас получилось!

### Обеспечим корректность работы нескольких реплик

В данный момент есть проблема: наше приложение использует БД SQLite в локальном файле, а мы добавили несколько реплик. Получается, что каждая реплика использует свою базу данных, а не общую. Это приведёт к тому, что запросы будут равномерно распределяться по репликам при создании и получении labels, а результаты запросов могут отличаться при перезапуске. Чтобы решить эту проблему, давайте переведем наше приложение на общую БД — MySQL.

Внесём все необходимые изменения в нашу рабочую директорию:

{% include chapter_prepare_repo_commands.md.liquid examples_to=page.examples_consistency from_scratch=false %}

Рассмотрим, какие изменения были внесены в приложение:
* Создан новый шаблон `database.yaml` для запуска MySQL в кластере.
{% include snippetcut_example path=".helm/templates/database.yaml" syntax="yaml" examples=page.examples_consistency %}
* В Deployment'е приложения миграции будут запускаться прямо перед стартом самого приложения. Также настройки SQLite более не нужны:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" examples=page.examples_consistency %}
* Приложение настроено на работу с MySQL:
{% include snippetcut_example path="config/database.yml" syntax="yaml" examples=page.examples_consistency %}
* Добавлен Ruby Gem для использования MySQL:
{% include snippetcut_example path="Gemfile" syntax="ruby" examples=page.examples_consistency %}
{% include snippetcut_example path="Gemfile.lock" syntax="gemfile.lock" examples=page.examples_consistency %}

Запустим деплой:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guided-rails
```

Дождёмся выполнения команды. Заметьте, что в процессе работы в логах могут появляться ошибки подключения к БД, потому что контейнер с ней ещё не успел запуститься. Это нормально - необходимо дождаться полного запуска приложения.

Проверим создание и получение labels и увидим консистентное поведение:
```shell
curl -X POST "http://URL:PORT/api/labels/?label=name"
curl "http://URL:PORT/api/labels/"
```

Так как у нас rails-приложение, мы когда-нибудь захотим зайти в работающий контейнер запущенного приложения и запустить там `rails console`. Это можно сделать вот так _(в общем случае это плохая практика, но допустимо на данном этапе для лучшей наглядности)_:
```shell
kubectl exec -ti deployment/basicapp -- bash
rails c
```

Получим список labels напрямую из базы данных:
```shell
Label.all
```
