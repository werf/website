---
title: Работа с электронной почтой
permalink: nodejs/200_real_apps/60_email.html
layout: development
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/secret-values.yaml
- .helm/values.yaml
- package.json
- app.js
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении работу с почтой.

Для того, чтобы использовать почту, мы предлагаем лишь один вариант - использовать внешний API. В нашем примере это [Mailgun](https://www.mailgun.com/), но есть и множество альтернатив, например: SendGrid, Amazon SES, Pepipost, Mailchimp, SparkPost.

Для того, чтобы Node.js-приложение могло работать с Mailgun, необходимо установить и сконфигурировать зависимость и начать её использовать. Установим через `npm` зависимость:

```shell
npm install mailgun-js
```

И сконфигурируем согласно [документации пакета](https://github.com/mailgun/mailgun-js#documentation). Подключение:

{% snippetcut name="app.js" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/260_email/app.js" %}
{% raw %}
```js
// Connection to Mailgun
let mg;
try {
  mg = mailgun({
    apiKey: process.env.MAILGUN_APIKEY,
    domain: process.env.MAILGUN_DOMAIN
  });
} catch (error) {
  console.error(error)
}
```
{% endraw %}
{% endsnippetcut %}

И использование:

{% snippetcut name="app.js" url="https://github.com/werf/werf-guides/blob/master/examples/nodejs/260_email/app.js" %}
{% raw %}
```js
//// Generate and send report to email ////
app.get('/api/send_report', function (req, res) {
  halt_if_errors(global_errors, res);

  try {
    email = {
      from: "Mailgun Sandbox <postmaster@" + process.env.MAILGUN_DOMAIN + ">",
      to: process.env.REPORT_RECIEVER,
      subject: "Report @ " + new Date().toISOString(),
      text: "Here is report: " + new Date().toISOString()
    }
    mg.messages().send(email, function(err, body){
      if (err) {
        console.error("error while sending e-mail")
        console.error(err)
        res.send(JSON.stringify({
          "result": "error",
          "comment": err.toString()
        }));
      } else {
        res.send(JSON.stringify({"result": true}));
      }
    });
  } catch (err) {
    res.send(JSON.stringify({
      "result": "error",
      "comment": err.toString()
    }));
  }
});
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
        command: ["node","/app/app.js"]
        image: {{ tuple "basicapp" . | werf_image}}
        workingDir: /app
        ports:
        - containerPort: 3000
          protocol: TCP
        env:
        - name: "SQLITE_FILE"
          value: "app.db"
        - name: "MAILGUN_APIKEY"
          value: "key-edec40bdee01d7c75cc25aeb1c09145e"
        - name: "MAILGUN_DOMAIN"
          value: "sandboxdad71610c342445aa4ab9bd5fe448eaf.mailgun.org"
        - name: "REPORT_RECIEVER"
          value: "example@gmail.com"
```
{% endraw %}
{% endsnippetcut %}

## Конфигурирование приложения

Оставлять переменные окружения и настройки mailgun прямо в helm-чарте — плохая практика. Чтобы иметь возможность конфигурировать разные значения для разных стендов (тестового, продакшн и т.п.) — воспользуйтесь подходами, описанными в главе "Конфигурирование инфраструктуры в виде кода": ключи от mailgun нужно унести в секретные значения (`secret-values.yaml`), а остальное — в настройки для чарта (`values.yaml`). 

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
        command: ["node","/app/app.js"]
        image: {{ tuple "basicapp" . | werf_image}}
        workingDir: /app
        ports:
        - containerPort: 3000
          protocol: TCP
        env:
        - name: "SQLITE_FILE"
          value: "app.db"
        - name: "MAILGUN_APIKEY"
          value: {{ pluck .Values.global.env .Values.app.mailgun.apikey | first | default .Values.app.mailgun.endpoint._default | quote }}
        - name: "MAILGUN_DOMAIN"
          value: {{ pluck .Values.global.env .Values.app.mailgun.domain | first | default .Values.app.mailgun.domain._default | quote }}
        - name: "REPORT_RECIEVER"
          value: {{ pluck .Values.global.env .Values.app.mailgun.report_reciever | first | default .Values.app.mailgun.report_reciever._default | quote }}
``` 
{% endraw %}
{% endsnippetcut %}

Несекретные значения — храним в `values.yaml`:

{% snippetcut name=".helm/values.yaml" url="#" %}
{% raw %}
```yaml
app:
  mailgun:
    domain:
      _default: "sandboxdad71610c342445aa4ab9bd5fe448eaf.mailgun.org"
      production: "sandboxdcd81847c3449873fghjk9bd5fe448eee.mailgun.org"
      staging: "sandboxdad71610c342445aa4ab9bd5fe448eaf.mailgun.org"
    report_reciever:
      _default: "kaylee@mycompany.io"
      production: "malcolm@mycompany.io"
      staging: "zoe@mycompany.io"
```
{% endraw %}
{% endsnippetcut %}

А секретные значения могут иметь, например, такой вид:

{% snippetcut name=".helm/secret-values.yaml (расшифрованный)" url="#" %}
{% raw %}
```yaml
app:
  mailgun:
    apikey:
      _default: "key-edec40bdee01d7c75cc25aeb1c09145e"
      production: "key-0742763e1c307236b0d05eb3e4f50d2c"
      staging: "key-f3a0df3239cb09688e99a6d217c8fa46"
```
{% endraw %}
{% endsnippetcut %}

{% endofftopic %}

## Деплой и тестирование

После того, как изменения в код внесены — нужно закоммитить изменения в git, задеплоить в кластер с помощью `werf converge` и протестировать метод `/api/send_report`.

Обратите внимание, что логика работы систем рассылки подразумевает, что вы можете отправить письма **только на подтверждённые** е-мэйл адреса. Таким образом, электронный адрес, указанный в переменной окружения `REPORT_RECIEVER` нужно первоначально подтвердить в Mailgun — иначе письма туда не будут уходить.

Также при тестировании обратите внимание, что пришедшее от сервиса письмо может оказаться в папке "Спам". Для отладки проблем используйте инструменты логирования отправки, предоставляемые почтовым сервисом.

<div id="go-forth-button">
    <go-forth url="80_database.html" label="Подключаем Managed PostgreSQL" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
