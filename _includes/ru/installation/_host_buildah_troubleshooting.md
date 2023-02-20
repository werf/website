## Устранение проблем

### Режим overlayfs

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

{% include "installation/__buildah_troubleshooting.md" %}