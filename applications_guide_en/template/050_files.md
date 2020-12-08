---
title: Работа с файлами
sidebar: applications_guide
guide_code: ____________
permalink: ____________/050_files.html
toc: false
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/secret-values.yaml
- ____________
- ____________
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении работу с пользовательскими файлами. Для этого потребуется персистентное (постоянное) хранилище.

В идеале — нужно добиться, чтобы приложение было stateless, а данные хранились в S3-совместимом хранилище — например, [MinIO](https://github.com/minio/minio) или AWS S3. Это обеспечивает простое масштабирование, работу в HA-режиме и высокую доступность.

{% offtopic title="А есть какие-то способы кроме S3?" %}
Первый и более общий способ — это использовать как [volume](https://kubernetes.io/docs/concepts/storage/volumes/) хранилище [NFS](https://kubernetes.io/docs/concepts/storage/volumes/#nfs), [CephFS](https://kubernetes.io/docs/concepts/storage/volumes/#cephfs) или [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath).

Мы не рекомендуем этот способ, потому что при возникновении неполадок с такими типами volume’ов они влияют на работоспособность контейнера и всего демона Docker в целом. Тогда могут пострадать приложения, не имеющие никакого отношения к вашему.

Более надёжный путь — пользоваться S3. Так мы используем отдельный сервис, который имеет возможность масштабироваться, работать в HA-режиме и иметь высокую доступность. Можно воспользоваться облачным решением вроде AWS S3, Google Cloud Storage, Microsoft Blobs Storage и т.д.

Природа Kubernetes (и оркестровки контейнеров в целом) такова, что если мы будем сохранять файлы в какой-либо директории у приложения (запущенного в Kubernetes), то после перезапуска контейнера все изменения пропадут.
{% endofftopic %}

Данная настройка производится полностью в рамках приложения. ____________

____________
____________
____________
____________
____________
____________

Для работы с S3 необходимо пробросить ключи доступа в приложение. Для этого стоит использовать [механизм секретных переменных]({{ site.docsurl }}/documentation/reference/deploy_process/working_with_secrets.html). *Вопрос работы с секретными переменными рассматривался подробнее, когда мы [делали базовое приложение](020_basic.html#secret-values-yaml).*

{% snippetcut name="secret-values.yaml (расшифрованный)" url="#" ignore-tests %}
{% raw %}
```yaml
app:
  ____________
  ____________
  ____________
```
{% endraw %}
{% endsnippetcut %}

Несекретные значения — храним в `values.yaml`:

{% snippetcut name="values.yaml" url="#" %}
{% raw %}
```yaml
____________
____________
____________
```
{% endraw %}
{% endsnippetcut %}

После того, как значения корректно прописаны и зашифрованы, можно пробросить соответствующие значения в Deployment:

{% snippetcut name="deployment.yaml" url="#" %}
{% raw %}
```yaml
        ____________
        ____________
        ____________
```
{% endraw %}
{% endsnippetcut %}

____________
____________
____________

<div id="go-forth-button">
    <go-forth url="060_email.html" label="Работа с электронной почтой" base-url="applications_guide_ru"></go-forth>
</div>
