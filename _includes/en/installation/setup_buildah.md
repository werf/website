{% offtopic title="Prepare your system for werf with Buildah backend" %}
1. Install a package providing `buildah` and it's dependencies.
2. Enable buildah backend by setting environment variable `WERF_BUILDAH_MODE=auto`:

```shell
export WERF_BUILDAH_MODE=auto
werf build
```
{% endofftopic %}
