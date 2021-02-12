---
title: Сборка образа
permalink: rails/100_basic/10_build.html
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
cp -r werf-guides/examples/rails/000_app ./
cd 000_app
git init
git add .
git commit -m "initial commit"
```

_Так вы скопируете себе код [приложения на Ruby on Rails](https://github.com/werf/werf-guides/tree/master/examples/rails/000_app) и инициируете Git в каталоге с ним._

werf следует принципам [гитерминизма]({{ site.docsurl }}/documentation/advanced/configuration/giterminism.html): опирается на состояние, описанное в Git-репозитории. Это означает, что файлы, не коммитнутые в Git-репозиторий, по умолчанию будут игнорироваться. Благодаря этому, имея исходные коды приложения, вы всегда можете реализовать его конкретное работоспособное состояние.

## Реализация сборки в Dockerfile

Конфигурация сборки нашего приложения будет состоять из следующих шагов:

- взять базовый образ с Ruby (`ruby:2.7.1` подойдёт);
- добавить в него код приложения;
- установить зависимости и gem'ы.

Реализуем это в `Dockerfile`:

{% snippetcut name="Dockerfile" url="https://github.com/werf/werf-guides/blob/master/examples/rails/010_build/Dockerfile" %}
{% raw %}
```Dockerfile
FROM ruby:2.7.1
WORKDIR /app
COPY . .

RUN apt-get update -qq && apt-get install -y \
 build-essential libpq-dev libxml2-dev libxslt1-dev curl
RUN bundle install

CMD ["bundle", "exec", "puma"]
```
{% endraw %}
{% endsnippetcut %}

## Интеграция Dockerfile с werf

Подключим готовый Dockerfile к werf. Для этого, создадим в корне репозитория файл `werf.yaml`, описывающий сборку всего проекта:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-rails
configVersion: 1
---
image: basicapp
dockerfile: Dockerfile
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="Что тут написано?" %}

`werf.yaml` начинается с обязательной **секции мета-информации**:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-rails
configVersion: 1
```
{% endraw %}
{% endsnippetcut %}

Здесь:
- **_project_** — поле, задающее уникальное имя проекта приложения. Имя проекта по умолчанию используется при генерации имени Helm-релиза и пространства имен, `namespace`, в Kubernetes. Изменение имени у активного проекта затруднительно и требует ряда действий, которые необходимо выполнить вручную (подробнее о возможных последствиях можно прочитать [здесь]({{ site.docsurl }}/documentation/reference/werf_yaml.html#последствия-смены-имени-проекта));
- **_configVersion_** определяет версию синтаксиса, используемую в `werf.yaml` (на данный момент поддерживается только версия `1`).

Следующая секция конфигурации, которая и будет основной для сборки: [**image config section**]({{ site.docsurl }}/documentation/reference/werf_yaml.html#секция-image).

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/rails/011_build_werf/werf.yaml" %}
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
│ │ basicapp/dockerfile  Sending build context to Docker daemon  11.64MB
│ │ basicapp/dockerfile  Step 1/14 : FROM ruby:2.7.1
│ │ basicapp/dockerfile   ---> d8ca85855516
│ │ basicapp/dockerfile  Step 2/14 : WORKDIR /app
│ │ basicapp/dockerfile   ---> Using cache
<..>
│ │ basicapp/dockerfile  Successfully built 3a4ede4e9556
│ │ basicapp/dockerfile  Successfully tagged 0041b344-efe4-416d-baff-5e50fbb712b0:latest
│ ├ Info
│ │       name: werf-guided-rails:31e0e7436c3055fa816fc770ebda185bacb7e8ef53775b8e5488a83f-1611855308907
│ │       size: 929.2 MiB
│ └ Building stage basicapp/dockerfile (94.47 seconds)
└ ⛵ image basicapp (96.07 seconds)

Running time 96.38 seconds
```

## Запуск

Запустим собранный образ с помощью [werf run]({{ site.docsurl }}/documentation/cli/main/run.html):

```shell
werf run --docker-options="--rm -p 3000:3000" -- bundle exec puma
```

Обратите внимание, что мы задаем [параметры docker](https://docs.docker.com/engine/reference/run/) опцией `--docker-options`, а саму команду запуска указываем после двух дефисов.

_Вы также можете заметить, что и вызов `werf run` осуществляет сборку, т.е. предварительно запускать сборку на самом деле не обязательно._

Теперь приложение доступно локально на порту 3000:

![](/guides/images/rails/100_10_app_in_browser.png)

_Как уже было сказано, используется база SQLite без сохранения данных между запусками. Поэтому при первом запросе Rails выдаст страницу с просьбой выполнить миграции. Это можно сделать, нажав на этой странице соответствующую кнопку._

## Внесение новых изменений

Мы будем постоянно дорабатывать приложение. Посмотрим, как это правильно делать на примере произвольных изменений в коде приложения:

{% snippetcut name="app/controllers/api/labels_controller.rb" url="#" %}
{% raw %}
```ruby
def index
  render plain: 'Our changes'
end
```
{% endraw %}
{% endsnippetcut %}

 1. Остановите ранее запущенный `werf run` (нажав Ctrl+C в консоли, где он запущен).
 2. Запустите его заново: 
    ```shell
    werf run --docker-options="--rm -p 3000:3000" -- bundle exec puma
    ```
 3. Произошла ошибка:
    ```
    Error: phase build on image basicapp stage dockerfile handler failed: the file "app/controllers/api/labels_controller.rb" must be committed

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
    werf run --docker-options="--rm -p 3000:3000" -- bundle exec puma
    ```
4. Посмотреть на результат в браузере.

Как мы уже упоминали в начале статьи, werf работает в режиме [гитерминизма]({{ site.docsurl }}/documentation/advanced/configuration/giterminism.html). Жёсткая связка с Git необходима для того, чтобы гарантировать воспроизводимость вашего решения. Подробнее о том, как работает эта механика _гитерминизма_, а также о режиме разработчика с флагом `--dev` мы расскажем в главе «Необходимо знать», а пока что — сфокусируемся на сборке и доставке до кластера.

{% endofftopic %}

<div id="go-forth-button">
    <go-forth url="20_cluster.html" label="Подготовка кластера" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
