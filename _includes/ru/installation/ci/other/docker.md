## Требования

- CI-система;

- Linux-хост для запуска CI-задач, имеющий:
  
  * Docker Runner вашей CI-системы;
  
  * [Docker Engine](https://docs.docker.com/engine/install/).

## Настройка Runner'а

На хосте для запуска CI задач создайте volume `werf`:

```shell
docker volume create werf
```

Настройте Runner вашей CI-системы так, чтобы создаваемые контейнеры имели следующие параметры:

* `--security-opt seccomp:unconfined`

* `--security-opt apparmor:unconfined`

* `--volume werf:/home/build/.werf`

Если хост для запуска CI-задач имеет версию ядра Linux 5.12 или ниже, то установите на хост `fuse` и настройте Runner так, чтобы создаваемые контейнеры имели дополнительный параметр `--device /dev/fuse`.

## Конфигурация CI/CD проекта

```yaml
# .pseudo-ci-config.yml:
image: "registry.werf.io/werf/werf:{{ include.version }}-{{ include.channel }}"
image_pull_policy: always

before_every_job:
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
