---
title: Организация локальной разработки
permalink: java_springboot/200_real_apps/10_local.html
layout: development
---

От локальной разработки хочется видеть:

- Быструю сборку/компиляцию изменений
- Инструменты для отладки ошибок
- Видеть результат в максимально приближенном к финальному виде

## Быстрая сборка изменений

При локальной разработке приходится: а) сделать коммит в git; б) запускать `werf converge ...`; в) дожидаться пересборки и передеплоя. Для ускорения этого процесса команды werf (такие, как `converge`, `build`, `run`, `compose`) поддерживают атрибут `--follow`. Если вы, например, запустите

```shell
werf converge --repo registry.mydomain.io/werf-guided-project --follow
```

то при коммите в git пересборка и деплой будет осуществляться автоматически.

## Отладка ошибок

При написании кода инфраструктуры иногда случаются ошибки и опечатки. Для их отладки в werf есть [инструмент интроспекции]({{ site.url }}/documentation/advanced/development_and_debug/stage_introspection.html).

{% offtopic title="А как пользоваться интроспекцией?" %}

Допустим, у вас есть `werf.yaml` с ошибкой: мы пытаемся в образе с `alpine` вызвать пакетный менеджер `apt`, которого там нет:

{% snippetcut name="werf.yaml" url="#" %}
{% raw %}
```yaml
project: some-project
configVersion: 1
---
image: mycontainer
from: alpine:latest
shell:
  beforeInstall:
  - apt update
  - apt install -y curl
  - curl ya.ru -v
```
{% endraw %}
{% endsnippetcut %}

Команда `werf build` в такой ситуации ожидаемо выдаст ошибку:

```
│ │ ┌ Building stage beforeInstall
│ │ │ mycontainer/beforeInstall  /.werf/shell/script.sh: line 3: apt: command not found
│ │ ├ Info
│ │ └ Building stage beforeInstall (1.62 seconds) FAILED
│ └ ⛵ image mycontainer (3.23 seconds) FAILED
└ Building stages (3.23 seconds) FAILED

Running time 7.01 seconds
Error: failed to build werf-stages-storage/some-project:96ae05afe886b927b8c131981591f38dfb834344e608fd5d66cd210f96867ecb: container run failed: Status: , Code: 127
```

Воспользуемся интроспекцией:

```shell
werf build --stages-storage :local --introspect-before-error
```

И вместо сообщения об ошибке вы окажетесь в консоли внутри контейнера:

```
│ │ ┌ Building stage beforeInstall
│ │ │ mycontainer/beforeInstall  /.werf/shell/script.sh: line 3: apt: command not found
│ │ │ mycontainer/beforeInstall  Launched command: /.werf/shell/script.sh
bash-4.3#
```

Например, если вы введёте в этот момент `ls -la` то увидите не хост-машину, а контейнер с его текущим состоянием:

```shell
bash-4.3# ls -la
total 68
drwxr-xr-x    1 root     root          4096 Nov 24 21:22 .
drwxr-xr-x    1 root     root          4096 Nov 24 21:22 ..
-rwxr-xr-x    1 root     root             0 Nov 24 21:22 .dockerenv
drwxr-xr-x    1 root     root          4096 Nov 24 21:22 .werf
drwxr-xr-x    2 root     root          4096 Oct 21 09:23 bin
drwxr-xr-x    5 root     root           360 Nov 24 21:22 dev
drwxr-xr-x    1 root     root          4096 Nov 24 21:22 etc
drwxr-xr-x    2 root     root          4096 Oct 21 09:23 home
drwxr-xr-x    7 root     root          4096 Oct 21 09:23 lib
drwxr-xr-x    5 root     root          4096 Oct 21 09:23 media
drwxr-xr-x    2 root     root          4096 Oct 21 09:23 mnt
drwxr-xr-x    2 root     root          4096 Oct 21 09:23 opt
dr-xr-xr-x  205 root     root             0 Nov 24 21:22 proc
drwx------    2 root     root          4096 Oct 21 09:23 root
drwxr-xr-x    2 root     root          4096 Oct 21 09:23 run
drwxr-xr-x    2 root     root          4096 Oct 21 09:23 sbin
drwxr-xr-x    2 root     root          4096 Oct 21 09:23 srv
dr-xr-xr-x   13 root     root             0 Nov 24 21:22 sys
drwxrwxrwt    2 root     root          4096 Oct 21 09:23 tmp
drwxr-xr-x    7 root     root          4096 Oct 21 09:23 usr
drwxr-xr-x   12 root     root          4096 Oct 21 09:23 var
```

Находясь внутри контейнера вы можете разобраться, что же пошло не так и провести необходимые эксперименты. В нашем случае мы нагуглим, что установка пакетов в `alpine` делается с помощью `apk`, установим его вручную. Убедимся, что это работает корректно вызвав `curl` в том виде, как это прописано в `werf.yaml`. После чего выйдем из интроспекции с помощью `exit`.

