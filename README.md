<p align="center">
  <img src="https://raw.githubusercontent.com/werf/werf-guides/master/assets/images/werf-logo.svg?sanitize=true" style="max-height:100%;" height="175">
</p>
___

## Development

### How do I run the guides site locally?

> Add `ru.localhost` record in your `/etc/hosts` to access the Russian version of the site. 

- Install [werf](https://werf.io/installation.html)
- Run:
  ```shell
  werf compose up --follow --docker-compose-command-options="-d" jekyll_base
  ```
- Check the English version is available on [https://localhost](https://localhost), and the Russian version on [https://ru.localhost](https://ru.localhost)
