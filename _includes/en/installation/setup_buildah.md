{% offtopic title="Prepare your system for werf with Buildah backend" %}
1. Install a package providing `buildah` and it's dependencies.
2. Enable buildah backend by setting environment variable `WERF_BUILDAH_MODE=auto`:

```shell
export WERF_BUILDAH_MODE=auto
werf build
```

### Troubleshooting

#### Overlayfs mode

If your Linux kernel version is 5.13+ (5.11+ for some distros) it is **recommended** to use native `overlay` kernel module:
- Make sure `overlay` kernel module is loaded with `lsmod | grep overlay`.
- Make sure `CONFIG_USER_NS=y` configuration flag enabled in your kernel with `grep CONFIG_USER_NS /boot/config-VERSION`.
- In debian based kernel command `sysctl kernel.unprivileged_userns_clone` should return `1`. Else execute:

    ```shell
    echo 'kernel.unprivileged_userns_clone = 1' | sudo tee -a /etc/sysctl.conf
    sudo sysctl -p
    ```

- Command `sysctl user.max_user_namespaces` should return at least `15000`. Else execute:

    ```shell
    echo 'user.max_user_namespaces = 15000' | sudo tee -a /etc/sysctl.conf
    sudo sysctl -p
    ```
If your kernel is older or if you can't activate `overlay` kernel module, then install `fuse-overlayfs`, which should be available in your distro package repos. As a last resort, `vfs` storage driver can be used by setting `WERF_BUILDAH_STORAGE_DRIVER=vfs`.

#### No buildah package in distro

If your distro does not have a package providing `buildah`, then following actions required:
- Install packages which provide `slirp4netns`, `newuidmap` and `newgidmap` (`uidmap`). Make sure that `newuidmap` and `newgidmap` has correct permissions:
   ```shell
   sudo setcap cap_setuid+ep /usr/bin/newuidmap
   sudo setcap cap_setgid+ep /usr/bin/newgidmap
   sudo chmod u-s,g-s /usr/bin/newuidmap /usr/bin/newgidmap
   ```
- If files `/etc/subuid` and `/etc/subgid` do not exist, then, in most distros, you should install a package that creates them. Current user should have at least `65536` subordinate UIDs/GUIDs reserved — this will look like a line `current_username:1000000:65536` in `/etc/subuid` and `/etc/subgid`. If there is no such a line you should add it yourself. After changing these files reboot might be necessary. More info: `man subuid`, `man subgid`.
- Path `~/.local/share/containers` should exist and the current user should have write/read permissions for it.

{% endofftopic %}
