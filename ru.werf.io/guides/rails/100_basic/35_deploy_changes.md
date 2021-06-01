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
    kubectl get pod
    ```
     
    =>

    ```
    NAME                      READY   STATUS    RESTARTS   AGE
    basicapp-57789b68-kxcb9   1/1     Running   0          72m
    ```

    - Давайте вручную поменяем на 2 реплики
    
    ```
    kubectl edit deployment basicapp
    ```

    В открывшемся редакторе выставляем `spec.replicas=2`, закрываем редактор.

    - Снова смотрим сколько реплик у нас запущено:

    ```
    kubectl get pod
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
    kubectl get pod
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
    kubectl get deployment
    kubectl get pod
    ```

    - Наш деплоймент `basicapp`, смотрим summary-информацию по нему:

    ```
    kubectl describe deployment basicapp
    ```

    - Деплоймент запускает поды, логи пишутся в подах, следующей командой можно запросить логи одного из запущенных подов:

    ```
    kubectl logs basicapp-57789b68-c2xlq
    ```

 - Проблема: наше приложение использует БД sqlite в локальном файле, а мы добавили несколько реплик, каждая реплика использует свою базу данных, а не общую. Это приведёт к тому, что запросы будут равномерно распределяться по репликам при создании и получении labels, результаты запросов могут отличаться при перезапуске. Давайте переедем на общий mysql, чтобы решить эту проблему.

    - создаём новый шаблон для запуска mysql в кластере: examples/rails/016_deploy_app_changes/.helm/templates/database.yml
        - вообще говоря, в реальных приложениях, mysql запустить несколько сложнее — нужен persistent volume, но в нашем случае для development-окружения бд будет терять все свои данные в случае перезапуска

    - адаптируем deployment приложения так, чтобы миграции запускались в init container и убираем настройки sqlite:
        - examples/rails/016_deploy_app_changes/.helm/templates/deployment.yaml

    - настраиваем приложение на работу с mysql:
        - examples/rails/016_deploy_app_changes/config/database.yml
        - examples/rails/016_deploy_app_changes/Gemfile
        - examples/rails/016_deploy_app_changes/Gemfile.lock
    
    - выкатываем через converge, дожидаемся выполнения команды, в процессе работы в логах могут быть ошибки подключения к бд, которая ещё не успела стартануть — это нормально, спокойно дожидаемся полного старта приложения

    - дальше проверяем создание и получение labels командами `curl -X POST "http://URL:PORT/api/labels/?label=name"` и `curl "http://URL:PORT/api/labels/"` и видим консистивное поведение

 - Ок, у нас rails приложение, хочу зайти в работающий контейнер запущенного приложения и запустить там rails console.
    - Вообще говоря заходить в поды не рекоммендуется, потому что ...
    - Но можно это сделать вот так:

    ```
    kubectl exec -ti basicapp-57789b68-c2xlq -- bash
    ```

    - Далее запускаем rails console:

    ```
    rails c
    ```

    - Далее можем получить список labels напрямую из базы данных:

    ```
    Label.all
    ```
