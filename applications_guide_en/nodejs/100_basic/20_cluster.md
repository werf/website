---
title: Preparing the cluster
permalink: nodejs/100_basic/20_cluster.html
---

You would need the following components to develop in Kubernetes:

- the locally installed werf tool — done in the previous chapters;
- the application code stored in Git — done in the previous chapters;
- a Kubernetes cluster;
- a container registry;
- the appropriate DNS records;

**Setting up an environment is a difficult task**. We do not recommend digging into these problems yourselves. It makes more sense to use the existing services available on the market, seek help from colleagues who have the relevant experience, or ask questions in the [werf's Telegram community](https://t.me/werf_ru).

**Настройка окружения — это сложная задача**. Не рекомендуем закапываться в решение проблем своими руками: используйте существующие на рынке услуги, позовите на помощь коллег, которые имеют практику настройки инфраструктуры, или задавайте вопросы в [Telegram-сообществе werf](https://t.me/werf_ru).

{% offtopic title="A brief checklist for self-checking the environment" %}

- You have a Kubernetes cluster (version 1.14 or later).
    - There is an ingress controller in the cluster with the respective domain (`example.com`) linking to it (in other words, you can visit it in your browser).
- You have a registry.
    - The registry is available at `registry.example.com`.
    - The cluster can pull images from your registry.
- On the local computer:
    - there is access to the cluster (`kubectl version` shows the version of K8s); 
    - there is access to the registry (`docker push` works fine);
    - the browser can access `example.com` (even if it shows the 404 error page by Ingress).

{% endofftopic %}

Now select your preferred method/location for setting up an environment:

<div style="display: flex; justify-content: space-between; margin: 0 10px 0 20px;">
<div class="button__blue button__blue_inline expand_columns_button" id="minikube_button"><a href="#">Minikube</a></div>
<div class="button__blue button__blue_inline expand_columns_button" id="docker_desktop_button"><a href="#">Docker Desktop</a></div>
<div class="button__blue button__blue_inline expand_columns_button" id="cloud_provider_button"><a href="#">Cloud provider</a></div>
<div class="button__blue button__blue_inline expand_columns_button" id="has_cluster_button"><a href="#">I already have a cluster</a></div>
</div>

{% expandonclick id="minikube_button__content" %}
{% include_relative 20_cluster_minikube.md %}
{% endexpandonclick %}

{% expandonclick id="docker_desktop_button__content" %}
{% include_relative 20_cluster_docker_desktop.md %}
{% endexpandonclick %}

{% expandonclick id="cloud_provider_button__content" %}
{% include_relative 20_cluster_cloud_provider.md %}
{% endexpandonclick %}

{% expandonclick id="has_cluster_button__content" %}
{% include_relative 20_cluster_has_cluster.md %}
{% endexpandonclick %}

## Final checks

After verifying the cluster functionality, you have to check if the registry and ingress work well. To do this, run the following commands:

```bash
docker tag ubuntu:18.04 registry.example.com/ubuntu:18.04
docker push registry.example.com/ubuntu:18.04
```

They should push the Ubuntu image to the registry without any errors. And the command below:

```bash
curl example.com
```

... should return an nginx ingress error page (if you haven't already deployed an application to the cluster).

<div id="go-forth-button">
    <go-forth url="30_deploy.html" label="Deploying the application" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
