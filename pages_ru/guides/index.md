---
title: "Самоучитель по Kubernetes и werf: научись K8s и деплоить в него"
layout: landing
permalink: index.html
---

<div class="landing">
    <div class="landing__header">
        <div class="landing__container">
            <a href="/" class="landing__header-title" data-proofer-ignore>
                {% asset werf-logo.svg alt="werf" width="167" height="70" %}
            </a>
            <a href="{{ site.site_urls['ru'] }}" class="landing__button" data-proofer-ignore>
                {% asset arrow.svg %}
                <span>ru.werf.io</span>
            </a>
        </div>
    </div>
    <div class="landing__content">
        <div class="landing__container">
            <div class="landing__section landing__section_first" data-sm-trigger="intro">
                <h1 class="landing__h1">
                    Быстро и&nbsp;просто<br>
                    деплоим в&nbsp;Kubernetes<br>
                    с&nbsp;помощью <b>werf</b>
                </h1>
                <div class="landing__text">
                    Мы разработали бесплатный самоучитель, который призван помочь разобраться,<br>
                    как&nbsp;подготовить свои приложения к&nbsp;Kubernetes и обеспечить их&nbsp;непрерывный деплой,<br>
                    реализуя CI/CD с&nbsp;помощью Open Source-инструментов.
                </div>
            </div>
            <div class="landing__section" data-sm-trigger="plan">
                <h1 class="landing__h2">
                    Что меня ждёт?
                </h1>
                <div class="landing__text">
                    <p>Самоучитель — достаточно подробное руководство, сочетающее <b>теорию и практику</b> разработки (Dev) и эксплуатации (Ops).</p>
                    <p>Учебные материалы ориентированы на разработчиков, стремящихся обрести базовые навыки DevOps-инженеров по организации непрерывной доставки приложений в Kubernetes. А также они будут полезны для DevOps-инженеров, которые хотят эффективнее решать свои задачи.</p>
                    <p>Последовательно рассматриваются <b>основы работы с Kubernetes</b>, <b>вопросы адаптации приложений для Kubernetes</b> (их «кубернетизация») и основные задачи, возникающие при&nbsp;разработке сервисов и организации CI/CD (сборка, деплой, работа с&nbsp;зависимостями и&nbsp;ассетами, с&nbsp;базами данных и файлами, хранение секретной и несекретной конфигурации приложения), а также <b>лучшие практики деплоя</b> в&nbsp;K8s.</p>
                    <p>В самоучителе учтена специфика языков/фреймворков и приложены примеры исходного кода приложения и инфраструктуры (IaC).</p>
                    <p><i>Самоучитель рассчитан на&nbsp;«ванильный» Kubernetes-кластер, но его должно быть просто адаптировать под&nbsp;кастомизированные сборки.</i></p>
                </div>
            </div>
            <div class="landing__section" data-sm-trigger="todo">
                <h1 class="landing__h2">
                    Что я получу?
                </h1>
                <div class="landing__text">
                    <p>Освоить навыки по самоучителю — непростая задача. Потребуется выделить время, приложить усилия и активировать внимание.</p>
                    <p>Те, кто найдёт в себе силы дойти до конца, смогут:</p>
                </div>
                <ul class="landing__list">
                    <li>
                        Стать более востребованными специалистами
                        <span>(в области пересечения разработки и эксплуатации, что обычно называют «DevOps-инженерами»)</span>
                    </li>
                    <li>
                        Принести в свою компанию Kubernetes
                    </li>
                    <li>
                        Научиться лучше решать задачи в текущей компании
                    </li>
                </ul>
                <div class="landing__text">
                    <p>Самоучитель поможет сфокусироваться только на том, что нужно именно вам, а растущее <a href="https://t.me/werf_ru">сообщество пользователей werf</a> — облегчит преодоление сложностей на пути.</p>
                </div>
            </div>
            <div class="landing__section" data-sm-trigger="select">
                <h1 class="landing__h2">
                    Начать погружение!
                </h1>

                <div class="landing__text">
                  <p>Выберите наиболее близкий вам фреймворк:</p>
                </div>

                {% include common/landing_tiles.html %}
            </div>
            <div class="landing__section">
                <div class="landing__credits">
                    Сделано с <span class="icon">любовью</span> в компании <a href="https://flant.ru/" target="_blank">Флант</a>
                </div>
            </div>
        </div>
    </div>
    <div class="landing__side">
        <div class="landing__captain" data-sm-captain="intro">
            {% asset landing/parrot.png class="landing__captain-img landing__captain-img_intro" %}
            {% asset landing/parrot2.png class="landing__captain-img landing__captain-img_plan" %}
            {% asset landing/parrot3.png class="landing__captain-img landing__captain-img_todo" %}
            {% asset landing/parrot4.png class="landing__captain-img landing__captain-img_select" %}
        </div>
    </div>
</div>
