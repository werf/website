---
title: Must haves
permalink: java_springboot/100_basic/40_must_haves.html
---

You need to know a few features and peculiarities of werf to use it confidently and comfortably:

- Firstly, you need to understand the **principles of working with source code and what giterminism is** — they ensure that the application can be reproduced in any environment.
- Secondly, you need to understand the **process for tagging images**. Do you need to add tags to them manually to build and deploy an image (no, with werf, you do not need to)?
- When deploying to various environments, you will learn about the concept of a **release** as well as the process of **debugging** them.
- After some time, due to continuous application improvements and regular releases, you may face the **shortage of space in the image storage** and the need to delete unneeded images. Werf provides mechanisms that simplify this task considerably.

All these issues are discussed in more detail below.

## Working with source code and giterminism

Often, some settings that affect the configuration of the deployed application depend on "external" data: files that exist on the runner, some dynamic measures emanating from external resources, etc.  This directly leads to the inability to **guarantee the reproducibility** of the application state.

{% offtopic title="Why do we need reproducibility?" %}
The ability to reproduce a specific application state at any moment makes the **debugging** process easier, allowing us to **deploy the copy of the project** for development or testing purposes. Plus, it makes the results of testing more credible and dependable thanks to a lower number of hidden parameters affecting the application state.

Reproducibility is fundamentally crucial to implement the [infrastructure as code](https://en.wikipedia.org/wiki/Infrastructure_as_code) and immutable infrastructure approach. {% endofftopic %}

To **ensure** reproducibility, werf, by default, considers only states that are committed to Git. In other words, Git is the primary and sufficient source of information about the application.

{% offtopic title="How is this implemented?" %}
Everything is defined in werf/Helm configs that are stored in Git. You cannot use environment variables or read a file that isn't committed to Git.

The `Files.Get` functions in werf and Helm allow for reading committed files only.

We plan to get rid of `--set` and `--values` CLI parameters in werf 1.3 and implement a mechanism for the strict definition and validation of allowable external parameters (OpenAPI) to make the usage of an external context as safe as possible.
{% endofftopic %}

However, we understand that committing every code change during the development isn't feasible for many reasons. That is why werf has a special `--dev` mode to ease the development and make it more convenient. Also, in some situations, the build must depend on the external context.  In these rare cases, you can configure the giterminism settings via the `werf-giterminism.yaml` file.

### werf-giterminism.yaml

Using this config file, you can eliminate some restrictions in a specific configuration (e.g., using uncommitted files, mounting directories and files during the build, passing environment variables, etc.). You can learn more about the configuration process [in the documentation](https://werf.io/documentation/advanced/configuration/giterminism.html#werf-giterminismyaml).

### dev mode

If you run `converge` (or any other command — e.g., `render`) with the `--dev` CLI parameter:

```shell
werf converge --repo registry.example.com/werf-guided-springboot --dev
```

... then werf will build and deploy not only committed but **tracked files** (added via `git add`) as well.


{% offtopic title="Why do I need to invoke git add every time?" %}
As you already know, werf (in "normal" mode) reads files exclusively from Git to ensure reproducibility.

It is not due to technical limitations (you can read files directly from the file system) but to eliminate the uncertainty. werf reads files exclusively from Git (even in the dev mode) to prevent bugs related to differences between files read from the file system and pulled from Git.

It creates an extra layer of protection and predictability of behavior. Plus, it ensures that you add files to the commit, while the final build will be the same as that made locally.
{% endofftopic %}

### follow mode
The **follow mode** is another useful special-purpose werf mode. If you run the `converge` or another command (e.g., `run` or `compose up`) with the CLI parameter `--follow`:

```shell
werf converge --repo registry.example.com/werf-guided-springboot --follow
```

... then the command will be automatically restarted with every new commit to Git.

If you combine `--follow` and `--dev` parameters, the command will be restarted in response to the `git add` command.

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
In short, werf calculates the checksum of the files added to the image and generates tags for the image based on that checksum. It uses an [MVCC-based](https://en.wikipedia.org/wiki/Multiversion_concurrency_control) approach and [optimistic locking](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) to make sure that there are no conflicts between various builds running simultaneously (including on different servers).

For more information about how data is stored in the registry, refer to the [werf documentation](https://werf.io/documentation/internals/stages_and_storage.html). It thoroughly describes the build process, how werf manages build stage dependencies, and how the naming is performed.
{% endofftopic %}

## Helm, releases, and more

[Helm-chart](https://helm.sh/docs/topics/charts/) is a collection of files that describe a related set of Kubernetes objects.  A _release_ is created when a chart is used to deploy an application to a specific environment.

{% offtopic title="I would like to learn more" %}
The [werf documentation](https://werf.io/documentation/v1.2/advanced/helm/releases/release.html) has plenty of information about working with releases, storing, and naming them.
{% endofftopic %}

When working with releases, helm implements the [3-way merge](https://helm.sh/docs/faq/#improved-upgrade-strategy-3-way-strategic-merge-patches) approach. Manual changes to the cluster that conflict with the state described in Git are corrected to conform to the latter. Note that manual changes that do not conflict with the state defined in Git remain outside the control of Helm and werf.

werf manages releases on its own, but if you want to do it the hard way, you can use the `werf helm <...>` commands.

### How do I list the installed components?

Releases provide information about **components installed in the cluster**, their state, and the **environment** they are running in.

To browse the list of releases or find out the name of the release you need, use the [`werf helm list -A` CLI command](https://werf.io/documentation/reference/cli/werf_helm_list.html).

### How do I delete the unneeded component?

Use [`werf helm uninstall`](https://werf.io/documentation/reference/cli/werf_helm_uninstall.html) to uninstall an application.

{% offtopic title="But what about the werf dismiss command?" %}
The [`werf dismiss`](https://werf.io/documentation/reference/cli/werf_dismiss.html) command also allows you to remove an application from Kubernetes. However, it requires access to the source code of the application (which can be inconvenient outside of the CI system).
{% endofftopic %}

## Debugging the installation process

Frequently, mistakes made in chart configurations lead to problems when rolling out a release. The [`werf render`](https://werf.io/documentation/reference/cli/werf_render.html) command can help you to debug such problems.

`werf render` performs all the actions related to building and generating charts, showing you the resulting charts instead of deploying the release to Kubernetes. It is a resource-intensive task that shows you the final result with all the required values filled in.

_Note that `werf render` only works with files committed to Git (like all other werf commands) but supports the `--dev` mode_.

## Storage space

Over time, a lot of data can accumulate in the storage (either local or in the registry). Werf has three built-in commands related to cleaning: `werf cleanup`, `werf purge`, `werf host purge`. They have different purposes, and we briefly discuss them below (the detailed information is available in the [documentation](https://werf.io/documentation/advanced/cleanup.html)).

### Regular registry cleanup

The [`werf cleanup`](https://werf.io/documentation/reference/cli/werf_cleanup.html) command performs a regular and safe cleanup. It safely deletes images that are no longer needed using an advanced algorithm that takes into account Git history, registry contents, and the cluster state.

We will configure scheduled registry cleanup using the CI system's tools in the chapter "Working with infrastructure". 

### Delete all

The following two service commands allow you to free up disk space by deleting all images and other data:

- [`werf host purge`](https://werf.io/documentation/reference/cli/werf_host_purge.html) deletes all host contents, leaving the registry intact.
- [`werf purge`](https://werf.io/documentation/reference/cli/werf_purge.html) (NOT SAFE) deletes all images, including those linked to applications running in the cluster!
