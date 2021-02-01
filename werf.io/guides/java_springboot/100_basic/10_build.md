---
title: Сборка образа
permalink: java_springboot/100_basic/10_build.html
---

{% filesused title="Files mentioned in the chapter" %}
- Dockerfile
- werf.yaml
{% endfilesused %}

[Dockerfile](https://docs.docker.com/engine/reference/builder/) is a classic method to build images. Probably, your applications are already using this mechanism for the assembly. That is why we will start with it and then learn how to use it with werf. In the next chapters, we will speed up the build using an alternative syntax for describing the assembly process, but in this chapter, we will focus on getting the result quickly.

## Preparing the workplace

We assume that you have already [installed werf]({{ site.docsurl }}/installation.html) and Docker.

Create a directory on your computer and follow these steps:

```shell
git clone git@github.com:werf/werf-guides.git
cp -r werf-guides/examples/springboot/000_app ./
cd 000_app 
git init
git add .
git commit -m "initial commit"
```

_This way you will copy the code of the [SpringBoot application](https://github.com/werf/werf-guides/tree/master/examples/springboot/000_app) to a local directory and initialize a Git repository in it._

Note that werf follows the principles of [giterminism]({{ site. docsurl }}/documentation/advanced/configuration/giterminism.html): it fully relies on the state described in the Git repository. This means that files not committed to the Git repository will be ignored by default. Thereby, if you have the source code, then you can turn an application to the specific operating condition at any time.

## Dockerfile-based build process

The build configuration of our application consists of the following steps:

- pull the OpenJDK image (e.g., `gradle:jdk8-openj9`);
- add the application code to it;
- build the application using gradle and move the resulting jar to the appropriate location;
- configure the application using environment variables.

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

{% offtopic title="Что тут написано?" %}

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

- **_project_** is a field containing the unique name of the application project. The default project name is used when generating the Helm release name and the namespace in Kubernetes. Changing the name of an active project is a tedious task that requires a number of manual actions (more information about possible consequences is available [here]({{ site.docsurl }}/documentation/reference/werf_yaml. html#warning-on-changing-project-name));
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

Now that you have successfully added `Dockerfile` and `werf.yaml` described above, it is necessary to commit changes to Git:

{% raw %}
```bash
git add .
git commit -m "work in progress"
```
{% endraw %}

The [`build` command]({{ site.docsurl }}/documentation/reference/cli/werf_build.html) starts the assembly process:


{% raw %}
```bash
werf build
```
{% endraw %}

_The sub-chapter "Speeding up the build" contains instructions on how to adapt the Dockerfile-based build process to an alternative werf syntax called `Stapel` and gain access to some advanced features, such as Git history-based incremental rebuilds, the usage of Ansible and inter-assembly cache, convenient diagnostic tools, and much more._

But even now, you may notice that werf outputs the build logs in the extended format:

```
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

## Running

Let's run the built image using the [werf run]({{ site.docsurl }}/documentation/cli/main/run.html) command:

```bash
werf run --docker-options="--rm -p 8080:8080" -- java -jar /app/demo.jar
```

Note that we set the [docker parameters](https://docs.docker.com/engine/reference/run/) via `--docker-options`, while the startup command is preceded by two hyphens.

_You may also notice that `werf run` also performs the build. In other words, it is not really necessary to run the build separately._

Now you can access the application locally on port 3000:

![](/guides/images/springboot/100_10_app_in_browser.png)

## Making changes

As you might guess, we are going to continually update our application. Let's see how to do this in the right way by making some arbitrary changes to the application code:

{% snippetcut name="/src/main/java/com/example/demo/mvc/controller/LabelController.java" url="#" %}
{% raw %}
```
    @GetMapping("/labels")
    public List<Labels> labels() {
        return "Our changes";
    }
```
{% endraw %}
{% endsnippetcut %}

1. Stop the running `werf run` (by pressing Ctrl+C in the console where it is running.
2. Start it again: 
    ```bash
    werf run --docker-options="-d -p 3000:3000 --restart=always" -- node /app/app.js
    ```
2. Watch as the application is being rebuilt and restarted, and then connect to the API: http://example.com:3000/labels
3. You probably expect to see the `Our changes` message, but it isn't there. **Everything is the same**... but why?

The thing is we **forgot to commit changes to Git prior to step 1** in the scenario above.

{% offtopic title="What is the correct way, and why go through all those troubles?" %}
1. Make changes to the code.
2. Commit them:
   ```shell
   git add .
   git commit -m "wip"
   ```
3. Restart `werf run`:
    ```shell
    werf run --docker-options="-d -p 8080:8080 --restart=always" -- java -jar /app/demo.jar
    ```
4. View the result in the browser.

A strict binding to Git ensures the reproducibility of each specific solution. More details about _giterminism_ mechanics are available in the "What you need to know" chapter. Until then, we'll focus on building and delivering an application to the cluster.
{% endofftopic %}


<div id="go-forth-button">
    <go-forth url="20_cluster.html" label="Preparing the cluster" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
