---
title: Must knows
permalink: nodejs/100_basic/40_must_haves.html
layout: "wip"
---

## Helm, releases, and more

[Chart](https://helm.sh/docs/topics/charts/) is a collection of files that describe a related set of Kubernetes objects.  A _release_ is created when a chart is used to deploy an application to a specific environment.

{% offtopic title="I would like to learn more" %}
The [werf documentation](https://werf.io/documentation/advanced/helm/basics.html#release) provides a detailed description of working with releases, storing, and naming them.
{% endofftopic %}

When working with releases, helm implements the [3-way merge](https://helm.sh/docs/faq/#improved-upgrade-strategy-3-way-strategic-merge-patches) strategy, i.e., manual changes to the cluster are also brought to the state described in Git.

Releases provide information about **components installed in the cluster**, their state, and the **environment** they are running in. You can manage those components using `werf helm <...>` commands.

### How do I list the installed components?
To browse the list of releases or find out the name of the release you need, use the [`werf helm list`](https://werf.io/documentation/reference/cli/werf_helm_list.html) command.

### How do I delete the unneeded component?
Use werf dismiss or [`werf helm uninstall`](https://werf.io/documentation/reference/cli/werf_helm_uninstall.html) to uninstall an application.

{% offtopic title="But what about the werf dismiss command?" %}
The [`werf dismiss`](https://werf.io/documentation/reference/cli/werf_dismiss.html) command also allows you to remove an application from Kubernetes. However, it requires access to the source code of an application.
{% endofftopic %}

## Debugging the install process

Frequently, mistakes made in chart configurations lead to problems when rolling out a release. The [`werf render`](https://werf.io/documentation/reference/cli/werf_render.html) command can help you to debug those mistakes.

`werf render` performs all the actions related to the building and generating charts, showing you the resulting charts instead of deploying the release to Kubernetes. It is a time-consuming task that shows you the final result with all the required values filled in.

_Note that `werf render` only works with files committed to git (like all other werf commands)._

## Tagging images

In the manual deploy process, you have to define strict rules for tagging images and follow them. On the other hand, you may have noticed that we use the {% raw %}`{{ tuple "basicapp" . | werf_image}}`{% endraw %} construct in werf charts: 

{% raw %}
```yaml
      - name: basicapp
        command: ["node","/app/app.js"]
        image: {{ tuple "basicapp" . | werf_image}}
```
{% endraw %}

Fortunately, **werf frees you of worries about tagging rules**: if there are changes in the code, it will rebuild the appropriate image, add service tags to it, push it to the registry, and insert the applicable image name into the templates.  For this, werf stores metadata in the registry and tracks file contents in the Git repository. 

{% offtopic title="How does it all work?" %}
You can learn more about how the data is stored in the registry in the [werf documentation](https://werf.io/documentation/internals/stages_and_storage.html). It thoroughly describes the build process, how werf manages build stage dependencies, and how the naming is performed.
{% endofftopic %}

All data can be stored either locally on the host or in the Docker Repo.

## Storage space

Over time, a lot of data can accumulate in the storage (either locally or in the registry). werf has three built-in commands related to cleaning: `werf cleanup`, `werf purge`, `werf host purge`. They have different purposes, and we briefly discuss them below (the detailed information is available in [the documentation](https://werf.io/documentation/advanced/cleanup.html)).

### Regular registry cleanup

The [`werf cleanup`](https://ru.werf.io/documentation/reference/cli/werf_cleanup.html) command performs a regular and safe cleanup. It does not delete images that are in use in the cluster.

We will configure scheduled registry cleanup using the CI system's tools in the chapter "Working with infrastructure". 

### Delete all

The following two service commands allow you to free up disk space by deleting all images:

- [`werf host purge`](https://werf.io/documentation/reference/cli/werf_host_purge.html) deletes all host contents, leaving the registry intact.
- [`werf purge`](https://werf.io/documentation/reference/cli/werf_purge.html) (NOT SAFE) deletes all images, including those linked to applications running in the cluster!

## Giterminism

As is often the case, some settings that affect the configuration of the deployed application depend on "external" data: files that exist on the runner, some dynamic measures emanating from external resources, etc. This directly leads to the inability to **guarantee the reproducibility** of the application state.

{% offtopic title="Why do we need reproducibility?" %}
The ability to reproduce a specific application state at any moment makes the **debugging** process easier and allows us to easily **deploy the copy of the project** for development or testing purposes. Plus, it makes the results of testing more credible and dependable thanks to a lower number of hidden parameters affecting the application state.
{% endofftopic %}

To **ensure** reproducibility, werf, by default, considers only states that are committed to Git. In other words, Git is the primary and sufficient source of information about the application.

{% offtopic title="How is this implemented?" %}
Everything is defined in werf/Helm configs that are stored in Git. You cannot use environment variables or read a file that isn't committed.

The `Files.Get` functions in werf and Helm allow for reading committed files only.

In version 1.3, `--set` and `--values` CLI parameters will be removed from werf so that no external context can be passed into the build (`--env` is the only parameter that will remain).
{% endofftopic %}

However, we understand that committing every code change during the development isn't feasible for many reasons. That is why werf has two special modes, `--dev` and `--follow`, to ease the development and make it more convenient.

### dev mode

If you run `converge` (or any other command — e.g., `render`) with the `--dev` CLI parameter:

```shell
werf converge --repo registry.example.com/werf-guided-nodejs --dev
```

… then werf will build and deploy not only committed but **added** (via the `git add` command) files as well.

### follow mode 

If you run `converge` (or any other command — e.g., `render`) with the `--follow` CLI parameter:

```shell
werf converge --repo registry.example.com/werf-guided-nodejs --follow
```

… then the command will be automatically restarted with every new commit to Git.

If you combine `--follow` and `--dev` parameters, the command will be restarted in response to the `git add` command.

<div id="go-forth-button">
    <go-forth url="40_optimize.html" label="Speding up the build" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
