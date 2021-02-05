---
title: Working with e-mail
permalink: java_springboot/200_real_apps/60_email.html
layout: "development"
---

TODO: Глава про файлы


{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/secret-values.yaml
- consumer/consumer.js
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении работу с почтой.

Для того, чтобы использовать почту, мы предлагаем лишь один вариант - использовать внешний API. В нашем примере это [Mailgun](https://www.mailgun.com/), но есть и множество альтернатив, например: SendGrid, Amazon SES, Pepipost, Mailchimp, SparkPost.

Для того, чтобы Node.js-приложение могло работать с Mailgun, необходимо установить и сконфигурировать зависимость и начать её использовать. Установим через `npm` зависимость:

```shell
npm install mailgun-js
```

И сконфигурируем согласно [документации пакета](https://github.com/mailgun/mailgun-js#documentation):

{% snippetcut name="consumer/consumer.js" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/060-email/consumer/consumer.js" %}
{% raw %}
```js
const mailgun = require("mailgun-js");
<...>
const mg = mailgun({apiKey: process.env.MAILGUN_APIKEY, domain: process.env.MAILGUN_DOMAIN, host: "api.eu.mailgun.net"});
```
{% endraw %}
{% endsnippetcut %}

В коде приложения подключение к API и отправка сообщения может выглядеть так:

{% snippetcut name="consumer/consumer.js" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/060-email/consumer/consumer.js" %}
{% raw %}
```js
const sendEmail = async (email) => {
  try {
    const mg = mailgun({apiKey: process.env.MAILGUN_APIKEY, domain: process.env.MAILGUN_DOMAIN, host: "api.eu.mailgun.net"});
    if (email != null) {
      email.from = "Mailgun Sandbox <postmaster@" + process.env.MAILGUN_DOMAIN + ">"
      email.subject = "Welcome to Chat!"
      email.html = "<p>" + email.text + "</p>"
      console.log("Email: " + JSON.stringify(email));
    }
    const sent = await mg.messages().send(email);
    return sent;
  } catch (error) {
    console.error(error)
  }
}
```
{% endraw %}
{% endsnippetcut %}

Для работы с Mailgun необходимо пробросить в ключи доступа в приложение. Для этого стоит использовать [механизм секретных переменных]({{ site.docsurl }}/documentation/reference/deploy_process/working_with_secrets.html). *Вопрос работы с секретными переменными рассматривался подробнее, когда мы [делали базовое приложение](020_basic.html#secret-values-yaml).*

{% snippetcut name=".helm/secret-values.yaml (расшифрованный)" url="#" ignore-tests %}
{% raw %}
```yaml
app:
  mailgun_apikey:
    _default: 192edaae18f13aaf120a66a4fefd5c4d-7fsaaa4e-kk5d08a5
  mailgun_domain:
    _default: sandboxf1b90123966447a0514easd0ea421rba.mailgun.org
```
{% endraw %}
{% endsnippetcut %}

После того, как значения корректно прописаны и зашифрованы, можно пробросить соответствующие значения в Deployment:

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-nodejs/060-email/.helm/templates/consumer.yaml" %}
{% raw %}
```yaml
        - name: MAILGUN_APIKEY
          value: {{ pluck .Values.global.env .Values.app.mailgun_apikey | first | default .Values.app.mailgun_apikey._default }}
        - name: MAILGUN_DOMAIN
          value: {{ pluck .Values.global.env .Values.app.mailgun_domain | first | default .Values.app.mailgun_domain._default | quote }}
```
{% endraw %}
{% endsnippetcut %}


<div id="go-forth-button">
    <go-forth url="80_database.html" label="Подключаем Managed PostgreSQL" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
