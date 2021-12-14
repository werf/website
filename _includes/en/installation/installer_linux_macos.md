Make sure you have [Docker](https://docs.docker.com/get-docker), git 2.18.0+ and gpg installed.

Download werf installer:
```shell
curl -sSLO https://werf.io/install.sh && chmod +x install.sh
```

To use werf on a workstation install werf and configure its automatic activation (open a new shell-session afterwards):
```shell
./install.sh --version {{ include.version }} --channel {{ include.channel }}
```

To use werf in CI install werf and activate it manually:
```shell
./install.sh --ci
source "$(~/bin/trdl use werf {{ include.version }} {{ include.channel }})"
```

List of installation options:
```shell
./install.sh --help
```

After activation `werf` should be available in the shell-session from which it was activated:
```shell
werf version
```
