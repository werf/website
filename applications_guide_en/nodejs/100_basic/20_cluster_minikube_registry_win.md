The command output will contain a string like the one below:

```
...
* Registry addon on with docker uses 32769 please use that instead of default 5000
...
```

Note the port number (e.g., `32769`)

Run the following port-forwarding command in a separate terminal, replacing `32769` with your port number:

```shell
docker run -ti --rm --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:host.docker.internal:32769"
```

Start the service at port 5000:

{% raw %}
```shell
kubectl -n kube-system expose rc/registry --type=ClusterIP --port=5000 --target-port=5000 --name=werf-registry --selector=actual-registry=true
```
{% endraw %}

Run the following port forwarding command in the terminal:

{% raw %}
```shell
kubectl port-forward --namespace kube-system service/werf-registry 5000
```
{% endraw %}
