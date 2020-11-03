---
title: Сборка
permalink: gitlab_nodejs/020_basic/10_build.html
---

{% offtopic title="Что делать, если образов и других констант станет много?" %}


{% endofftopic %}








Добавим исходный код нашего приложения в контейнер с помощью [**директивы git**]({{ site.docsurl }}/documentation/configuration/stapel_image/git_directive.html):

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/020-basic/werf.yaml" %}
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


{% endofftopic %}

{% offtopic title="Реальная практика" %}



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

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/020-basic/werf.yaml" %}
{% raw %}
```yaml
shell:
  beforeInstall:
  - apt update
  - apt install -y tzdata locales
  install:
  - cd /app && npm ci
```
{% endraw %}
{% endsnippetcut %}

Чтобы при запуске приложения по умолчанию использовалась директория `/app`, воспользуемся **[указанием Docker-инструкций]({{ site.docsurl }}/documentation/configuration/stapel_image/docker_directive.html)**:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/020-basic/werf.yaml" %}
{% raw %}
```yaml
docker:
  WORKDIR: /app
```
{% endraw %}
{% endsnippetcut %}

Также мы должны прописать связь файла `package.json` со стадией `install` внутри блока `git`:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/020-basic/werf.yaml" %}
```yaml
git:
  <...>
  stageDependencies:
    install:
    - package.json
```
{% endsnippetcut %}

Данная конструкция отвечает за отслеживание изменений в файле `package.json` и пересборки стадии `install` в случае нахождения таковых.

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
$ werf run --stages-storage :local --docker-options="-d -p 3000:3000 --restart=always" -- node /app/src/js/index.js
```

Первая часть команды очень похожа на `build`, а во второй — мы задаем [параметры docker](https://docs.docker.com/engine/reference/run/) и после двойного дефиса команду, с которой запустить образ.

Теперь приложение доступно локально на порту 3000:

![](/applications_guide_ru/images/applications-guide/020-hello-world-in-browser.png)

Как только мы убедились в том, что всё корректно, необходимо **загрузить образ в Registry**. Сборка с последующей загрузкой в Registry делается [командой `build-and-publish`]({{ site.docsurl }}/documentation/cli/main/build_and_publish.html). Когда werf запускается внутри CI-процесса, werf узнаёт реквизиты для доступа к Registry [из переменных окружения](https://docs.gitlab.com/ee/ci/variables/predefined_variables.html).

Если же werf запускается вне Gitlab CI, потребуется:

* Вручную подключиться к GitLab Registry [с помощью `docker login`](https://docs.docker.com/engine/reference/commandline/login/);
* Установить переменную окружения `WERF_IMAGES_REPO` с путём до Registry (вида `registry.mydomain.io/myproject`);
* Выполнить сборку и загрузку в Registry: `werf build-and-publish`.

Если вы всё правильно сделали, собранный образ появится в Registry. В случае использования Registry от GitLab собранный образ можно увидеть через веб-интерфейс GitLab.

<div id="go-forth-button">
    <go-forth url="20_iac.html" label="Конфигурирование инфраструктуры в виде кода" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
