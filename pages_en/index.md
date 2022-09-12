---
title: Giterministic CLI tool
permalink: /
layout: default
sidebar: none
---

<div class="news-and-updates">
    <div class="news-and-updates__block">
        <div class="news-and-updates__title">
            Development news
        </div>
        <div class="news-and-updates__item" id="nau-news"></div>
    </div>
    <div class="news-and-updates__block">
        <div class="news-and-updates__title">
            Latest releases
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
                <div class="intro__subtitle">What is it?</div>
                <h1 class="intro__title">Consistent<br/>delivery tool</h1>
                <ul class="intro__features">
                    <li>The CLI tool gluing Git, Docker, Helm & Kubernetes<br/>
                    with any CI system to implement CI/CD and Giterminism.</li>
                </ul>
                <div class="intro__btns page__btn-group">
                    <a href="{{ "how_it_works.html" | true_relative_url }}" target="_blank" class="page__btn page__btn_b page__btn_small">
                        How it works
                    </a>
                    <a href="{{ "documentation/quickstart.html" | true_relative_url }}" target="_blank" class="page__btn page__btn_b page__btn_small">
                        Quickstart
                    </a>
                    <a href="{{ "documentation/index.html" | true_relative_url }}" target="_blank" class="page__btn page__btn_b page__btn_small">
                        Documentation
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
                <div class="intro__subtitle">Why do I need it?</div>
                <h1 class="intro__title">Fast & efficient<br/>CI/CD</h1>
                <ul class="intro__features">
                    <li>Establish and benefit from efficient,<br/>
                    robust and integrated CI/CD pipelines<br/>
                    on top of proven technologies.</li>
                    <li>With werf, it’s easy to start, to apply best practices<br/>
                    and to avoid reinventing the wheel.</li>
                </ul>
            </div>
        </div>
    </div>
</div>

<div class="presentation" id="presentation">
    <div class="page__container presentation__container">
        <div class="presentation__row">
            <div class="presentation__row-item" id="presentation-title">
                <div class="presentation__subtitle">How does it work?</div>
                <h1 class="presentation__title">What you Git<br/> is what you get!</h1>
                <ul class="presentation__features">
                    <li>Git is treated as a single source of truth<br/>
                    (we call it «Giterminism»).</li>
                    <li>werf not only builds & deploys but also continuously<br/>
                    syncs current Kubernetes state with changes made in Git.</li>
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
                Desired state<br>
                is defined in Git
            </div>
            <div class="presentation-notes__item-text"></div>
        </div>
        <div class="presentation-notes__item" id="presentation-notes-2">
            <div class="presentation-notes__item-num">2</div>
            <div class="presentation-notes__item-title">
                Syncs the registry<br>
                to the defined state
            </div>
            <div class="presentation-notes__item-text">
                <ol>
                    <li>
                        Builds images (if anything changed or missing)
                    </li>
                    <li>
                        Pushes images (if needed)
                    </li>
                </ol>
            </div>
        </div>
        <div class="presentation-notes__item" id="presentation-notes-3">
            <div class="presentation-notes__item-num">3</div>
            <div class="presentation-notes__item-title">
                Syncs Kubernetes<br>
                to the defined state
            </div>
            <div class="presentation-notes__item-text">
                <ol>
                    <li>
                        Applies the Helm chart (with appropriate images)
                    </li>
                    <li>
                        Tracks deployment progress till success (and provides feedback)
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
                <div class="features__list-item-title">Two ways to deploy</div>
                <div class="features__list-item-text">
                    werf supports 2 ways to deploy an application:
                    <ol>
                        <li><b>converge</b> application from git commit into the Kubernetes;</li>
                        <li><b>publish</b> application from git commit into the Container Registry as a <b>bundle</b>, then <b>deploy bundle</b> into the Kubernetes.</li>
                    </ol>
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_usage"></div>
                <div class="features__list-item-title">Flexible integration</div>
                <div class="features__list-item-text">
                    werf allows usage in a different ways:
                    <ol>
                        <li>manually;</li>
                        <li>by CI/CD system;</li>
                        <li>by Kubernetes operator <small>(coming soon)</small>;</li>
                        <li>by heroku-like git push approach <small>(coming soon)</small>.</li>
                    </ol>
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_lifecycle"></div>
                <div class="features__list-item-title">It’s a CLI tool</div>
                <div class="features__list-item-text">
                    werf is not a SaaS, it is an Open Source, self-contained client-side CLI&nbsp;tool. Use this single tool for <b>local development</b> or <b>embed it into any CI/CD system</b> (like GitLab CI/CD, GitHub Actions, Jenkins, CircleCI, etc.) using its main commands as building blocks:
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
                <div class="features__list-item-title">Easy to use</div>
                <div class="features__list-item-text">
                    werf just works out of the box with a minimal configuration. You don't even need to be a DevOps/SRE engineer to use werf. <a href="/guides.html"><b>Many guides</b></a> are provided to quickly deploy your app into Kubernetes.
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_config"></div>
                <div class="features__list-item-title">Combining the best</div>
                <div class="features__list-item-text">
                    werf glues well-established software forming a transparent, <b>integrated CI/CD platform</b>. Benefit from a conveniently controlled, smooth interaction of Git, Docker, your container registry &amp; existing CI system, Helm, Kubernetes!
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_fast"></div>
                <div class="features__list-item-title">Distributed building</div>
                <div class="features__list-item-text">
                    werf implements an advanced builder boasting a distributed algorithm that <b>makes your pipelines really fast</b> thanks to distributed caching.
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_debug"></div>
                <div class="features__list-item-title">Built-in cleaning</div>
                <div class="features__list-item-text">
                    werf implements a sophisticated algorithm to <b>clean unused Docker images</b> that is based on the Git history of your application.
                </div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_helm"></div>
                <div class="features__list-item-title">Extended Helm</div>
                <div class="features__list-item-text">
                    werf uses a built-in <code>helm</code> binary to implement a Helm-compatible deployment with additional features. It doesn't require to have <code>helm</code> separately installed. It provides you a descriptive and sharp <b>deploy logging</b>, fast <b>failures detection</b> during deploy process, secrets support and other extensions making deploy process <b>clear, robust and reliable</b>.
                </div>
            </li>
            <li class="features__list-item features__list-item_special">
                <div class="features__list-item-title">Open Source</div>
                <div class="features__list-item-description">
                    <a href="https://github.com/werf/werf" target="_blank">Open Source</a> project since its launch in 2016. Written in Go. Proud of its strong &amp; growing community of users.
                </div>
            </li>
        </ul>
    </div>
