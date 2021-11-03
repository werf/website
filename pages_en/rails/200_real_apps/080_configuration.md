---
title: Configuration and secrets
permalink: rails/200_real_apps/080_config.html
examples_initial: examples/rails/050_s3
examples: examples/rails/080_configuration
base_url: https://github.com/werf/werf-guides/blob/master/
description: |
  In this chapter, we will show how to use and store sensitive and non-sensitive application configurations properly. 

  In previous chapters, we added the configuration directly to containers during the build or used the container's environment variables to pass parameters during the deployment.
  
  Now, you will learn how to store application parameters in ConfigMaps and Secrets for security and flexibility. We will show how you can use Helm chart values and werf secrets and discuss parameterization and configuration reuse approaches. In addition, you will learn how to store sensitive data along with the code in the application's Git repository.
---

## ConfigMap and Secret

ConfigMap and Secret Kubernetes resources allow you to separate environment-dependent and context-specific configuration from container images. 

Both objects store data in key-value pairs and provide them to Pods as environment variables, command-line arguments, or configuration files mounted in a selected container.

ConfigMap stores non-confidential data, while Secret stores confidential data (various secret types).

You can learn more about these resource types in the official Kubernetes documentation ([ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/), [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)). This chapter discusses various use cases for ConfigMaps and Secrets using our application as an example.

## Storing application's config files in the ConfigMap

Currently, the [nginx.conf]({{ page.base_url | append: page.examples_initial | append: "/.werf/nginx.conf" }}) config file gets copied into the image during the build. As a result, the image is rebuilt, and Pods are restarted each time the file is modified. In addition, there is no easy way to benefit from templating in our `nginx.conf`.

You can get around these obstacles by moving `.werf/nginx.conf` into a dedicated ConfigMap. This way, you can mount `nginx.conf` when deploying instead of adding the file during the build:
{% include snippetcut_example path=".helm/templates/configmap-nginx.yaml" syntax="yaml" examples=page.examples %}
Now, let's add our ConfigMap to the Deployment — we will mount it as a file in the `frontend` container:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="nginx_conf" examples=page.examples %}
We also need to delete the no longer needed `.werf/nginx.conf` as well as a command to copy this file to an image during the build. The Dockerfile for the `frontend` image will look as follows:
{% include snippetcut_example path="Dockerfile" syntax="dockerfile" snippet="frontend" examples=page.examples %}

## Handling ConfigMap and Secret changes

