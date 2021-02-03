---
title: Storing files
permalink: java_springboot/200_real_apps/50_files.html
layout: "development"
---

TODO: Глава про файлы


{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/secret-values.yaml
- .helm/values.yaml
- src/server/server.js
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении работу с пользовательскими файлами. Для этого потребуется персистентное (постоянное) хранилище.

В идеале — нужно добиться, чтобы приложение было stateless, а данные хранились в S3-совместимом хранилище — например, [MinIO](https://github.com/minio/minio) или AWS S3. Это обеспечивает простое масштабирование, работу в HA-режиме и высокую доступность.

{% offtopic title="А есть какие-то способы кроме S3?" %}
Первый и более общий способ — это использовать как [volume](https://kubernetes.io/docs/concepts/storage/volumes/) хранилище [NFS](https://kubernetes.io/docs/concepts/storage/volumes/#nfs), [CephFS](https://kubernetes.io/docs/concepts/storage/volumes/#cephfs) или [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath).

Мы не рекомендуем этот способ, потому что при возникновении неполадок с такими типами volume’ов они влияют на работоспособность контейнера и всего демона Docker в целом. Тогда могут пострадать приложения, не имеющие никакого отношения к вашему.

Более надёжный путь — пользоваться S3. Так мы используем отдельный сервис, который имеет возможность масштабироваться, работать в HA-режиме и иметь высокую доступность. Можно воспользоваться облачным решением вроде AWS S3, Google Cloud Storage, Microsoft Blobs Storage и т.д.

Природа Kubernetes (и оркестровки контейнеров в целом) такова, что если мы будем сохранять файлы в какой-либо директории у приложения (запущенного в Kubernetes), то после перезапуска контейнера все изменения пропадут.
{% endofftopic %}

Данная настройка производится полностью в рамках приложения. Рассмотрим подключение к S3 на примере MinIO.

Для этого подключим пакет `minio` в `npm`:

```bash
$ npm install minio --save
```

… и настроим работу с MinIO S3 в приложении:

{% snippetcut name="src/server/server.js" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/050-files/src/server/server.js" %}
```js
const Minio = require("minio");
const S3_ENDPOINT = process.env.S3_ENDPOINT || "127.0.0.1";
const S3_PORT = Number(process.env.S3_PORT) || 9000;
const TMP_S3_SSL = process.env.S3_SSL || "true";
const S3_SSL = TMP_S3_SSL.toLowerCase() == "true";
const S3_ACCESS_KEY = process.env.S3_ACCESS_KEY || "SECRET";
const S3_SECRET_KEY = process.env.S3_SECRET_KEY || "SECRET";
const S3_BUCKET = process.env.S3_BUCKET || "avatars";
const CDN_PREFIX = process.env.CDN_PREFIX || "http://127.0.0.1:9000";

// S3 client
var s3Client = new Minio.Client({
  endPoint: S3_ENDPOINT,
  port: S3_PORT,
  useSSL: S3_SSL,
  accessKey: S3_ACCESS_KEY,
  secretKey: S3_SECRET_KEY,
});
```
{% endsnippetcut %}

Для работы с S3 необходимо пробросить ключи доступа в приложение. Для этого стоит использовать [механизм секретных переменных]({{ site.docsurl }}/documentation/reference/deploy_process/working_with_secrets.html). *Вопрос работы с секретными переменными рассматривался подробнее, когда мы [делали базовое приложение](020_basic.html#secret-values-yaml).*

{% snippetcut name=".helm/secret-values.yaml (расшифрованный)" url="#" ignore-tests %}
{% raw %}
```yaml
app:
  s3:
    access_key:
      _default: bNGXXCF1GF
    secret_key:
      _default: zpThy4kGeqMNSuF2gyw48cOKJMvZqtrTswAQ
```
{% endraw %}
{% endsnippetcut %}

Несекретные значения — храним в `values.yaml`:

{% snippetcut name=".helm/values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/050-files/.helm/values.yaml" %}
{% raw %}
```yaml
app:
<...>
  s3:
    host:
      _default: minio
    port:
      _default: 9000
    bucket:
      _default: 'avatars'
    ssl:
      _default: 'false'
```
{% endraw %}
{% endsnippetcut %}

После того, как значения корректно прописаны и зашифрованы, можно пробросить соответствующие значения в Deployment:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/050-files/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        - name: CDN_PREFIX
          value: {{ printf "%s%s" (pluck .Values.global.env .Values.app.cdn_prefix | first | default .Values.app.cdn_prefix._default) (pluck .Values.global.env .Values.app.s3.bucket | first | default .Values.app.s3.bucket._default) | quote }}
        - name: S3_SSL
          value: {{ pluck .Values.global.env .Values.app.s3.ssl | first | default .Values.app.s3.ssl._default | quote }}
        - name: S3_ENDPOINT
          value: {{ pluck .Values.global.env .Values.app.s3.host | first | default .Values.app.s3.host._default | quote  }}
        - name: S3_PORT
          value: {{ pluck .Values.global.env .Values.app.s3.port | first | default .Values.app.s3.port._default | quote }}
        - name: S3_ACCESS_KEY
          value: {{ pluck .Values.global.env .Values.app.s3.access_key | first | default .Values.app.s3.access_key._default | quote }}
        - name: S3_SECRET_KEY
          value: {{ pluck .Values.global.env .Values.app.s3.secret_key | first | default .Values.app.s3.secret_key._default | quote }}
        - name: S3_BUCKET
          value: {{ pluck .Values.global.env .Values.app.s3.bucket | first | default .Values.app.s3.bucket._default | quote }}
```
{% endraw %}
{% endsnippetcut %}


<div id="go-forth-button">
    <go-forth url="60_email.html" label="Работа с электронной почтой" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
