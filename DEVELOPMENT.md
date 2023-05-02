<p align="center">
  <img src="https://raw.githubusercontent.com/werf/website/main/assets/images/werf-logo.svg?sanitize=true" style="max-height:100%;" height="175">
</p>
___

## Development

### Run the site locally

#### Variant 1 (faster)

- Install [werf](http://werf.io/documentation/index.html). 
- Run:
  ```shell
  werf compose up jekyll_base --dev
  ```
- Or run with specific architecture (e.g. ARM-based Macbooks):
  ```shell
  werf compose up jekyll_base --dev --platform='linux/amd64'
  ```
- Check the English version is available on [https://localhost](http://localhost), and the Russian version on [http://ru.localhost](https://ru.localhost) (add `ru.localhost` record in your `/etc/hosts` to access the Russian version of the site).
- Configurator part can be added with:
  ```shell
  task local:gen:configurator
  ```
  Run `task clean` to cleanup working directory.

#### Variant 2 (slower)

- Install [werf](http://werf.io/documentation/index.html). 
- Run (add `--follow --docker-compose-command-options="-d"` if necessary):
  ```shell
  werf compose up --docker-compose-options="-f docker-compose-slow.yml" --dev
  ```
- Check the English version is available on [https://localhost](http://localhost), and the Russian version on [http://ru.localhost](https://ru.localhost) (add `ru.localhost` record in your `/etc/hosts` to access the Russian version of the site).

#### Variant 3 (the fastest but only one language)

- Install ruby and bundler
- To run English version:
  ```shell
  jekyll serve -t --disable-disk-cache --config _config.yml,_config_dev.yml,_config_en.yml -d /tmp/_site
  ```
- To run Russian version:
  ```shell
  jekyll serve -t --disable-disk-cache --config _config.yml,_config_dev.yml,_config_ru.yml -d /tmp/_site
  ```
- Site content (with guides but without documentation section) is available on http://localhost:4000
 
