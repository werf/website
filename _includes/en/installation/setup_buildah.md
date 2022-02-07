{% offtopic title="Activating experimental Buildah-backend in werf" %}
> For now Buildah mode can be activated only if `werf.yaml` contains only Dockerfile builds, but not Stapel builds.

If you want to run werf in containers/Kubernetes, then follow these [Docker](https://werf.io/documentation/v{{ include.version }}/advanced/ci_cd/run_in_container/use_docker_container.html) or [Kubernetes](https://werf.io/documentation/v{{ include.version }}/advanced/ci_cd/run_in_container/use_kubernetes.html) instructions. But if you want to run werf in Buildah mode outside of containers or want to build images bundled with werf in Buildah mode from scratch, then:
* If your Linux kernel version is 5.13+ (5.11+ for some distros), make sure `overlay` kernel module is loaded with `lsmod | grep overlay`. If your kernel is older or if you can't activate `overlay` kernel module, then install `fuse-overlayfs`, which should be available in your distro package repos.
* Command `sysctl kernel.unprivileged_userns_clone` should return `1`. Else execute:
  ```shell
  echo 'kernel.unprivileged_userns_clone = 1' | sudo tee -a /etc/sysctl.conf
  sudo sysctl -p
  ```
* Command `sysctl user.max_user_namespaces` should return at least `15000`. Else execute:
  ```shell
  echo 'user.max_user_namespaces = 15000' | sudo tee -a /etc/sysctl.conf
  sudo sysctl -p
  ```
* If files `/etc/subuid` and `/etc/subgid` do not exist, then, in most distros, you should install a package that creates them. Current user should have at least `65536` subordinate UIDs/GUIDs reserved — this will look like a line `current_username:1000000:65536` in `/etc/subuid` and `/etc/subgid`. If there is no such a line you should add it yourself. After changing these files reboot might be necessary. More info: `man subuid`, `man subgid`.
* Path `~/.local/share/containers` should exist and the current user should have write/read permissions for it.
* To use `werf` outside of containers, install `crun`, which is usually available in distro package repos.
* Install package which provides `newuidmap` and `newgidmap` binaries.

Now activate the Buildah backend and try to build your project:
```shell
export WERF_BUILDAH_MODE=auto
werf build
```

If there were errors running werf, try:
* Fix permissions for `newuidmap` and `newgidmap`:
  ```shell
  sudo setcap cap_setuid+ep /usr/bin/newuidmap
  sudo setcap cap_setgid+ep /usr/bin/newgidmap
  sudo chmod u-s,g-s /usr/bin/newuidmap /usr/bin/newgidmap
  ```
* If you experience problems with OverlayFS try VFS:
  ```shell
  export WERF_BUILDAH_STORAGE_DRIVER=vfs
  ```
* When werf is not containerized then consider switching from native rootless mode to less isolated chroot mode:
  ```shell
  export WERF_BUILDAH_MODE=native-chroot
  ```
{% endofftopic %}