</div>

<div class="stats">
    <div class="page__container">
        <div class="stats__content">
            <div class="stats__title">Active development & adoption</div>
            <ul class="stats__list">
                <li class="stats__list-item">
                    <div class="stats__list-item-num">4</div>
                    <div class="stats__list-item-title">releases per week</div>
                    <div class="stats__list-item-subtitle">on average during the last year</div>
                </li>
                <li class="stats__list-item">
                    <div class="stats__list-item-num">2000+</div>
                    <div class="stats__list-item-title">installations</div>
                    <div class="stats__list-item-subtitle">for large and small projects</div>
                </li>
                <li class="stats__list-item">
                    <div class="stats__list-item-num gh_counter">2100</div>
                    <div class="stats__list-item-title">stars on GitHub</div>
                    <div class="stats__list-item-subtitle">let’s make it more ;)</div>
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
                    werf is a mature,<br>
                    reliable tool you can trust
                </div>
                <a href="{{ "about/release_channels.html" | true_relative_url }}" class="page__btn page__btn_b page__btn_small page__btn_inline">
                    Read about release channels
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

{%- assign publications = site.data.publications.articles | sort: 'created' | reverse %}
{%- assign publications_by_year = publications | group_by_exp: "publication", "publication.created | date: '%Y'" %}

<div class="publications">
  <div class="page__container">
    <div class="publications__content">
      <div class="publications__title">Latest publications about werf</div>
      <div class="publications__subtitle">Learn more about tool's features and hands-on<br />experience from new publications</div>
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
            </li>
          {%- endfor %}
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

<div class="page__container">
    <div class="documentation">
        <div class="documentation__image">
        </div>
        <div class="documentation__info">
            <div class="documentation__info-title">
                Detailed documentation
            </div>
            <div class="documentation__info-text">
                werf documentation comprises over 100 articles covering typical use cases (getting started, deploying to Kubernetes, CI/CD integration, and more), CLI, commands, and providing a thorough description of functions & architecture.
            </div>
        </div>
        <div class="documentation__btns">
            <a href="{{ "how_it_works.html" | true_relative_url }}" class="page__btn page__btn_b documentation__btn">
                How it works
            </a>
            <a href="{{ "documentation/quickstart.html" | true_relative_url }}" class="page__btn page__btn_o documentation__btn">
                Quickstart
            </a>
            <a href="/guides.html" class="page__btn page__btn_o documentation__btn">
                Guides
            </a>
        </div>
    </div>
</div>
