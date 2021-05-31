---
title: Must haves
permalink: java_springboot/100_basic/40_must_haves.html
---

You need to know a few features and peculiarities of werf to use it confidently and comfortably:

- Firstly, you need to understand the **principles of working with source code and what giterminism is** — these ensure that the application can be reproduced in any environment and unify all processes.
- Secondly, you need to understand the **process for tagging images**. Do you need to add tags to them manually to build and deploy an image (no, with werf, you do not need to)?
- When deploying to various environments, you will learn about the concept of a **release** as well as the process of **debugging** it.
- After some time, due to continuous application improvements and regular releases, you may face a shortage of **space in the image storage** and the need to clean up the unneeded images.

All these issues are discussed in more detail below.

## Working with source code and giterminism

Usually, some settings that affect the configuration of the deployed application are based on various "external" data: the environment in which building and deploying are performed, preset environment variables, generated files, external resources, etc. This directly leads to the inability to **guarantee the reproducibility** of the application state.

{% offtopic title="Why do we need reproducibility?" %}
The ability to reproduce a specific application state at any moment makes the **debugging process** easier and allows you to **deploy the copy of the application** for development or testing purposes. Plus, it makes the results of testing more credible and dependable thanks to a lower number of hidden parameters affecting the application state.

Reproducibility is fundamentally crucial to implement the [infrastructure as code](https://en.wikipedia.org/wiki/Infrastructure_as_code) and immutable infrastructure approach. 
{% endofftopic %}

werf follows the principles of giterminism to **guarantee** reproducibility. This way, the application state is entirely determined by the current git state (the commit that HEAD points to). By default, werf doesn't allow using uncommitted files in the configuration and the build context of the images being built. Also, it excludes functionality that potentially has external dependencies.

We strongly suggest following this approach. However, if necessary, you can ease giterminism restrictions as well as enable the functionality that requires an in-depth understanding of the process using the [werf-giterminism.yaml]({{ site.url }}/documentation/reference/werf_giterminism_yaml.html) configuration file.

When debugging and developing, changing application files can be inconvenient because of the need to make interim commits. In these cases, we recommend you switch to the [development mode](#dev-mode).

Please, refer to the [werf documentation]({{ site.url }}/documentation/advanced/giterminism.html) to learn more about giterminism.

### dev mode

You can activate the dev mode via the `--dev` option or `WERF_DEV` environment variable. This mode supports two operating models that you can switch between using the `--dev-mode` option (by default, the _simple_ mode is enabled):

- _simple_: for working with the worktree state of the application's Git repository;
- _strict_: for working with the [index state](http://shafiul.github.io/gitbook/1_the_git_index.html) in the application's Git repository.

Note that you must add new (_untracked_) files to the index manually via the `git add` command (a one-time task that you must perform for each newly added file).

### follow mode

With the follow mode enabled, werf restarts the command in response to the git state changes:

- the command is restarted when a new commit is created;
- if the [dev mode](#dev-mode) (`--dev`) is also enabled, then the state of the index of the git repository is monitored for changes; you just need to add the new changes to the index (`git add`) to restart the command.

The mode is activated via the `--follow` option or related `WERF_FOLLOW` environment variable.

## Tagging images

In the manual deployment process (i.e., without werf), you have to define the strict rules for tagging images and follow them (it is [no easy task](https://www.youtube.com/watch?v=oh4N2wBJCc8), frankly). However, you may have noticed that we use the {% raw %}`image: {{ .Values.werf.image.basicapp }}`{% endraw %} construct in werf charts: 

{% raw %}
```yaml
      - name: basicapp
        command: ["java","-jar","/app/demo.jar"]
        image: {{ .Values.werf.image.basicapp }}
```
{% endraw %}

**werf frees you of worries about tagging rules**: if there are changes in the code, it will rebuild the appropriate image, add service tags to it, push it to the registry, and insert the applicable image name into the templates.  For this, werf stores metadata in the registry and tracks the file contents in the Git repository. 

{% offtopic title="How does it work?" %}
werf implements the so-called content-based tagging (in other words, tags depend on the content): it automatically builds and deploys images in response to changes in Git contents. What's essential is that images do not change after each commit (and there are no unnecessary re-deployments) as long as the content state in Git stays the same.

Without werf, you usually have to implement some formal rules for naming images in the registry. In our case, you do not have to worry about that: werf can automatically tag images for you.

In short, werf calculates the checksum of the files added to the image and generates tags for the image based on that checksum. It uses an [MVCC-based](https://en.wikipedia.org/wiki/Multiversion_concurrency_control) approach and [optimistic locking](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) to make sure that there are no conflicts between various builds running simultaneously (including on different servers).

For more information about how data are stored in the registry, refer to the [werf documentation]({{ site.url }}/documentation/internals/stages_and_storage.html). It thoroughly describes the build process, how werf manages build stage dependencies, and how the naming is performed.
{% endofftopic %}

## Helm, releases, and more

[Helm-chart](https://helm.sh/docs/topics/charts/) is a collection of files that describe a related set of Kubernetes objects.  A _release_ is created when a chart is used to deploy an application to a specific environment.

{% offtopic title="I would like to learn more" %}
The [werf documentation]({{ site.url }}/documentation/v1.2/advanced/helm/releases/release.html) has plenty of information about working with releases, storing, and naming them.
{% endofftopic %}

When working with releases, helm implements the [3-way merge](https://helm.sh/docs/faq/#improved-upgrade-strategy-3-way-strategic-merge-patches) approach. Manual changes to the cluster that conflict with the state described in Git are corrected to conform to the latter. Note that manual changes that do not conflict with the state defined in Git remain outside the control of Helm and werf.

werf manages releases on its own, but if you want to do it the hard way, you can use the `werf helm <...>` commands.

### How do I list the installed components?

Releases provide information about **components installed in the cluster**, their state, and the **environment** they are running in.

To browse the list of releases or find out the name of the release you need, use the [`werf helm list -A` CLI command]({{ site.url }}/documentation/reference/cli/werf_helm_list.html).

### How do I delete the unneeded component?

Use [`werf dismiss`]({{ site.url }}/documentation/reference/cli/werf_dismiss.html) to uninstall an app release from Kubernetes.

## Debugging the installation process

Frequently, mistakes made in chart configurations lead to problems when rolling out a release. The [`werf render`]({{ site.url }}/documentation/reference/cli/werf_render.html) command can help you with debugging this kind of issues.

`werf render` performs all the actions related to building and generating charts, showing you the resulting charts instead of deploying the release to Kubernetes. It is a resource-intensive task that shows you the final result with all the required values filled in.

_Note that `werf render` only works with files committed to Git (like all other werf commands) but supports the `--dev` mode_.

## Storage space

Over time, a lot of data can accumulate in the storage (either local or in the registry). Werf has three built-in commands related to cleaning: `werf cleanup`, `werf purge`, `werf host purge`. They have different purposes, and we briefly discuss them below (the detailed information is available in the [documentation]({{ site.url }}/documentation/advanced/cleanup.html)).

### Regular registry cleanup

The [`werf cleanup`]({{ site.url }}/documentation/reference/cli/werf_cleanup.html) command performs a regular and safe cleanup. It safely deletes images that are no longer needed using an advanced algorithm that takes into account Git history, registry contents, and the cluster state.

We will configure scheduled registry cleanup using the CI system's tools in the chapter "Working with infrastructure". 

### Delete all

The following two service commands allow you to free up disk space by deleting all images and other data:

- [`werf host purge`]({{ site.url }}/documentation/reference/cli/werf_host_purge.html) deletes all host contents, leaving the registry intact;
- [`werf purge`]({{ site.url }}/documentation/reference/cli/werf_purge.html) deletes all images (CAUTION: NOT SAFE!).
