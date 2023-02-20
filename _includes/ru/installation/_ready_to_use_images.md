## Готовые образы

Рекомендуемым вариантом запуска werf в контейнерах является использование официальных образов werf. Для выбора предлагаются следующие образы, основанные на различных дистрибутивах:

- `registry.werf.io/werf/werf:{{ include.version }}-{{ include.channel }}-alpine` (или `registry.werf.io/werf/werf:{{ include.version }}-{{ include.channel }}`);
- `registry.werf.io/werf/werf:{{ include.version }}-{{ include.channel }}-ubuntu`;
- `registry.werf.io/werf/werf:{{ include.version }}-{{ include.channel }}-fedora`.