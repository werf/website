---
title: Эффективный и консистентный CI/CD с Kubernetes
permalink: /
layout: default
sidebar: none
---

<div class="intro-banner" id="intro-banner">
    <div class="page__container">
        <div class="intro-banner__background-shapes">
            <img class="left" src="/assets/images/backgrounds/intro-banner-left.svg" alt="">
            <img class="right" src="/assets/images/backgrounds/intro-banner-right.svg" alt="">
        </div>
        <div class="intro-banner__wrap">
            <div class="intro-banner__title">Недостающая часть вашей CI/CD-системы</div>
            <div class="building-utility__pic">
                <img src="{% asset landing/building-utility-scheme.svg @path %}" alt="">
            </div>
            <div class="intro-banner__links">
                <div class="intro-banner__links-github">
                    <span class="page__icon page__icon_github-white intro-banner__github-icon"></span>
                    <a href="https://github.com/werf/werf" class="intro-banner__github-counter">
                        <span class="intro-banner__github-counter-num"><span class="gh_counter">3600</span></span>
                    </a>
                </div>
                <div class="intro-banner__link-cncf">
                    <a href="https://www.cncf.io/projects/werf/" target="_blank">
                        <img src="/assets/images/cncf-logo-small.svg" alt="">
                        <div class="link-cncf__text">werf — проект категории sandbox <br>в Cloud Native Computing Foundation</div>
                    </a>
                </div>
            </div>
        </div>
    </div>
</div>

<div class="features-card" id="werf-features">
    <div class="page__container">
        <div class="features-card__wrap">
            <ul class="features-card__list">
                <li class="features-card__item">
                    <div class="card__item-icon">
                        <img src="{% asset icons/finger.svg @path %}" alt="">
                    </div>
                    <div class="card__item-title">Простота <br>использования</div>
                    <div class="card__item-text">Предоставьте Dockerfile и<br>Helm-чарт — остальное werf возьмет на себя</div>
                </li>
                <li class="features-card__item">
                    <div class="card__item-icon">
                        <img src="{% asset icons/box-icon.svg @path %}" alt="">
                    </div>
                    <div class="card__item-title">Один <br>инструмент</div>
                    <div class="card__item-text">Используйте одно решение для сборки образов, запуска тестов, дистрибуции релизных артефактов и развёртывания приложения в Kubernetes</div>
                </li>
                <li class="features-card__item">
                    <div class="card__item-icon">
                        <img src="{% asset icons/weightlifter-icon.svg @path %}" alt="">
                    </div>
                    <div class="card__item-title">Продвинутые <br>возможности</div>
                    <div class="card__item-text">Получайте выгоду от автоматического кэширования, тегирования на основе содержимого, отслеживания ресурсов в Helm и другого</div>
                </li>
                <li class="features-card__item">
                    <div class="card__item-icon">
                        <img src="{% asset icons/triangel.svg @path %}" alt="">
                    </div>
                    <div class="card__item-title">Интеграция <br>стандартных технологий</div>
                    <div class="card__item-text">Положитесь на привычные Git, Buildah, Helm, Kubernetes и любимую CI-систему</div>
                </li>
            </ul>
        </div>
    </div>
</div>

<div class="what-problems-solve">
    <div class="page__container">
        <div class="what-problems__title title-h3">Что приносит werf в ваш CI/CD</div>
        <div class="what-problems__subtitle">
            <span class="subtitle-problem-number">1</span>
            <div class="subtitle-problem-text">Предсказуемость и надежность процесса доставки</div>
        </div>
        <div class="what-problems__wrap">
            <div class="what-problems__grid">
                <ul class="what-problems__list">
                    <li class="features-card__item">
                        <div class="card__item-title">Детерминированность CI/CD</div>
                        <div class="card__item-text">What you Git is what you get. werf предлагает подход гитерминизм, который cтимулирует IaC и Git как единый источник истины</div>
                    </li>
                    <li class="features-card__item">
                        <div class="card__item-title">Воспроизводимость сборок</div>
                        <div class="card__item-text">Cинхронизация с container registry. Однажды собранный образ неизменяем, а работа с container registry скоординирована между сборщиками</div>
                    </li>
                    <li class="features-card__item">
                        <div class="card__item-title">Предсказуемость выката</div>
                        <div class="card__item-text">
                            <p>Успешный выкат – это готовое к использованию приложение</p>
                            <p>Перевыкат только изменившихся компонентов. Теги образов приложения независимы (решение вопроса тегирования образов в монорепозиториях)</p>
                        </div>
                    </li>
                    <li class="features-card__item">
                        <div class="card__item-title">Диагностика проблем</div>
                        <div class="card__item-text">
                            <p>Связывание релизных артефактов с Git и CI/CD</p>
                            <p>Мгновенная и богатая обратная связь при развёртывании в Kubernetes</p>
                        </div>
                    </li>
                </ul>
            </div>
            <div class="what-problems__pic">
                <img src="/assets/images/werf-schema_ru.svg" alt="">
            </div>
        </div>
        <div class="what-problems__subtitle">
            <span class="subtitle-problem-number">2</span>
            <div class="subtitle-problem-text">Эффективность</div>
        </div>
        <div class="what-problems__wrap">
            <div class="what-problems__pic">
                <img src="/assets/images/werf-time.svg" alt="">
            </div>
            <div class="what-problems__grid">
                <ul class="what-problems__list">
                    <li class="features-card__item">
                        <div class="card__item-title">Инкрементальность сборок</div>
                        <div class="card__item-text">
                            <p>Пересборка только изменившихся образов</p>
                            <p>Переиспользование раннее собранных слоёв из container registry</p>
                        </div>
                    </li>
                    <li class="features-card__item">
                        <div class="card__item-title">Время развёртывания</div>
                        <div class="card__item-text">
                            <p>Перевыкат только изменившихся компонентов</p>
                            <p>Моментальное завершение проблемного выката</p>
                        </div>
                    </li>
                </ul>
            </div>
        </div>
    </div>
</div>

<div class="demo-block" id="how-it-works">
    <div class="page__container">
        <div class="demo-block__wrap">
            <div class="demo-block__title title-h3">Демонстрация</div>
            <div class="demo-block__tabs-wrap">
                <div class="demo-block__code">
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
                                <span>FROM node</span>
                                <br>
                                <span>WORKDIR /app</span>
                                <span>COPY . .</span>
                                <span>RUN npm ci</span>
                                <br>
                                <span>CMD ["node", "server.js"]</span>
                            </div>
                        </li>
                        <li class="code-item">
                            <div class="code-item-title">.helm/templates/deployment.yaml</div>
                            <div class="code-item-text">
                                <span>apiVersion: apps/v1</span>
                                <span>kind: Deployment</span>
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
                <div class="demo-block__tabs">
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

{%- assign publications = site.data.publications.articles | sort: 'created' | reverse %}

<div class="publications" id="publications">
    <div class="page__container">
        <div class="publications__content">
            <div class="publications__title">Последние публикации</div>
            <div class="publications__subtitle">Узнайте больше о фичах и практике<br />использования утилиты из новых статей</div>
            <div class="publications__cards">
                <ul class="publications__cards--list">
                    {%- for publication in publications limit: 4 %}
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

<div class="community" id="community">
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
