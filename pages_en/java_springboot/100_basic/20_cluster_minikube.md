## Minikube

Minikube is a lightweight Kubernetes version that can run on a local machine.

### Installation

[Install minikube](https://minikube.sigs.k8s.io/docs/start/) and run it:

```shell
minikube start --driver=docker
```

**IMPORTANT:** You need to make sure that Minikube uses a docker driver if it is already running on your system. Otherwise, you have to restart Minikube using the `minikube delete` command and start it using the command provided above.

The `.kube/config` file containing keys for connecting to the local cluster will be automatically created, allowing werf to connect to the local registry.

As a result, you must be able to access the cluster using the `kubectl` utility (you may need to install this utility separately). For example, the command below:

```shell
kubectl get ns
```

… will print a list of all namespaces in the cluster (and not an error message).

### Ingress

You need to enable the appropriate [addon](https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/#enable-the-ingress-controller) in Minikube.

{% offtopic title="How do I verify that Ingress operates as expected?" %}

You can assume that everything works fine by default. If you have difficulties in understanding the points below, please, return to this sub-chapter later.

- The load balancer is installed.
- The pod containing the load balancer is up and running fine (use this command for checking: kubectl -n ingress-nginx get po).
- The application is listening on port 80 (check it with lsof -n | grep LISTEN).

{% endofftopic %}

### Registry

Enable the minikube registry addon:

```shell
minikube addons enable registry
```

And then perform the following actions depending on your operating system:

{% offtopic title="Windows" %}
{% include_relative 20_cluster_minikube_registry_win.md %}
{% endofftopic %}
{% offtopic title="MacOS" %}
{% include_relative 20_cluster_minikube_registry_macos.md %}
{% endofftopic %}
{% offtopic title="Linux" %}
{% include_relative 20_cluster_minikube_registry_linux.md %}
{% endofftopic %}

### Hosts

This tutorial assumes that the cluster (in fact, its Nginx Ingress controller) is located at `example.com`, and its registry is accessible at `registry.example.com`. We will be using this domain and its subdomains in the configs. In the case of using different addresses, adjust the configuration accordingly.

Enter the following lines in `/etc/hosts`:

```
127.0.0.1           example.com
127.0.0.1           registry.example.com
```

### Authorizing in the registry

The Minikube registry does not require authorization. No further action is required.
