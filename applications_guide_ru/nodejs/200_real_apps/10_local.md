---
title: Организация локальной разработки
permalink: nodejs/200_real_apps/10_local.html
layout: "development"
---

TODO: тру вэй - миникуб, но разработчики могут уже уметь в докер композ. И он может быть тупо удобнее или есть какие-то наработки.

???

Тут же — про отладку






можно сделать werf run --follow и при коммите делается перезапуск (аналогично он есть у конверджа и компоуза)





                  
version: "3.8"
services:
  web:
    image: "$WERF_IMAGE_B_NAME"




werf compose up --help




$ werf build --report-format=envfile --report-path=.env
$ cat .env
WERF_APP_DOCKER_IMAGE_NAME=REPO:TAG
$ docker-compose --env-file .env





