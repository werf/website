{% include "installation/_ready_to_use_images.md" %}

## Пример использования

Запуск werf в контейнере выглядит следующим образом:

```shell
docker run \
    --security-opt seccomp=unconfined \
    --security-opt apparmor=unconfined \
    registry.werf.io/werf/werf:{{ include.version }}-{{ include.channel }} werf help
```

{% include "installation/_buildah_settings.md" %}

{% include "installation/_buildah_troubleshooting.md" %}