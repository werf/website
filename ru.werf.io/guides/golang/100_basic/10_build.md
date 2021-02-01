---
title: Сборка образа
permalink: golang/100_basic/10_build.html
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
git clone git@github.com:werf/werf-guides.git
cp -r werf-guides/examples/golang/000_app ./
cd 000_app 
git init
git add .
git commit -m "initial commit"
```

_Так вы скопируете себе код [приложения на Go](https://github.com/werf/werf-guides/tree/master/examples/golang/000_app) и инициируете Git в каталоге с ним._

werf следует принципам [гитерминизма]({{ site.docsurl }}/documentation/advanced/configuration/giterminism.html): опирается на состояние, описанное в Git-репозитории. Это означает, что файлы, не коммитнутые в Git-репозиторий, по умолчанию будут игнорироваться. Благодаря этому, имея исходные коды приложения, вы всегда можете реализовать его конкретное работоспособное состояние.

## Реализация сборки в Dockerfile

Конфигурация сборки нашего приложения будет состоять из следующих шагов:

- взять официальный образ Go (`golang:1.15-alpine` подойдёт);
- добавить в него код приложения;
- установить модули-зависимости;
- скомпилировать бинарный файл приложения;
- взять базовый образ alpine (например, `alpine:3.13`);
- скопировать в него бинарный файл с приложением из сборочного образа;


Реализуем это в `Dockerfile`:

{% snippetcut name="Dockerfile" url="https://github.com/werf/werf-guides/blob/master/examples/golang/010_build/Dockerfile" %}
{% raw %}
```Dockerfile
FROM golang:1.15-alpine as builder
RUN apk add --update gcc musl-dev
WORKDIR /app
COPY go.* /app/
RUN go mod download
COPY /cmd /app/cmd
RUN go build ./cmd/demoapp

FROM alpine:3.13
COPY --from=builder /app/demoapp /app/demoapp
COPY /app.db /app/app.db

CMD ["/app/demoapp"]
```
{% endraw %}
{% endsnippetcut %}

## Интеграция Dockerfile с werf

Подключим готовый Dockerfile к werf. Для этого, создадим в корне репозитория файл `werf.yaml`, описывающий сборку всего проекта:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/golang/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-golang
configVersion: 1
---
image: basicapp
dockerfile: Dockerfile
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="Что тут написано?" %}

`werf.yaml` начинается с обязательной **секции мета-информации**:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/golang/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-golang
configVersion: 1
```
{% endraw %}
{% endsnippetcut %}

Здесь:
- **_project_** — поле, задающее уникальное имя проекта приложения. Имя проекта по умолчанию используется при генерации имени Helm-релиза и пространства имен, `namespace`, в Kubernetes. Изменение имени у активного проекта затруднительно и требует ряда действий, которые необходимо выполнить вручную (подробнее о возможных последствиях можно прочитать [здесь]({{ site.docsurl }}/documentation/reference/werf_yaml.html#последствия-смены-имени-проекта));
- **_configVersion_** определяет версию синтаксиса, используемую в `werf.yaml` (на данный момент поддерживается только версия `1`).

Следующая секция конфигурации, которая и будет основной для сборки: [**image config section**]({{ site.docsurl }}/documentation/reference/werf_yaml.html#секция-image).

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/golang/011_build_werf/werf.yaml" %}
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
```bash
git add .
git commit -m "work in progress"
```
{% endraw %}

Для того, чтобы запустить сборку, воспользуемся [командой `build`]({{ site.docsurl }}/documentation/reference/cli/werf_build.html):

{% raw %}
```bash
werf build
```
{% endraw %}

_В подглаве «Ускорение сборки» мы переведём сборку с Dockerfile на альтернативный синтаксис werf под названием Stapel и получим расширенные возможности: инкрементальную пересборку с учетом истории Git, возможность применения Ansible, использование кэша между сборками, удобные инструменты диагностики и многое другое._

Но уже сейчас вы можете заметить, что werf делает расширенный вывод логов сборки:

```
┌ ⛵ image basicapp
│ ┌ Building stage basicapp/dockerfile
│ │ basicapp/dockerfile  Sending build context to Docker daemon  69.63kB
│ │ basicapp/dockerfile  Step 1/16 : FROM golang:1.15-alpine as builder
│ │ basicapp/dockerfile   ---> 1de1afaeaa9a
│ │ basicapp/dockerfile  Step 2/16 : RUN apk add --update gcc musl-dev
│ │ basicapp/dockerfile   ---> Using cache
<..>
│ │ basicapp/dockerfile  Successfully built f2a49a68a666
│ │ basicapp/dockerfile  Successfully tagged c89c09f6-411f-4d4c-9c9e-43f4bdfad074:latest
│ ├ Info
│ │       name: werf-guided-golang:f276ddd4d73aafb69d657234505e718f78284bbdd816863f1540a912-1611832304608
│ │       size: 17.8 MiB
│ └ Building stage basicapp/dockerfile (55.42 seconds)
└ ⛵ image basicapp (55.75 seconds)


Running time 55.82 seconds
```

## Запуск

Запустим собранный образ с помощью [werf run]({{ site.docsurl }}/documentation/cli/main/run.html):

```shell
werf run --docker-options="--rm -p 3000:3000" -- /app/demoapp
```

Обратите внимание, что мы задаем [параметры docker](https://docs.docker.com/engine/reference/run/) опцией `--docker-options`, а саму команду запуска указываем после двух дефисов.

_Вы также можете заметить, что и вызов `werf run` осуществляет сборку, т.е. предварительно запускать сборку на самом деле не обязательно._

Теперь приложение доступно локально на порту 3000:

![](/guides/images/golang/100_10_app_in_browser.png)

## Внесение новых изменений

Мы будем постоянно дорабатывать приложение. Посмотрим, как это правильно делать на примере произвольных изменений в коде приложения:

{% snippetcut name="cmd/demoapp/main.go" url="#" %}
{% raw %}
```golang

func main() {
    ...
	r.Get("/labels", labelsHandler)
	...
}

...

func listLabels(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Our changes"))
}
```
{% endraw %}
{% endsnippetcut %}

1. Остановите ранее запущенный `werf run` (нажав Ctrl+C в консоли, где он запущен).
2. Запустите его заново: 
    ```bash
    werf run --docker-options="--rm -p 3000:3000" -- /app/demoapp
    ```
2. Посмотрите, как произойдёт пересборка и запуск, а затем обратитесь к API: http://example.com:3000/labels
3. Вы ожидали увидеть сообщение `Our changes`, но увидите старый результат. **Ничего не изменилось**, почему?

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
    werf run --docker-options="--rm -p 3000:3000" -- /app/demoapp
    ```
4. Посмотреть на результат в браузере.

Жёсткая связка с Git необходима для того, чтобы гарантировать воспроизводимость вашего решения. Подробнее о том, как работает эта механика _гитерминизма_, мы расскажем в главе «Необходимо знать», а пока что — сфокусируемся на сборке и доставке до кластера.
{% endofftopic %}


<div id="go-forth-button">
    <go-forth url="20_cluster.html" label="Подготовка кластера" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
