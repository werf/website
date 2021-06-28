---
title: Логирование
permalink: rails/200_real_apps/25_logging.html
examples: examples/rails/100_logging
layout: wip
---

В этой статье мы добавим в наше приложение новый метод API и рассмотрим тему логирования в приложении и просмотра логов запущенного приложения, и при этом  сконфигурируем приложение так, чтобы писать логи, которые подхватит Kubernetes.

## Подготовка окружения
Перейдём к состоянию приложения, в котором добавлен новый метод API:

```shell
cd werf-guides/examples/rails/020_logging/
git init --separate-git-dir ~/werf-guides-repo
```

## Настройка logs и окружения

По умолчанию rails-приложение предоставляет несколько окружений: development, test и production.
В окружении production при установке переменной окружения `RAILS_LOG_TO_STDOUT=1` стандартный логер по умолчанию будет автоматически писать логи в stdout, поэтому для всех окружений он будет переделан чтобы писать еще и в stderr — это наиболее правильный способ логирования приложений, запускаемых в Kubernetes-кластере.

В нашем случае дальнейший сбор и хранение логов, при необходимости, будет осуществляться отдельными решениями и это основной принципиальный момент - мы об этом не думаем и просто пишем в stdout/stderr.

Конфигурацию логера приложения для всех окружений сделана в файле {% include snippetcut_example path="config/application.rb" syntax="ruby" examples=page.examples %}.

а стандартная конфигурация логера для production-окружения удалена из файла {% include snippetcut_example path="config/environments/production.rb" syntax="ruby" examples=page.examples %}.


## Логирование в приложении

В наше приложение был добавлен новый метод POST `/api/generate-image`, использующий стандартный rails-логер через конструкцию `logger.debug`. На вход метода передаётся параметр `text` а на выходе генерируется png-картинка с переданным `text`.

Полный листинг метода можно посмотреть в файле {% include snippetcut_example path="app/controllers/api/images_controller.rb" syntax="ruby" examples=page.examples %}

## Деплоим приложение

Мы готовы задеплоить новую версию приложения. Сделаем коммит изменения:

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
curl -v -X POST "example.com/api/generate-image?text=do%20it" -o /tmp/out
```

- В файле /tmp/out должно появится изображение с указанным тестом.
- Если вдруг не получилось — не беда, сейчас будем смотреть логи приложения, также в случае ошибки в файле /tmp/out будет её текст.

## Как посмотреть логи

За API нашего приложения отвечает backend, который запущен в Deployment. Задача deployment — создать указанное количество реплик pod'ов, при этом каждый pod это запущенный экземпляр нашего приложения, соответственно, чтобы увидеть логи backend, необходимо запросить их для одного из pod.

Получим список pod'ов:

```shell
kubectl -n werf-guided-rails get pod
```

В ответ отобразится следующее:
```shell
NAME                      READY   STATUS    RESTARTS   AGE
backend-7964b6b68-n28v5   1/1     Running   0          21m
backend-7964b6b68-wzgw4   1/1     Running   0          21m
mysql-666f76d7cb-967xb    1/1     Running   3          4d6h
```

Выберем один из pod с префиксом `backend` и посмотрим его логи следующей командой:

```shell
kubectl -n werf-guided-rails logs -f backend-7964b6b68-n28v5
```

В ответе увидим:
```shell
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
Ура! У нас всё получилось.