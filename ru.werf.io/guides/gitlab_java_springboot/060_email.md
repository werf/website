---
title: Работа с электронной почтой
permalink: gitlab_java_springboot/060_email.html
---

{% filesused title="Файлы, упомянутые в главе" %}
- .helm/templates/deployment.yaml
- .helm/secret-values.yaml
- .helm/values.yaml
- src/main/java/com/example/demo/SendGridClient.java
- pom.xml
{% endfilesused %}

В этой главе мы настроим в нашем базовом приложении работу с почтой.

Для того чтобы использовать почту мы предлагаем лишь один вариант - использовать внешнее API. В нашем примере это [sendgrid](https://sendgrid.com/).

Для того, чтобы Java приложение могло работать с sendgrid необходимо установить и сконфигурировать зависимость `sendgrid` и начать её использовать. Пропишем зависимости в `pom.xml`, чтобы они устаналивались:

{% snippetcut name="pom.xml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/060-email/pom.xml" %}
{% raw %}
```xml
		<dependency>
				<groupId>com.sendgrid</groupId>
				<artifactId>sendgrid-java</artifactId>
				<version>4.1.2</version>
		</dependency>
```
{% endraw %}
{% endsnippetcut %}

В коде приложения подключение к API и отправка сообщения может выглядеть так:

{% snippetcut name="src/main/java/com/example/demo/SendGridClient.java" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/060-email/src/main/java/com/example/demo/SendGridClient.java" %}
{% raw %}
```java
@Service
public class SendGridEmailService implements EmailService {
    private SendGrid sendGridClient;
    @Autowired
    public SendGridEmailService(SendGrid sendGridClient) {
        this.sendGridClient = sendGridClient;
    }
    @Override
    public void sendText(String from, String to, String subject, String body) {
        Response response = sendEmail(from, to, subject, new Content("text/plain", body));
        System.out.println("Status Code: " + response.getStatusCode() + ", Body: " + response.getBody() + ", Headers: "
                + response.getHeaders());
    }
    @Override
    public void sendHTML(String from, String to, String subject, String body) {
        Response response = sendEmail(from, to, subject, new Content("text/html", body));
        System.out.println("Status Code: " + response.getStatusCode() + ", Body: " + response.getBody() + ", Headers: "
                + response.getHeaders());
    }
    private Response sendEmail(String from, String to, String subject, Content content) {
        Mail mail = new Mail(new Email(from), subject, new Email(to), content);
        mail.setReplyTo(new Email("abc@gmail.com"));
        Request request = new Request();
        Response response = null;
        try {
            request.setMethod(Method.POST);
            request.setEndpoint("mail/send");
            request.setBody(mail.build());
            this.sendGridClient.api(request);
        } catch (IOException ex) {
            System.out.println(ex.getMessage());
        }
        return response;
    }
}
```
{% endraw %}
{% endsnippetcut %}

Для работы с `sendgrid` необходимо пробросить в ключи доступа в приложение. Для этого стоит использовать [механизм секретных переменных]({{ site.docsurl }}/v1.1-stable/documentation/reference/deploy_process/working_with_secrets.html). *Вопрос работы с секретными переменными рассматривался подробнее, [когда мы делали базовое приложение](020_basic/20_iac.html#secret-values-yaml)*

{% snippetcut name=".helm/secret-values.yaml (расшифрованный)" url="#" ignore-tests %}
{% raw %}
```yaml
app:
  sendgrid:
    apikey: 
      _default: sendgridapikey
    password:
      _default: sendgridpassword
```
{% endraw %}
{% endsnippetcut %}

А не секретные значения — храним в `values.yaml`

{% snippetcut name=".helm/values.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/060-email/.helm/values.yaml" %}
{% raw %}
```yaml
app:
<...>
  sendgrid:
    username:
      _default: sendgridusername
```
{% endraw %}
{% endsnippetcut %}

После того, как значения корректно прописаны и зашифрованы — мы можем пробросить соответствующие значения в Deployment.

{% snippetcut name=".helm/templates/deployment.yaml" url="https://github.com/werf/werf-guides/blob/master/examples/gitlab-java-springboot/060-email/.helm/templates/deployment.yaml" %}
{% raw %}
```yaml
        - name: SGAPIKEY
          value: {{ pluck .Values.global.env .Values.app.sendgrid.apikey | first | default .Values.app.sendgrid.apikey._default | quote }}
        - name: SGUSERNAME
          value: {{ pluck .Values.global.env .Values.app.sendgrid.username | first | default .Values.app.sendgrid.username._default | quote }}
        - name: SGPASSWORD
          value: {{ pluck .Values.global.env .Values.app.sendgrid.password | first | default .Values.app.sendgrid.password._default | quote }}
```
{% endraw %}
{% endsnippetcut %}

В интернете можно найти много разных примеров настройки почты через `sendgrid` используя spring. Имплементация может отличаться, но нам важно понять, что нужно параметризировать `application.properties`, чтобы java узнавала о значения из переменных окружения уже во время выполнения в кластере, а не на этапе сборки.

<div id="go-forth-button">
    <go-forth url="070_redis.html" label="Подключаем redis" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
