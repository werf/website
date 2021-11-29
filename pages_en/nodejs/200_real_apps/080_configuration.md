---
title: Configuration and secrets
layout: wip
permalink: nodejs/200_real_apps/080_config.html
examples_initial: examples/nodejs/050_s3
examples: examples/nodejs/080_configuration
base_url: https://github.com/werf/werf-guides/blob/master/
description: |
  In this chapter, we will show how to use and store sensitive and non-sensitive application configurations properly.

  In previous chapters, we added the configuration directly to containers during the build or used the container's environment variables to pass parameters during the deployment.

  Now, you will learn how to store application parameters in ConfigMaps and Secrets for security and flexibility. We will show how you can use Helm chart values and werf secrets and discuss parameterization and configuration reuse approaches. In addition, you will learn how to store sensitive data along with the code in the application's Git repository.
---

{% include 200_real_apps/080_configuration.md.liquid %}
