---
title: Работа с файлами
permalink: gitlab_java_springboot/050_files.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/secret-values.yaml
- .helm/values.yaml
- src/main/resources/application.properties
- pom.xml
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении работу с пользовательскими файлами. Для этого потребуется персистентное (постоянное) хранилище.

В идеале — нужно добиться, чтобы приложение было stateless, а данные хранились в S3-совместимом хранилище — например, [MinIO](https://github.com/minio/minio) или AWS S3. Это обеспечивает простое масштабирование, работу в HA-режиме и высокую доступность.

{% offtopic title="А есть какие-то способы кроме S3?" %}
Первый и более общий способ — это использовать как [volume](https://kubernetes.io/docs/concepts/storage/volumes/) хранилище [NFS](https://kubernetes.io/docs/concepts/storage/volumes/#nfs), [CephFS](https://kubernetes.io/docs/concepts/storage/volumes/#cephfs) или [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath).

Мы не рекомендуем этот способ, потому что при возникновении неполадок с такими типами volume’ов они влияют на работоспособность контейнера и всего демона Docker в целом. Тогда могут пострадать приложения, не имеющие никакого отношения к вашему.

Более надёжный путь — пользоваться S3. Так мы используем отдельный сервис, который имеет возможность масштабироваться, работать в HA-режиме и иметь высокую доступность. Можно воспользоваться облачным решением вроде AWS S3, Google Cloud Storage, Microsoft Blobs Storage и т.д.

Природа Kubernetes (и оркестровки контейнеров в целом) такова, что если мы будем сохранять файлы в какой-либо директории у приложения (запущенного в Kubernetes), то после перезапуска контейнера все изменения пропадут.
{% endofftopic %}

Данная настройка производится полностью в рамках приложения. Нужно подключить `aws-java-sdk`, сконфигурировать его и начать использовать.

Пропишем использование `aws-java-sdk` как зависимость:

{% snippetcut name="pom.xml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/050-files/pom.xml" %}
{% raw %}
```xml
		<dependency>
				<groupId>com.amazonaws</groupId>
				<artifactId>aws-java-sdk</artifactId>
				<version>1.11.163</version>
		</dependency>
```
{% endraw %}
{% endsnippetcut %}

И сконфигурируем:

{% snippetcut name="src/main/resources/application.properties" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/050-files/src/main/resources/application.properties" %}
{% raw %}
```yaml
amazonProperties.endpointUrl=${S3ENDPOINT}
amazonProperties.login=${S3LOGIN}
amazonProperties.password=${S3PASSWORD}
amazonProperties.bucketName=${S3BUCKET}
amazonProperties.zone=${S3ZONE}
```
{% endraw %}
{% endsnippetcut %}

Для работы с S3 необходимо пробросить ключи доступа в приложение. Для этого стоит использовать [механизм секретных переменных]({{ site.docsurl }}/documentation/reference/deploy_process/working_with_secrets.html). *Вопрос работы с секретными переменными рассматривался подробнее, когда мы [делали базовое приложение](020_basic.html#secret-values-yaml).*

{% snippetcut name=".helm/secret-values.yaml (расшифрованный)" url="#" ignore-tests %}
{% raw %}
```yaml
app:
  s3:
    login:
      _default: mys3keyidstage
    password:
      _default: mys3keysecretstage
```
{% endraw %}
{% endsnippetcut %}

Несекретные значения — храним в `values.yaml`:

{% snippetcut name=".helm/values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/050-files/.helm/values.yaml" %}
{% raw %}
```yaml
app:
<...>
  s3:
    endpoint:
      _default: "https://s3.selcdn.ru"
    bucket:
      _default: "test-container"
    zone:
      _default: "ru-1a"
```
{% endraw %}
{% endsnippetcut %}

После того, как значения корректно прописаны и зашифрованы, можно пробросить соответствующие значения в Deployment:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/050-files/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        env:
<...>
        - name: S3ENDPOINT
          value: {{ pluck .Values.global.env .Values.app.s3.endpoint | first | default .Values.app.s3.endpoint._default | quote }}
        - name: S3LOGIN
          value: {{ pluck .Values.global.env .Values.app.s3.login | first | default .Values.app.s3.login._default | quote }}
        - name: S3PASSWORD
          value: {{ pluck .Values.global.env .Values.app.s3.password | first | default .Values.app.s3.password._default | quote }}
        - name: S3BUCKET
          value: {{ pluck .Values.global.env .Values.app.s3.bucket | first | default .Values.app.s3.bucket._default | quote }}
        - name: S3ZONE
          value: {{ pluck .Values.global.env .Values.app.s3.zone | first | default .Values.app.s3.zone._default | quote }}
{{ tuple "basicapp" . | include "werf_container_env" | indent 8 }}
```
{% endraw %}
{% endsnippetcut %}

Об особенностях использования aws-java-sdk можно подробно почитать в [документации](https://cloud.spring.io/spring-cloud-aws/spring-cloud-aws.html)

<div id="go-forth-button">
    <go-forth url="060_email.html" label="Работа с электронной почтой" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
