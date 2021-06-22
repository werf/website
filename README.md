<p align="center">
  <img src="https://raw.githubusercontent.com/werf/werf-guides/master/assets/images/werf-logo.svg?sanitize=true" style="max-height:100%;" height="175">
</p>
___

## Development

### Running the site locally

Up site and synchronize the state in case of local changes:

```shell
werf compose up --env development --dev --docker-compose-command-options='-d'
```

Now you can access the site:
```shell
curl localhost
```

> To access the Russian version of the site add `ru.localhost` record in your `/etc/hosts`
