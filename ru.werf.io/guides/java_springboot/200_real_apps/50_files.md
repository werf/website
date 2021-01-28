---
title: Работа с файлами
permalink: java_springboot/200_real_apps/50_files.html
layout: "wip"
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/secret-values.yaml
- .helm/values.yaml
- build.gradle
- src/main/java/com/example/demo/mvc/controller/S3Controller.java
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении работу с пользовательскими файлами. Для этого потребуется персистентное (постоянное) хранилище.

В идеале — нужно добиться, чтобы приложение было stateless, а данные хранились в S3-совместимом хранилище — например, AWS S3, Selectel S3 или [MinIO](https://github.com/minio/minio). Это обеспечивает простое масштабирование, работу в HA-режиме и высокую доступность.

{% offtopic title="А есть какие-то способы кроме S3?" %}
Первый и более общий способ — это использовать как [volume](https://kubernetes.io/docs/concepts/storage/volumes/) хранилище [NFS](https://kubernetes.io/docs/concepts/storage/volumes/#nfs), [CephFS](https://kubernetes.io/docs/concepts/storage/volumes/#cephfs) или [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath).

Мы не рекомендуем этот способ, потому что при возникновении неполадок с такими типами volume’ов они влияют на работоспособность контейнера и всего демона Docker в целом. Тогда могут пострадать приложения, не имеющие никакого отношения к вашему.

Более надёжный путь — пользоваться S3. Так мы используем отдельный сервис, который имеет возможность масштабироваться, работать в HA-режиме и иметь высокую доступность. Можно воспользоваться облачным решением вроде AWS S3, Google Cloud Storage, Microsoft Blobs Storage и т.д.

Природа Kubernetes (и оркестровки контейнеров в целом) такова, что если мы будем сохранять файлы в какой-либо директории у приложения (запущенного в Kubernetes), то после перезапуска контейнера все изменения пропадут.
{% endofftopic %}

## Подключение S3

Данная настройка производится полностью в рамках приложения. Рассмотрим подключение к S3 на примере пакета `com.amazonaws:aws-java-sdk`.

{% snippetcut name="src/server/server.js" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/250_files/build.gradle" %}
{% raw %}
```gradle
plugins {
    id 'org.springframework.boot' version '2.4.0'
    id 'io.spring.dependency-management' version '1.0.10.RELEASE'
    id 'java'
}

group = 'com.example'
version = '0.0.1-SNAPSHOT'
sourceCompatibility = '8'

repositories {
    mavenCentral()
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-data-jpa'
    implementation 'org.springframework.boot:spring-boot-starter-web'
    compile group: 'org.xerial', name: 'sqlite-jdbc', version: '3.28.0'
    runtimeOnly 'com.h2database:h2'
    compile "com.amazonaws:aws-java-sdk:1.11.163"
    compile "javax.xml.bind:jaxb-api:2.2.4"
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
}

test {
    useJUnitPlatform()
}
```
{% endraw %}
{% endsnippetcut %}


… и настроим работу с S3 в приложении. Подключение:

{% snippetcut name="src/server/server.js" url="https://github.com/werf/werf-guides/blob/master/examples/springboot/250_files/- src/main/java/com/example/demo/mvc/controller/S3Controller.java" %}
{% raw %}
```java
package com.example.demo.mvc.controller;

import com.amazonaws.auth.AWSCredentials;
import com.amazonaws.auth.AWSStaticCredentialsProvider;
import com.amazonaws.auth.BasicAWSCredentials;
import com.amazonaws.client.builder.AwsClientBuilder;
import com.amazonaws.services.s3.AmazonS3;
import com.amazonaws.services.s3.AmazonS3ClientBuilder;
import com.amazonaws.services.s3.model.ObjectListing;
import com.amazonaws.services.s3.model.S3ObjectSummary;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(value = "/api")
public class S3Controller {
    AmazonS3 s3client;

    S3Controller() {
        AWSCredentials credentials = new BasicAWSCredentials(
                System.getenv("S3_ACCESS_KEY"),
                System.getenv("S3_SECRET_KEY")
        );
        this.s3client = AmazonS3ClientBuilder
                .standard()
                .withEndpointConfiguration(new AwsClientBuilder.EndpointConfiguration(System.getenv("S3_ENDPOINT"), System.getenv("S3_ZONE")))
                .withCredentials(new AWSStaticCredentialsProvider(credentials))
                .withPathStyleAccessEnabled(true)
                .build();
    }

    @GetMapping("/generate_report")
    public String report(){
        long date = System.currentTimeMillis();
        s3client.putObject(
                System.getenv("S3_BUCKET"),
                String.format("report_%o_.txt", date),
                String.valueOf(date));
        ObjectListing objectListing = s3client.listObjects(System.getenv("S3_BUCKET"));
        return "{\"result\": true}";
    }
}
```
{% endraw %}
{% endsnippetcut %}

Переменные окружения нужно будет прописать в объекте Deployment. В итоге файл `deployment.yaml` должен будет принять примерно следующий вид:

{% snippetcut name=".helm/templates/deployment.yaml" url="#" %}
{% raw %}
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: basicapp
spec:
  selector:
    matchLabels:
      app: basicapp
  revisionHistoryLimit: 3
  strategy:
    type: RollingUpdate
  replicas: 1
  template:
    metadata:
      labels:
        app: basicapp
    spec:
      imagePullSecrets:
      - name: "registrysecret"
      containers:
      - name: basicapp
        command: ["java","-jar","/app/demo.jar"]
        image: {{ .Values.werf.image.basicapp }}
        workingDir: /app
        ports:
        - containerPort: 8080
          protocol: TCP
        env:
        - name: "SQLITE_FILE"
          value: "app.db"
        - name: "S3_ENDPOINT"
          value: "s3.selcdn.ru"
        - name: "S3_PORT"
          value: "9000"
        - name: "S3_SSL"
          value: "true"
        - name: "S3_ZONE"
          value: "ru-1a"
        - name: "S3_BUCKET"
          value: "test-container"
        - name: "S3_ACCESS_KEY"
          value: "47951_Clarisse"
        - name: "S3_SECRET_KEY"
          value: "i{z^e9WX"
```
{% endraw %}
{% endsnippetcut %}

## Конфигурирование приложения

Оставлять переменные окружения и настройки S3 прямо в helm-чарте — плохая практика. Чтобы иметь возможность конфигурировать разные значения для разных стендов (тестового, продакшн и т.п.) — воспользуйтесь подходами, описанными в главе "Конфигурирование инфраструктуры в виде кода": ключи от S3 нужно унести в секретные значения (`secret-values.yaml`), а остальное — в настройки для чарта (`values.yaml`). 

{% offtopic title="А можно подробнее?" %}

Можно. У вас может получиться что-то вроде:

{% snippetcut name=".helm/templates/deployment.yaml" url="#" %}
{% raw %}
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: basicapp
spec:
  selector:
    matchLabels:
      app: basicapp
  revisionHistoryLimit: 3
  strategy:
    type: RollingUpdate
  replicas: 1
  template:
    metadata:
      labels:
        app: basicapp
    spec:
      imagePullSecrets:
      - name: "registrysecret"
      containers:
      - name: basicapp
        command: ["java","-jar","/app/demo.jar"]
        image: {{ .Values.werf.image.basicapp }}
        workingDir: /app
        ports:
        - containerPort: 8080
          protocol: TCP
        env:
        - name: "SQLITE_FILE"
          value: "app.db"
        - name: "S3_ENDPOINT"
          value: {{ pluck .Values.global.env .Values.app.s3.endpoint | first | default .Values.app.s3.endpoint._default | quote }}
        - name: "S3_PORT"
          value: {{ pluck .Values.global.env .Values.app.s3.port | first | default .Values.app.s3.port._default | quote }}
        - name: "S3_SSL"
          value: {{ pluck .Values.global.env .Values.app.s3.ssl | first | default .Values.app.s3.ssl._default | quote }}
        - name: "S3_ZONE"
          value: {{ pluck .Values.global.env .Values.app.s3.zone | first | default .Values.app.s3.zone._default | quote }}
        - name: "S3_BUCKET"
          value: {{ pluck .Values.global.env .Values.app.s3.bucket | first | default .Values.app.s3.bucket._default | quote }}
        - name: "S3_ACCESS_KEY"
          value: {{ pluck .Values.global.env .Values.app.s3.login | first | default .Values.app.s3.login._default | quote }}
        - name: "S3_SECRET_KEY"
          value: {{ pluck .Values.global.env .Values.app.s3.password | first | default .Values.app.s3.password._default | quote }}
``` 
{% endraw %}
{% endsnippetcut %}

Несекретные значения — храним в `values.yaml`:

{% snippetcut name=".helm/values.yaml" url="#" %}
{% raw %}
```yaml
app:
  s3:
    endpoint:
      _default: "s3.selcdn.ru"
    port:
      _default: "9000"
    ssl:
      _default: "true"
    zone:
      _default: "ru-1a"
    bucket:
      _default: "test-container"
      production: "production-myapp-reports"
      staging: "staging-myapp-reports"
    login:
      _default: "47951_Clarisse"
      production: "61235_Malcolm"
      staging: "54563_Kaylee"
```
{% endraw %}
{% endsnippetcut %}

А секретные значения могут иметь, например, такой вид:

{% snippetcut name=".helm/secret-values.yaml (расшифрованный)" url="#" %}
{% raw %}
```yaml
app:
  s3:
    password:
      _default: "i{z^e9WX"
      production: "&Brc_Rzn4/7f2E]u"
      staging: ".^At8wE,k<kMk+x""
```
{% endraw %}
{% endsnippetcut %}

{% endofftopic %}

<div id="go-forth-button">
    <go-forth url="60_email.html" label="Работа с электронной почтой" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
