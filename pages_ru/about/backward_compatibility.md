---
title: Обратная совместимость
permalink: about/backward_compatibility.html
layout: default
sidebar: none
description: Гарантии обратной совместимости
---
{%- asset backward-compatibility.css %}

<div class="page__container page_installation">
  <div class="installation-compatibility">
    <h2 class="installation-compatibility__title" id="backward-compatibility-promise">Гарантии обратной совместимости</h2>
    <div markdown="1" class="docs">
werf использует [семантическое версионирование](https://semver.org). Мажорные релизы могут ломать обратную совместимость (совместимость между 1.0 и 2.0, между 2.0 и 3.0 и т.д.). Для werf следующий мажорный релиз _может_ потребовать полный перезапуск приложений или других ручных действий, которые не могут быть выполнены автоматически.

Минорные релизы (1.1, 1.2, и т.д.) могут содержать новые глобальные возможности. Изменения, вносимые в минорных релизах не могут значительно менять обратную совместимость в рамках текущей мажорной ветки версий (например, 1.x). Для werf это означает, что обновление между минорными версиями происходит незаметно в большинстве случаев. Однако такое обновление может потребовать от пользователя запуска некоторых инструкций, которые должны быть предоставлены при выходе новой версии werf.

Патч релизы (1.1.0, 1.1.1, 1.1.2) могут содержать как исправления, так и новые возможности. Однако данные изменения обязаны быть полностью обратно совместимыми в рамках текущей минорной ветки версий (например, 1.1.х). Для werf это означает, что обновление на следующий патч релиз всегда происходит незаметно и может производится автоматически.

 - Мы **не гарантируем** обратную совместимость между:
    - `alpha` релизами;
    - `beta` релизами;
    - `ea` релизами.
 - Мы **гарантируем** обратную совместимость между:
    - `stable` релизами в рамках минорной ветки версий (например 1.1.x);
    - `rock-solid` релизами в рамках минорной ветки версий (например 1.1.x).
</div>
  </div>
</div>
