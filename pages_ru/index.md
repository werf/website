---
title: Гитерминированная CLI утилита
permalink: /
layout: default
sidebar: none
---

<!--<div class="news-and-updates">-->
<!--    <div class="news-and-updates__block">-->
<!--        <div class="news-and-updates__title">-->
<!--            Новости разработки-->
<!--        </div>-->
<!--        <div class="news-and-updates__item" id="nau-news"></div>-->
<!--    </div>-->
<!--    <div class="news-and-updates__block">-->
<!--        <div class="news-and-updates__title">-->
<!--            Последние релизы-->
<!--        </div>-->
<!--        <div class="news-and-updates__row" id="nau-releases"></div>-->
<!--    </div>-->
<!--</div>-->

<div class="intro-banner">
    <div class="page__container">
        <div class="intro-banner__wrap">
            <div class="intro-banner__title">Лучшие практики CI/CD из коробки</div>
            <div class="intro-banner__github">
                <span class="page__icon page__icon_github intro-banner__github-icon"></span>
                <a href="https://github.com/werf/werf" class="intro-banner__github-counter">
                    <span class="intro-banner__github-counter-num"><span class="gh_counter">3414</span></span>
                </a>
            </div>
            <div class="intro-banner__tags">
                <ul class="tags__list">
                    <li class="tags__item">#разработка</li>
                    <li class="tags__item">#тестирование</li>
                    <li class="tags__item">#сборка</li>
                </ul>
                <ul class="tags__list">
                    <li class="tags__item">#дистрибуция</li>
                    <li class="tags__item">#развертывание</li>
                </ul>
                <ul class="tags__list">
                    <li class="tags__item">#очистка</li>
                </ul>
            </div>
        </div>
    </div>
</div>

<div class="building-utility">
    <div class="page__container">
        <div class="building-utility__wrap">
            <div class="building-utility__title">Утилита для построения полного цикла<br> доставки CI/CD c Kubernetes</div>
            <div class="building-utility__pic">
                <img src="{% asset landing/building-utility-scheme.svg @path %}" alt="">
            </div>
        </div>
    </div>
</div>

<div class="features-card">
    <div class="page__container">
        <div class="features-card__wrap">
            <ul class="features-card__list">
                <li class="features-card__item">
                    <div class="card__item-icon">
                        <img src="{% asset icons/gear.svg @path %}" alt="">
                    </div>
                    <div class="card__item-title">Простота <br>использования</div>
                    <div class="card__item-text">Предоставьте Dockerfile и Helm-чарт, остальное werf возьмет на себя</div>
                </li>
                <li class="features-card__item">
                    <div class="card__item-icon">
                        <img src="{% asset icons/puzzles.svg @path %}" alt="">
                    </div>
                    <div class="card__item-title">Интеграция стандартных технологий</div>
                    <div class="card__item-text">Git, Buildah (tooltip: Dockerfile-сборщик от RedHat), Helm, Kubernetes и ваша CI-система</div>
                </li>
                <li class="features-card__item">
                    <div class="card__item-icon">
                        <img src="{% asset icons/certificate.svg @path %}" alt="">
                    </div>
                    <div class="card__item-title">Улучшение стандартных технологий</div>
                    <div class="card__item-text">Продвинутое кеширование, тегирование на основе содержимого, продвинутое отслеживание ресурсов в Helm и другое</div>
                </li>
            </ul>
        </div>
    </div>
</div>

<div class="three-files-you-need">
    <div class="page__container">
        <div class="three-files-you-need__wrap">
            <div class="three-files-you-need__title">Демонстрация работы werf</div>
            <div class="three-files-you-need__tabs-wrap">
                <div class="three-files-you-need__code">
                    <ul class="code-list">
                        <li class="code-item">
                            <div class="code-item-title">werf.yaml</div>
                            <div class="code-item-text">
                                <span>configVersion: 1</span>
                                <span>project: hello</span>
                                <span>---</span>
                                <span>image: hello</span>
                                <span>dockerfile: ./Dockerfile</span>
                            </div>
                        </li>
                        <li class="code-item">
                            <div class="code-item-title">Dockerfile</div>
                            <div class="code-item-text">
                                <span>FROM mode</span>
                                <br>
                                <span>WORKDIR /app</span>
                                <span>COPY . .</span>
                                <span>RUN npm ci</span>
                                <br>
                                <span>CMD [“node”, “server.js”]</span>
                            </div>
                        </li>
                        <li class="code-item">
                            <div class="code-item-title">.helm/templates/deployment.yaml</div>
                            <div class="code-item-text">
                                <span>apiVersion: apps/v1</span>
                                <span>kind: Deployment</span>
                                <span>WORKDIR /app</span>
                                <span>metadata:
                                    <span>name: hello</span>
                                </span>
                                <span>spec:
                                <span>selector:
                                    <span>matchLabels:
                                        <span>app: hello</span>
                                    </span>
                                </span>
                                <span>template:
                                    <span>metadata:
                                        <span>labels:
                                            <span>app: hello</span>
                                        </span>
                                    </span>
                                    <span>spec:
                                        <span>containers:<br>
                                            - image: {% raw %}{{ .Values.werf.image.hello }}{% endraw %}
                                        </span>
                                    </span>
                                </span>
                            </span>
                            </div>
                        </li>
                    </ul>
                </div>
                <div class="three-files-you-need__tabs">
                    <div class="tabs__nav">
                        <ul class="tabs__nav-list">
                            <li class="tabs__nav-item"><button class="button active" data-tabs-button="build">build</button></li>
                            <li class="tabs__nav-item"><button class="button" data-tabs-button="test">test</button></li>
                            <li class="tabs__nav-item"><button class="button" data-tabs-button="deploy">deploy</button></li>
                            <li class="tabs__nav-item"><button class="button" data-tabs-button="distribute">distribute</button></li>
                            <li class="tabs__nav-item"><button class="button" data-tabs-button="cleanup">cleanup</button></li>
                        </ul>
                    </div>
                    <div class="tabs__video-wrap">
                        <ul class="tabs__video-list">
                            <li class="tabs__video-item" data-tabs-video="build">
                                <div class="tabs__video-item-terminal">
                                    <div class="tabs__video-item-container" id="demo"></div>
                                </div>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<div class="page__container documentation-block">
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
