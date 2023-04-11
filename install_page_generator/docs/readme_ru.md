Генератор страницы установки
===

Гененирует конфигуратор для страницы установки.

## Конфигурационный файл

В качестве единого источника для всех создаваемых данных используется конфигурационный файл `config.yml`.

Делится на две глобальные секции:
* **Options**, из которых создаются блоки кнопок и дерево состояний,
* **Combinations**, из которых формируются паршелы для вкладок с информацией.

Пример конфигурационно файла:

```yaml
options:
- name: "Version"
  groupId: "version"
  help: "Some help information"
  values:
    - name: "1.2"
      tabName: "1.2"
    - name: "1.1"
      tabName: "1.1"
- name: "Channel"
  groupId: "channel"
  help: "Some help information"
  values:
    - name: "Stable"
      tabName: "stable"
    - name: "Early-Access"
      tabName: "ea"
    - name: "Beta"
      tabName: "beta"
    - name: "Alpha"
      tabName: "alpha"
- name: "Setup"
  groupId: "setup"
  help: "Some help information"
  values:
    - name: "Shell"
      tabName: "shell"
    - name: "Docker"
      tabName: "docker"
    - name: "Kubernetes"
      tabName: "kubernetes"
combinations:
- tabs:
  - name: "Setup 1"
    linkName: "setup_1"
    templateName: "default"
  - name: "View 1"
    linkName: "view_1"
    templateName: "default"
    params:
      - type: "type"
  options:
  - name: "Version"
    value: "1.1"
  - name: "Channel"
    value: "Stable"
  - name: "Setup"
    value: "Shell"
- tabs:
  - name: "Setup 2"
    linkName: "setup_2"
    templateName: "default"
  - name: "View 2"
    linkName: "view_2"
    templateName: "default"
    params:
      - type: "type"
  options:
  - name: "Version"
    value: "1.1"
  - name: "Channel"
    value: "Early-Access"
  - name: "Setup"
    value: "Shell"
```

### Набор Options

Пример одной секции:

```yaml
- name: "Version"
  groupId: "version"
  help: "Some help information"
  values:
    - name: "1.2"
      tabName: "1.2"
    - name: "1.1"
      tabName: "1.1"
```

В нем использованы следующе параметры:

* **name** — имя группы кнопок (опций).
* **groupId** — имя группы кнопок для JS.
* **help** – текст справки, показываемой по наведению курсора на знак вопроса возле имени группы кнопок.
  * Если оставлен пустым — при рендере не отображается.
* **values** – массив значений, которые присваиваются кнопкам.
  * **name** — текст кнопки.
  * **tabName** — имя кнопки для JS.

### Набор Combinations

Пример одной секции:

```yaml
- tabs:
  - name: "Setup 1"
    linkName: "setup_1"
    templateName: "default"
  - name: "View 1"
    linkName: "view_1"
    templateName: "default"
    params:
      - type: "type"
  options:
  - name: "Version"
    value: "1.1"
  - name: "Channel"
    value: "Stable"
  - name: "Setup"
    value: "Shell"
```

Состоит из двух частей:
* **tabs** – вкладки, соответствующие набору опций из options;
* **options** — набор опций, соответствующий этим вкладкам.
  * При рендере читаются подряд сверху внизу (аналогично дереву).

## Генерация дерева состояний

Дерево генерируется в формате JSON. Сгенерированное дерево лежит в файле `js_conf.json` в создаваемом каталоге `generated`.

При генерации дерева состояний Options читаются сверху внизу в прямом порядке.

Дерево генерируется в минимизированном виде.

Пример дерева, получающегося из приведенного конфигурационного файла:

