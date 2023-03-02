## Требования

- CI-система;

- Kubernetes для запуска CI-задач с установленным Kubernetes Runner вашей CI-системы.

## Настройка Runner'а

Настройте Kubernetes Runner вашей CI-системы так, чтобы создаваемые им Pod'ы имели следующую конфигурацию:

```yaml
apiVersion: v1
kind: Pod
metadata:
  namespace: ci
  annotations:
    "container.apparmor.security.beta.kubernetes.io/build": unconfined
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    runAsGroup: 1000
    fsGroup: 1000
```

Если нужно кеширование `.werf` (рекомендуется), то Pod'ы должны иметь дополнительную конфигурацию:

```yaml
spec:
  containers:
  - volumeMounts:
    - name: werf-cache
      mountPath: /home/build/.werf
  volumes:
  - name: werf-cache
    persistentVolumeClaim:
      claimName: ci-kubernetes-runner-werf-cache
```

... а если кеширование не нужно, то пусть имеют эту конфигурацию:

```yaml
spec:
  containers:
  - volumeMounts:
    - name: werf-cache
      mountPath: /home/build/.werf
  volumes:
  - name: werf-cache
    emptyDir: {}
```

Если будете развертывать приложения с werf в тот же кластер, в котором запускается Kubernetes Runner, то добавьте эту конфигурацию:

```yaml
spec:
  serviceAccountName: ci-kubernetes-runner
```

## Настройка Kubernetes

Если в конфигурации Kubernetes Runner вы включили кеширование `.werf`, то создайте в кластере соответствующий PersistentVolumeClaim, например:

```shell
$ kubectl create -f - <<EOF
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: ci-kubernetes-runner-werf-cache
  namespace: ci
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 10Gi
EOF
```

Если будете развертывать приложения в тот же кластер, в котором запускается Kubernetes Runner, то в кластере для запуска CI-задач настройте RBAC следующей командой:

```shell
$ kubectl create -f - <<EOF
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ci-kubernetes-runner
  namespace: ci
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ci-kubernetes-runner
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: ci-kubernetes-runner
    namespace: ci
EOF
```

> Для большей безопасности можете создать более ограниченную ClusterRole/Role и использовать её вместо кластерной роли `cluster-admin` выше.

Если Kubernetes-ноды для Kubernetes Runner имеют ядро Linux версии 5.12 или ниже, то разрешите в этом кластере использование FUSE для Kubernetes Runner следующей командой:

```shell
$ kubectl create -f - <<EOF
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fuse-device-plugin
  namespace: kube-system
spec:
  selector:
    matchLabels:
      name: fuse-device-plugin
  template:
    metadata:
      labels:
        name: fuse-device-plugin
    spec:
      hostNetwork: true
      containers:
      - image: soolaugust/fuse-device-plugin:v1.0
        name: fuse-device-plugin-ctr
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop: ["ALL"]
        volumeMounts:
          - name: device-plugin
            mountPath: /var/lib/kubelet/device-plugins
      volumes:
        - name: device-plugin
          hostPath:
            path: /var/lib/kubelet/device-plugins
---
apiVersion: v1
kind: LimitRange
metadata:
  name: enable-fuse
  namespace: ci
spec:
  limits:
  - type: "Container"
    default:
      github.com/fuse: 1
EOF
```

## Конфигурация CI/CD проекта

```yaml
# .pseudo-ci-config.yml:
image: "registry.werf.io/werf/werf:{{ include.version }}-{{ include.channel }}"
image_pull_policy: always

before_every_job:
- source "$(werf ci-env gitlab --as-file)"

jobs:
  prod:
    commands:
    - werf converge
    environment: prod
    on: master
    how: manually

  images:cleanup:
    commands:
    - werf cleanup
    on: master
    how: daily
```

TODO: больше листингов

TODO: настройка очистки (авторизация, ...)
