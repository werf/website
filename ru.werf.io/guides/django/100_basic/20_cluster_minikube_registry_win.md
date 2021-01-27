Вывод команды должен содержать похожую строку:

```
...
* Registry addon on with docker uses 32769 please use that instead of default 5000
...
```

Запомните порт (например, `32769`).

Запустите следующий проброс портов в отдельном терминале, заменив порт `32769` вашим портом:

```shell
docker run -ti --rm --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:host.docker.internal:32769"
```

Запустите сервис с привязкой к порту 5000:

{% raw %}
```shell
kubectl -n kube-system expose rc/registry --type=ClusterIP --port=5000 --target-port=5000 --name=werf-registry --selector=actual-registry=true
```
{% endraw %}

Запустите следующий проброс портов в отдельном терминале:

{% raw %}
```shell
kubectl port-forward --namespace kube-system service/werf-registry 5000
```
{% endraw %}
