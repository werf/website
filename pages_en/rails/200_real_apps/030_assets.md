---
title: Delivering assets
permalink: rails/200_real_apps/030_assets.html
examples_initial: examples/rails/020_logging
examples: examples/rails/030_assets
base_url: https://github.com/werf/werf-guides/blob/master/
description: |
  In this chapter, you will learn how to work with static files and deliver them to the client correctly.
---

## Adding an `/image` page to the application

Let's add a new `/image` endpoint to our application. It will return a page generated using a set of static files. We will use Webpacker (instead of Sprockets) to bundle all JS, CSS, and media files.

Currently, our application provides a basic API with no means to manage static files and frontend code. To turn this API into a simple web application, we have generated a backbone of the new Rails application (without the `--api` and `--skip-javascript` flags):
```shell
rails new --skip-keeps --skip-action-mailer --skip-action-mailbox --skip-action-text --skip-active-record --skip-active-job --skip-active-storage --skip-action-cable --skip-sprockets --skip-spring --skip-listen --skip-turbolinks --skip-jbuilder --skip-test --skip-system-test --skip-bootsnap .
```

This backbone includes all the functionality we need for managing JS code and static files. The following changes have been made to our application:
1. Adding `webpacker` to [Gemfile]({{ page.base_url | append: page.examples | append: "/Gemfile" }}).
1. Generating standard Webpacker configuration files:
  * [config/webpacker.yml]({{ page.base_url | append: page.examples | append: "/config/webpacker.yml" }})
  * [config/webpack/environment.js]({{ page.base_url | append: page.examples | append: "/config/webpack/environment.js" }})
  * [config/webpack/production.js]({{ page.base_url | append: page.examples | append: "/config/webpack/production.js" }})
1. Creating a set of directories to host the source static files:
  * [app/javascript/styles]({{ page.base_url | append: page.examples | append: "/app/javascript/styles" }})
  * [app/javascript/src]({{ page.base_url | append: page.examples | append: "/app/javascript/src" }})
  * [app/javascript/images]({{ page.base_url | append: page.examples | append: "/app/javascript/images" }})

Now let's add an HTML page template (available at `/image`) with the _Get image_ button:
{% include snippetcut_example path="app/views/layouts/image.html.erb" syntax="erb" examples=page.examples %}

Clicking the _Get image_ button must result in an Ajax request that pulls and displays an [SVG image]({{ page.base_url | append: page.examples | append: "/app/javascript/images/werf-logo.svg" }}):
{% include snippetcut_example path="app/javascript/src/image.js" syntax="js" examples=page.examples %}

Our page will also use a [app/javascript/styles/site.css]({{ page.base_url | append: page.examples | append: "/app/javascript/styles/site.css" }}) CSS file.

JS and CSS files as well as an SVG image will be bundled with Webpack and placed to `/packs/...`:
{% include snippetcut_example path="app/javascript/packs/application.js" syntax="js" examples=page.examples %}

Let's update the controller and routes by adding a new `/image` endpoint that takes care of the HTML template created above:
{% include snippetcut_example path="app/controllers/application_controller.rb" syntax="rb" examples=page.examples %}
{% include snippetcut_example path="config/routes.rb" syntax="rb" examples=page.examples %}

The application now has a new endpoint called `/image` in addition to the `/ping` endpoint we made in previous chapters. This new endpoint displays a page that uses different types of static files.