Note that by default, changes in ConfigMaps or Secrets mounted in a Deployment, StatefulSet, or DaemonSet will not trigger Pod restarts for the new configuration to take effect. To trigger Pod restart, you need to add annotations with checksums of all ConfigMaps and Secrets used by this Pod. In this case, the annotations will change in response to ConfigMap and Secret modifications, leading to the Pod restart. Below is an example of an annotation that includes the checksum of the `nginx.conf` ConfigMap:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="nginx_conf_checksum" examples=page.examples %}
Note that you will need a unique annotation for each mounted ConfigMap/Secret.
>_The authors of this tutorial prefer to use operators like [stakater/Reloader](https://github.com/stakater/Reloader) instead of checksum-containing annotations because they are simpler, more flexible, and easier to work with._

## Values

Using a Helm chart to deploy an application has many advantages — for example, you can create manifest templates with Helm charts. `Values` is the key built-in templating object. With it, you can access values passed to the chart from various sources.

When using werf, all data passed to the chart can be grouped into several categories:
- [User-defined regular values]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/values.html#user-defined-regular-values): `values.yaml` parameters and corresponding options. 
- [User-defined secret values]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/values.html#user-defined-secret-values): parameters from the `secret-values.yaml` file.
- [Service values]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/values.html#service-values): information about the project, release, built images, etc.

## Reusing configuration with Values and ConfigMaps

Frequently used parameters can be put into `.helm/values.yaml`:
{% include snippetcut_example path=".helm/values.yaml" syntax="yaml" snippet="mysql" examples=page.examples %}
... to make manifest configuration easier:
{% include snippetcut_example path=".helm/templates/database.yaml" syntax="yaml" snippet="volume_claim_templates" examples=page.examples %}
This way, you can adjust the parameters to the specific environment (we will discuss this in more detail in the following tutorial chapters).

Probably, the most useful feature is the ability to place repetitive configuration snippets into `.helm/values.yaml`. For example, you can store in this file environment variables used in multiple places:
{% include snippetcut_example path=".helm/values.yaml" syntax="yaml" snippet="app" examples=page.examples %}
The environment variables from `.Values.app.envs` can be inserted either into the container manifest as `env` values (as we did before), or you can put them into the ConfigMap and load it into the container via `envFrom`.

The first option is easier to do, but the ConfigMap option is more convenient if you have a large number of shared environment variables. This way, you will avoid duplicating them between controllers. Your ConfigMap may look like this:
{% include snippetcut_example path=".helm/templates/configmap-app-envs.yaml" syntax="yaml" examples=page.examples %}

Now you can use `envFrom` in the Deployment to define ConfigMap's data as container environment variables:

{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="backend_conf" examples=page.examples %}

Don't forget to add checksum-containing annotations to trigger Pod re-deployment in response to ConfigMap changes:

{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="configmap_app_envs_checksum" examples=page.examples %}

Similarly, you can load this ConfigMap into other resources that use the same environment variables.

## Handling confidential data using Values and Secrets

To start working with secrets, you first need to generate a symmetric encryption key. You can use the following command to do this: `werf helm secret generate-secret-key`. However, since we have already prepared/encrypted the secrets for you, the encryption key that will be used in this chapter is generated as well. This encryption key is stored in the repository in the `.werf_secret_key` file; werf uses it automatically.
>_Caution! You MUST NOT store the key in the repository when working with real-life applications. We recommend keeping it in a safe place and passing via the `WERF_SECRET_KEY` environment variable. You can read more about using encryption keys in the [werf documentation]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/secrets.html#encryption-key)._

The application configuration contains data that should not be stored in unencrypted form in the repository (e.g., the database login and password). Thus, we won't store the login and password in plain text in the configuration file (as we currently do):
{% include snippetcut_example path="config/database.yml" syntax="yaml" examples=page.examples_initial %}
Instead, we will store them encrypted in the `.helm/secret-values.yaml` file along with other confidential parameters:
{% include snippetcut_example path=".helm/secret-values.yaml" syntax="yaml" examples=page.examples %}

Use the command below to decrypt the secrets stored in `.helm/secret-values.yaml`:
```bash
werf helm secret values decrypt .helm/secret-values.yaml
```

You should see the following output:
```yaml
app:
  secretEnvs:
    S3_USERNAME: minioadmin
    S3_PASSWORD: minioadmin
    DB_USERNAME: root
    DB_PASSWORD: password
    SECRET_KEY_BASE: something
mysql:
  secretEnvs:
    MYSQL_ROOT_PASSWORD: password
minio:
  secretEnvs:
    MINIO_ROOT_USER: minioadmin
    MINIO_ROOT_PASSWORD: minioadmin
```

Now, you need to pass the secrets contained in `.helm/secret-values.yaml` back to the application config. To do this, let's first pass them to the Secret resource:
{% include snippetcut_example path=".helm/templates/secret-app-envs.yaml" syntax="yaml" examples=page.examples %}

... and then load this resource's data as a set of environment variables by mounting it in containers:
{% include snippetcut_example path=".helm/templates/deployment.yaml" syntax="yaml" snippet="backend_secret" examples=page.examples %}

You have to make similar changes to the rest of the `.helm/templates` files (we won't delve into details here; you can take a look at their contents in the [repository]({{ page.base_url | append: page.examples | append: "/.helm/templates/" }}).)

After passing environment variables to the container, you have to substitute them in the application's configuration files:
{% include snippetcut_example path="config/database.yml" syntax="yaml" examples=page.examples %}
{% include snippetcut_example path="config/storage.yml" syntax="yaml" examples=page.examples %}
{% include snippetcut_example path="config/secrets.yml" syntax="yaml" examples=page.examples %}

You can also use Secrets to store and mount complete secret configuration files. This is similar to using ConfigMap to mount non-secret configuration files (as [described above](#storing-applications-config-files-in-the-configmap)], except that the Secret contents must be stored encrypted in `.helm/secret/...` or `.helm/secret-values.yaml`. You can read more about this in the [werf documentation]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/secrets.html#secret-files).

Note: learn more about working with secrets in the [werf documentation]({{ site.url }}/documentation/v1.2/advanced/helm/configuration/secrets.html).

## Checking whether the application runs as expected

Let's make sure the configuration changes did not affect the application performance:
```shell
werf converge --repo <DOCKER HUB USERNAME>/werf-guide-app
```

You should see the following output:
```
...
┌ ⛵ image backend
│ Use cache image for backend/dockerfile
│      name: .../werf-guide-app:8df39a540816b5851314934d932f0152a8fa409369a58dbfeccce4be-1632852281567
│        id: 0825950bfd70
│   created: 2021-09-28 21:04:41 +0300 MSK
│      size: 366.3 MiB
└ ⛵ image backend (7.19 seconds)

┌ ⛵ image frontend
│ Use cache image for frontend/dockerfile
│      name: .../werf-guide-app:18b3a25a107798890200664b7de037c998bef8ed9d9bbe4d66efaf53-1632852281351
│        id: 7d029b84655e
│   created: 2021-09-28 21:04:41 +0300 MSK
│      size: 9.5 MiB
└ ⛵ image frontend (7.60 seconds)

┌ Waiting for release resources to become ready
...
│ ┌ Status progress
│ │ DEPLOYMENT      REPLICAS  AVAILABLE  UP-TO-DATE
│ │ werf-guide-app  1/1       1          1
│ │ │    POD                         READY  RESTARTS  STATUS
│ │ ├──  guide-app-85bd8c68f5-w9nff  2/2    0         Running
│ │ └──  guide-app-865b6d68bc-bkt9l  2/2    0         Terminating
│ │ STATEFULSET  REPLICAS  READY  UP-TO-DATE
│ │ minio        1/1       1      1
│ │ │    POD  READY  RESTARTS  STATUS
│ │ └──  0    1/1    0         Running
│ │ mysql        1/1       1      1
│ │ │    POD  READY  RESTARTS  STATUS
│ │ └──  0    1/1    0         Running
│ │ JOB                        ACTIVE  DURATION  SUCCEEDED/FAILED
│ │ setup-and-migrate-db-rev2  0       24s       1/0
│ │ │    POD                        READY  RESTARTS  STATUS
│ │ └──  and-migrate-db-rev2-lhnq7  0/1    0         Completed
│ │ setup-minio-rev2           0       33s       0->1/0
│ │ │    POD               READY  RESTARTS  STATUS
│ │ └──  minio-rev2-wr7tw  0/1    0         Running  ->
│ │ Completed
│ └ Status progress
└ Waiting for release resources to become ready (33.05 seconds)

Release "werf-guide-app" has been upgraded. Happy Helming!
NAME: werf-guide-app
LAST DEPLOYED: Tue Sep 28 21:26:48 2021
NAMESPACE: werf-guide-app
STATUS: deployed
REVISION: 2
TEST SUITE: None
Running time 43.88 seconds
```

Check if the application is available:
```shell
curl http://werf-guide-app/ping
```

You should see the following output:
```
pong
```

In this chapter, you learned how to securely store the confidential application parameters, minimize duplicate configuration, and turn application configs into templates (thus, no re-builds are necessary when changes are made).
