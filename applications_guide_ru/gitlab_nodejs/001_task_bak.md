

<!--

Мы против бездумного копирования файлов, но для того, чтобы вам было проще разобраться — в некоторых точках гайда есть **milestone-ы с полным кодом файлов на текущий момент**. Таким образом вы можете свериться с образцом:

<div class="tabs">
  <a href="javascript:void(0)" class="tabs__btn active" onclick="openTab(event, 'tabs__btn', 'tabs__content', 'example_1')">example_1.yaml</a>
  <a href="javascript:void(0)" class="tabs__btn" onclick="openTab(event, 'tabs__btn', 'tabs__content', 'example_2')">example_2.yaml</a>
</div>

<div id="example_1" class="tabs__content active">
{% snippetcut name="example_1.yaml" url="#" limited="true" %}
{% raw %}
```yaml
      containers:
      - name: example_1
        command:
         - java
         - -jar
         - /app/target/demo-1.0.jar $JAVA_OPT
      containers:
      - name: example_1
        command:
         - java
         - -jar
         - /app/target/demo-1.0.jar $JAVA_OPT
               containers:
      - name: example_1
        command:
         - java
         - -jar
         - /app/target/demo-1.0.jar $JAVA_OPT
      containers:
      - name: example_1
        command:
         - java
         - -jar
         - /app/target/demo-1.0.jar $JAVA_OPT
      containers:
      - name: example_1
        command:
         - java
         - -jar
         - /app/target/demo-1.0.jar $JAVA_OPT
               containers:
      - name: example_1
        command:
         - java
         - -jar
         - /app/target/demo-1.0.jar $JAVA_OPT
{{ tuple "basicapp" . | include "werf_container_image" | indent 8 }}
```
{% endraw %}
{% endsnippetcut %}
</div>
<div id="example_2" class="tabs__content">
{% snippetcut name="example_2.yaml" url="#" limited="true" %}
{% raw %}
```yaml
      containers:
      - name: example_2
        command:
         - java
         - -jar
         - /app/target/demo-1.0.jar $JAVA_OPT
      containers:
      - name: example_2
        command:
         - java
         - -jar
         - /app/target/demo-1.0.jar $JAVA_OPT
{{ tuple "basicapp" . | include "werf_container_image" | indent 8 }}
```
{% endraw %}
{% endsnippetcut %}
</div>

-->