>_The commands provided at the [beginning of the chapter](#preparing-a-repository) allow you to view a complete, exhaustive list of the changes made to the application in the current chapter._

## Delivering static files

By default, Rails does not even try to distribute static files in production environments. Instead, Rails developers suggest using a reverse proxy like NGINX for this task. This is because the reverse proxy distributes static files much more efficiently than Rails and Puma.

In practice, you can do without the reverse proxy running in front of an application only during development. In addition to efficiently serving static files, the reverse proxy nicely complements the application server (Puma) by providing many additional features not available at the application server level. In addition, it allows for fast and flexible request routing.

There are several ways to deploy the application server behind the reverse proxy in Kubernetes. We will use a popular and straightforward method that nevertheless scales well. As part of it, the NGINX container is deployed into each Pod with the Rails/Puma container. This auxiliary container serves as a proxy for all Rails/Puma requests except for static file requests. These astatic files are delivered by the NGINX container.

Now, let's get to the implementation.

## Making changes to the build and deploy process

First of all, we have to make changes to the application build process. Now we need to build an NGINX image in addition to the Rails/Puma one. This NGINX image must have all the static files to serve to the client:
{% include snippetcut_example path="Dockerfile" syntax="dockerfile" examples=page.examples %}

The NIGINX config will be added to the NGINX image during the build:
{% include snippetcut_example path=".werf/nginx.conf" syntax="nginx" examples=page.examples %}

Let's update `werf.yaml` so that werf can build two images (backend, frontend) instead of one:
{% include snippetcut_example path="werf.yaml" syntax="yaml" examples=page.examples %}

Add a new NGINX container to the application's Deployment:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" examples=page.examples %}

Services and Ingresses should now connect to port 80 instead of 3000 so that all requests are routed through the NGINX proxy instead of sending them directly to Rails/Puma:
{% include snippetcut_example path=".helm/templates/service.yaml" syntax="yaml" examples=page.examples %}
{% include snippetcut_example path=".helm/templates/ingress.yaml" syntax="yaml" examples=page.examples %}

## Checking that everything works as expected

Now we have to re-deploy our application:
```shell
werf converge --repo <DOCKER HUB USERNAME>/werf-guide-app
```

You should see the following output:
```shell
...
┌ ⛵ image backend
│ ┌ Building stage backend/dockerfile
│ │ backend/dockerfile  Sending build context to Docker daemon    320kB
│ │ backend/dockerfile  Step 1/24 : FROM ruby:2.7 as base
│ │ backend/dockerfile   ---> 1faa5f2f8ca3
...
│ │ backend/dockerfile  Step 24/24 : LABEL werf-version=v1.2.12+fix2
│ │ backend/dockerfile   ---> Running in d6126083b530
│ │ backend/dockerfile  Removing intermediate container d6126083b530
│ │ backend/dockerfile   ---> 434a85e3df9d
│ │ backend/dockerfile  Successfully built 434a85e3df9d
│ │ backend/dockerfile  Successfully tagged 145074fa-9024-4bcb-9144-e245e1bbaf87:latest
│ │ ┌ Store stage into .../werf-guide-app
│ │ └ Store stage into .../werf-guide-app (14.27 seconds)
│ ├ Info
│ │      name: .../werf-guide-app:7d52fed109415751ad6f28dd259d7a201dc2ffe34863e70f0a404adf-1629129180941
│ │        id: 434a85e3df9d
│ │   created: 2021-08-16 18:53:00 +0300 MSK
│ │      size: 335.1 MiB
│ └ Building stage backend/dockerfile (23.70 seconds)
└ ⛵ image backend (30.13 seconds)

┌ ⛵ image frontend
│ ┌ Building stage frontend/dockerfile
│ │ frontend/dockerfile  Sending build context to Docker daemon    320kB
│ │ frontend/dockerfile  Step 1/28 : FROM ruby:2.7 as base
│ │ frontend/dockerfile   ---> 1faa5f2f8ca3
...
│ │ frontend/dockerfile  Step 28/28 : LABEL werf-version=v1.2.12+fix2
│ │ frontend/dockerfile   ---> Running in 29323e1512c7
│ │ frontend/dockerfile  Removing intermediate container 29323e1512c7
│ │ frontend/dockerfile   ---> 8f637e588672
│ │ frontend/dockerfile  Successfully built 8f637e588672
│ │ frontend/dockerfile  Successfully tagged 215f9544-2637-4ecf-9b81-8cbbcd6777fe:latest
│ │ ┌ Store stage into .../werf-guide-app
│ │ └ Store stage into .../werf-guide-app (9.83 seconds)
│ ├ Info
│ │      name: .../werf-guide-app:bdb809c26315846927553a069f8d6fe72c80e33d2a599df9c6ed2be0-1629129180706
│ │        id: 8f637e588672
│ │   created: 2021-08-16 18:53:00 +0300 MSK
│ │      size: 9.4 MiB
│ └ Building stage frontend/dockerfile (19.02 seconds)
└ ⛵ image frontend (29.28 seconds)

┌ Waiting for release resources to become ready
│ ┌ Status progress
│ │  DEPLOYMENT      REPLICAS                    AVAILABLE  UP-TO-DATE
│ │  werf-guide-app  2/1                         1          1
│ │  │               POD                         READY      RESTARTS    STATUS             ---
│ │  ├──             guide-app-566dbd7cdc-n8hxj  2/2        0           Running            Waiting  for:  replicas  2->1
│ │  └──             guide-app-c9f746755-2m98g   0/2        0           ContainerCreating
│ └ Status progress
│
│ ┌ Status progress
│ │  DEPLOYMENT      REPLICAS                    AVAILABLE  UP-TO-DATE
│ │  werf-guide-app  2->1/1                      1          1
│ │  │               POD                         READY      RESTARTS    STATUS
│ │  ├──             guide-app-566dbd7cdc-n8hxj  2/2        0           Running            ->  Terminating
│ │  └──             guide-app-c9f746755-2m98g   2/2        0           ContainerCreating  ->  Running
│ └ Status progress
└ Waiting for release resources to become ready (9.11 seconds)

Release "werf-guide-app" has been upgraded. Happy Helming!
NAME: werf-guide-app
LAST DEPLOYED: Mon Aug 16 18:53:19 2021
NAMESPACE: werf-guide-app
STATUS: deployed
REVISION: 23
TEST SUITE: None
Running time 41.55 seconds
```

Go to [http://werf-guide-app/image](http://werf-guide-app/image) in your browser and click on the _Get image_ button. You should see the following output:

{% asset guides/rails/030_assets_success.png %}

Note which resources were requested and which links were used (the last resource was retrieved via the Ajax request):

{% asset guides/rails/030_assets_resources.png %}

Now our application not only provides an API but boasts a set of tools to manage static and JavaScript files effectively.

Furthermore, our application can handle a high load since many requests for static files will not affect the operation of the application as a whole. Note that you can quickly scale Puma (serves dynamic content) and NGINX (serves static content) by increasing the number of `replicas` in the application's Deployment.
