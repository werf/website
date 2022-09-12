---
title: Гитерминированная CLI утилита
permalink: /
layout: default
sidebar: none
---

<div class="news-and-updates">
    <div class="news-and-updates__block">
        <div class="news-and-updates__title">
            Новости разработки
        </div>
        <div class="news-and-updates__item" id="nau-news"></div>
    </div>
    <div class="news-and-updates__block">
        <div class="news-and-updates__title">
            Последние релизы
        </div>
        <div class="news-and-updates__row" id="nau-releases"></div>
    </div>
</div>

<div class="intro-scheme" id="intro-scheme">
    <div class="page__container intro-scheme__container">
        {% include common/intro.md %}
    </div>
</div>

<div class="intro">
    <div class="intro__bg" id="intro-bg"></div>
    <div class="page__container intro__container">
        <div class="intro__row">
            <div class="intro__row-item" id="intro-title">
                <div class="intro__subtitle">Что это?</div>
                <h1 class="intro__title">Инструмент<br/>консистентной<br/>доставки</h1>
                <ul class="intro__features">
                    <li>CLI-утилита, «склеивающая» Git, Docker, Helm и Kubernetes<br>
                    с любой CI-системой для реализации CI/CD и подхода гитерминизм.</li>
                </ul>
                <div class="intro__btns page__btn-group">
                    <a href="{{ "how_it_works.html" | true_relative_url }}" target="_blank" class="page__btn page__btn_b page__btn_small">
                        Как это работает
                    </a>
                    <a href="{{ "documentation/quickstart.html" | true_relative_url }}" target="_blank" class="page__btn page__btn_b page__btn_small">
                        Быстрый старт
                    </a>
                    <a href="{{ "documentation/index.html" | true_relative_url }}" target="_blank" class="page__btn page__btn_b page__btn_small">
                        Документация
                    </a>
                </div>
            </div>
        </div>
    </div>
</div>

<div class="intro">
    <div class="page__container intro__container">
        <div class="intro__row">
            <div class="intro__row-item" id="intro-title">
                <div class="intro__subtitle">Зачем нужно?</div>
                <h1 class="intro__title">Быстрый<br/>и эффективный<br/>CI/CD</h1>
                <ul class="intro__features">
                    <li>Создавайте эффективные, предсказуемые и целостные<br>
                    CI/CD-процессы на базе устоявшихся технологий.</li>
                    <li>С werf быстро начать работу, легко применять<br>
                    лучшие практики и не нужно переизобретать колесо.</li>
                </ul>
            </div>
        </div>
    </div>
</div>

<div class="presentation" id="presentation">
    <div class="page__container presentation__container">
        <div class="presentation__row">
            <div class="presentation__row-item" id="presentation-title">
                <div class="presentation__subtitle">Как работает?</div>
                <h1 class="presentation__title"
                    data-toggle="tooltip" title="Что ты Git'ишь, то и видишь!">
                    What you Git<br/> is what you get!
                </h1>
                <ul class="presentation__features">
                    <li>Git используется как единый источник истины («гитерминизм»).</li>
                    <li>werf не только собирает и деплоит, но и непрерывно<br/>
                    синхронизирует изменения в Git с состоянием Kubernetes.</li>
                </ul>
            </div>
            <div class="presentation__row-item presentation__row-item_scheme">
                {% include common/scheme.md %}
            </div>
        </div>
    </div>
</div>

<div class="page__container">
    <div class="presentation-notes">
        <div class="presentation-notes__item" id="presentation-notes-1">
            <div class="presentation-notes__item-num">1</div>
            <div class="presentation-notes__item-title">
                Желаемое состояние<br>
                определяется в Git
            </div>
            <div class="presentation-notes__item-text"></div>
        </div>
        <div class="presentation-notes__item" id="presentation-notes-2">
            <div class="presentation-notes__item-num">2</div>
            <div class="presentation-notes__item-title">
                Приводит registry<br>
                к желаемому состоянию
            </div>
            <div class="presentation-notes__item-text">
                <ol>
                    <li>
                        Собирает образы (при&nbsp;изменениях)
                    </li>
                    <li>
                        Публикует образы (при&nbsp;необходимости)
                    </li>
                </ol>
            </div>
        </div>
        <div class="presentation-notes__item" id="presentation-notes-3">
            <div class="presentation-notes__item-num">3</div>
            <div class="presentation-notes__item-title">
                Приводит Kubernetes<br>
                к желаемому состоянию
            </div>
            <div class="presentation-notes__item-text">
                <ol>
                    <li>
                        Выкатывает Helm-чарт (с&nbsp;необходимыми образами)
                    </li>
                    <li>
                        Отслеживает прогресс вплоть до успешного выката (и&nbsp;дает обратную связь)
                    </li>
                </ol>
            </div>
        </div>
    </div>
