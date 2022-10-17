{% offtopic title="Активация экспериментального Buildah-backend'а в werf" %}
Если вы хотите запускать werf в контейнерах/Kubernetes, то следуйте этим инструкциям для [Docker](https://ru.werf.io/documentation/v{{ include.version }}/advanced/ci_cd/run_in_container/use_docker_container.html) или [Kubernetes](https://ru.werf.io/documentation/v{{ include.version }}/advanced/ci_cd/run_in_container/use_kubernetes.html). Если же вы хотите запускать werf в Buildah-режиме вне контейнеров или хотите собрать образы с запакованным в них werf в Buildah-режиме с нуля, то:
* Если ваше ядро Linux версии 5.13+ (в некоторых дистрибутивах 5.11+), то убедитесь, что модуль ядра `overlay` загружен с `lsmod | grep overlay`. Если ядро более старое или у вас не получается активировать модуль ядра `overlay`, то установите `fuse-overlayfs`, который обычно доступен в репозиториях вашего дистрибутива.
* Команда `sysctl user.max_user_namespaces` должна вернуть по меньшей мере `15000`. В ином случае выполните:
  ```shell
  echo 'user.max_user_namespaces = 15000' | sudo tee -a /etc/sysctl.conf
  sudo sysctl -p
  ```
* В ядре должна быть включена опция CONFIG_USER_NS. Зачастую это можно проверить так:
  ```shell
  grep CONFIG_USER_NS /boot/config-<текущая версия ядра>
  ```
* Если файлы `/etc/subuid` и `/etc/subgid` не существуют, то, в большинстве дистрибутивов, вам потребуется установить пакет, который их создаст. Текущий пользователь должен иметь по крайней мере `65536` выделенных для него subordinate UIDs/GIDs — это может выглядеть как строка вида `current_username:1000000:65536` в `/etc/subuid` и `/etc/subgid`. Если в этих файлах нет подобной строки, вам потребуется добавить её самостоятельно. Изменение этих файлов может потребовать перезагрузки. Подробнее: `man subuid`, `man subgid`.
* Путь `~/.local/share/containers` должен быть создан и у текущего пользователя должны иметься права на чтение и запись в него.
* Для использования `werf` вне контейнеров установите один из `crun`/`runc`/`kata`/`runsc`, обычно доступные в репозиториях вашего дистрибутива.
* Установите пакет, предоставляющий программы `newuidmap` и `newgidmap`.

Теперь активируйте Buildah backend и попробуйте запустить сборку вашего проекта:
```shell
export WERF_BUILDAH_MODE=auto
werf build
```

Если возникли ошибки при запуске werf:
* Исправьте права на `newuidmap` и `newgidmap`:
  ```shell
  sudo setcap cap_setuid+ep /usr/bin/newuidmap
  sudo setcap cap_setgid+ep /usr/bin/newgidmap
  sudo chmod u-s,g-s /usr/bin/newuidmap /usr/bin/newgidmap
  ```
* При проблемах с OverlayFS попробуйте VFS:
  ```shell
  export WERF_BUILDAH_STORAGE_DRIVER=vfs
  ```
* При запуске werf вне контейнеров можно переключиться с rootless-режима на менее изолированный chroot-режим:
  ```shell
  export WERF_BUILDAH_MODE=native-chroot
  ```
{% endofftopic %}
