---
title: Efficient and consistent CI/CD with Kubernetes
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
            <div class="intro-banner__title">Missing part of your CI/CD system</div>
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
                        <!-- spell-delete -->
                        <img src="/assets/images/cncf-logo-small.svg" alt="">
                        <!-- end-spell-delete -->
                        <div class="link-cncf__text">werf is a Cloud Native Computing Foundation<br> sandbox project</div>
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
                    <div class="card__item-title">Ease of use</div>
                    <div class="card__item-text">Create your Dockerfile and Helm chart, and let werf handle all the rest</div>
                </li>
                <li class="features-card__item">
                    <div class="card__item-icon">
                        <img src="{% asset icons/box-icon.svg @path %}" alt="">
                    </div>
                    <div class="card__item-title">A single, all-in-one tool</div>
                    <div class="card__item-text">Build images, run tests, distribute release artifacts, and deploy the application to Kubernetes using a unified, all-in-one tool</div>
                </li>
                <li class="features-card__item">
                    <div class="card__item-icon">
                        <img src="{% asset icons/weightlifter-icon.svg @path %}" alt="">
                    </div>
                    <div class="card__item-title">Advanced features</div>
                    <div class="card__item-text">Take advantage of automatic caching, content-based tagging, resource tracking in Helm, and more</div>
                </li>
                <li class="features-card__item">
                    <div class="card__item-icon">
                        <img src="{% asset icons/triangel.svg @path %}" alt="">
                    </div>
                    <div class="card__item-title">Gluing technologies</div>
                    <div class="card__item-text">Rely on Git, <span class="tooltip-text" title="Dockerfile builder from Red Hat">Buildah</span>, Helm, Kubernetes, and your<br>CI system of choice</div>
                </li>
            </ul>
        </div>
    </div>
</div>

<div class="what-problems-solve">
    <div class="page__container">
        <div class="what-problems__title title-h3">What werf brings to your CI/CD</div>
        <div class="what-problems__subtitle">
            <span class="subtitle-problem-number">1</span>
            <div class="subtitle-problem-text">Predictable and reliable delivery process</div>
        </div>
        <div class="what-problems__wrap">
            <div class="what-problems__grid">
                <ul class="what-problems__list">
                    <li class="features-card__item">
                        <div class="card__item-title">Deterministic CI/CD</div>
                        <div class="card__item-text">What you Git is what you get. werf introduces Giterminism that encourages an IaC approach and the use of Git as a single source of truth</div>
                    </li>
                    <li class="features-card__item">
                        <div class="card__item-title">Reproducible builds</div>
                        <div class="card__item-text">Synchronization with the container registry. Once built, the image is immutable, while builders use the container registry in a coordinated fashion</div>
                    </li>
                    <li class="features-card__item">
                        <div class="card__item-title">Predictable deployment</div>
                        <div class="card__item-text">
                            <p>A successful deployment means that an application is up and running</p>
                            <p>Redeploying only the components that have been modified. Application image tags are independent (this solves the issue of image tagging in monorepos)</p>
                        </div>
                    </li>
                    <li class="features-card__item">
                        <div class="card__item-title">Easier troubleshooting</div>
                        <div class="card__item-text">
                            <p>Linking release artifacts to Git and CI/CD</p>
                            <p>Instant and detailed feedback when deploying to Kubernetes</p>
                        </div>
                    </li>
                </ul>
            </div>
            <div class="what-problems__pic">
                <img src="/assets/images/werf-schema.svg" alt="">
            </div>
        </div>
        <div class="what-problems__subtitle">
            <span class="subtitle-problem-number">2</span>
            <div class="subtitle-problem-text">Efficiency</div>
        </div>
        <div class="what-problems__wrap">
            <div class="what-problems__pic">
                <img src="/assets/images/werf-time.svg" alt="">
            </div>
            <div class="what-problems__grid">
                <ul class="what-problems__list">
                    <li class="features-card__item">
                        <div class="card__item-title">Incremental builds</div>
                        <div class="card__item-text">
                            <p>Rebuilding only the components that have been modified</p>
                            <p>Reusing the existing layers found in the container registry</p>
                        </div>
                    </li>
                    <li class="features-card__item">
                        <div class="card__item-title">Deployment time</div>
                        <div class="card__item-text">
                            <p>Redeploying only the components that have been modified</p>
                            <p>Instant termination of a troubled deployment process</p>
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
            <div class="demo-block__title title-h3">How it works</div>
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
      <div class="publications__title">Latest publications</div>
      <div class="publications__subtitle">Learn more about tool's features and hands-on<br />experience from new publications</div>
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
              {%- elsif publication.medium_url %}
                <a href="{{ publication.medium_url }}" class="publications__cards--link" target="_blank">
                  <span class="publications__cards--pic" style="background-image: url('{{ publication.img | true_relative_url }}')"></span>
                  <div class="publications__cards--title">{{ publication.title }}</div>
                  <div class="publications__cards--date">{{ publication.created | date: "%d-%m-%Y" }}</div>
                </a>
              {%- elsif publication.blog_url %}
                <a href="{{ publication.blog_url }}" class="publications__cards--link" target="_blank">
                  <span class="publications__cards--pic" style="background-image: url('{{ publication.img | true_relative_url }}')"></span>
                  <div class="publications__cards--title">{{ publication.title }}</div>
                  <div class="publications__cards--date">{{ publication.created | date: "%d-%m-%Y" }}</div>
                </a>
              {%- elsif publication.habr_url %}
                <a href="{{ publication.habr_url }}" class="publications__cards--link" target="_blank">
                  <span class="publications__cards--pic" style="background-image: url('{{ publication.img | true_relative_url }}')"></span>
                  <div class="publications__cards--title">{{ publication.title }}</div>
                  <div class="publications__cards--date">{{ publication.created | date: "%d-%m-%Y" }}</div>
                </a>
              {%- elsif publication.custom_urls %}
              {%- for custom in publication.custom_urls %}
                <a href="{{ custom.url }}" class="publications__cards--link" target="_blank">
                  <span class="publications__cards--pic" style="background-image: url('{{ publication.img | true_relative_url }}')"></span>
                  <div class="publications__cards--title">{{ publication.title }}</div>
                  <div class="publications__cards--date">{{ publication.created | date: "%d-%m-%Y" }}</div>
                </a>
              {% break %}
              {%- endfor -%}
              {%- elsif publication.youtube_url %}
                <a href="{{ publication.youtube_url }}" class="publications__cards--link" target="_blank">
                  <span class="publications__cards--pic" style="background-image: url('{{ publication.img | true_relative_url }}')"></span>
                  <div class="publications__cards--title">{{ publication.title }}</div>
                  <div class="publications__cards--date">{{ publication.created | date: "%d-%m-%Y" }}</div>
                </a>
              {%- elsif publication.cloud_yuga_url %}
                <a href="{{ publication.cloud_yuga_url }}" class="publications__cards--link" target="_blank">
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
              <span>Read more</span>
          </a>
      </div>
    </div>
  </div>
</div>

<div class="community" id="community">
    <div class="page__container">
        <div class="community__content">
            <div class="community__title">Friendly and rapidly growing community</div>
            <div class="community__subtitle">werf’s developers are always in contact with the community.<br/> You can reach us in Telegram, Twitter and GitHub Discussions.</div>
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
