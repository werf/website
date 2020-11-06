---
title: Подготовка кластера
permalink: gitlab_nodejs/210_cluster.html
---

Для дальнейшей работы нам понадобится: kubernetes кластер и хранилище образов (registry).

**Настройка окружения — это сложная задача**. Не рекомендуем закапываться в решение проблем своими руками: используйте существующие на рынке услуги, позовите на помощь коллег, которые имеют практику настройки инфраструктуры, или задавайте вопросы в [Telegram-сообществе werf](https://t.me/werf_ru).

В целях обучения мы рекомендуем воспользоваться одним из двух вариантов:

- Если у вас достаточно мощный компьютер — установить и настроить minikube
- В противном случае — воспользоваться предложением одного из облачных провайдеров


Ремарка:

- позже мы будем брать CI-систему и регистри будем юзать оттуда
- и кластер будем брать взрослый







## Выбираем Minikube

Установить миникуб, включить регистри. Всё сделается само.

TODO: тут как-то дать пошаговку что делать с миникубом, как его ставить https://github.com/kubernetes/minikube#installation

TODO: запускаем команду в кубецтл и танцуемы


## Выбираем один из cloud provider-ов

<!--  CLOUDPROVIDER -->

Поднять кластер и регистри в клауде.


По-любому добываем из регистри ключики — нам туда пушить надо будет

И надо будет кубконфиг добыть
Если вы ещё ни разу не устанавливали Kubernetes самостоятельно и/или не обладаете опытом системного администрирования — не стоит пытаться освоить такую объемную тему в рамках этого самоучителя. Это может быть слишком сложно.

Существует огромное количество услуг, позиционирующихся как [Managed Kubernetes](https://www.google.com/search?q=managed+kubernetes). Часть из них включает серверные мощности, часть — нет.

Проще всего взять услугу в AWS (EKS) или Google Cloud (GKE). При первой регистрации они представляют бонус, которого должно хватить на неделю-другую работы с кластером. Однако для регистрации понадобится ввести данные банковской карты.

Альтернативой может стать использование Yandex.Cloud: здесь также предоставляется Managed Kubernetes и тоже требуется ввести данные карты. Но в качестве карты можно воспользоваться, например, виртуальной картой от Яндекс.Денег, позволяющей проводить «подписочные» платежи внутри России.

Также можно попробовать развернуть Kubernetes самостоятельно [в Hetzner на основании их статьи](https://community.hetzner.com/tutorials/install-kubernetes-cluster) — это один из самых дешёвых облачных провайдеров. Однако надо понимать, что вам придётся самостоятельно разбираться с большим пластом работ по администрированию платформы. Если решите переводить свой production — найдите надёжного провайдера или команду, которая будет администрировать платформу.

TODO: как добыть .kube/config?

Поддерживаниемые докер регистри https://ru.werf.io/documentation/advanced/supported_registry_implementations.html

<!--  .kube/config -->
{% offtopic title="Что за .kube/config?" %}
Это файл, хранящий реквизиты доступа к кластеру. Если разместить его по пути `~/.kube/config` - утилиты, работающие с kubernetes смогут подключиться к кластеру. TODO: ссылка на более подробное.

Файл выглядит примерно так:

{% snippetcut name=".kube/config" url="#" %}
{% raw %}
```yaml
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ALotOfNumbersAndLettersAndSoOnAveryVERYveryLongStringInBase64=
    server: https://127.0.0.1:6445
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: ManyNumbersAndLettersAVERYveryVerylongString=
    client-key-data: ManyLettersAndNumbersAveryVeryVERYlongString=
```
{% endraw %}
{% endsnippetcut %}

Обратите внимание, что кроме ключей для подключения, здесь содержится путь до API kubernetes, например, `server: https://127.0.0.1:6445`. Важно, чтобы этот путь был корректным.

{% endofftopic %}
<!-- / .kube/config -->



TODO: .kube/config и кладём куда надо

TODO: запускаем команду в кубецтл и танцуемы





Затем **зашифруйте `.kube/config` с помощью base64**. Если лень пользоваться консолью, можно воспользоваться каким-то из [веб-сервисов](https://www.base64encode.org/). Результатом станет что-то вроде:

{% snippetcut name=".kube/config (base64)" url="#" %}
{% raw %}
```yaml
YXBpVmVyc2lvbjogdjEKY2x1c3RlcnM6Ci0gY2x1c3RlcjoKICAgIGNlcnRpZmljYXRlLWF1dGhvcml0eS1kYXRhOiBBTG90T2ZOdW1iZXJzQW5kTGV0dGVyc0FuZFNvT25BdmVyeVZFUll2ZXJ5TG9uZ1N0cmluZ0luQmFzZTY0PQogICAgc2VydmVyOiBodHRwczovLzEyNy4wLjAuMTo2NDQ1CiAgbmFtZToga3ViZXJuZXRlcwpjb250ZXh0czoKLSBjb250ZXh0OgogICAgY2x1c3Rlcjoga3ViZXJuZXRlcwogICAgdXNlcjoga3ViZXJuZXRlcy1hZG1pbgogIG5hbWU6IGt1YmVybmV0ZXMtYWRtaW5Aa3ViZXJuZXRlcwpjdXJyZW50LWNvbnRleHQ6IGt1YmVybmV0ZXMtYWRtaW5Aa3ViZXJuZXRlcwpraW5kOiBDb25maWcKcHJlZmVyZW5jZXM6IHt9CnVzZXJzOgotIG5hbWU6IGt1YmVybmV0ZXMtYWRtaW4KICB1c2VyOgogICAgY2xpZW50LWNlcnRpZmljYXRlLWRhdGE6IE1hbnlOdW1iZXJzQW5kTGV0dGVyc0FWRVJZdmVyeVZlcnlsb25nU3RyaW5nPQogICAgY2xpZW50LWtleS1kYXRhOiBNYW55TGV0dGVyc0FuZE51bWJlcnNBdmVyeVZlcnlWRVJZbG9uZ1N0cmluZz0=
```
{% endraw %}
{% endsnippetcut %}





## Домен

Мы предполагаем, что к началу прохождения самоучителя вы привязали к своёму кластеру какой-то домен.

{% offtopic title="Что это значит?" %}
Ожидается, что у вас есть домен, к которому вы привязываете DNS-сервер, позволяющий настраивать [ресурсные записи DNS](https://ru.wikipedia.org/wiki/%D0%A2%D0%B8%D0%BF%D1%8B_%D1%80%D0%B5%D1%81%D1%83%D1%80%D1%81%D0%BD%D1%8B%D1%85_%D0%B7%D0%B0%D0%BF%D0%B8%D1%81%D0%B5%D0%B9_DNS). Он может быть предоставлен вашим регистратором доменов или облачной платформой, а может, вы самостоятельно его установили и настраиваете (см. также [как работает DNS](https://firstwiki.ru/index.php/%D0%9A%D0%B0%D0%BA_%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D0%B0%D0%B5%D1%82_DNS)).
Это также означает, что у вашего кластера есть IP-адреса, по которым он доступен.

Требуется прописать правильные A-записи в DNS-сервере, чтобы нужные домены были направлены на кластер.
{% endofftopic %}

В самоучителе мы будем использовать следующие домены:

- mydomain.io — направлен на кластер;
- *.mydomain.io — направлены на кластер;
- registry.mydomain.io — для GitLab Registry;
- gitlab.mydomain.io — для GitLab.














---
TODO то, что ниже входит в историю про клауды но очень сложно

## Доступ куба к регистри

TODO: перед тем как пушить нам надо registrysecret прокидывать с локали 


Для того, чтобы кластер имел доступ к собранным образам, необходимо создать ключ доступа к Registry и прописать его в кластер. Этот ключ мы назовём `registrysecret`.

Сперва нужно создать API-ключ в GitLab: зайдите в настройки пользователя (`Settings`) и в разделе `Personal Access Tokens` создайте API-ключ с правами на `read_registry`. Корректнее всего создать отдельного служебного пользователя, чтобы не завязываться на персональный аккаунт.

Полученный ключ должен быть прописан в **каждом** пространстве имён в Kubernetes, куда осуществляется деплой, в виде объекта Secret. Сделать это можно, выполнив на master-узле кластера команду:

```bash
kubectl create secret docker-registry registrysecret -n <namespace> --docker-server=<registry_domain> --docker-username=<account_email> --docker-password=<account_password> --docker-email=<account_email>
```

Здесь:

- `<namespace>` — название пространства имён в Kubernetes (например, `werf-guided-project-production`);
- `<registry_domain>` — домен Registry (например, `registry.gitlab.com`);
- `<account_email>` — email вашей учётной записи в GitLab;
- `<account_password>` — созданный API-ключ.

В каждом Deployment'е также указывается имя секрета:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/020-basic-1/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
    spec:
      imagePullSecrets:
      - name: "registrysecret"
```
{% endraw %}
{% endsnippetcut %}

{% offtopic title="Это каждый раз копировать руками?" %}

В идеале проблема копирования секретов должна решаться на уровне платформы, но если это не сделано — можно прописать нужные команды в CI-процесс.

Вариант решения — завести секрет один раз в пространство имён `kube-system`, а затем в `.gitlab-ci.yaml` при деплое копировать этот секрет:

```bash
kubectl get ns ${CI_ENVIRONMENT_CUSTOM_SLUG:-${CI_ENVIRONMENT_SLUG}} || kubectl create namespace ${CI_ENVIRONMENT_CUSTOM_SLUG:-${CI_ENVIRONMENT_SLUG}}
kubectl get secret registrysecret -n kube-system -o json |
                      jq ".metadata.namespace = \"${CI_ENVIRONMENT_CUSTOM_SLUG:-${CI_ENVIRONMENT_SLUG}}\"|
                      del(.metadata.annotations,.metadata.creationTimestamp,.metadata.resourceVersion,.metadata.selfLink,.metadata.uid)" |
                      kubectl apply -f -
```
{% endofftopic %}


<a name="iac-debug-deploy" />

---

TODO: а вот гитлаб нам не нужен пока что. Об этом надо сделать ремарку чуть ли не в самом начале.

--------



------

В каждом репозитории с кодом, с которым вы будете работать, вам надо будет **прописать в GitLab'е специальную переменную**. Зайдите в репозиторий и в левой панели нажмите `Settings` -> `CI/CD`. В этом разделе, в главе `Variables`, нужно прописать переменную окружения `KUBECONFIG_BASE64` с нашим `.kube/config`, зашифрованным base64.

Так доступ к кластеру попадет в werf, и с его помощью мы сможем строить CI-процесс.



<div id="go-forth-button">
    <go-forth url="210_cluster.html" label="Сборка" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
