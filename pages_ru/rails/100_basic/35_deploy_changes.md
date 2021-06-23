---
title: Внесение изменений
permalink: rails/100_basic/35_deploy_changes.html
---
Внесём изменения в уже выкаченное приложение и его инфраструктуру. Продемонстрируем, как работает подход infrastructure-as-code (IaC).

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

```shell
cp ../werf-guides/examples/rails/016_scale/.helm/templates/deployment.yaml .helm/templates/deployment.yaml
```

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/016_scale/.helm/templates/deployment.yaml" %}
{% include_file "examples/rails/016_scale/.helm/templates/deployment.yaml" %}
{% endsnippetcut %}

Закоммитим изменения:
```shell
git add .
git commit -m go
```

Запустим выкат:
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

Вернём состояние в рабочей директории к начальному:

```shell
cp ../werf-guides/examples/rails/017_unscale/.helm/templates/deployment.yaml .helm/templates/deployment.yaml
```

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/017_unscale/.helm/templates/deployment.yaml" %}
{% include_file "examples/rails/017_unscale/.helm/templates/deployment.yaml" %}
{% endsnippetcut %}

Закоммитим изменения:
```shell
git add .
git commit -m go
```

Запустим деплой:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guided-rails
```

## Меняем приложение

Наше приложение ­— это [CRUD](https://ru.wikipedia.org/wiki/CRUD) для создания labels. Ниже рассмотрим некоторые действия, которые можно выполнить с приложением.

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

```shell
cp ../werf-guides/examples/rails/018_add_fields/db/migrate/20210526202700_add_timestamps_to_labels.rb db/migrate/20210526202700_add_timestamps_to_labels.rb
cp ../werf-guides/examples/rails/018_add_fields/db/schema.rb db/schema.rb
cp ../werf-guides/examples/rails/018_add_fields/app/views/api/labels/_label.json.jbuilder app/views/api/labels/_label.json.jbuilder
```

Внесённые изменения включают в себя:
- новый файл миграций `20210526202700_add_timestamps_to_labels.rb`:

    {% snippetcut name="db/migrate/20210526202700_add_timestamps_to_labels.rb" url="https://github.com/werf/werf-guides/blob/master/examples/rails/018_add_fields/db/migrate/20210526202700_add_timestamps_to_labels.rb" %}
    {% include_file "examples/rails/018_add_fields/db/migrate/20210526202700_add_timestamps_to_labels.rb" %}
    {% endsnippetcut %}

- обновлённую схему БД `schema.rb`:

    {% snippetcut name="db/schema.rb" url="https://github.com/werf/werf-guides/blob/master/examples/rails/018_add_fields/db/schema.rb" %}
    {% include_file "examples/rails/018_add_fields/.helm/templates/deployment.yaml" %}
    {% endsnippetcut %}

- правки в `_label.json.jbuilder`, чтобы API выдавал не только `id` и имя `label`, но и поля `created_at` и `updated_at`:

    {% snippetcut name="app/views/api/labels/_label.json.jbuilder" url="https://github.com/werf/werf-guides/blob/master/examples/rails/018_add_fields/app/views/api/labels/_label.json.jbuilder" %}
    {% include_file "examples/rails/018_add_fields/app/views/api/labels/_label.json.jbuilder" %}
    {% endsnippetcut %}

Сделаем коммит изменений:
```shell
git add .
git commit -m go
```

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

```shell
cp ../werf-guides/examples/rails/019_fixup_consistency/.helm/templates/database.yaml .helm/templates/database.yaml
cp ../werf-guides/examples/rails/019_fixup_consistency/.helm/templates/deployment.yaml .helm/templates/deployment.yaml
cp ../werf-guides/examples/rails/019_fixup_consistency/config/database.yml config/database.yml
cp ../werf-guides/examples/rails/019_fixup_consistency/Gemfile Gemfile
cp ../werf-guides/examples/rails/019_fixup_consistency/Gemfile.lock Gemfile.lock
```

Рассмотрим, какие изменения были внесены в приложение:
 - Создан новый шаблон `database.yaml` для запуска MySQL в кластере. В реальных приложениях MySQL запустить несколько 
   сложнее, т.к. нужен persistent volume. Но в нашем случае для development-окружения это не требуется: БД будет терять все свои данные в случае перезапуска.

    {% snippetcut name=".helm/templates/database.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/019_fixup_consistency/.helm/templates/database.yaml" %}
    {% include_file "examples/rails/019_fixup_consistency/.helm/templates/database.yaml" %}
    {% endsnippetcut %}

 - Deployment приложения адаптирован так, чтобы миграции запускались в init-контейнере. Также настройки SQLite более не нужны:

    {% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/019_fixup_consistency/.helm/templates/deployment.yaml" %}
    {% include_file "examples/rails/019_fixup_consistency/.helm/templates/deployment.yaml" %}
    {% endsnippetcut %}

 - Приложение настроено на работу с MySQL:
    - `database.yaml`:

        {% snippetcut name=".helm/templates/database.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/019_fixup_consistency/.helm/templates/database.yaml" %}
        {% include_file "examples/rails/019_fixup_consistency/.helm/templates/database.yaml" %}
        {% endsnippetcut %}
    
    - `Gemfile`:

        {% snippetcut name="Gemfile" url="https://github.com/werf/werf-guides/blob/master/examples/rails/019_fixup_consistency/Gemfile" %}
        {% include_file "examples/rails/019_fixup_consistency/Gemfile" %}
        {% endsnippetcut %}

    - `Gemfile.lock`:

        {% snippetcut name="Gemfile.lock" url="https://github.com/werf/werf-guides/blob/master/examples/rails/019_fixup_consistency/Gemfile.lock" %}
        {% include_file "examples/rails/019_fixup_consistency/Gemfile.lock" %}
        {% endsnippetcut %}

Сделаем коммит изменений:
```shell
git add .
git commit -m go
```

Запустим деплой:
```shell
werf converge --repo <имя пользователя Docker Hub>/werf-guided-rails
```

Дождёмся выполнения команды. Заметьте, что в процессе работы в логах могут появляться ошибки подключения к БД, потому что контейнер с ней ещё не успел запуститься. Это нормально - необходимо дождаться полного запуска приложения.

Проверим создание и получение labels и увидим консистивное поведение:
```shell
curl -X POST "http://URL:PORT/api/labels/?label=name"
curl "http://URL:PORT/api/labels/" 
```

Так как у нас rails-приложение, мы когда-нибудь захотим зайти в работающий контейнер запущенного приложения и запустить там `rails console`. Это можно сделать вот так _(в общем случае это плохая практика, но допустимо на данном этапе для лучшей наглядности)_:
```shell
kubectl exec -ti basicapp-57789b68-c2xlq -- bash
rails c
```

Получим список labels напрямую из базы данных:
```shell
Label.all
```
