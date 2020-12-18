---
title: Доработка приложения
permalink: nodejs/100_basic/40_changing_app.html
---

TODO: follow и работа с гитерменизмом


## Быстрая сборка изменений

При локальной разработке приходится: а) сделать коммит в git; б) запускать `werf converge ...`; в) дожидаться пересборки и редеплоя. Для ускорения этого процесса команды werf (такие, как `converge`, `build`, `run`, `compose`) поддерживают атрибут `--follow`. Если вы, например, запустите

```bash
werf converge --repo registry.example.com/werf-guided-project --follow
```

то при коммите в git пересборка и деплой будет осуществляться автоматически.



<div id="go-forth-button">
    <go-forth url="40_optimize.html" label="Ускорение сборки" framework="{{ page.label_framework }}" ci="{{ page.label_ci }}" guide-code="{{ page.guide_code }}" base-url="{{ site.baseurl }}"></go-forth>
</div>
