---
title: First steps
permalink: nodejs/100_basic.html
---

In this chapter, you will acquire the basic skills required to deliver an application to Kubernetes:
- You will learn how to containerize an application. For this, we will define the assembly process, build a container, and run the application in Docker.
- You will deploy an application. We will configure the application architecture by describing its components/their relationships and deploy an application to Kubernetes.

Moreover, we will do it all on the local machine without using any CI system.

But before we start, you need to [install werf]({{ site.docsurl }}/installation.html) locally.

## The application

This chapter deals with an [elementary Node.js example application](https://github.com/werf/werf-guides/tree/master/examples/nodejs/000_app). It provides the basic HTTP REST API (in the JSON format) that performs CRUD operations on some objects (`/api/labels`). We will use SQLite as a database. We will specify key parameters using environment variables (in our case, it is the `SQLITE_FILE` variable - the path to the SQLite database).

For the sake of simplicity, we will not preserve the data in the SQLite database. Thus, the data will be deleted after each redeployment of the application. We will cover the issues of data storage and preservation in the chapter "Full-fledged applications".

