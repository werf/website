## I have a cluster already

If you already have a cluster, then you probably have some experience in configuring it. Check your cluster using the checklist provided at the beginning of the article.

Regardless of how the cluster is installed, you will need to get the keys for accessing the Kubernetes cluster in the form of `.kube/config file`.

<!--  .kube/config -->
{% offtopic title="What is .kube/config?" %}
This file stores the cluster access details. Using this file (located at `~/.kube/config`), tools that work with Kubernetes can connect to the cluster. 

Here is the example:

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

Note that in addition to connection keys, there is also a path to the Kubernetes API (e.g., `server: https://127.0.0.1:6445`). Important: this path must be correct! Some `.kube/config`-generating utilities do not insert a public IP into the file. In this case, you need to adjust it manually.

{% endofftopic %}
<!-- / .kube/config -->

### Checking the operability and accessibility of the cluster

Check your cluster using the checklist provided at the beginning of the article.

As a result, you must be able to access the cluster using the `kubectl` utility (you may need to install it separately). For example, the command below:

```bash
kubectl get ns
```

… must print a list of all namespaces in the cluster (and not an error message).

### Ingress

Nginx Ingress must be installed in your cluster. The installation guide is available in the [Kubernetes documentation](https://kubernetes.github.io/ingress-nginx/deploy/).

### Registry

werf supports [various registry implementations]({{ site.docsurl }}/documentation/advanced/supported_registry_implementations.html).

Make sure that there is network connectivity between the cluster and the registry.

### Hosts

This tutorial assumes that the cluster (in fact, its Nginx Ingress controller) is located at `example.com`, and its registry is accessible at `registry.example.com`. You will have to choose and configure the necessary DNS records yourself and substitute your domains in the code throughout this tutorial.

### Authorizing in the registry

You need to log in to the local machine using the `docker login` command so that werf can push assembled images to the registry:

```bash
docker login <registry_domain> -u <account_login> -p <account_password>
```
