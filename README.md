# werf website

This repository contains files required to build main part of the werf.io website including [werf guides](https://werf.io/guides/) — are in-depth free online tutorials to learn
how you can use [werf CI/CD tool](https://werf.io/) to deploy your apps to Kubernetes.

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

# werf website repo

This repo serves as a source for all the content generated in werf.io excluding the documentation (all under the /documentation path). The source for the werf documentation (https://werf.io/documentation/) is [here...](https://github.com/werf/werf/tree/main/docs)  

The website is based on [Jekyll](https://jekyllrb.com/). We have also
implemented our own [jekyll_include_plugin](https://github.com/flant/jekyll_include_plugin).

# Contributions

Any feedback is much appreciated!

You are also welcome to propose your ideas for them using issues and PRs in this repo.
If you're interested in running werf website locally, please refer
to [DEVELOPMENT.md](DEVELOPMENT.md).

Feel free to reach us in [this Telegram chat](https://t.me/werf_io) if you have
any questions. _(There is a Russian-speaking [Telegram chat](https://t.me/werf_ru)
as well.)_
