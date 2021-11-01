<p align="center">
  <img src="https://raw.githubusercontent.com/werf/werf-guides/master/assets/images/werf-logo.svg?sanitize=true" style="max-height:100%;" height="175">
</p>
___

## Development

### How do I run the guides site locally?

- Install [werf](http://werf.io/installation.html)
- Run:
  ```shell
  werf compose up jekyll_base
  ```
- Check the English version is available on [https://localhost](http://localhost), and the Russian version on [http://ru.localhost](https://ru.localhost) (add `ru.localhost` record in your `/etc/hosts` to access the Russian version of the site.)
