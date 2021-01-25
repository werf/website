The command output will contain a string like this:
```
...
* Registry addon on with docker uses 32769 please use that instead of default 5000
...
```

Note the port number (e.g., `32769`).

Run the following port-forwarding command in a separate terminal, replacing `32769` with the number of your port: 

```shell
docker run -ti --rm --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:host.docker.internal:32769"
```

Run the following port-forwarding command in a separate terminal, replacing `32769` with the number of your port: 

```shell
brew install socat
socat TCP-LISTEN:5000,reuseaddr,fork TCP:host.docker.internal:32769
```
