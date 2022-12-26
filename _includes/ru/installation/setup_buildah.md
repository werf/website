{% offtopic title="Подготовка системы к использованию werf с бэкендом Buildah" %}
1. Установить пакет предоставляющий `buildah` и его зависимости.
2. Активировать buildah бэкенд переменной окружения `WERF_BUILDAH_MODE=auto`:

```shell
export WERF_BUILDAH_MODE=auto
werf build
```

### Устранение проблем

#### Режим overlayfs

Если ваше ядро Linux версии 5.13+ (в некоторых дистрибутивах 5.11+) **рекомендуется** режим работы через модуль ядра `overlay`:
- Убедитесь, что модуль ядра `overlay` загружен с `lsmod | grep overlay`.
- Убедитесь, что настройка ядра `CONFIG_USER_NS=y` включена в вашем ядре с помощью `grep CONFIG_USER_NS /boot/config-VERSION`.
- При использовании ядра в debian-системах команда `sysctl kernel.unprivileged_userns_clone` должна вернуть `1`. В ином случае выполните:
    ```shell
    echo 'kernel.unprivileged_userns_clone = 1' | sudo tee -a /etc/sysctl.conf
    sudo sysctl -p
    ```
- Команда `sysctl user.max_user_namespaces` должна вернуть по меньшей мере `15000`. В ином случае выполните:
    ```shell
    echo 'user.max_user_namespaces = 15000' | sudo tee -a /etc/sysctl.conf
    sudo sysctl -p
    ```

Если ядро более старое или у вас не получается активировать модуль ядра `overlay`, то установите `fuse-overlayfs`, который обычно доступен в репозиториях вашего дистрибутива. В крайнем случае может быть использован драйвер хранилища `vfs` с помощью установки переменной окружения `WERF_BUILDAH_STORAGE_DRIVER=vfs`.

#### Отсутствует пакет buildah

При отсутствии пакета buildah для вашего дистрибутива необходимы следующие действия:
- Установить пакеты предоставляющие: `slirp4netns`, `newuidmap` и `newgidmap` (`uidmap`). Удостоверится что `newuidmap` и `newgidmap` имеют корректные права:
    ```shell
    sudo setcap cap_setuid+ep /usr/bin/newuidmap
    sudo setcap cap_setgid+ep /usr/bin/newgidmap
    sudo chmod u-s,g-s /usr/bin/newuidmap /usr/bin/newgidmap
    ```
-  Если файлы `/etc/subuid` и `/etc/subgid` не существуют, то, в большинстве дистрибутивов, вам потребуется установить пакет, который их создаст. Текущий пользователь должен иметь по крайней мере `65536` выделенных для него subordinate UIDs/GIDs — это может выглядеть как строка вида `current_username:1000000:65536` в `/etc/subuid` и `/etc/subgid`. Если в этих файлах нет подобной строки, вам потребуется добавить её самостоятельно. Изменение этих файлов может потребовать перезагрузки. Подробнее: `man subuid`, `man subgid`.
- Путь `~/.local/share/containers` должен быть создан и у текущего пользователя должны иметься права на чтение и запись в него.

{% endofftopic %}
