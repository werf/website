---
title: Сборка образа
permalink: java_springboot/100_basic/10_build.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- Dockerfile
- werf.yaml
{% endfilesused %}

Классическим способом сборки является использование [Dockerfile](https://docs.docker.com/engine/reference/builder/). Возможно, в ваших приложениях уже реализована сборка с помощью этого механизма, поэтому мы начнём с него, а затем подключим его в werf. В следующих главах мы ускорим сборку, воспользовавшись альтернативным синтаксисом описания сборки, а сейчас — сфокусируемся на быстром получении результата.

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

werf следует принципам [гитерминизма]({{ site.docsurl }}/documentation/advanced/configuration/giterminism.html): опирается на состояние, описанное в Git-репозитории. Это означает, что некоммитнутые в Git-репозиторий файлы по умолчанию будут игнорироваться. Благодаря этому, имея исходные коды приложения, вы всегда можете реализовать его конкретное работоспособное состояние.

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

После того, как мы добавили описанные выше файлы `Dockerfile` и `werf.yaml`, надо обязательно закоммитить изменения в Git:

{% raw %}
```shell
git add .
git commit -m "work in progress"
```
{% endraw %}

Для того, чтобы запустить сборку, воспользуемся [командой `build`]({{ site.docsurl }}/documentation/reference/cli/werf_build.html):

{% raw %}
```shell
werf build
```
{% endraw %}

_В подглаве «Ускорение сборки» мы переведём сборку с Dockerfile на альтернативный синтаксис werf под названием Stapel и получим расширенные возможности: инкрементальную пересборку с учетом истории Git, возможность применения Ansible, использование кэша между сборками, удобные инструменты диагностики и многое другое._

Но уже сейчас вы можете заметить, что werf делает расширенный вывод логов сборки:

```
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

## Запуск

Запустим собранный образ с помощью [werf run]({{ site.docsurl }}/documentation/cli/main/run.html):

```shell
werf run --docker-options="--rm -p 8080:8080" -- java -jar /app/demo.jar
```

Обратите внимание, что мы задаем [параметры docker](https://docs.docker.com/engine/reference/run/) опцией `--docker-options`, а саму команду запуска указываем после двух дефисов.

_Вы также можете заметить, что и вызов `werf run` осуществляет сборку, т.е. предварительно запускать сборку на самом деле не обязательно._

Теперь приложение доступно локально на порту 8080:

![](/guides/images/springboot/100_10_app_in_browser.png)

## Внесение новых изменений

Мы будем постоянно дорабатывать приложение. Посмотрим, как это правильно делать на примере произвольных изменений в коде приложения:

{% snippetcut name="/src/main/java/com/example/demo/mvc/controller/LabelController.java" url="#" %}
{% raw %}
```
    @GetMapping("/labels")
    public String labels() {
        return "Our changes";
    }
```
{% endraw %}
{% endsnippetcut %}

 1. Остановите ранее запущенный `werf run` (нажав Ctrl+C в консоли, где он запущен).
 2. Запустите его заново: 
    ```shell
    werf run --docker-options="--rm -p 8080:8080" -- java -jar /app/demo.jar
    ```
 3. Произошла ошибка:
    ```
    Error: phase build on image basicapp stage dockerfile handler failed: the file "src/main/java/com/example/demo/mvc/controller/LabelController.java" must be committed

    You might also be interested in developer mode (activated with --dev option) that allows you to work with staged changes without doing redundant commits. Just use "git add <file>..." to include the changes that should be used.

    To provide a strong guarantee of reproducibility, werf reads the configuration and build's context files from the project git repository and eliminates external dependencies. We strongly recommend to follow this approach. But if necessary, you can allow the reading of specific files directly from the file system and enable the features that require careful use. Read more about giterminism and how to manage it here: https://werf.io/v1.2-ea/documentation/advanced/configuration/giterminism.html.
    ```

В описанном сценарии **перед шагом 1 забыли сделать коммит** в Git.

{% offtopic title="А как правильно и зачем такие сложности?" %}
1. Внести изменения в код.
2. Сделать коммит:
   ```shell
   git add .
   git commit -m "wip"
   ```
3. Перезапустить `werf run`:
    ```shell
    werf run --docker-options="--rm -p 8080:8080" -- java -jar /app/demo.jar
    ```
4. Посмотреть на результат в браузере.

Как мы уже упоминали в начале статьи, werf работает в режиме [гитерминизма]({{ site.docsurl }}/documentation/advanced/configuration/giterminism.html). Жёсткая связка с Git необходима для того, чтобы гарантировать воспроизводимость вашего решения. Подробнее о том, как работает эта механика _гитерминизма_, а также о режиме разработчика с флагом `--dev` мы расскажем в главе «Необходимо знать», а пока что — сфокусируемся на сборке и доставке до кластера.

{% endofftopic %}

<div id="go-forth-button">
    <go-forth url="20_cluster.html" label="Подготовка кластера" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
