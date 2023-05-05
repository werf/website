{% offtopic title="Подготовка системы к использованию werf с бэкендом Buildah" %}
1. Установить пакет предоставляющий `buildah` и его зависимости.
2. Активировать buildah бэкенд переменной окружения `WERF_BUILDAH_MODE=auto`:

```shell
export WERF_BUILDAH_MODE=auto
werf build
```
{% endofftopic %}
