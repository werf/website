## Docker Desktop

### Installing

Install Docker Desktop for [Windows](https://docs.docker.com/docker-for-windows/install/) or [macOS](https://docs.docker.com/docker-for-mac/install/).

Start Kubernetes and configure the resources allocated to it.

{% offtopic title="Why are resources important?" %}

Lack of resources in the cluster may prevent your application, Ingress, or even some system components required by the orchestrator from starting. You must be really experienced in cluster administration to figure out what is happening.

The authors of this tutorial tested it using the following configuration: 6 CPU cores, 6 GB RAM, 1.5 GB of swap space, 24 GB of disk space. 

{% endofftopic %}

The `.kube/config` file containing keys for connecting to the local cluster will be automatically created, allowing werf to connect to the local registry.

As a result, you must be able to access the cluster using the `kubectl` utility (you may need to install it separately). For example, the command below:

```shell
kubectl get ns
```

… will print the list of all namespaces in the cluster (and not an error message).

### Ingress

You have to manually [install the Nginx Ingress controller](https://kubernetes.github.io/ingress-nginx/deploy/).

In the case of Docker Desktop, the Ingress controller may be inaccessible in some cases. That is because ports are not forwarded to the host machine. To make sure that Ingress works as intended, check that:

- The pod containing the Ingress controller is up and running fine (use this command for checking in the case of nginx-ingress: `kubectl -n ingress-nginx get po`).
- Kubernetes is listening on port 80 (check it with `lsof -n | grep LISTEN` or [similar way](https://www.google.com/search?q=check+used+ports&oq=check+used+ports))

If there is no Kubernetes on port 80 (HTTP), you may need to forward the port manually. For this, find out the name of the Ingress controller's pod:

```shell
kubectl -n ingress-nginx get po
```

… and then run port forwarding:

```shell
kubectl port-forward --address 0.0.0.0 -n ingress-nginx ingress-nginx-controller-<random_alphanumeric_sequence> 80:80
```

After running the above commands, verify that an application is listening on port 80 (`lsof -n | grep LISTEN`).

### Registry

Docker Desktop does not have a built-in registry. The easiest way is to run it manually as a Docker image:

```shell
docker run -d -p 5000:5000 --restart=always --name registry registry:2
```

Note that in this case, the registry does not use SSL encryption by default. Thus, you have to add the `--insecure-registry=true` [parameter](https://werf.io/documentation/reference/cli/werf_managed_images_add.html#options) when using werf.

### Hosts

This tutorial assumes that the cluster (in fact, its Nginx Ingress controller) is located at `example.com`, and its registry is accessible at `registry.example.com`. We will be using this domain and its subdomains in the configs. If you are using different addresses, adjust the configuration accordingly.

Insert the following lines in `/etc/hosts`:

```
127.0.0.1           example.com
127.0.0.1           registry.example.com
```

### Authorizing in the registry

When running the registry in a local mode (as shown in the example above), no password is necessary, and you do not need to do anything.

{% offtopic title="I use an external registry or username/password" %}
You need to log in to the local machine using the `docker login` command for werf to be able to push assembled images to the registry:

```shell
docker login <registry_domain> -u <account_login> -p <account_password>
```
{% endofftopic %}
