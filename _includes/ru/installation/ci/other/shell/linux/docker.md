## Требования

* CI-система;

* Linux-хост для запуска CI-задач, имеющий:
  
  * Shell Runner вашей CI-системы;
  - Bash;
  
  - Git версии 2.18.0 или выше;
  
  - GPG;
  
  - [Docker Engine](https://docs.docker.com/engine/install/).

## Установка werf

Для установки werf, на хосте для запуска CI-задач выполните:

```shell
curl -sSL https://werf.io/install.sh | bash -s -- --ci
```

## Конфигурация CI/CD проекта

```yaml
# .pseudo-ci-config.yml:
before_every_job:
- source "$(~/bin/trdl use werf {{ include.version }} {{ include.channel }})"
- source "$(werf ci-env gitlab --as-file)"

jobs:
  prod:
    commands:
    - werf converge
    environment: prod
    on: master
    how: manually

  images:cleanup:
    commands:
    - werf cleanup
    on: master
    how: daily
```

TODO: больше листингов

TODO: настройка очистки (авторизация, ...)
