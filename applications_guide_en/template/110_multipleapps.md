---
title: Несколько приложений в одном репозитории
sidebar: applications_guide
guide_code: ____________
permalink: ____________/110_multipleapps.html
toc: false
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/service-app.yaml
- .helm/templates/app.yaml
- .helm/templates/ingress.yaml
- werf.yaml
{% endfilesused %}

В этой главе мы добавим к нашему базовому приложению ещё одно, находящееся в том же репозитории. Это корректная ситуация для:

* сложных случаев с двумя приложениями на двух разных языках;
* ситуации, когда есть более одного запускаемого процесса (например, сервис, отвечающий на HTTP-запросы, и worker);
* ситуации, когда в одном репозитории находятся frontend и backend.

Рекомендуем также посмотреть [доклад Дмитрия Столярова](https://www.youtube.com/watch?v=g9cgppj0gKQ) о том, почему и в каких ситуациях это хороший путь для микросервисов. Также вы можете посмотреть [аналогичную статью]({{ site.docsurl }}/documentation/guides/advanced_build/multi_images.html) о приложении с несколькими образами.

Добавим к нашему приложению Single Page Application-приложением на react которое отображает публичную часть. Наш подход будет очень похож на то, что делалось в главе [Генерируем и раздаем ассеты](040_assets.html), с одним существенным отличием: изменения в коде react сильно отделены от изменений в ____________ приложении. Как следствие — мы разнесём их в разные папки, а также в различные Pod-ы.

Мы рассмотрим вопрос организации структуры файлов и папок, соберём два образа: для ____________ приложения и для react приложения и сконфигурируем запуск этих образов в Kubernetes.

## Структура файлов и директорий

Структура каталогов будет организована следующим образом:

```
├── .helm/
│   ├── templates/
│   └── values.yaml
├── backend/
├── frontend/
└── werf.yaml
```

Для одного репозитория рекомендуется использовать один файл `werf.yaml` и одну папку `.helm` с конфигурацией инфраструктуры. Такой подход делает работу над кодом прозрачнее и помогает избегать рассинхронизации в двух частях одного проекта.

{% offtopic title="А если получится слишком много информации в одном месте и станет сложно ориентироваться?" %}
Helm обрабатывает все файлы, которые находятся в папке `.helm/templates`, а значит их может быть столько, сколько удобно вам. Для упрощения кода можно использовать [общие блоки](https://helm.sh/docs/chart_template_guide/named_templates/).

Кроме того `werf.yaml` также поддерживает [Описание конфигурации в нескольких файлах]({{ site.docsurl }}/documentation/configuration/introduction.html#%D0%BE%D0%BF%D0%B8%D1%81%D0%B0%D0%BD%D0%B8%D0%B5-%D0%BA%D0%BE%D0%BD%D1%84%D0%B8%D0%B3%D1%83%D1%80%D0%B0%D1%86%D0%B8%D0%B8-%D0%B2-%D0%BD%D0%B5%D1%81%D0%BA%D0%BE%D0%BB%D1%8C%D0%BA%D0%B8%D1%85-%D1%84%D0%B0%D0%B9%D0%BB%D0%B0%D1%85) и вынесение части кода в общие блоки.
{% endofftopic %}

## Сборка приложений

На стадии сборки приложения нам необходимо правильно организовать структуру файла `werf.yaml`, описав в нём сборку двух приложений на разном стеке.

Мы соберём два образа: `backend` c ____________-приложением и `frontend` c React-приложением. Для сборки последнего мы воспользуемся механизмом артефактов (и соберём артефакт `frontend-build`) — мы использовали подобное в главе [Генерируем и раздаем ассеты](040_assets.html).

{% offtopic title="Как конкретно?" %}

Сборка образа `backend` аналогична ранее описанному [базовому приложению](020_basic.html) с [зависимостями](030_dependencies.html), за исключением того, откуда берётся исходный код:

{% snippetcut name="werf.yaml" url="#" %}
{% raw %}
```yaml
git:
- add: /backend
  to: /app
```
{% endraw %}
{% endsnippetcut %}

Мы добавляем в собираемый образ только исходные коды, относящиеся к ____________ приложению. Таким образом, пересборка этой части проекта не будет срабатывать, когда изменился только React-код.

Сборка для frontend приложения описана в файле `werf.yaml` как отдельный образ. Поскольку nodejs нужен только для сборки - соберем артефакт:

{% snippetcut name="werf.yaml" url="#" %}
{% raw %}
```yaml
artifact: frontend-build
from: node:{{ .NODE_MAJOR }}
git:
- add: /frontend
  to: /app
```
{% endraw %}
{% endsnippetcut %}

И затем импортируем необходимые файлы в образ с nginx:

{% snippetcut name="deployment-frontend.yaml" url="#" %}
{% raw %}
```yaml
---
image: frontend
from: nginx:alpine
import:
- artifact: frontend-build
  add: /app/build
  to: /www
  after: setup
```
{% endraw %}
{% endsnippetcut %}

{% endofftopic %}

## Конфигурация инфраструктуры в Kubernetes

Подготовленные приложения мы будем запускать отдельными объектами Deployment: таким образом в случае изменений только в одной из частей приложения будет перевыкатываться только эта часть. Создадим два отдельных файла для описания объектов: `frontend.yaml` и `backend.yaml`. В условиях, когда в одном сервисе меньше 15-20 объектов — удобно следовать принципу максимальной атомарности в шаблонах.

При деплое нескольких Deployment крайне важно правильно прописать `selector`-ы в Service и Deployment:

{% snippetcut name="service-backend.yaml" url="#" %}
{% raw %}
```yaml
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .Chart.Name }}-backend
spec:
  selector:
    app: {{ .Chart.Name }}-backend
```
{% endraw %}
{% endsnippetcut %}

{% snippetcut name="deployment-backend.yaml" url="#" %}
{% raw %}
```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Chart.Name }}-backend
spec:
  selector:
    matchLabels:
      app: {{ .Chart.Name }}-backend
  template:
    metadata:
      labels:
        app: {{ .Chart.Name }}-backend
```
{% endraw %}
{% endsnippetcut %}

Маршрутизация запросов будет осуществляться через Ingress:

{% snippetcut name="ingress.yaml" url="#" %}
{% raw %}
```yaml
  rules:
  - host: {{ .Values.global.ci_url }}
    http:
      paths:
      - path: /
        backend:
          serviceName: {{ .Chart.Name }}-frontend
          servicePort: 80
      - path: /api
        backend:
          serviceName: {{ .Chart.Name }}-backend
          servicePort: ____________
```
{% endraw %}
{% endsnippetcut %}

<div id="go-forth-button">
    <go-forth url="120_dynamicenvs.html" label="Динамические окружения" base-url="applications_guide_ru"></go-forth>
</div>
