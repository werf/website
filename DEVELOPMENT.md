<p align="center">
  <img src="https://raw.githubusercontent.com/werf/website/main/assets/images/werf-logo.svg?sanitize=true" style="max-height:100%;" height="175">
</p>
___

## Development

- Install [werf](http://werf.io/docs/index.html).
- Install [task](https://taskfile.dev/installation/).

### Local development setup

Run `jekyll serve` with --watch option to test changes in "real time":

- Run:
  ```shell
  task dev:setup:website
  ```
- Check the English version is available on [https://localhost](http://localhost), and the Russian version on [http://ru.localhost](https://ru.localhost) (add `ru.localhost` record in your `/etc/hosts` to access the Russian version of the site).

Add configurator content:
- Run:
  ```shell
  task local:gen:configurator
  ```
  Run `task dev:cleanup:website` to cleanup working directory.

Optionally serve documentation content:

- Run in werf/werf repository:
  ```shell
  cd ../werf/werf
  task dev:setup:website
  ```

### Cleanup

- `task dev:cleanup:website` - stop all containers and delete shared docker network.

### Spell-check management

We use a `./scripts/docs/spelling/wordlist.txt` file as the base for creating a custom spell-check dictionary.

- Add any new words directly to `wordlist.txt` in any order.
- Run:
  ```shell
  task site:generate-special-dictionary
  ```