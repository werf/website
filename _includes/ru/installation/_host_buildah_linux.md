Установить пакет предоставляющий `buildah` и его зависимости для вашего дистрибутива.

{% offtopic title="Отсутствует пакет buildah" %}

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