---
title: Быстрый старт разработчика
permalink: rails/100_basic.html
---

В этой главе мы покажем, как с помощью werf собрать образ с приложением и запустить контейнер на основе собранного образа локально. После этого мы развернём локальный тестовый кластер Kubernetes и покажем, как разворачивать приложение уже в нём.

## Используемое приложение

В качестве тестового приложения мы будем использовать [простое REST JSON API на Ruby on Rails](https://github.com/werf/werf-guides/tree/master/examples/rails/000_app), хранящее данные в SQLite. Для упрощения задачи, мы не будем обеспечивать сохранность данных в SQLite: при перевыкате приложения данные будут удаляться. Вопросы хранения данных будут рассмотрены в главе "Полноценные приложения".

<div id="go-forth-button">
    <go-forth url="100_basic/10_build.html" label="Сборка образа" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
