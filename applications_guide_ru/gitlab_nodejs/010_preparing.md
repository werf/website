---
title: Подготовка к работе
sidebar: applications_guide
guide_code: gitlab_nodejs
permalink: gitlab_nodejs/010_preparing.html
toc: false
---
В этой главе мы рассмотрим, как подготовить приложение и инфраструктуру для того, чтобы можно было развернуть это приложение в Kubernetes.

Итак, последующее прохождение самоучителя подразумевает, что у вас:

- Установлен Kubernetes-кластер;
- На Kubernetes направлен трафик;
- Есть GitLab;
- Настроен GitLab Runner, где установлена утилита werf и с которого есть доступ к кластеру Kubernetes;
- Приложение соответствует требованиям для работы в Kubernetes.

Ниже мы максимально кратко объясним, что подразумевается под каждым из этих пунктов.

**Настройка окружения — это сложная задача**. Не рекомендуем закапываться в решение проблем своими руками: используйте существующие на рынке услуги, позовите на помощь коллег, которые имеют практику настройки инфраструктуры, или задавайте вопросы в [Telegram-сообществе werf](https://t.me/werf_ru).

## Требования к Kubernetes-кластеру

Никаких специальных требований к кластеру не предъявляется: вы можете развернуть его самостоятельно, воспользоваться предложением одного из облачных провайдеров или установить локальный Minikube.

{% offtopic title="Где взять Kubernetes, если нет опыта самостоятельной установки кластера?" %}

Если вы ещё ни разу не устанавливали Kubernetes самостоятельно и/или не обладаете опытом системного администрирования — не стоит пытаться освоить такую объемную тему в рамках этого самоучителя. Это может быть слишком сложно.

Существует огромное количество услуг, позиционирующихся как [Managed Kubernetes](https://www.google.com/search?q=managed+kubernetes). Часть из них включает серверные мощности, часть — нет.

Проще всего взять услугу в AWS (EKS) или Google Cloud (GKE). При первой регистрации они представляют бонус, которого должно хватить на неделю-другую работы с кластером. Однако для регистрации понадобится ввести данные банковской карты.

Альтернативой может стать использование Yandex.Cloud: здесь также предоставляется Managed Kubernetes и тоже требуется ввести данные карты. Но в качестве карты можно воспользоваться, например, виртуальной картой от Яндекс.Денег, позволяющей проводить «подписочные» платежи внутри России.

Также можно попробовать развернуть Kubernetes самостоятельно [в Hetzner на основании их статьи](https://community.hetzner.com/tutorials/install-kubernetes-cluster) — это один из самых дешёвых облачных провайдеров. Однако надо понимать, что вам придётся самостоятельно разбираться с большим пластом работ по администрированию платформы. Если решите переводить свой production — найдите надёжного провайдера или команду, которая будет администрировать платформу.
{% endofftopic %}

**После того, как развернёте кластер**, необходимо достать с него файл `.kube/config` с ключами доступа к кластеру. Он выглядит примерно так:

{% snippetcut name=".kube/config" url="#" %}
{% raw %}
```yaml
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ALotOfNumbersAndLettersAndSoOnAveryVERYveryLongStringInBase64=
    server: https://127.0.0.1:6445
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: ManyNumbersAndLettersAVERYveryVerylongString=
    client-key-data: ManyLettersAndNumbersAveryVeryVERYlongString=
```
{% endraw %}
{% endsnippetcut %}

Этот файл необходим, чтобы мы могли строить CI-процесс и деплоиться в созданный кластер.

Обратите внимание, что объект Ingress описывается в нашем самоучителе в расчёте на [nginx-ingress](https://kubernetes.github.io/ingress-nginx/). Если вы используете другой контроллер — вам предстоит адаптировать конфигурацию самостоятельно.

{% offtopic title="Что делать с Ingress, если у вас нет опыта администрирования?" %}
Вы можете задать вопрос в [Telegram-сообществе werf](https://t.me/werf_ru), указав, где и как вы развернули свой Kubernetes-кластер. Мы постараемся подсказать, как подойти к решению проблемы.
{% endofftopic %}

## Домен

Мы предполагаем, что к началу прохождения самоучителя вы привязали к своёму кластеру какой-то домен.

{% offtopic title="Что это значит?" %}
Ожидается, что у вас есть домен, к которому вы привязываете DNS-сервер, позволяющий настраивать [ресурсные записи DNS](https://ru.wikipedia.org/wiki/%D0%A2%D0%B8%D0%BF%D1%8B_%D1%80%D0%B5%D1%81%D1%83%D1%80%D1%81%D0%BD%D1%8B%D1%85_%D0%B7%D0%B0%D0%BF%D0%B8%D1%81%D0%B5%D0%B9_DNS). Он может быть предоставлен вашим регистратором доменов или облачной платформой, а может, вы самостоятельно его установили и настраиваете (см. также [как работает DNS](https://firstwiki.ru/index.php/%D0%9A%D0%B0%D0%BA_%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D0%B0%D0%B5%D1%82_DNS)).
Это также означает, что у вашего кластера есть IP-адреса, по которым он доступен.

Требуется прописать правильные A-записи в DNS-сервере, чтобы нужные домены были направлены на кластер.
{% endofftopic %}

В самоучителе мы будем использовать следующие домены:

- mydomain.io — направлен на кластер;
- *.mydomain.io — направлены на кластер;
- registry.mydomain.io — для GitLab Registry;
- gitlab.mydomain.io — для GitLab.

## Требования к GitLab

Предполагается, что в GitLab будет включён GitLab CI/CD и Registry.
Для прохождения самоучителя не обязательна собственная установка GitLab — достаточно ограничиться заведением учётной записи на [gitlab.com](https://about.gitlab.com/).

## Требования к GitLab Runner

{% offtopic title="Что такое Runner?" %}
Под термином «раннер» зачастую подразумевают три вещи:

- CLI-утилиту, которая устанавливается куда-то в кластере и осуществляет запуск сценария сборки и деплоя;
- Настройки, которые делаются в GitLab CI/CD в связи с установкой CLI-утилиты на какой-то из серверов;
- Выделенный сервер (если такой есть), на котором запускается CLI-утилита.
{% endofftopic %}

Самый простой вариант установки и настройки раннера — поднять отдельную виртуальную машину.

{% offtopic title="А другие варианты?" %}

На практике вариантов настройки раннеров используется гораздо больше. Можно использовать:

- несколько виртуальных машин;
- GitLab Shared Runners;
- динамически создаваемые узлы в Kubernetes.

Последние два варианта требуют работы с Docker внутри Docker с учётом монтирования томов. Реализация этих фич в werf [стоит](https://github.com/werf/werf/issues/1926) в ближайших планах.

_Примечание: werf уже умеет нативно работать с GitHub Runners._
{% endofftopic %}

Для настройки виртуальной машины необходимо:

- Установить на неё [CLI-утилиту gitlab-runner](https://docs.gitlab.com/runner/install/) (в качестве тега создаваемого раннера укажите `werf`).
- Установить на неё `docker`, `git` и [multiwerf]({{ site.docsurl }}/documentation/guides/installation.html#installing-multiwerf).
- Обеспечить связь с API Kubernetes.

{% offtopic title="Как обеспечить связь с API Kubernetes?" %}

Есть несколько способов реализовать это. Рассмотрим наиболее простой с точки зрения управления.

Сперва нужно убедиться, что в `.kube/config` указан корректный IP-адрес вашего кластера в поле `server`. Возможно, там указано некорректное значение:

{% snippetcut name=".kube/config" url="#" %}
{% raw %}
```yaml
- cluster:
<...>
    server: https://127.0.0.1:6445
```
{% endraw %}
{% endsnippetcut %}

Затем **зашифруйте `.kube/config` с помощью base64**. Если лень пользоваться консолью, можно воспользоваться каким-то из [веб-сервисов](https://www.base64encode.org/). Результатом станет что-то вроде:

{% snippetcut name=".kube/config (base64)" url="#" %}
{% raw %}
```yaml
YXBpVmVyc2lvbjogdjEKY2x1c3RlcnM6Ci0gY2x1c3RlcjoKICAgIGNlcnRpZmljYXRlLWF1dGhvcml0eS1kYXRhOiBBTG90T2ZOdW1iZXJzQW5kTGV0dGVyc0FuZFNvT25BdmVyeVZFUll2ZXJ5TG9uZ1N0cmluZ0luQmFzZTY0PQogICAgc2VydmVyOiBodHRwczovLzEyNy4wLjAuMTo2NDQ1CiAgbmFtZToga3ViZXJuZXRlcwpjb250ZXh0czoKLSBjb250ZXh0OgogICAgY2x1c3Rlcjoga3ViZXJuZXRlcwogICAgdXNlcjoga3ViZXJuZXRlcy1hZG1pbgogIG5hbWU6IGt1YmVybmV0ZXMtYWRtaW5Aa3ViZXJuZXRlcwpjdXJyZW50LWNvbnRleHQ6IGt1YmVybmV0ZXMtYWRtaW5Aa3ViZXJuZXRlcwpraW5kOiBDb25maWcKcHJlZmVyZW5jZXM6IHt9CnVzZXJzOgotIG5hbWU6IGt1YmVybmV0ZXMtYWRtaW4KICB1c2VyOgogICAgY2xpZW50LWNlcnRpZmljYXRlLWRhdGE6IE1hbnlOdW1iZXJzQW5kTGV0dGVyc0FWRVJZdmVyeVZlcnlsb25nU3RyaW5nPQogICAgY2xpZW50LWtleS1kYXRhOiBNYW55TGV0dGVyc0FuZE51bWJlcnNBdmVyeVZlcnlWRVJZbG9uZ1N0cmluZz0=
```
{% endraw %}
{% endsnippetcut %}

В каждом репозитории с кодом, с которым вы будете работать, вам надо будет **прописать в GitLab'е специальную переменную**. Зайдите в репозиторий и в левой панели нажмите `Settings` -> `CI/CD`. В этом разделе, в главе `Variables`, нужно прописать переменную окружения `KUBECONFIG_BASE64` с нашим `.kube/config`, зашифрованным base64.

Так доступ к кластеру попадет в werf, и с его помощью мы сможем строить CI-процесс.
{% endofftopic %}

## Требования к приложению

В самоучителе используется приложение, которое уже соответствует необходимым требованиям (см. ниже). Однако, когда вы будете реализовывать собственные приложения, возможно, потребуется немного их дорабатывать.

{% offtopic title="Что за требования?" %}

Наилучшим образом приложения будут работать в Kubernetes, если они соответствуют [методологии 12 факторов](https://12factor.net/ru/). В таком случае в Kubernetes работают stateless-приложения, которые не зависят от среды. Это важно, так как кластер может самостоятельно переносить приложения с одного узла на другой, заниматься масштабированием и т.п. — и мы не указываем, где конкретно запускать приложение, а лишь формируем правила, на основании которых кластер принимает свои решения.

Если приложение не соответствует этим требованиям полностью или частично, можно найти обходные пути, однако это непростая задача, рассмотрение которой уходит за рамки данного самоучителя.

{% endofftopic %}

<div>
    <a href="020_basic.html" class="nav-btn">Далее: Базовые настройки</a>
</div>
