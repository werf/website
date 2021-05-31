---
title: Сборка образа
permalink: java_springboot/100_basic/10_build.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- Dockerfile
- werf.yaml
{% endfilesused %}
  
В этой главе мы начнём работу с werf: контейнизируем тестовое приложение, используя werf и [Dockerfile](https://docs.docker.com/engine/reference/builder/), и запустим конечный образ локально в Docker.

## Подготовка рабочего места

Мы предполагаем, что вы уже [установили werf]({{ site.docsurl }}/installation.html) и Docker.

Создайте директорию на своём компьютере и выполните там следующие шаги:

```shell
git clone https://github.com/werf/werf-guides.git
cp -r werf-guides/examples/springboot/000_app ./
cd 000_app 
git init
git add .
git commit -m "initial commit"
```

_Так вы скопируете себе код [приложения на SpringBoot](https://github.com/werf/werf-guides/tree/master/examples/springboot/000_app) и инициируете Git в каталоге с ним._

## Реализация сборки в Dockerfile

Конфигурация сборки нашего приложения будет состоять из следующих шагов:

- взять базовый образ с OpenJDK (`gradle:jdk8-openj9` подойдёт);
- добавить в него код приложения;
- собрать приложение с помощью gradle и переместить полученный jar в подходящее место.

Реализуем это в `Dockerfile`:

{% snippetcut name="Dockerfile" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/010_build/Dockerfile" %}
{% raw %}
```Dockerfile
FROM gradle:jdk8-openj9
WORKDIR /app
COPY . .

RUN gradle build --no-daemon
RUN cp /app/build/libs/*.jar /app/demo.jar

CMD ["java","-jar","/app/demo.jar"]
```
{% endraw %}
{% endsnippetcut %}

## Интеграция Dockerfile с werf

Подключим готовый Dockerfile к werf. Для этого, создадим в корне репозитория файл `werf.yaml`, описывающий сборку всего проекта:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-springboot
configVersion: 1
---
image: basicapp
dockerfile: Dockerfile
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="Что тут написано?" %}

`werf.yaml` начинается с обязательной **секции мета-информации**:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-springboot
configVersion: 1
```
{% endraw %}
{% endsnippetcut %}

Здесь:
- **_project_** — поле, задающее уникальное имя проекта приложения. Имя проекта по умолчанию используется при генерации имени Helm-релиза и пространства имен, `namespace`, в Kubernetes. Изменение имени у активного проекта затруднительно и требует ряда действий, которые необходимо выполнить вручную (подробнее о возможных последствиях можно прочитать [здесь]({{ site.docsurl }}/documentation/reference/werf_yaml.html#последствия-смены-имени-проекта));
- **_configVersion_** определяет версию синтаксиса, используемую в `werf.yaml` (на данный момент поддерживается только версия `1`).

Следующая секция конфигурации, которая и будет основной для сборки: [**image config section**]({{ site.docsurl }}/documentation/reference/werf_yaml.html#секция-image).

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
---
image: basicapp
dockerfile: Dockerfile
```
{% endraw %}
{% endsnippetcut %}

В строке `image: basicapp` определяется идентификатор образа, который впоследствии может использоваться при конфигурации выката, а также для вызова команд werf для определённого образа из `werf.yaml` (к примеру, `werf build basicapp`, `werf run basicapp` и т.д.). 

Строка `dockerfile: Dockerfile` указывает, что сборочная конфигурация описана в существующем файле, расположенном по пути `Dockerfile`. 

Также доступны и другие директивы, с которыми можно ознакомиться [по ссылке]({{ site.docsurl }}/documentation/reference/werf_yaml.html#dockerfile-image-section-image).

В одном `werf.yaml` может быть описано произвольное количество образов.

{% endofftopic %}

## Сборка

Перед выполнением сборки необходимо добавить изменения в git-репозиторий проекта:

```shell
git add werf.yaml Dockerfile
git commit -m "Add build configuration"
```

> Почему изменения должны добавляться в git-репозиторий, что такое гитерминизм и режим разработчика, а также другие особенности работы с файлами проекта будут разобраны далее в главе «Необходимо знать»

Сборка выполняется командой [`werf build`]({{ site.docsurl }}/documentation/reference/cli/werf_build.html):

{% raw %}
```shell
werf build
┌ ⛵ image basicapp
│ ┌ Building stage basicapp/dockerfile
│ │ basicapp/dockerfile  Sending build context to Docker daemon  116.7kB
│ │ basicapp/dockerfile  Step 1/15 : FROM gradle:jdk8-openj9
│ │ basicapp/dockerfile   ---> 2fb781988fa5
│ │ basicapp/dockerfile  Step 2/15 : WORKDIR /app
│ │ basicapp/dockerfile   ---> Using cache
<..>
│ │ basicapp/dockerfile  Successfully built e0d6df14df8b
│ │ basicapp/dockerfile  Successfully tagged ee51ea7f-c498-45a5-a435-0fd830fbb576:latest
│ ├ Info
│ │       name: werf-guided-springboot:50558f3f54d2ebbbd817824c6d7194aabe725bff6d7beae4df9c5e29-1606128099580
│ │       size: 738.6 MiB
│ └ Building stage basicapp/dockerfile (86.12 seconds)
└ ⛵ image basicapp (86.32 seconds)

Running time 86.37 seconds
```
{% endraw %}

## Запуск

Запуск контейнера выполняется командой [werf run]({{ site.docsurl }}/documentation/cli/main/run.html):

```shell
werf run --docker-options="--rm -p 8080:8080" -- java -jar /app/demo.jar
```

Обратите внимание, что [параметры docker](https://docs.docker.com/engine/reference/run/) задаются опцией `--docker-options`, а команда запуска после двух дефисов.

В `werf.yaml` может описываться произвольное количество образов. Для запуска контейнера определённого `image` из werf.yaml необходимо использовать позиционный аргумент команды (`werf run basicapp ...`).

_Можно заметить, что вызов `werf run` осуществляет сборку, т.е. предварительная сборка не требуется._

Теперь приложение доступно локально на порту 8080:

{% asset guides/springboot/100_10_app_in_browser.png %}

