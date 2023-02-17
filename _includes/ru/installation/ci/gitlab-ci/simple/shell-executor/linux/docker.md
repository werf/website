## Требования

- GitLab;

- Linux-хост для установки GitLab Runner, имеющий:
  
  - Bash;
  
  - Git версии 2.18.0 или выше;
  
  - GPG;
  
  - [Docker Engine](https://docs.docker.com/engine/install/).Bash;
  
  - Git версии 2.18.0 или выше;
  
  - GPG;
  
  - [Docker Engine](https://docs.docker.com/engine/install/).

## Установка GitLab Runner

Установите GitLab Runner на выделенный для него хост, следуя [официальным инструкциям](https://docs.gitlab.com/runner/install/linux-repository.html).

## Установка werf

Для установки werf, на хосте для GitLab Runner выполните:

```
curl -sSL https://werf.io/install.sh | bash -s -- --ci
```

## Регистрация GitLab Runner

Для регистрации GitLab Runner в GitLab следуйте [официальным инструкциям](https://docs.gitlab.com/runner/register/index.html), указав Shell в качестве executor'а. При желании после регистрации произведите [дополнительную конфигурацию](https://docs.gitlab.com/runner/configuration/advanced-configuration.html) GitLab Runner'а.

## Настройка проекта GitLab

* [Создайте и сохраните access token](https://docs.gitlab.com/ee/user/project/settings/project_access_tokens.html#create-a-project-access-token) для очистки ненужных образов из container registry со следующей конфигурацией:
  
  * **Token name:** `werf-images-cleanup`
  
  * **Role:** `developer`
  
  * **Scopes:** `api`

* [В переменных проекта](https://docs.gitlab.com/ee/ci/variables/#for-a-project) добавьте следующие переменные:
  
  * Access token для очистки ненужных образов:
    
    * **Key:** `WERF_IMAGES_CLEANUP_PASSWORD`
    
    * **Value:** `<сохранённый "werf-images-cleanup" access token>`
    
    * **Protect variable:** `yes`
    
    * **Mask variable:** `yes`

* [Добавьте плановое задание](https://docs.gitlab.com/ee/ci/pipelines/schedules.html#add-a-pipeline-schedule) на каждую ночь для очистки ненужных образов в container registry, указав ветку `main`/`master` в качестве **Target branch**.

## Конфигурация CI/CD проекта

```yaml
# .gitlab-ci.yml:
stages:
- prod
- cleanup

default:
  before_script:
  - source "$(~/bin/trdl use werf {{ include.version }} {{ include.channel }})"
  - source "$(werf ci-env gitlab --as-file)"
  tags: ["<тег зарегистрированного GitLab Runner>"]

prod:
  stage: prod
  script:
  - werf converge
  environment:
    name: prod
  rules:
  - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH && $CI_PIPELINE_SOURCE != "schedule"
    when: manual

images:cleanup:
  stage: cleanup
  script:
    - werf cr login -u nobody -p "${WERF_IMAGES_CLEANUP_PASSWORD:?}" "${WERF_REPO:?}"
    - werf cleanup
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH && $CI_PIPELINE_SOURCE == "schedule"
```

TODO: добавить другие файлы проекта