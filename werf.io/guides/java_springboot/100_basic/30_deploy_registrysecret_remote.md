In most cases, the remote registry is protected with the username and password. We have already discussed it in the chapter "Preparing the environment".

We pass this username/password pair to the cluster using an object of the Secret type called `registrysecret`. Let us suppose that your username is `admin`, your password is also `admin`, and the repository is located at `registry.example.com`. Let's create the `registry-secret.yaml` file:

{% snippetcut name=".helm/templates/registry-secret.yaml" url="#" %}
{% raw %}
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: registrysecret
  annotations:
    "helm.sh/hook": pre-install
type: kubernetes.io/dockerconfigjson
data:
  {{- $registry := "registry.example.com" }}
  {{- $login := "admin" }}
  {{- $password := "admin" }}
  .dockerconfigjson: {{ printf "{\"auths\": {\"%s\": {\"auth\": \"%s\"}}}" $registry (printf "%s:%s" $login $password | b64enc) | b64enc }}
```
{% endraw %}
{% endsnippetcut %}

_Note: In our example, access keys are stored unencrypted. This is not safe. The nuances of the secure key storage will be discussed later in the chapter "Organizing the non-local development"._

Now you need to specify this Secret in each Deployment object in the `imagePullSecrets` attribute so that the cluster can pull the image:

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
      imagePullSecrets:                           # Added line
      - name: "registrysecret"                    # Added line
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
```
{% endraw %}
{% endsnippetcut %}
