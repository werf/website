---
title: How it works
permalink: how_it_works.html
layout: plain
banner: guides
breadcrumbs: none
---

{% asset introduction.css %}
{% asset introduction.js %}

<div id="introduction-presentation" class="introduction-presentation">
    <div id="introduction-presentation-controls" class="introduction-presentation__controls">
        <a href="javascript:void(0)" class="introduction-presentation__controls-nav">
            <img src="{% asset introduction/nav.svg @path %}" />
        </a>
        <div class="introduction-presentation__controls-stage">
            Define desired state
        </div>
        <div class="introduction-presentation__controls-step">
            0. Pre-configuration
        </div>
        <div class="introduction-presentation__controls-selector">
            <div class="introduction-presentation__controls-selector-stage">
                Define desired state
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="0"
                    data-presentation-selector-stage="Define desired state">
                    0. Pre-configuration
                </a>
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="1"
                    data-presentation-selector-stage="Define desired state">
                    1. Put Dockerfiles into the application repository
                </a>
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="2"
                    data-presentation-selector-stage="Define desired state">
                    2. Create the werf.yaml configuration file
                </a>
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="3"
                    data-presentation-selector-stage="Define desired state">
                    3. Describe Helm chart templates to deploy your app
                </a>
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="4"
                    data-presentation-selector-stage="Define desired state">
                    4. Commit
                </a>
            </div>
            <div class="introduction-presentation__controls-selector-stage">
                Converge
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="5"
                    data-presentation-selector-stage="Converge">
                    1. Calculate configured images based on the current Git commit state.
                </a>
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="6"
                    data-presentation-selector-stage="Converge">
                    2. Read images that are already available in the container registry.
                </a>
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="7"
                    data-presentation-selector-stage="Converge">
                    3. Calculate the difference between images required for the current Git commit state and those that are already available in the container registry.
                </a>
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="8"
                    data-presentation-selector-stage="Converge">
                    4. Build and publish only those images that do not exist in the container registry (if any).
                </a>
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="9"
                    data-presentation-selector-stage="Converge">
                    5. Read the target configuration of Kubernetes resources based on the current Git commit state.
                </a>
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="10"
                    data-presentation-selector-stage="Converge">
                    6. Read the configuration of existing Kubernetes resources from the cluster.
                </a>
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="11"
                    data-presentation-selector-stage="Converge">
                    7. Calculate the difference between the target configuration of Kubernetes resources defined for the current Git commit state and the configuration of existing Kubernetes resources in the cluster.
                </a>
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="12"
                    data-presentation-selector-stage="Converge">
                    8. Apply changes (if any) to Kubernetes resources so that they conform to the state defined in the Git commit.
                </a>
            </div>
            <div class="introduction-presentation__controls-selector-step">
                <a href="javascript:void(0)"
                    data-presentation-selector-option="13"
                    data-presentation-selector-stage="Converge">
                    9. Make sure all resources are ready to operate and report any errors immediately (with an option to roll back to the previous state if an error occurred).
                </a>
            </div>
        </div>
    </div>
    <div class="introduction-presentation__container">
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text">
                werf configuration should be stored in the application Git repository along with the application code.
            </div>
            <img src="{% asset introduction/s-1.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text"></div>
            <img src="{% asset introduction/s-2.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text">
<div markdown="1">
Pay attention to the crucial `project` parameter that stores a _project name_ since werf will be using it  extensively during the _converge process_. Changing this name later, when your application is up and running, involves downtime, and requires re-deploying the application.
</div>
            </div>
            <img src="{% asset introduction/s-3.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text">
<div markdown="1">
`werf_image` is a special template function for generating the full name of the image that was built. This function has a name parameter that corresponds to the image defined in the `werf.yaml` (`"frontend"` or `"backend"` in our case).
</div>
            </div>
            <img src="{% asset introduction/s-4.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text"></div>
            <img src="{% asset introduction/s-5.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text">
<div markdown="1">
At this step, werf calculates the target image names. Image names may change or stay the same after the new commit is made depending on the changes in the Git repository. Note that each Git commit is associated with deterministic and reproducible image names.
</div>
            </div>
            <img src="{% asset introduction/s-6.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text"></div>
            <img src="{% asset introduction/s-7.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text"></div>
            <img src="{% asset introduction/s-8.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text"></div>
            <img src="{% asset introduction/s-9.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text"></div>
            <img src="{% asset introduction/s-10.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text"></div>
            <img src="{% asset introduction/s-11.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text"></div>
            <img src="{% asset introduction/s-12.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text">
                <div class="introduction-presentation__slide-title"></div>
            </div>
            <img src="{% asset introduction/s-13.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
        <div class="introduction-presentation__slide">
            <div class="introduction-presentation__slide-text">
                <div class="introduction-presentation__slide-title"></div>
            </div>
            <img src="{% asset introduction/s-14.svg @path %}"
            class="introduction-presentation__slide-img" />
        </div>
    </div>
</div>
<div markdown="1">

## What's next?

Deploy your first demo application with this [quickstart]({{ "documentation/quickstart.html" | true_relative_url }}) or check the [guides]({{ "guides.html" | true_relative_url }}) – they cover the configuration of a wide variety of applications based on different programming languages and frameworks. We recommend finding a manual suitable for your application and follow instructions.

Refer to [this article]({{ "documentation/advanced/ci_cd/ci_cd_workflow_basics.html" | true_relative_url }}) if you feel like you are ready to dig deeper into the general overview of CI/CD workflows that can be implemented with werf.
</div>