</div>

<div class="features">
    <div class="page__container">
        <ul class="features__list">
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_deploy"></div>
                <div class="features__list-item-title">Два варианта деплоя</div>
                <div class="features__list-item-text">
                    werf даёт 2 варианта деплоя приложений:
                    <ol>
                        <li>
                            <b>converge</b> приложения из git-коммита в кластер Kubernetes;
                        </li>
                        <li>
                            <b>публикация</b> приложения из git-коммита в Container Registry в виде <b>бандла</b>, затем <b>деплой</b> этого <b>бандла</b> в Kubernetes.
                        </li>
                    </ol>
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_usage"></div>
                <div class="features__list-item-title">Гибкое взаимодействие</div>
                <div class="features__list-item-text">
                    Разные способы взаимодействия c werf:
                    <ol>
                        <li>
                            вручную;
                        </li>
                        <li>
                            из CI/CD системы;
                        </li>
                        <li>
                            как оператор Kubernetes <small>(пока недоступно...)</small>;
                        </li>
                        <li>
                            через git push как в heroku <small>(пока недоступно...)</small>.
                        </li>
                    </ol>
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_lifecycle"></div>
                <div class="features__list-item-title">Консольная утилита</div>
                <div class="features__list-item-text">
                    werf — это не SaaS, а самодостаточная CLI-утилита с открытым кодом, запускаемая на стороне клиента. werf можно использовать как для <b>локальной разработки</b>, так и для <b>встраивания в любую CI/CD-систему</b> (GitLab CI/CD, GitHub Actions, Jenkins, CircleCI и т.д.), оперируя основными командами как составляющими пайплайна:
                    <ul>
                        <li><code>werf converge</code>;</li>
                        <li><code>werf dismiss</code>;</li>
                        <li><code>werf cleanup</code>;</li>
                        <li><code>werf bundle publish</code>;</li>
                        <li><code>werf bundle apply</code>.</li>
                    </ul>
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_kubernetes"></div>
                <div class="features__list-item-title">Простая в использовании</div>
                <div class="features__list-item-text">
                    werf работает «из коробки» с минимальной конфигурацией. Вам не нужно быть DevOps/SRE-инженером, чтобы использовать werf. Доступно <a href="/guides.html"><b>множество гайдов</b></a>, которые помогут быстро организовать деплой приложений в Kubernetes.
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_config"></div>
                <div class="features__list-item-title">Объединяет лучшее</div>
                <div class="features__list-item-text">
                    werf связывает привычные инструменты, превращая их в понятую, целостную, <b>интегрированную CI/CD-платформу</b>. werf делает хорошо контролируемым и удобным взаимодействие Git, Docker, вашего container registry и существующей CI-системы, Helm и Kubernetes.
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_fast"></div>
                <div class="features__list-item-title">Распределенная сборка</div>
                <div class="features__list-item-text">
                    В werf реализован продвинутый сборщик, среди возможностей которого — алгоритм распределенной сборки. Благодаря нему и его распределенному кэшированию <b>ваши пайплайны становятся по-настоящему быстрыми</b>.
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_debug"></div>
                <div class="features__list-item-title">Встроенная очистка</div>
                <div class="features__list-item-text">
                    Продуманный алгоритм <b>очистки неиспользуемых Docker-образов</b> в werf основан на анализе Git-истории собираемых приложений.
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_helm"></div>
                <div class="features__list-item-title">Расширенный Helm</div>
                <div class="features__list-item-text">
                    В werf встроен бинарник <code>helm</code>, который реализует процесс деплоя, совместимый с Helm, и расширяет его возможности. С ним не требуется отдельная установка <code>helm</code>, а его дополнения обеспечивают детальные и понятные <b>логи при деплое</b>, быстрое <b>определение сбоев</b> во время деплоя, поддержку секретов и другие фичи, превращающие деплой в <b>понятный и надежный процесс</b>.
                </div>
            </li>
            <li class="features__list-item features__list-item_special">
                <div class="features__list-item-title">Open Source</div>
                <div class="features__list-item-description">
                    <a href="https://github.com/werf/werf" target="_blank">Код открыт</a> и написан на Go. За годы развития проекта у него сформировалось большое сообщество пользователей.
                </div>
            </li>
        </ul>
    </div>
