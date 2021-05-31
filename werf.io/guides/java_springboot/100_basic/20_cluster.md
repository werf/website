---
title: Preparing the environment
permalink: java_springboot/100_basic/20_cluster.html
---

You will need the following components to develop in Kubernetes:

- the locally installed werf tool — done in the previous chapters;
- the application code stored in Git — done in the previous chapters;
- a Kubernetes cluster;
- a container registry;
- the appropriate DNS records;

**Setting up an environment is a difficult task**. We do not recommend digging into these problems yourselves. It makes more sense to use the existing services available on the market, seek help from colleagues who have the relevant experience, or ask questions in the [werf's Telegram community](https://t.me/werf_io).

{% offtopic title="A brief checklist for self-checking the environment" %}

- You have a Kubernetes cluster (version 1.14 or later).
    - There is an ingress controller in the cluster with the respective domain (`example.com`) linking to it (in other words, you can visit it in your browser).
- You have a registry.
    - The registry is available at `registry.example.com`.
    - The cluster can pull images from your registry.
- On the local computer:
    - there is access to the cluster (`kubectl version` shows K8s version); 
    - there is access to the registry (`docker push` works fine);
    - the browser can access `example.com` (even if it shows the 404 error page by Ingress).

{% endofftopic %}

Now select your preferred method/location for setting up an environment:

<div class="tabs">
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn active" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__docker')">Docker Desktop</a>
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__minikube')">Minikube</a>
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__cloud')">Cloud provider</a>
<a href="javascript:void(0)" class="tabs__btn tabs__install__btn" onclick="openTab(event, 'tabs__install__btn', 'tabs__install__content', 'tab__install__ihave')">I already have a cluster</a>
</div>

<div id="tab__install__docker" class="tabs__content tabs__install__content active" markdown="1">
{% include_relative 20_cluster_docker_desktop.md %}
</div>

<div id="tab__install__minikube" class="tabs__content tabs__install__content" markdown="1">
{% include_relative 20_cluster_minikube.md %}
</div>

<div id="tab__install__cloud" class="tabs__content tabs__install__content" markdown="1">
{% include_relative 20_cluster_cloud_provider.md %}
</div>

<div id="tab__install__ihave" class="tabs__content tabs__install__content" markdown="1">
{% include_relative 20_cluster_has_cluster.md %}
</div>

## Final checks

After verifying the cluster functionality, you have to check if the registry and ingress work well. To do this, run the following commands:

```shell
docker tag ubuntu:18.04 registry.example.com/ubuntu:18.04
docker push registry.example.com/ubuntu:18.04
```

They should push the Ubuntu image to the registry without any errors. And the command below:

```shell
curl example.com
```

... should return an nginx ingress error page (if you haven't already deployed an application to the cluster).

<div id="go-forth-button">
    <go-forth url="30_deploy.html" label="Deploying the application" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
