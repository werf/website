Вывод команды должен содержать похожую строку:

```
...
* Registry addon on with docker uses 32769 please use that instead of default 5000
...
```

Запоминаем порт (например, `32769`).

Запустите следующий проброс портов в отдельном терминале, заменив порт `32769` вашим портом:

```shell
docker run -ti --rm --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:host.docker.internal:32769"
```

Запустите следующий проброс портов в отдельном терминале, заменив порт `32769` вашим портом:

```shell
brew install socat
socat TCP-LISTEN:5000,reuseaddr,fork TCP:host.docker.internal:32769
```