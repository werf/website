---
title: Logging
permalink: rails/200_real_apps/020_logging.html
examples_initial: examples/rails/010_basic_app
examples: examples/rails/020_logging
description: |
  In this chapter, we will learn about application logging in Kubernetes and implement it. Additionally, we will introduce a structured logging format to make it ready for parsing by log collection and analysis systems.
---

## Redirecting logs to stdout

All applications deployed to a Kubernetes cluster must write logs to [stdout/stderr](https://en.wikipedia.org/wiki/Standard_streams). Sending logs to standard streams makes them accessible by Kubernetes and log collection systems. Such an approach keeps logs from being deleted when containers are recreated and prevents consuming all available storage on Kubernetes nodes if the output is sent to log files in containers.

By default, Rails saves logs to a file instead of stdout/stderr. Set the appropriate environment variable (`RAILS_LOG_TO_STDOUT=1`) to enable redirecting logs (including errors) to stdout.

Since streaming logs to stdout is all we need, they will be directed to stdout **in all cases** regardless of the environment variables:

{% include snippetcut_example path="config/environments/production.rb" syntax="ruby" snippet="logger" examples=page.examples %}

## Formatting logs

By default, Rails-based application generate plain text logs:
```shell
=> Booting Puma
=> Rails 6.1.4 application starting in production
=> Run "bin/rails server --help" for more startup options
Puma starting in single mode...
* Puma version: 5.3.2 (ruby 2.7.4-p191) ("Sweetnighter")
*  Min threads: 5
*  Max threads: 5
*  Environment: production
*          PID: 1
* Listening on http://0.0.0.0:3000
Use Ctrl-C to stop
I, [2021-06-07T16:47:28.498897 #1]  INFO -- : Started GET \"/ping\" for 192.168.49.1 at 2021-07-23 15:08:16 +0000
I, [2021-06-07T16:47:28.498972 #1]  INFO -- : Processing by ApplicationController#ping as */*
I, [2021-06-07T16:47:28.498897 #1]  INFO -- : Completed 200 OK in 0ms (Views: 0.1ms | Allocations: 166)
```

Note how the application logs end up mixed with Rails and Puma logs. All these logs have different formats. Regular log collection and analysis systems will probably struggle with parsing all this gibberish.

You can solve this problem by shipping logs in a structured JSON-like format. Most log collection systems can both, easily parse JSON and correctly process logs/messages in other (unexpected, unstructured) formats when they sneak in between JSON-formatted logs.

Let's define an `ActiveSupport::Logger::SimpleFormatter` class that converts plain text logs into JSON-formatted ones:
{% include snippetcut_example path="lib/json_simple_formatter.rb" syntax="ruby" examples=page.examples %}

> _Log tagging support is beyond the scope of this guide; however, you can implement it if necessary._

Now let's make the new `JSONSimpleFormatter` class the default one to convert all logs to JSON: 
{% include snippetcut_example path="config/environments/production.rb" syntax="ruby" snippet="log_formatter" examples=page.examples %}

## Managing the logging level

By default, the application logging level for the production environment is set to `info`. However, sometimes you may wish to change that.

For example, switching the logging level from `info` to `debug` can provide additional debugging information and help in troubleshooting problems in production. However, if the application has many replicas, switching them all to the `debug` level may not be the best idea. It may affect security and significantly increase the load on log collection, storage, and analysis components.

You can use an environment variable to set the logging level, thus solving the above problem. Using this approach, you can run a single Deployment replica with the `debug` logging level next to the existing replicas with the `info` logging level enabled. Moreover, you can disable centralized log collection for this new Deployment (if any). Together, all these measures prevent overloading log collection systems while keeping debug logs containing potentially sensitive information from being streamed to them.

Here is how you can set the logging level using the `RAILS_LOG_LEVEL` environment variable:
{% include snippetcut_example path="config/environments/production.rb" syntax="ruby" snippet="log_level" examples=page.examples %}

If the environment variable is omitted, the `info` level is used by default.

## Log filtering

Let's enable log filtering to filter out secrets from the log file:
{% include snippetcut_example path="config/initializers/filter_parameter_logging.rb" syntax="ruby" examples=page.examples %}

> _Keep in mind that you need to update this list after adding new secrets._

## Displaying logs in werf-based deploys

By default, when deploying, werf prints logs of all application containers until they become `Ready`. 
You can filter these logs using custom werf annotations. In this case, werf will only print lines that match the specified patterns.
Additionally, you can enable log output for specific containers.

For example, here is how you can disable log output for the `container_name` container during the deployment:
```yaml
annotations:
  werf.io/skip-logs-for-containers: container_name
```

The example below shows how you can enable printing lines that match a pre-defined regular expression:
```yaml
annotations:
  werf.io/log-regex: ".*ERROR.*"
```

A list of all available annotations can be found [here](https://werf.io/documentation/v1.2/reference/deploy_annotations.html).

> _Note that these annotations only influence the way logs are shown during the werf-based deployment process. They do not affect the application being deployed or its configuration in any way. You can still use stdout/stderr streams of the container to view raw logs._

## Checking whether the application is running

Let's deploy our application:
```shell
werf converge --repo <Docker Hub username>/werf-guide-app
```

You should see the following output:
```shell
...
┌ ⛵ image app
│ ┌ Building stage app/dockerfile
│ │ app/dockerfile  Sending build context to Docker daemon  30.72kB
│ │ app/dockerfile  Step 1/13 : FROM ruby:2.7
│ │ app/dockerfile   ---> 1faa5f2f8ca3
...
│ │ app/dockerfile  Step 13/13 : LABEL werf-version=v1.2.11+fix10
│ │ app/dockerfile   ---> Running in db6e76c3f427
│ │ app/dockerfile  Removing intermediate container db6e76c3f427
│ │ app/dockerfile   ---> d7f69fbfdedb
│ │ app/dockerfile  Successfully built d7f69fbfdedb
│ │ app/dockerfile  Successfully tagged eb50cb50-1191-4d0b-8bf2-a4d5bba18ecf:latest
│ ├ Info
│ │      name: .../werf-guide-app:...
│ │        id: d7f69fbfdedb
│ │   created: 2021-07-23 18:00:19 +0300 MSK
│ │      size: 327.3 MiB
│ └ Building stage app/dockerfile (12.72 seconds)
└ ⛵ image app (17.81 seconds)

┌ Waiting for release resources to become ready
│ ┌ Status progress
│ │ DEPLOYMENT      REPLICAS  AVAILABLE  UP-TO-DATE
│ │ werf-guide-app  1/1       1          1
│ │ │   POD                         READY  RESTARTS  STATUS
│ │ ├── guide-app-5f97776488-vwqfg  1/1    0         Terminating
│ │ └── guide-app-fcf7c4ff5-wvb62   1/1    0         Running
│ └ Status progress
└ Waiting for release resources to become ready (4.89 seconds)

Release "werf-guide-app" has been upgraded. Happy Helming!
NAME: werf-guide-app
LAST DEPLOYED: Fri Jul 23 18:00:34 2021
NAMESPACE: werf-guide-app
STATUS: deployed
REVISION: 14
TEST SUITE: None
Running time 26.67 seconds
```

Make several requests in order to generate some logging data:
```shell
curl http://werf-guide-app/ping       # returns "pong" + 200 OK status code
curl http://werf-guide-app/not_found  # no response; returns 404 Not Found
```

Now let's take a look at the logs:
```shell
kubectl logs deploy/werf-guide-app
```

You should see the following output:
```shell
=> Booting Puma
=> Rails 6.1.4 application starting in production
=> Run `bin/rails server --help` for more startup options
Puma starting in single mode...
* Puma version: 5.3.2 (ruby 2.7.4-p191) ("Sweetnighter")
*  Min threads: 5
*  Max threads: 5
*  Environment: production
*          PID: 1
* Listening on http://0.0.0.0:3000
Use Ctrl-C to stop
{"type":"INFO","time":"2021-07-23T15:11:36+00:00","message":"Started GET \"/ping\" for 192.168.49.1 at 2021-07-23 15:11:36 +0000"}
{"type":"INFO","time":"2021-07-23T15:11:36+00:00","message":"Processing by ApplicationController#ping as */*"}
{"type":"INFO","time":"2021-07-23T15:11:36+00:00","message":"Completed 200 OK in 0ms (Views: 0.1ms | Allocations: 135)"}
{"type":"INFO","time":"2021-07-23T15:11:30+00:00","message":"Started GET \"/not_found\" for 192.168.49.1 at 2021-07-23 15:11:30 +0000"}
{"type":"FATAL","time":"2021-07-23T15:11:30+00:00","message":"  \nActionController::RoutingError (No route matches [GET] \"/not_found\"):\n  "}
```

Note that application logs are now rendered in JSON format, and most log processing systems can easily parse them. At the same time, Rails and Puma logs are streamed in plain text just like before. The main advantage of this approach is that log processing systems will no longer try to parse application logs and Rails/Puma logs as if they have the same format. JSON logs will be stored separately, letting you perform searching/filtering based on the selected fields.
