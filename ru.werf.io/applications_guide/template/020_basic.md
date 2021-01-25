---
title: Базовые настройки
sidebar: applications_guide
guide_code: ____________
permalink: ____________/020_basic.html
toc: false
---

В этой главе мы возьмём приложение, которое будет выводить сообщение «Hello world!» по HTTP, и опубликуем его в Kubernetes с помощью werf. Сперва разберёмся со сборкой и добьёмся того, чтобы образ оказался в Registry, затем — разберёмся с деплоем собранного приложения в Kubernetes, после чего, наконец, организуем CI/CD-процесс силами GitLab CI.

Если у вас мало опыта с описанием объектов Kubernetes — учтите, что это может занять у вас больше времени, чем организация сборки и CI-процесса с помощью werf. Это нормально, и мы постарались облегчить данный процесс.

Приложение будет состоять из одного Docker-образа, собранного с помощью werf. В этом образе будет работать один основной процесс, который запустит ____________. Управлять маршрутизацией запросов к приложению будет Ingress в Kubernetes кластере. Мы реализуем два стенда: [production]({{ site.docsurl }}/documentation/reference/ci_cd_workflows_overview.html#production) и [staging]({{ site.docsurl }}/documentation/reference/ci_cd_workflows_overview.html#staging).

<div id="go-forth-button">
    <go-forth url="020_basic/10_build.html" label="Сборка" base-url="applications_guide"></go-forth>
</div>
