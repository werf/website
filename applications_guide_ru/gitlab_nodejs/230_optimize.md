---
title: Ускорение сборки
permalink: gitlab_nodejs/230_optimize.html
---

TODO: Переходим на стапель
TODO: Оптимизация с учётом особенностей фреймворка (кэши, отделение компиляций и т.п)

## Оптимизация сборки

Чтобы не хранить ненужные кэши пакетного менеджера в образе, можно при сборке смонтировать директорию, в которой будет храниться кэш.

Для того, чтобы оптимизировать работу с этим кешем при сборке, мы добавим специальную конструкцию в `werf.yaml`:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/030-deps/werf.yaml" %}
{% raw %}
```yaml
mount:
- from: build_dir
  to: /var/cache/apt
```
{% endraw %}
{% endsnippetcut %}

Теперь при каждом запуске сборки эта директория будет монтироваться с сервера, где запускается `build`, и она не будет очищаться между сборками. Таким образом, кэш будет сохраняться между сборками.

<div id="go-for
<div id="go-forth-button">
    <go-forth url="210_cluster.html" label="Сборка" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
