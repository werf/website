Скачайте установщик werf:

```shell
curl -sSLO https://werf.io/install.sh && chmod +x install.sh
```

Для использования на рабочей машине установите werf и настройте его автоматическую активацию (после чего откройте новую shell-сессию):

```shell
./install.sh --version {{ include.version }} --channel {{ include.channel }}
```

Для использования werf в CI установите werf и активируйте его вручную:

```shell
./install.sh --ci
source "$(~/bin/trdl use werf {{ include.version }} {{ include.channel }})"
```

Список опций установщика:

```shell
./install.sh --help
```

После активации `werf` должен быть доступен в той же shell-сессии, в которой он был активирован:

```shell
werf version
```