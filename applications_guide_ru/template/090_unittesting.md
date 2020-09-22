---
title: Юнит-тесты и линтеры
sidebar: applications_guide
guide_code: ____________
permalink: ____________/090_unittesting.html
toc: false
---

{% filesused title="Файлы, упомянутые в главе" %}
- .gitlab-ci.yml
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении выполнение тестов/линтеров. Запуск тестов и линтеров - это отдельная стадия в пайплайне Gitlab CI, для выполнения которой может требоваться соблюдение определенных условий. Рассмотрим на примере ____________.

Требуется добавить эту зависимость  ____________ и прописать выполнение задания отдельной стадией на GitLab Runner командой [werf run]({{ site.docsurl }}/documentation/cli/main/run.html).

{% snippetcut name=".gitlab-ci.yml" url="#" %}
{% raw %}
```yaml
____________:
  stage: test
  script:
    - werf run ____________ -- ____________
  except:
    - schedules
  tags:
    - werf
  needs: ["Build"]
```
{% endraw %}
{% endsnippetcut %}

Созданную стадию нужно добавить в список стадий:

{% snippetcut name=".gitlab-ci.yml" url="#" %}
{% raw %}
```yaml
stages:
  - build
  - test
  - deploy
```
{% endraw %}
{% endsnippetcut %}

Обратите внимание, что процесс будет выполняться на runner'е, внутри собранного контейнера, но без доступа к базе данных и каким-либо ресурсам Kubernetes-кластера.

{% offtopic title="А если нужно больше?" %}
Если нужен доступ к ресурсам кластера или база данных — это уже не линтер: потребуется собрать отдельный образ и прописать сложный сценарий деплоя объектов Kubernetes. Эти случаи выходят за рамки нашего гайда для начинающих.
{% endofftopic %}

<div>
    <a href="110_multipleapps.html" class="nav-btn">Далее: Несколько приложений в одном репозитории</a>
</div>
