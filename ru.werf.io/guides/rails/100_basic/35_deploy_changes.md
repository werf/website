---
title: Внесение изменений
permalink: rails/100_basic/35_deploy_changes.html
---

 - Наше приложение ­— это crud для создания labels, демонстрируем это:
    - Получаем список labels из консоли:

        ```
        curl "http://URL:PORT/api/labels"
        ```

    - Создаём новые labels из консоли:

        ```
        curl -X POST "http://URL:PORT/api/labels/?label=red"
        curl -X POST "http://URL:PORT/api/labels/?label=hot"
        curl -X POST "http://URL:PORT/api/labels/?label=blue"
        curl -X POST "http://URL:PORT/api/labels/?label=cold"
        ```

    - Получаем обновлённый список labels:

        ```
        curl "http://URL:PORT/api/labels"
        ```

 - Давайте добавим новые timestamp поля для label: время его создания и обновления.
    - добавляем новые поля created_at и updated_at в таблицу labels:
        - новый файл миграций: examples/rails/016_deploy_app_changes/db/migrate/20210526202700_add_timestamps_to_labels.rb
        - обновлённый файл схемы бд: examples/rails/016_deploy_app_changes/db/schema.rb
    - правим examples/rails/016_deploy_app_changes/app/views/api/labels/_label.json.jbuilder, чтобы api выдавал не только id и имя label, но и поля created_at и updated_at
    - делаем git add . ; git commit
    - запускаем werf converge

    - Проверяем результат:
    
      ```
      curl -X POST "http://URL:PORT/api/labels/?label=black"
      curl -X POST "http://URL:PORT/api/labels/?label=white"
      curl "http://URL:PORT/api/labels"
      ```

      Видим новые поля — получилось!

 - Как поскейлить вручную
    - Как мы знаем веб сервер запущен в deployment basicapp.
    - Смотрим сколько реплик у нас запущено:

    ```
    kubectl -n werf-guided-rails get pod
    ```
     
    =>

    ```
    NAME                      READY   STATUS    RESTARTS   AGE
    basicapp-57789b68-kxcb9   1/1     Running   0          72m
    ```

    - Давайте вручную поменяем на 2 реплики
    
    ```
    kubectl -n werf-guided-rails edit deployment basicapp
    ```

    В открывшемся редакторе выставляем `spec.replicas=2`, закрываем редактор.

    - Снова смотрим сколько реплик у нас запущено:

    ```
    kubectl -n werf-guided-rails get pod
    ```
     
    =>

    ```
    NAME                      READY   STATUS    RESTARTS   AGE
    basicapp-57789b68-c2xlq   1/1     Running   0          6s
    basicapp-57789b68-kxcb9   1/1     Running   0          72m
    ```

 - Ок, мы поскейлили вручную, а что будет если снова запустить werf converge:
    - werf converge ...
    - Снова смотрим сколько реплик у нас запущено:

    ```
    kubectl -n werf-guided-rails get pod
    ```
     
    =>

    ```
    NAME                      READY   STATUS    RESTARTS   AGE
    basicapp-57789b68-kxcb9   1/1     Running   0          72m
    ```

    - Видим опять одну реплику, в чём дело? Дело в том, что werf привёл состояние кластера к состоянию описанному в текущем git коммите.
    - Как правильно?
        - Меняем тот же spec.replicas, но уже в git: examples/rails/016_deploy_app_changes/.helm/templates/deployment.yaml
        - Запускаем converge.
        
 - А как вообще посмотреть на состояние запущенного приложения в кубах и его логи?
    - Список запущенных деплойментов и подов

    ```
    kubectl -n werf-guided-rails get deployment
    kubectl -n werf-guided-rails get pod
    ```

    - Наш деплоймент `basicapp`, смотрим summary-информацию по нему:

    ```
    kubectl -n werf-guided-rails describe deployment basicapp
    ```

    - Деплоймент запускает поды, логи пишутся в подах, следующей командой можно запросить логи одного из запущенных подов:

    ```
    kubectl -n werf-guided-rails logs basicapp-57789b68-c2xlq
    ```

 - Ок, у нас rails приложение, хочу зайти в работающий контейнер запущенного приложения и запустить там rails console.
    - Вообще говоря заходить в поды не рекоммендуется, потому что ...
    - Но можно это сделать вот так:

    ```
    kubectl -n werf-guided-rails exec -ti basicapp-57789b68-c2xlq -- bash
    ```

    - Далее запускаем rails console:

    ```
    rails c
    ```

    - Далее можем получить список labels напрямую из базы данных:

    ```
    Label.all
    ```