```shell
bash-4.3# apk --no-cache add curl
fetch http://dl-cdn.alpinelinux.org/alpine/v3.12/main/x86_64/APKINDEX.tar.gz
fetch http://dl-cdn.alpinelinux.org/alpine/v3.12/community/x86_64/APKINDEX.tar.gz
(1/4) Installing ca-certificates (20191127-r4)
(2/4) Installing nghttp2-libs (1.41.0-r0)
(3/4) Installing libcurl (7.69.1-r1)
(4/4) Installing curl (7.69.1-r1)
Executing busybox-1.31.1-r19.trigger
Executing ca-certificates-20191127-r4.trigger
OK: 7 MiB in 18 packages
bash-4.3#
bash-4.3#
bash-4.3# curl ya.ru -v
*   Trying 87.250.250.242:80...
* Connected to ya.ru (87.250.250.242) port 80 (#0)
> GET / HTTP/1.1
> Host: ya.ru
> User-Agent: curl/7.69.1
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 302 Found
< Cache-Control: no-cache,no-store,max-age=0,must-revalidate
< Content-Length: 0
< Date: Tue, 24 Nov 2020 21:30:33 GMT
< Expires: Tue, 24 Nov 2020 21:30:33 GMT
< Last-Modified: Tue, 24 Nov 2020 21:30:33 GMT
< Location: https://ya.ru/
< P3P: policyref="/w3c/p3p.xml", CP="NON DSP ADM DEV PSD IVDo OUR IND STP PHY PRE NAV UNI"
< Set-Cookie: yandexuid=5379887481606253433; Path=/; Domain=.ya.ru; Expires=Thu, 24 Nov 2022 21:30:33 GMT
< Set-Cookie: is_gdpr=1; Path=/; Domain=.ya.ru; Expires=Thu, 24 Nov 2022 21:30:33 GMT
< Set-Cookie: is_gdpr_b=CMmFQhC4DhgB; Path=/; Domain=.ya.ru; Expires=Thu, 24 Nov 2022 21:30:33 GMT
< X-Content-Type-Options: nosniff
<
* Connection #0 to host ya.ru left intact
bash-4.3#
bash-4.3# exit
```

С полученным пониманием правильного синтаксиса, который реально сработает в контейнере, исправим `werf.yaml`:

{% snippetcut name="werf.yaml" url="#" %}
{% raw %}
```yaml
project: some-project
configVersion: 1
---
image: mycontainer
from: alpine:latest
shell:
  beforeInstall:
  - apk --no-cache add curl
  - curl ya.ru -v
```
{% endraw %}
{% endsnippetcut %}

И перезапустим сборку, которая в этот раз успешно отработает без интроспекции:

```shell
werf build --stages-storage :local --introspect-before-error
```

Вы успешно отладили и исправили ошибку в сборочном сценарии!

{% endofftopic %}

## Видеть результат

Если registry и стенд для разворачивания результата находится на локальном компьютере — мы экономим время на перекачку собранных образов по сети. В качестве площадки для запуска контейнеров можно использовать:

- minikube — более требователен к ресурсам, хотя и ближе к "настоящему" кластеру
- docker compose — требует меньше ресурсов и меньшей квалификации

Настройку minikube мы рассмотрели в главе ["Подготовка кластера"](../100_basic/20_cluster.html). Однако, werf позволяет напрямую работать с docker compose, дополняя его возможности.

Если вы переводите существующее приложение на werf и в рамках главы "Ускорение сборки" уже перешли целиком на `werf.yaml` - возможно вы захотите задействовать готовый `docker-compose.yml`.

Допустим, ваш `docker-compose.yml` запускает два образа. Один — ранее собиравшийся через `Dockerfile`, второй — с бд Redis:

{% snippetcut name="docker-compose.yml" url="#" %}
{% raw %}
```yaml
version: "3.8"
services:
  basicapp:
    build: .
    ports:
      - "5000:5000"
    depends_on:
      - redis
  redis:
    image: "redis:alpine"
```
{% endraw %}
{% endsnippetcut %}

Скорректируйте описание того, откуда берётся образ для `basicapp`, подставив туда переменную окружения, которую сгенерирует werf:

{% snippetcut name="docker-compose.yml" url="#" %}
{% raw %}
```yaml
version: "3.8"
services:
  basicapp:
    image: "$WERF_IMAGE_BASICAPP"
    ports:
      - "5000:5000"
    depends_on:
      - redis
  redis:
    image: "redis:alpine"
```
{% endraw %}
{% endsnippetcut %}

Для каждого `image` в `werf.yaml` будет доступно по одной переменной вида `$WERF_IMAGE_{НАЗВАНИЕ ОБРАЗА КАПСОМ}`.

Закоммитьте изменения в git и запустите compose:

```shell
werf compose up --follow
```

Образы запустятся и будут пересобираться и перезапускаться при каждом коммите.

