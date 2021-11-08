---
title: Конфигурация и секреты
permalink: framework_id/200_real_apps/080_config.html
examples_initial: examples/framework_id/050_s3
examples: examples/framework_id/080_configuration
base_url: https://github.com/werf/werf-guides/blob/master/
description: |
  В этой главе мы покажем, как правильно использовать и хранить секретную и несекретную конфигурацию приложения. 

  В предыдущих главах конфигурация добавлялась прямо в контейнеры при сборке или использовалась как есть в переменных окружения контейнеров при выкате.
  
  Теперь для безопасности и гибкости конфигурация будет сохраняться в ConfigMap и Secret. А в дополнение к параметрам Helm-чарта (Values) и секретам werf будут продемонстрированы подходы параметризации и переиспользования конфигурации, а также хранения конфиденциальных данных вместе с кодом в Git-репозитории проекта.
---

{% include 200_real_apps/080_configuration.md.liquid %}