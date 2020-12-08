---
title: Работа с электронной почтой
sidebar: applications_guide
guide_code: ____________
permalink: ____________/060_email.html
toc: false
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/secret-values.yaml
- ____________
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении работу с почтой.

Для того, чтобы использовать почту, мы предлагаем лишь один вариант - использовать внешний API. В нашем примере это [Mailgun](https://www.mailgun.com/), но есть и множество альтернатив, например: SendGrid, Amazon SES, Pepipost, Mailchimp, SparkPost.

Для того, чтобы ____________ приложение могло работать с Mailgun, необходимо установить и сконфигурировать зависимость и начать её использовать. Установим через `____________` зависимость:

```bash
____________
```

И сконфигурируем согласно [документации пакета](____________):

{% snippetcut name="____________" url="#" %}
{% raw %}
```____________
____________
____________
____________
```
{% endraw %}
{% endsnippetcut %}

В коде приложения подключение к API и отправка сообщения может выглядеть так:

{% snippetcut name="____________" url="#" %}
{% raw %}
```____________
____________
____________
____________
```
{% endraw %}
{% endsnippetcut %}

Для работы с Mailgun необходимо пробросить в ключи доступа в приложение. Для этого стоит использовать [механизм секретных переменных]({{ site.docsurl }}/documentation/reference/deploy_process/working_with_secrets.html). *Вопрос работы с секретными переменными рассматривался подробнее, когда мы [делали базовое приложение](020_basic.html#secret-values-yaml).*

{% snippetcut name="secret-values.yaml (расшифрованный)" url="#" ignore-tests %}
{% raw %}
```yaml
app:
  ____________
  ____________
```
{% endraw %}
{% endsnippetcut %}

А не секретные значения — храним в `values.yaml`

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
____________TODO: надо дать отсылку на какой-то гайд, где описано, как конкретно использовать ____________. Мало же просто его установить — надо ещё как-то юзать в коде.
____________



<div id="go-forth-button">
    <go-forth url="070_redis.html" label="Подключаем redis" base-url="applications_guide_ru"></go-forth>
</div>
