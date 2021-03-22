---
title: Подключение зависимостей
permalink: gitlab_java_springboot/030_dependencies.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- werf.yaml
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении работу с зависимостями. Важно корректно вписать зависимости в [стадии сборки]({{ site.docsurl }}/documentation/internals/stages_and_storage.html), что позволит не тратить время на пересборку зависимостей тогда, когда зависимости не изменились.

{% offtopic title="Что за стадии?" %}
werf подразумевает, что лучшей практикой будет разделить сборочный процесс на этапы, у каждого из которых есть свои четкие функции и назначение. Каждый такой этап соответствует промежуточному образу — подобно слоям в Docker. В werf такой этап называется стадией, и конечный образ в итоге состоит из набора собранных стадий. Все стадии хранятся в хранилище стадий, которое можно рассматривать как кэш сборки приложения, хотя по сути это скорее часть контекста сборки.

Стадии — это этапы сборочного процесса, кирпичи, из которых в итоге собирается конечный образ. Стадия собирается из группы сборочных инструкций, указанных в конфигурации. Причем группировка этих инструкций не случайна, имеет определенную логику и учитывает условия и правила сборки. С каждой стадией связан конкретный Docker-образ.

Подробнее о том, какие стадии для чего предполагаются, можно посмотреть [здесь]({{ site.docsurl }}/v1.1-stable/documentation/reference/stages_and_images.html). Если вкратце, то werf предлагает использовать для стадий следующую стратегию:

*   использовать стадию `beforeInstall` для инсталляции системных пакетов;
*   `install` — для инсталляции системных зависимостей и зависимостей приложения;
*   `beforeSetup` — для настройки системных параметров и установки приложения;
*   `setup` — для настройки приложения.

Другие подробности о стадиях описаны в [документации]({{ site.docsurl }}/v1.1-stable/documentation/configuration/stapel_image/assembly_instructions.html).

Одно из основных преимуществ использования стадий в том, что мы можем перезапускать сборку не с нуля, а только с той стадии, которая зависит от изменений в определенных файлах.
{% endofftopic %}

В Java, в частности в spring, в качестве менеджера зависимостей может использоваться maven, gradle. Мы будем, как и ранее использовать maven, но для gradle принципиальных отличий не будет. Пропишем его использование в файле `werf.yaml` и затем оптимизируем его использование.

## Подключение менеджера зависимостей

Пропишем разрешение зависимостей в нужную стадию сборки в `werf.yaml`:

{% snippetcut name="werf.yaml" url="#" ignore-tests %}
{% raw %}
```yaml
shell:
  beforeSetup:
    - mvn -B -f pom.xml package dependency:resolve
```
{% endraw %}
{% endsnippetcut %}

Однако, если оставить всё так — стадия не будет запускаться при изменении `pom.xml` и любого кода в `src/`. Подобная зависимость пользовательской стадии от изменений [указывается]({{ site.docsurl }}/v1.1-stable/documentation/configuration/stapel_image/assembly_instructions.html#%D0%B7%D0%B0%D0%B2%D0%B8%D1%81%D0%B8%D0%BC%D0%BE%D1%81%D1%82%D1%8C-%D0%BE%D1%82-%D0%B8%D0%B7%D0%BC%D0%B5%D0%BD%D0%B5%D0%BD%D0%B8%D0%B9-%D0%B2-git-%D1%80%D0%B5%D0%BF%D0%BE%D0%B7%D0%B8%D1%82%D0%BE%D1%80%D0%B8%D0%B8) с помощью параметра `git.stageDependencies`:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/030-deps/werf.yaml" %}
{% raw %}
```yaml
git:
- add: /
  to: /app
  stageDependencies:
    beforeSetup:
    - pom.xml
    setup:
    - src
```
{% endraw %}
{% endsnippetcut %}

Теперь при изменении файла `pom.xml` или любого из файлов в `src/` стадия `setup` будет запущена заново.

## Оптимизация сборки

Сборка занимает много времени, поэтому оптимизировать её — важная задача. Применим два приёма:

* Уменьшим объём скачиваемых файлов благодаря улучшенному использованию кэша maven.
* Усовершенствуем использование пользовательских стадий

Даже в пустом проекте сборщику нужно скачать приличное количество файлов. Скачивать эти файлы раз за разом выглядит нецелесообразным, поэтому разумно **переиспользовать кэш в `.m2/repository` между сборками**. С помощью директивы `mount` будем хранить кэш на раннере:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/030-deps/werf.yaml" %}
{% raw %}
```yaml
mount:
- from: build_dir
  to: /root/.m2/repository
```
{% endraw %}
{% endsnippetcut %}

**Усовершенствуем использование пользовательских стадий**: отделим resolve зависимостей от сборки jar — таким образом те коммиты, в которых правится исходный код, но не меняются зависимости, будут собираться быстрее.

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/030-deps/werf.yaml" %}
{% raw %}
```yaml
shell:
  beforeSetup:
    - cd /app
    - mvn -B -f pom.xml dependency:resolve
  setup:
  - cd /app
  - mvn -B -f pom.xml package
```
{% endraw %}
{% endsnippetcut %}


<div id="go-forth-button">
    <go-forth url="040_assets.html" label="Генерируем и раздаем ассеты" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
