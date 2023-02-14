## Buildah

Для сборки без Docker-сервера werf использует Buildah в rootless-режиме. Buildah встроен в бинарник werf.

Buildah включается установкой переменной окружения `WERF_BUILDAH_MODE` в один из вариантов: `auto`, `native-chroot`, `native-rootless`. Большинству пользователей для включения режима Buildah достаточно установить `WERF_BUILDAH_MODE=auto`.

### Уровень изоляции сборки

По умолчанию в режиме `auto` werf автоматически выбирает уровень изоляции в зависимости от платформы и окружения.

Поддерживается 2 варианта изоляции сборочных контейнеров:
1. Chroot — вариант без использования container runtime, включается переменной окружения `WERF_BUILDAH_MODE=native-chroot`.
2. Rootless — вариант с использованием container runtime (в системе должен быть установлен crun, или runc, или kata, или runsc), включается переменной окружения `WERF_BUILAH_MODE=native-rootless`.

### Драйвер хранилища

werf может использовать драйвер хранилища `overlay` или `vfs`:

* `overlay` позволяет использовать файловую систему OverlayFS. Можно использовать либо встроенную в ядро Linux поддержку OverlayFS (если она доступна), либо реализацию fuse-overlayfs. Это рекомендуемый выбор по умолчанию.
* `vfs` обеспечивает доступ к виртуальной файловой системе вместо OverlayFS. Эта файловая система уступает по производительности и требует привилегированного контейнера, поэтому ее не рекомендуется использовать. Однако в некоторых случаях она может пригодиться.

Как правило, достаточно использовать драйвер по умолчанию (`overlay`). Драйвер хранилища можно задать с помощью переменной окружения `WERF_BUILDAH_STORAGE_DRIVER`.

### Ulimits

По умолчанию Buildah-режим в werf наследует системные ulimits при запуске сборочных контейнеров. Пользователь может переопределить эти параметры с помощью переменной окружения `WERF_BUILDAH_ULIMIT`.

Формат `WERF_BUILDAH_ULIMIT=type=softlimit[:hardlimit][,type=softlimit[:hardlimit],...]` — конфигурация лимитов, указанные через запятую:
* "core": maximum core dump size (ulimit -c);
* "cpu": maximum CPU time (ulimit -t);
* "data": maximum size of a process's data segment (ulimit -d);
* "fsize": maximum size of new files (ulimit -f);
* "locks": maximum number of file locks (ulimit -x);
* "memlock": maximum amount of locked memory (ulimit -l);
* "msgqueue": maximum amount of data in message queues (ulimit -q);
* "nice": niceness adjustment (nice -n, ulimit -e);
* "nofile": maximum number of open files (ulimit -n);
* "nproc": maximum number of processes (ulimit -u);
* "rss": maximum size of a process's (ulimit -m);
* "rtprio": maximum real-time scheduling priority (ulimit -r);
* "rttime": maximum amount of real-time execution between blocking syscalls;
* "sigpending": maximum number of pending signals (ulimit -i);
* "stack": maximum stack size (ulimit -s).
