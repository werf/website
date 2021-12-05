---
title: Deploying the application
permalink: guides/java_springboot/100_basic/30_deploy.html
examples: examples/basic/002_deploy
examples_initial: examples/basic/001_build
description: |
    In the previous chapters, we built an image of the application and set up the environment to deploy it. Now let's deploy the application to the Kubernetes cluster we configured.

    Kubernetes manifests are used for deploying to Kubernetes. They describe the resources (Kubernetes objects) required for the application to run. These resources include Deployment (it is responsible for running applications in containers) and Service/Ingress (these are responsible for accessing running applications from inside and outside of the cluster).
---

{% include guides/100_basic/30_deploy.md.liquid %}
