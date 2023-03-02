## Требования

* GitHub Actions;

* GitHub-hosted Linux Runner.

## Настройка проекта GitHub

Сохраните kubeconfig-файл для доступа к кластеру Kubernetes в [зашифрованный секрет](https://docs.github.com/en/actions/security-guides/encrypted-secrets) `KUBECONFIG_BASE64`, предварительно закодировав его в Base64.

## Конфигурация CI/CD проекта

```yaml
# .github/workflows/prod.yml:
name: prod
on:
  push:
    branches:
      - main

jobs:
  prod:
    name: prod
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0

      - uses: werf/actions/install@v{{ include.version }}
        with:
          channel: {{ include.channel }}

      - run: |
          source "$(werf ci-env github --as-file)"
          werf converge
        env:
          WERF_ENV: prod
          WERF_KUBECONFIG_BASE64: ${{ secrets.KUBECONFIG_BASE64 }}
```

```yaml
# .github/workflows/cleanup.yml:
name: cleanup
on:
  schedule:
    - cron: '0 3 * * *'

jobs:
  cleanup:
    name: cleanup
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: git fetch --prune --unshallow

      - uses: werf/actions/install@v{{ include.version }}
        with:
          channel: {{ include.channel }}

      - run: |
          source "$(werf ci-env github --as-file)"
          werf cleanup
        env:
          WERF_KUBECONFIG_BASE64: ${{ secrets.KUBECONFIG_BASE64 }}
```

TODO: больше листингов
