---
title: Сборка
sidebar: applications_guide
guide_code: ____________
permalink: ____________/020_basic/10_build.html
toc: false
---

{% filesused title="Файлы, упомянутые в главе" %}
- werf.yaml
{% endfilesused %}

В этой главе мы научимся писать конфигурацию сборки, отлаживать её в случае ошибок и загружать полученный образ в Registry.

{% offtopic title="Локальная разработка или на сервере?" %}
Существует 2 варианта ведения разработки: локально на компьютере либо на сервере. У обоих вариантов есть свои плюсы и минусы.

Локальная разработка подойдет, если у вас достаточно мощный компьютер и важна скорость разработки. В остальных случаях, как правило, лучше выбирать разработку на сервере.

Подробнее о плюсах и минусах каждого подхода:

**Локальная разработка**

Основной плюс этого метода — скорость разработки и проверки изменений. Результат правок можно увидеть уже через несколько секунд (хотя это зависит от используемого языка программирования).

К минусам можно отнести необходимость в более мощном железе, чтобы запускать приложение, а также то, что локальное окружение будет значительно отличаться от production.

_В скором времени werf [позволит](https://github.com/werf/werf/issues/1940) разработчикам запускать локальное окружение, идентичное production._

**Разработка на сервере**

Главный плюс этого метода — однотипность production с остальными окружениями. При возникновении проблем со сборкой мы узнаем об этом моментально. Отбрасывается необходимость в дополнительных ресурсах для локального запуска.

Один из минусов — это отзывчивость: на процесс от push'а кода до появления результата может потребоваться несколько минут.
{% endofftopic %}

Возьмите исходный код приложения [из репозитория на GitHub](https://github.com/werf/werf-guides/tree/master/examples/____________/000-app):

```bash
git clone git@github.com:werf/werf-guides.git guide
cd examples/____________/000-app
```

… и скопируйте его в свой проект в GitLab. Далее мы будем работать с исходным кодом проекта в GitLab.

____________
____________

Для того, чтобы werf смогла собрать Docker-образ с приложением, необходимо в корне нашего репозитория создать файл `werf.yaml`, в котором будут описаны инструкции по сборке.

{% offtopic title="Варианты синтаксиса werf.yaml" %}

werf имеет два механизма для сборки: Stapel (поддерживающий синтаксис Ansible и Shell) и Dockerfile (использующий синтаксис Dockerfile).

**Ansible**

Поскольку werf поддерживает почти все модули из Ansible, если у вас имеются Ansible-плейбуки для сборки, их можно легко адаптировать под werf.

Пример установки curl:
{% raw %}
```yaml
- name: "Install additional packages"
  apt:
    state: present
    update_cache: yes
    pkg:
    - curl
```
{% endraw %}
Полный список поддерживаемых модулей Ansible в werf доступен в [документации]({{ site.docsurl }}/documentation/configuration/stapel_image/assembly_instructions.html#supported-modules).

**Shell**

Также werf поддерживает использование обычных shell-команд — как будто мы запускаем Bash-скрипт.

Пример установки curl:
{% raw %}
```yaml
shell:
  beforeInstall:
  - apt-get update
  - apt-get install -y curl
```
{% endraw %}
Для простоты и удобства в гайде рассматривается именно shell, однако для своего проекта можно выбрать любой из вариантов.

Прочитать подробнее про виды синтаксисов в werf можно в документации: [**syntax section**]({{ site.docsurl }}/documentation/guides/advanced_build/first_application.html#%D1%88%D0%B0%D0%B3-1-%D0%BA%D0%BE%D0%BD%D1%84%D0%B8%D0%B3%D1%83%D1%80%D0%B0%D1%86%D0%B8%D1%8F-werfyaml).

{% endofftopic %}

Начнём создание `werf.yaml` с обязательной [**секции мета-информации**]({{ site.docsurl }}/documentation/configuration/introduction.html#секция-мета-информации):

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-project
configVersion: 1
```
{% endraw %}
{% endsnippetcut %}

Здесь:
- **_project_** — поле, задающее имя для проекта. Под проектом понимается набор из Docker-образов, которые мы будем собирать. Имя проекта по умолчанию используется в названии Helm-релиза и пространства имен в Kubernetes, куда будет выкатываться приложение. Данное имя не рекомендуется изменять (или подходить к таким изменениям с должным уровнем ответственности), так как после изменений уже имеющиеся ресурсы, которые выкачены в кластер, не будут автоматически переименованы.
- **_configVersion_** — в данном случае определяет версию синтаксиса, используемую в `werf.yaml` (на данный момент мы всегда используем `1`).

Перейдем к следующей секции конфигурации, которая и будет основной для сборки: [**image config section**]({{ site.docsurl }}/documentation/configuration/introduction.html#%D1%81%D0%B5%D0%BA%D1%86%D0%B8%D1%8F-%D0%BE%D0%B1%D1%80%D0%B0%D0%B7%D0%B0).

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/werf.yaml" %}
{% raw %}
```yaml
---
image: basicapp
from: ____________
```
{% endraw %}
{% endsnippetcut %}

В строке `image: basicapp` дано название для образа, который соберёт werf. Это имя будет впоследствии указываться при запуске контейнера. Строка `from: ____________` определяет, что берётся за основу: в данном случае это официальный публичный образ с нужной нам версией ____________.

{% offtopic title="Что делать, если образов и других констант станет много?" %}

В werf поддерживаются [**Go templates**]({{ site.docsurl }}/documentation/configuration/introduction.html#%D1%88%D0%B0%D0%B1%D0%BB%D0%BE%D0%BD%D1%8B-go), поэтому легко определять переменные и записывать в них константы и часто используемые образы.

Например, сделаем 2 образа, используя один базовый образ `golang:1.11-alpine`:

{% raw %}
```yaml
{{ $base_image := "golang:1.11-alpine" }}

project: my-project
configVersion: 1
---

image: gogomonia
from: {{ $base_image }}
---
image: saml-authenticator
from: {{ $base_image }}
```
{% endraw %}

Подробнее почитать про Go-шаблоны в werf можно в документации: [**werf go templates**]({{ site.docsurl }}/documentation/configuration/introduction.html#%D1%88%D0%B0%D0%B1%D0%BB%D0%BE%D0%BD%D1%8B-go).

{% endofftopic %}

Добавим исходный код нашего приложения в контейнер с помощью [**директивы git**]({{ site.docsurl }}/documentation/configuration/stapel_image/git_directive.html):

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/werf.yaml" %}
{% raw %}
```yaml
git:
- add: /
  to: /app
```
{% endraw %}
{% endsnippetcut %}

При наличии директивы `git` werf использует репозиторий проекта, в котором расположен `werf.yaml`. Добавляемые директории и файлы проекта указываются списком относительно корня репозитория.

* `add: /` — директория, которая будет добавлена внутрь Docker-образа. Мы указываем, что это корень, т.е. весь репозиторий будет склонирован внутрь Docker-образа;
* `to: /app` — куда клонируется репозиторий внутри Docker-образа. Важно заметить, что директорию назначения werf создаст сам.

Как такового репозитория в образе не создаётся: директория `.git` отсутствует. На одном из первых этапов сборки все исходники пакуются в архив (с учётом всех пользовательских параметров) и распаковываются внутри контейнера. При последующих сборках [накладываются патчи с разницей]({{ site.docsurl }}/documentation/configuration/stapel_image/git_directive.html#подробнее-про-gitarchive-gitcache-gitlatestpatch).

{% offtopic title="Можно ли использовать несколько репозиториев?" %}

Директива `git` также позволяет добавлять код из удалённых Git-репозиториев, используя параметр `url`.

Пример:
{% raw %}
```yaml
git:
- add: /src
  to: /app
- url: https://github.com/ariya/phantomjs
  add: /
  to: /src/phantomjs
```
{% endraw %}
Детали и особенности можно почитать в [документации]({{ site.docsurl }}/documentation/configuration/stapel_image/git_directive.html#работа-с-удаленными-репозиториями).

{% endofftopic %}

{% offtopic title="Реальная практика" %}

В реальной практике нужно добавлять файлы, **фильтруя по названию или пути**.

К примеру, в данном варианте добавляются все файлы `.php` и `.js` из каталога `/src`, исключая файлы с суффиксом `-dev.` и `-test.` в имени файла.
{% raw %}
```yaml
git:
- add: /src
  to: /app
  includePaths:
  - '**/*.php'
  - '**/*.js'
  excludePaths:
  - '**/*-dev.*'
  - '**/*-test.*'
```
{% endraw %}
В Linux-системах важно не забывать про установку **владельца файлов**.

Например, вот так изменяется владелец файла `index.php` на `www-data`:

```yaml
git:
- add: /src/index.php
  to: /app/index.php
  owner: www-data
```

Подробнее о всех опциях директивы `git` можно прочитать в  [документации]({{ site.docsurl }}/documentation/configuration/stapel_image/git_directive.html#Изменение-владельца).

{% endofftopic %}

Следующим этапом необходимо описать **правила [сборки для приложения]({{ site.docsurl }}/documentation/configuration/stapel_image/assembly_instructions.html)**. werf позволяет кэшировать сборку образа подобно слоям в Docker, но с более явной конфигурацией. Этот механизм называется [стадиями]({{ site.docsurl }}/documentation/configuration/stapel_image/assembly_instructions.html#пользовательские-стадии). Для текущего приложения будут описаны 2 стадии, в которых сначала устанавливаются пакеты, а потом ­— производятся действия над исходными кодами приложения.

{% offtopic title="Что нужно знать про стадии" %}

Стадии — это очень важный и в то же время непростой инструмент, который сокращает время сборки приложения. Вплотную работа со стадиями будет разобрана при рассмотрении вопроса зависимостей и ассетов.

Пока же достаточно знать следующее:

* При правильном использовании стадий мы сильно уменьшаем время сборки приложения.
* На стадии `install` устанавливаются системные пакеты, необходимое для работы программы.
* На стадии `setup` осуществляется работа с исходным кодом приложения.

{% endofftopic %}

Добавим в `werf.yaml` следующий блок, используя shell-синтаксис:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/werf.yaml" %}
{% raw %}
```yaml
shell:
  ____________
  ____________
```
{% endraw %}
{% endsnippetcut %}

Чтобы при запуске приложения по умолчанию использовалась директория `/app`, воспользуемся **[указанием Docker-инструкций]({{ site.docsurl }}/documentation/configuration/stapel_image/docker_directive.html)**:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/____________/____________/werf.yaml" %}
{% raw %}
```yaml
docker:
  WORKDIR: /app
```
{% endraw %}
{% endsnippetcut %}

____________
____________

Когда `werf.yaml` готов (или кажется таковым) — пробуем запустить сборку:

```bash
$  werf build --stages-storage :local
```

{% offtopic title="Что делать, если что-то пойдёт не так?" %}

werf предоставляет удобные способы отладки.

Если во время сборки что-то пошло не так, можно мгновенно оказаться внутри процесса сборки и изучить, что конкретно пошло не так. Для этого есть механизм [интроспекции стадий]({{ site.docsurl }}/documentation/reference/development_and_debug/stage_introspection.html). Если сборка запущена с флагом `--introspect-before-error`, вы окажетесь в сборочном контейнере перед тем, как сработала ошибка.

Например, ошибка произошла на команде `jekyll build`. Если вы воспользуетесь интроспекцией, то окажетесь внутри контейнера прямо перед тем, как должна была выполняться эта команда. А значит, выполнив её самостоятельно, можно увидеть сообщение об ошибке и посмотреть на состояние других файлов в контейнере на этот момент, чтобы понять, что в сценарии сборки пошло не так.

{% endofftopic %}

{% offtopic title="Что за stages-storage?" %}
Стадии, которые используются для ускорения сборки, должны где-то храниться. werf подразумевает, что сборка может производиться на нескольких раннерах, поэтому кэш сборочных стадий нужно хранить в каком-то едином хранилище.

`stages-storage` [позволяет настраивать]({{ site.docsurl }}/documentation/reference/stages_and_images.html#%D1%85%D1%80%D0%B0%D0%BD%D0%B8%D0%BB%D0%B8%D1%89%D0%B5-%D1%81%D1%82%D0%B0%D0%B4%D0%B8%D0%B9), где будет храниться кэш сборочных стадий: на локальном сервере или в Registry.
{% endofftopic %}

Если всё написано правильно, то сборка завершится примерно так:

```
...
│ ┌ Building stage basicapp/beforeInstall
│ ├ Info
│ │     repository: werf-stages-storage/werf-guided-project
│ │       image_id: 2743bc56bbf7
│ │        created: 2020-05-26T22:44:26.0159982Z
│ │            tag: 7e691385166fc7283f859e35d0c9b9f1f6dc2ea7a61cb94e96f8a08c-1590533066068
│ │           diff: 0 B
│ │   instructions: WORKDIR /app
│ └ Building stage basicapp/beforeInstall (0.82 seconds)
└ ⛵ image basicapp (239.56 seconds)
```

В конце werf выдает информацию о готовом репозитории и о теге для этого образа для дальнейшего использования.

```
werf-stages-storage/werf-guided-project:981ece3acc63d57d5ab07f45fd0c0c477088649523822c2bff033df4-1594911680806
```

Запустим собранный образ с помощью [werf run]({{ site.docsurl }}/documentation/cli/main/run.html):

```bash
$ werf run --stages-storage :local --docker-options="-d -p ____________:____________ --restart=always" -- ____________
```

Первая часть команды очень похожа на `build`, а во второй — мы задаем [параметры docker](https://docs.docker.com/engine/reference/run/) и после двойного дефиса команду, с которой запустить образ.

Теперь приложение доступно локально на порту ____________:

![](/applications_guide_ru/images/applications-guide/020-hello-world-in-browser.png)

Как только мы убедились в том, что всё корректно, необходимо **загрузить образ в Registry**. Сборка с последующей загрузкой в Registry делается [командой `build-and-publish`]({{ site.docsurl }}/documentation/cli/main/build_and_publish.html). Когда werf запускается внутри CI-процесса, werf узнаёт реквизиты для доступа к Registry [из переменных окружения](https://docs.gitlab.com/ee/ci/variables/predefined_variables.html).

Если же werf запускается вне Gitlab CI, потребуется:

* Вручную подключиться к GitLab Registry [с помощью `docker login`](https://docs.docker.com/engine/reference/commandline/login/);
* Установить переменную окружения `WERF_IMAGES_REPO` с путём до Registry (вида `registry.example.com/myproject`);
* Выполнить сборку и загрузку в Registry: `werf build-and-publish`.

Если вы всё правильно сделали, собранный образ появится в Registry. В случае использования Registry от GitLab собранный образ можно увидеть через веб-интерфейс GitLab.

____________
____________
____________

<div id="go-forth-button">
    <go-forth url="20_iac.html" label="Конфигурирование инфраструктуры в виде кода" base-url="applications_guide_ru"></go-forth>
</div>
