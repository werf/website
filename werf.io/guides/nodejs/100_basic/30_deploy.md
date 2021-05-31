---
title: Deploying the application
permalink: nodejs/100_basic/30_deploy.html
---

{% filesused title="Files mentioned in the chapter" %}
- .helm/templates/deployment.yaml
- .helm/templates/registry-secret.yaml
- .helm/templates/ingress.yaml
- .helm/templates/service.yaml
{% endfilesused %}

In the previous chapter, we defined IaC for building, now we need to define IaC for deploying the application to Kubernetes. You will need the following Kubernetes objects: Deployment, Service, and Ingress.

{% offtopic title="What if I know nothing about Kubernetes?" %}
The tutorial contains the source code of the infrastructure, and you will be able to guess what is going on intuitively. You can refer to tutorials/guides and [official Kubernetes documentation](https://kubernetes.io/docs/tutorials/kubernetes-basics/) to learn how to write code like this yourself. Also, there are many video courses and books available that describe various Kubernetes objects and their settings.
{% endofftopic %}

werf supports all the features of the Helm template engine (Helm is compiled into werf and helps it to perform deployments) and provides [additional functionality]({{ site. docsurl }}/documentation/v1.2/advanced/helm/overview.html). We will discuss the peculiarities of templating, and the nuances of creating Kubernetes objects later in the "Configuring the infrastructure as code" chapter. Meantime, we will deploy the application to a real cluster and make sure that it is running smoothly.

## Deployment

The Deployment object allows you to create a Pod object that contains application containers and manages them. Our Pod will contain just one container called `basicapp`.

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/015_deploy_app/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: basicapp
spec:
  selector:
    matchLabels:
      app: basicapp
  revisionHistoryLimit: 3
  strategy:
    type: RollingUpdate
  replicas: 1
  template:
    metadata:
      labels:
        app: basicapp
    spec:
      containers:
      - name: basicapp
        command: ["node","/app/app.js"]
        image: {{ .Values.werf.image.basicapp }}
        workingDir: /app
        ports:
        - containerPort: 3000
          protocol: TCP
        env:
        - name: "SQLITE_FILE"
          value: "app.db"
```
{% endraw %}
{% endsnippetcut %}

Note the {% raw %}`image: {{ .Values.werf.image.basicapp }}`{% endraw %} construct — it puts in the actual image name.

werf rebuilds containers only when necessary: if there are changes in the related Git files or in `werf.yaml`. When an image changes, its tag also changes. As a result, the Pod gets redeployed with the new image.

## Registry Secret

The Kubernetes cluster uses images stored in the registry to run the application. Therefore, it is important that the cluster can log in to the registry. Generally, the situation varies for local and remote registries.

<div class="tabs">
<a href="javascript:void(0)" class="tabs__btn tabs__secret__btn active" onclick="openTab(event, 'tabs__secret__btn', 'tabs__secret__content', 'tab__secret__local')">Local registry</a>
<a href="javascript:void(0)" class="tabs__btn tabs__secret__btn" onclick="openTab(event, 'tabs__secret__btn', 'tabs__secret__content', 'tab__secret__remote')">Remote registry</a>
</div>

<div id="tab__secret__local" class="tabs__content tabs__secret__content active" markdown="1">
{% include_relative 30_deploy_registrysecret_local.md %}
</div>

<div id="tab__secret__remote" class="tabs__content tabs__secret__content" markdown="1">
{% include_relative 30_deploy_registrysecret_remote.md %}
</div>


## Service

The Service object allows applications in the cluster to discover each other. Let's define it:

{% snippetcut name=".helm/templates/service.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/015_deploy_app/.helm/templates/service.yaml" %}
{% raw %}
```yaml
apiVersion: v1
kind: Service
metadata:
  name: basicapp
spec:
  selector:
    app: basicapp
  ports:
  - name: http
    port: 3000
    protocol: TCP
```
{% endraw %}
{% endsnippetcut %}

## Ingress

The Ingress object routes traffic to our Service for the appropriate domain (`example.com` in our case).

{% offtopic title="What is the Ingress object?" %}
Here, the conflict of terminology is possible:

* First, there is the so-called [NGINX Ingress Controller](https://github.com/kubernetes/ingress-nginx). It runs inside the cluster and accepts incoming requests.
* Second, there is an [Ingress API object](https://kubernetes.io/docs/concepts/services-networking/ingress/) that, in fact, describes the configuration of the NGINX Ingress Controller.

Both of them are often called "Ingress" in articles and everyday speech, so you have to guess from the context. 
{% endofftopic %}

{% snippetcut name=".helm/templates/ingress.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/015_deploy_app/.helm/templates/ingress.yaml" %}
{% raw %}
```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: nginx
  name: basicapp
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /
        backend:
          serviceName: basicapp
          servicePort: 3000
```
{% endraw %}
{% endsnippetcut %}

## Deploying to the cluster

We will use the [`werf converge`]({{ site.docsurl }}/documentation/reference/cli/werf_converge.html) command to build an image, push it to the registry, and deploy our application to the Kubernetes cluster. The repository to store images (`--repo registry.example.com/werf-guided-nodejs`) serves as its only parameter.

{% offtopic title="How are images stored in the repository?" %}
When implementing the deployment process without werf, you have to take care of the scheme for naming images in the registry. In our case, you do not have to worry about that: werf can automatically tag images for you.

werf implements the so-called `content-based tagging` (in other words, tags depend on the content): it automatically builds and deploys images in response to changes in Git contents.

You can learn more in the [documentation]({{ site.docsurl }}/documentation/internals/stages_and_storage.html#stage-naming) and in the [article](https://medium.com/flant-com/content-based-tagging-in-werf-eb96d22ac509).
{% endofftopic %}

Commit changes to the code repository and then run the command below:

```shell
werf converge --repo registry.example.com/werf-guided-nodejs
```

As a result, you should see output similar to this:

```
│ │ basicapp/dockerfile  Successfully built 7e38465ee6de
│ │ basicapp/dockerfile  Successfully tagged cbb1cef2-a03a-432f-b13d-b95f0f0cb4e9:latest
│ ├ Info
│ │       name: localhost:5005/werf-guided-nodejs:017ce9df8dbd7d3505546c95557f1c1f39ce1e6666aaae29e8c12608-1605619646009
│ │       size: 375.8 MiB
│ └ Building stage basicapp/dockerfile (209.48 seconds)
└ ⛵ image basicapp (213.60 seconds)

Release "werf-guided-nodejs" does not exist. Installing it now.

┌ Waiting for release resources to become ready
│ ┌ Status progress
│ │ DEPLOYMENT                                                                                                                                                      REPLICAS                      AVAILABLE                        UP-TO-DATE
│ │ basicapp                                                                                                                                                        1/1                           1                                1
│ │ │   POD                                                           READY                  RESTARTS                       STATUS
│ │ └── 6cf5b444bc-6rh4j                                              1/1                    0                              Running
│ └ Status progress
└ Waiting for release resources to become ready (4.02 seconds)

NAME: werf-guided-nodejs
LAST DEPLOYED: Tue Nov 17 16:29:16 2020
NAMESPACE: werf-guided-nodejs
STATUS: deployed
REVISION: 1
TEST SUITE: None
Running time 222.54 seconds
```

And the application should be accessible via the browser:

![](/guides/images/template/100_30_app_in_browser.png)

<div id="go-forth-button">
    <go-forth url="40_optimize.html" label="Speeding up the build" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
