# werf website

This repository contains files required to build the [werf CI/CD tool website](https://werf.io/)
and run it locally.

It includes all of website pages with the only exception for the official
[werf documentation](https://werf.io/documentation/) which can be found
in the [main werf](https://github.com/werf/werf/tree/main/docs) repo since it's
versioned and provided together with werf releases themselves (e.g. we have
the documentation for werf v1.2.51, v1.1.23, and so on).

The website is based on [Jekyll](https://jekyllrb.com/). We have also implemented our
own [jekyll_include_plugin](https://github.com/flant/jekyll_include_plugin).

## werf guides

This repo is also as a source for the content displayed in our [werf guides](https://werf.io/guides/)
which are in-depth free online tutorials to learn how you can use werf to deploy
your apps to Kubernetes.

werf guides will work even for those who have no experience with Kubernetes. Combining
the theory and practice of development (Dev) and operation (Ops), they will guide you
step by step from learning the basics of werf and Kubernetes to fully benefit
from running highly available apps in K8s clusters.

werf guides provide instructions for different programming languages/frameworks
including:
* JavaScript (Node.js),
* Java (Spring Boot),
* Python (Django),
* Go,
* Ruby (Ruby on Rails),
* PHP (Laravel).

All guides' articles are currently available in two languages, [English](https://werf.io/guides/)
and [Russian](https://ru.werf.io/guides/).

# Contributions

Any fixes, improvements, and feedback are much appreciated!

Please propose your ideas using issues and PRs in this repo. If you're interested
in running werf website locally, please refer to [DEVELOPMENT.md](DEVELOPMENT.md).

Feel free to reach us in [this Telegram chat](https://t.me/werf_io) if you have
any questions. _(There is a Russian-speaking [Telegram chat](https://t.me/werf_ru)
as well.)_