```json
{"Option":"version","Values":[{"1.2":[{"Option":"channel","Values":[{"stable":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"ea":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"beta":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"alpha":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]}]},{"Option":"channel","Values":[{"stable":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"ea":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"beta":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"alpha":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]}]},{"Option":"channel","Values":[{"stable":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"ea":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"beta":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"alpha":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]}]}]},{"1.1":[{"Option":"channel","Values":[{"stable":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"ea":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"beta":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"alpha":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]}]},{"Option":"channel","Values":[{"stable":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"ea":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"beta":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"alpha":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]}]},{"Option":"channel","Values":[{"stable":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"ea":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"beta":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]},{"alpha":[{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null},{"Option":"setup","Values":null}]}]}]}]}
```

В развернутом виде оно соответствует графу зависимостей Options друг от друга:

```json
{
  "Option": "version",
  "Values": [
    {
      "1.2": [
        {
          "Option": "channel",
          "Values": [
            {
              "stable": [
                {
                  "Option": "setup",
                  "Values": null
                },
                {
                  "Option": "setup",
                  "Values": null
                },
                {
                  "Option": "setup",
                  "Values": null
                },
                {
                  "Option": "setup",
                  "Values": null
                }
              ]
            },
            {
              "ea": [
                {
                  "Option": "setup",
                  "Values": null
                },
                {
                  "Option": "setup",
                  "Values": null
                },
                # ...
```

## Генерация

Все шаблоны лежат в каталоге `templates`.

### Блок с кнопками

Для генерации блока с кнопками используется шаблон `configurator_buttons.html`:

```html
<div>
    <div class="page__container">
        <div class="configurator">
            <div class="configurator__buttons">
                {{ range .Groups }}
                <div class="button__wrap button--{{ .GroupName }}">
                    <div class="configurator__buttons-title">
                        {{ .Name }}
                        {{ if .Help }}
                        <span title="{{ .Help }} <a href='#'>Узнать подробнее</a>"></span>
                        {{ end }}
                    </div>
                    <div class="buttons-group">
                        {{ range .Buttons }}
                            <a href="#" class="btn btn_o" data-key="{{ .TabName }}" data-value="{{ .Name }}">{{ .Name }}</a>
                        {{ end }}
                    </div>
                </div>
                {{ end }}
            </div>
        </div>
    </div>
</div>
```

Он генерирует поочередно все блоки с кнопками, подставляя названия и служебные элементы их конфигурационного файла.

Сгенерированный результат называется `configurator.html`.

### Паршелы табов

Все табы описаны в разделе `configurations` конфигурационного файла. Для каждого таба можно задать список параметров, которые будут переданы в шаблон.

Пример описания одного таба:

```yaml
- tabs:
  - name: "Setup 1"
    linkName: "setup_1"
    templateName: "default"
  - name: "View 1"
    linkName: "view_1"
    templateName: "default"
    params:
      - type: "type"
  options:
  - name: "Version"
    value: "1.1"
  - name: "Channel"
    value: "Stable"
  - name: "Setup"
    value: "Shell"
```

Здесь описаны две вкладки, соответствующие трем выбранным кнопкам конфигуратора.

* **name** – название вкладки;
* **linkName** – путь в пермалинке;
* **templateName** – имя используемого шаблона;
* **params** – массив передаваемых в шаблон параметров.

Шаблоны хранятся в каталоге `templates/tabs` и содержат всю необходимое для генерации контента.

Пример шаблона `default`:

```html
---
title: {{ .Name }}
permalink: {{ .Permalink }}
---

DEFAULT
```

Permalink формируется по принципу «все options в порядке появления + имя из конфига». Например, для писанного примера permalink будет иметь следующий вид: `/1.1_stable_shell_setup_1.html`.

Массив данных передается под названием `.data` и содержит все передаваемые параметрв в формате `name - data`.

## Пример генерации

После генерации приведенного примера конфигурационного файла получится следующее содержимое:

```
generated
├── 1.1_ea_shell_setup_2.html
├── 1.1_ea_shell_view_2.html
├── 1.1_stable_shell_setup_1.html
├── 1.1_stable_shell_view_1.html
├── configurator.html
└── js_conf.json
```
