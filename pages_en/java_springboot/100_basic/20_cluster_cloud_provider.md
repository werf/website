## Choosing a cloud provider

There is a large number of [Managed Kubernetes](https://www.google.com/search?q=managed+kubernetes) offerings. Some of them include server capacity, while others don't.

The easiest way is to take advantage of AWS (EKS) or Google (GKE) offerings.  During the initial registration, they provide a bonus large enough for working with a cluster for a couple of weeks. However, you will need to enter your bank card details to register.

Yandex.Cloud is a viable alternative: the Managed Kubernetes service is available here as well, and you also must provide your card details. However, you can use some variation of a local bank card (such as the YouMoney card that allows for subscription payments within Russia).

You can also try to deploy Kubernetes yourself to the Hetzner cloud using their [tutorial](https://community.hetzner.com/tutorials/install-kubernetes-cluster) (Hetzner is one of the cheapest cloud providers out there). However, you must realize that you will have to deal with a large volume of platform administration tasks. You must find a reliable provider or support team for your production environment. 
If you have never installed Kubernetes yourself and/or do not have experience in system administration, then you should not try to master this rather voluminous topic as part of this tutorial.

Regardless of the cloud provider chosen, you will need to get the keys for accessing the Kubernetes cluster in the form of `.kube/config` file.

<!--  .kube/config -->
{% offtopic title="What is .kube/config?" %}
This file stores the cluster access details. Using this file (located at `~/.kube/config`), tools that work with Kubernetes can connect to the cluster.

Here is the example of that file:

{% snippetcut name=".kube/config" url="#" %}
{% raw %}
```yaml
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ALotOfNumbersAndLettersAndSoOnAveryVERYveryLongStringInBase64=
    server: https://127.0.0.1:6445
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: ManyNumbersAndLettersAVERYveryVerylongString=
    client-key-data: ManyLettersAndNumbersAveryVeryVERYlongString=
```
{% endraw %}
{% endsnippetcut %}

Note that in addition to access keys, there is also a path to the Kubernetes API (e.g., `server: https://127.0.0.1:6445`). Important: this path must be correct! Some `.kube/config`-generating utilities do not insert a public IP into the file. In this case, you need to adjust it manually.

{% endofftopic %}
<!-- / .kube/config -->

### Checking the operability and access to the cluster

Check your cluster using the checklist provided at the beginning of the article.

As a result, you must be able to access the cluster using the `kubectl` utility (you may need to install this utility separately). For example, the command below:

```shell
kubectl get ns
```

… will print a list of all namespaces in the cluster (and not an error message).

### Ingress

You also have to install Nginx Ingress in your cluster. There are [detailed instructions](https://kubernetes.github.io/ingress-nginx/deploy/) available for all major providers on the NGINX Ingress Controller's website.

### Registry

In most cases, cloud providers provide a registry as one of the services (however, this is not always the case).

{% offtopic title="My provider doesn't provide registry, what should I do?" %}

You can deploy Docker Registry by yourself on a separate virtual machine or use some cloud-based solution. werf supports [various implementations]({{ site.url }}/documentation/advanced/supported_registry_implementations.html) of container registries.

{% endofftopic %}

Ensure the network connectivity between the cluster and the registry: again, in most cases, the cloud provider solves this problem on its own or provides specific manuals on the subject.

### Hosts

This tutorial assumes that the cluster (in fact, its Nginx Ingress controller) is located at `example.com`, and its registry is accessible at `registry.example.com`. You will have to choose and configure the necessary DNS records yourself and substitute your domains in the code throughout this tutorial.

### Authorizing in the registry

You need to log in on the local machine using the `docker login` command so that werf can push assembled images to the registry:

```shell
docker login registry.example.com -u <account_login> -p <account_password>
```

where `<account_login>` is the login to the registry, and `<account_password>` is the password.