</div>

<div class="stats">
    <div class="page__container">
        <div class="stats__content">
            <div class="stats__title">Активная разработка</div>
            <ul class="stats__list">
                <li class="stats__list-item">
                    <div class="stats__list-item-num">4</div>
                    <div class="stats__list-item-title">релиза в неделю</div>
                    <div class="stats__list-item-subtitle">в среднем за прошлый год</div>
                </li>
                <li class="stats__list-item">
                    <div class="stats__list-item-num">2000+</div>
                    <div class="stats__list-item-title">инсталляций</div>
                    <div class="stats__list-item-subtitle">в больших и маленьких проектах</div>
                </li>
                <li class="stats__list-item">
                    <div class="stats__list-item-num gh_counter">2100</div>
                    <div class="stats__list-item-title">звезд на GitHub</div>
                    <div class="stats__list-item-subtitle">поддержите проект ;)</div>
                </li>
            </ul>
        </div>
    </div>
</div>

<div class="reliability">
    <div class="page__container">
        <div class="reliability__content">
            <div class="reliability__column">
                <div class="reliability__title">
                    werf — это зрелый, надежный<br>
                    инструмент, которому можно доверять
                </div>
                <a href="{{ "about/release_channels.html" | true_relative_url }}" class="page__btn page__btn_b page__btn_small page__btn_inline">
                    Подробнее о каналах обновлений
                </a>
            </div>
            <div class="reliability__column reliability__column_image">
                <div class="reliability__image"></div>
            </div>
        </div>
    </div>
</div>

<div class="community">
    <div class="page__container">
        <div class="community__content">
            <div class="community__title">Растущее дружелюбное сообщество</div>
            <div class="community__subtitle">Мы всегда на связи с сообществом<br/> в Telegram, Twitter и GitHub Discussions.</div>
            <div class="community__btns">
                <a href="{{ site.social_links[page.lang].telegram }}" target="_blank" class="page__btn page__btn_w community__btn">
                    <span class="page__icon page__icon_telegram"></span>
                    Telegram
                </a>
                <a href="{{ site.social_links[page.lang].twitter }}" target="_blank" class="page__btn page__btn_w community__btn">
                    <span class="page__icon page__icon_twitter"></span>
                    Twitter
                </a>
                <a href="https://github.com/werf/werf/discussions" rel="noopener noreferrer" target="_blank" class="page__btn page__btn_w community__btn">
                    <span class="page__icon page__icon_github"></span>
                    GitHub Discussions
                </a>
            </div>
        </div>
    </div>
</div>

{%- assign publications = site.data.publications.articles | sort: 'created' | reverse %}
{%- assign publications_by_year = publications | group_by_exp: "publication", "publication.created | date: '%Y'" %}

