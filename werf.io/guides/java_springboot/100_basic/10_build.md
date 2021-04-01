---
title: Building an image
permalink: java_springboot/100_basic/10_build.html
---

{% filesused title="Files mentioned in the chapter" %}
- Dockerfile
- werf.yaml
{% endfilesused %}

In this chapter, we'll start using werf. We will containerize the test application using werf and [Dockerfile](https://docs.docker.com/engine/reference/builder/) and run the final image locally in Docker.

## Preparing the workplace

We assume that you have already [installed werf]({{ site.docsurl }}/installation.html) and Docker.

Create a directory on your computer and follow these steps:

```shell
git clone https://github.com/werf/werf-guides.git
cp -r werf-guides/examples/springboot/000_app ./
cd 000_app 
git init
git add .
git commit -m "initial commit"
```

_This way you will copy the code of the [Spring Boot application](https://github.com/werf/werf-guides/tree/master/examples/springboot/000_app) to a local directory and initialize a Git repository in it._

## Dockerfile-based build process

The build configuration of our application consists of the following steps:

- pull the OpenJDK image (e.g., `gradle:jdk8-openj9`);
- add the application code to it;
- build the application using gradle and move the resulting jar to the appropriate location.

Now let's insert all these steps into a `Dockerfile`:

{% snippetcut name="Dockerfile" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/010_build/Dockerfile" %}
{% raw %}
```Dockerfile
FROM gradle:jdk8-openj9
WORKDIR /app
COPY . .

RUN gradle build --no-daemon
RUN cp /app/build/libs/*.jar /app/demo.jar

CMD ["java","-jar","/app/demo.jar"]
```
{% endraw %}
{% endsnippetcut %}

## Integrating Dockerfile into werf

Let's add our Dockerfile to werf. To do this, create a `werf.yaml` file in the repository root  (it describes the build process of the entire project):

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-springboot
configVersion: 1
---
image: basicapp
dockerfile: Dockerfile
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="What do all those lines mean?" %}

`werf.yaml` starts with a mandatory **meta-information section**:

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
project: werf-guided-springboot
configVersion: 1
```
{% endraw %}
{% endsnippetcut %}

Here:

- **_project_** is a field containing the unique name of the application project. The default project name is used when generating the Helm release name and the namespace in Kubernetes. Changing the name of an active project is a tedious task that requires a number of manual actions (more information about possible consequences is available [here]({{ site.docsurl }}/documentation/v1.2/reference/werf_yaml.html#warning-on-changing-project-name));
- **_configVersion_** specifies the syntax version used in the `werf.yaml` (currently, only version `1` is supported).

The next section, called the [**image config section**]({{ site.docsurl }}/documentation/reference/werf_yaml.html#image-section), defines the build process.

{% snippetcut name="werf.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/011_build_werf/werf.yaml" %}
{% raw %}
```yaml
---
image: basicapp
dockerfile: Dockerfile
```
{% endraw %}
{% endsnippetcut %}

The `image: basicapp` line sets the image ID that will be used in the rollout configuration as well as for invoking werf commands for a specific image described in `werf.yaml` (for example, `werf build basicapp`, `werf run basicapp`, etc.).

The `dockerfile: Dockerfile` line specifies that the build configuration is defined in the existing file located at the `Dockerfile` path.

You can also make use of other directives (their description is available [here]({{ site.docsurl }}/documentation/reference/werf_yaml.html#dockerfile-image-section-image)).

The single `werf.yaml` file can contain the definitions of an arbitrary number of images.

{% endofftopic %}

## Building

Before starting the build, you have to add changes to the project's Git repository:

```shell
git add werf.yaml Dockerfile
git commit -m "Add build configuration"
```

> Why should changes be added to the git repository? What are giterminism and the dev mode? You can find answers to these and other questions as well as nuances of working with project files in the chapter "Need to know"

The build is performed by the [`werf build`] ({{site.docsurl}} / documentation / reference / cli / werf_build.html) command:


{% raw %}
```shell
werf build
┌ ⛵ image basicapp
│ ┌ Building stage basicapp/dockerfile
│ │ basicapp/dockerfile  Sending build context to Docker daemon  116.7kB
│ │ basicapp/dockerfile  Step 1/15 : FROM gradle:jdk8-openj9
│ │ basicapp/dockerfile   ---> 2fb781988fa5
│ │ basicapp/dockerfile  Step 2/15 : WORKDIR /app
│ │ basicapp/dockerfile   ---> Using cache
<..>
│ │ basicapp/dockerfile  Successfully built e0d6df14df8b
│ │ basicapp/dockerfile  Successfully tagged ee51ea7f-c498-45a5-a435-0fd830fbb576:latest
│ ├ Info
│ │       name: werf-guided-springboot:50558f3f54d2ebbbd817824c6d7194aabe725bff6d7beae4df9c5e29-1606128099580
│ │       size: 738.6 MiB
│ └ Building stage basicapp/dockerfile (86.12 seconds)
└ ⛵ image basicapp (86.32 seconds)

Running time 86.37 seconds
```
{% endraw %}

## Running

The container is run using the [werf run]({{ site.docsurl }}/documentation/cli/main/run.html) command:

```shell
werf run --docker-options="--rm -p 8080:8080" -- java -jar /app/demo.jar
```

Note that the [docker parameters](https://docs.docker.com/engine/reference/run/) are set via the `--docker-options` option, while the startup command is preceded by two hyphens.

`werf.yaml` can describe any number of images. Use a positional argument of the command (`werf run basicapp ...`) to run a container related to the specific `image` described in `werf.yaml`.

_You may notice that `werf run` also performs the build, i.e., no pre-build is needed._

Now you can access the application locally on port 3000:

![](/guides/images/springboot/100_10_app_in_browser.png)

<div id="go-forth-button">
    <go-forth url="20_cluster.html" label="Preparing the cluster" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