<div class="publications">
  <div class="page__container">
    <div class="publications__content">
      <div class="publications__title">Последние публикации про werf</div>
      <div class="publications__subtitle">Узнайте больше о фичах и практике<br />использования утилиты из новых статей</div>
      <div class="publications__cards">
        <ul class="publications__cards--list">
        {%- for year in publications_by_year limit: 1 %}
          {%- for publication in year.items limit: 4 %}
            <li class="publications__cards--item">
              {%- if publication.url %}
                <a href="{{ publication.url }}" class="publications__cards--link" target="_blank">
                  <span class="publications__cards--pic" style="background-image: url('{{ publication.img | true_relative_url }}')"></span>
                  <div class="publications__cards--title">{{ publication.title }}</div>
                  <div class="publications__cards--date">{{ publication.created | date: "%d-%m-%Y" }}</div>
                </a>
              {%- endif %}
              {%- if publication.medium_url %}
                <a href="{{ publication.medium_url }}" class="publications__cards--link" target="_blank">
                  <span class="publications__cards--pic" style="background-image: url('{{ publication.img | true_relative_url }}')"></span>
                  <div class="publications__cards--title">{{ publication.title }}</div>
                  <div class="publications__cards--date">{{ publication.created | date: "%d-%m-%Y" }}</div>
                </a>
              {%- endif %}
              {%- if publication.blog_url %}
                <a href="{{ publication.blog_url }}" class="publications__cards--link" target="_blank">
                  <span class="publications__cards--pic" style="background-image: url('{{ publication.img | true_relative_url }}')"></span>
                  <div class="publications__cards--title">{{ publication.title }}</div>
                  <div class="publications__cards--date">{{ publication.created | date: "%d-%m-%Y" }}</div>
                </a>
              {%- endif %}
              {%- if publication.habr_url %}
                <a href="{{ publication.habr_url }}" class="publications__cards--link" target="_blank">
                  <span class="publications__cards--pic" style="background-image: url('{{ publication.img | true_relative_url }}')"></span>
                  <div class="publications__cards--title">{{ publication.title }}</div>
                  <div class="publications__cards--date">{{ publication.created | date: "%d-%m-%Y" }}</div>
                </a>
              {%- endif %}
              {%- for custom in publication.custom_urls %}
                <a href="{{ custom.url }}" class="publications__cards--link" target="_blank">
                  <span class="publications__cards--pic" style="background-image: url('{{ publication.img | true_relative_url }}')"></span>
                  <div class="publications__cards--title">{{ publication.title }}</div>
                  <div class="publications__cards--date">{{ publication.created | date: "%d-%m-%Y" }}</div>
                </a>
              {%- endfor -%}
              {%- if publication.youtube_url %}
                <a href="{{ publication.youtube_url }}" class="publications__cards--link" target="_blank">
                  <span class="publications__cards--pic" style="background-image: url('{{ publication.img | true_relative_url }}')"></span>
                  <div class="publications__cards--title">{{ publication.title }}</div>
                  <div class="publications__cards--date">{{ publication.created | date: "%d-%m-%Y" }}</div>
                </a>
              {%- endif %}
              {%- if publication.tproger_url %}
                <a href="{{ publication.tproger_url }}" class="publications__cards--link" target="_blank">
                  <span class="publications__cards--pic" style="background-image: url('{{ publication.img | true_relative_url }}')"></span>
                  <div class="publications__cards--title">{{ publication.title }}</div>
                  <div class="publications__cards--date">{{ publication.created | date: "%d-%m-%Y" }}</div>
                </a>
              {%- endif %}
            </li>
          {%- endfor %}
        {%- endfor %}
        </ul>
      </div>
      <div class="community__btns">
          <a href="/publications.html" class="page__btn page__btn_o publications__btn">
              <span>Ещё публикации</span>
          </a>
      </div>
    </div>
  </div>
</div>

<div class="page__container">
    <div class="documentation">
        <div class="documentation__image">
        </div>
        <div class="documentation__info">
            <div class="documentation__info-title">
                Исчерпывающая документация
            </div>
            <div class="documentation__info-text">
              <a href="{{ "documentation/index.html" | true_relative_url }}">Документация</a> содержит более 100 статей, включающих описание частых случаев (первые шаги, деплой в Kubernetes, интеграция с CI/CD-системами и другое), полное описание функций, архитектуры и CLI-команд.
            </div>
        </div>
        <div class="documentation__btns">
            <a href="{{ "how_it_works.html" | true_relative_url }}" class="page__btn page__btn_b documentation__btn">
                Как это работает
            </a>
            <a href="{{ "documentation/quickstart.html" | true_relative_url }}" class="page__btn page__btn_o documentation__btn">
                Быстрый старт
            </a>
            <a href="/guides.html" class="page__btn page__btn_o documentation__btn">
                Руководства
            </a>
        </div>
    </div>
</div>
