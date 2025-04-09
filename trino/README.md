# Trino

Основные настройки для Trino стандартные https://trino.io/docs/current/installation/deployment.html#configuring-trino

Для подключения источника BigQuery в /etc/catalog нужно поместить `<filename>.properties` файл следующего содержания:
```
connector.name=bigquery
bigquery.project-id=<project-id>
bigquery.credentials-file=<path-to-credentials.json>
```

Вместо json файла можно использовать base64 строку где этот же файл закодирован. 
Подробнее: https://trino.io/docs/current/connector/bigquery.html

В этой папке приведёт минимальный пример содержимого всех необходимых файлов для запуска Trino с BigQuery в качестве источника

<pre>
├── docker-compose.yml          
├── etc                             // volume для trino
    ├── catalog                     // каталог с источниками
    │   ├── bigquery.properties     // конфигурация источника bigquery (название файла может быть любым)
    │   └── credentials.json        // токен BigQuery (может быть в другом месте)
    ├── config.properties           // <a href="https://trino.io/docs/current/installation/deployment.html#config-properties">конфигурация сервера</a>
    ├── jvm.config                  // <a href="https://trino.io/docs/current/installation/deployment.html#jvm-config">конфигурация jvm для trino</a>
    └── node.properties             // <a href="https://trino.io/docs/current/installation/deployment.html#node-properties">конфигурация инстанса trino</a>
</pre>

Эмулятор BigQuery к сожалению [пока что](https://github.com/trinodb/trino/issues/25539) не поддерживается в trino.

## Запуск и использование

1. Запустить `docker compose up` или просто одиночный контейнер аналогичными аргументаи
2. Подключится к контейнеру: `docker exec -it trino trino` (первый trino у меня - название контейнера, второй - команда для вызова cli trino)
3. Для вывода доступных источников `show catalogs;`. Если мы назвали наш файл с источником `bigquery.properties` - то наш каталог будет называться также (`bigquery`).
4. Посмотреть датасеты: `show schemas from bigquery;`
5. Посмотреть таблицы в датасете: `show tables from bigquery.<dataset>;`
6. Посмотреть колонки в таблице: `show tables from bigquery.<dataset>.<table>;`
7. Остальные запросы аналогично. Например взять из таблицы: `select * from bigquery.<dataset>.<table>;`

Также в web-интерфейсе можно посмотреть историю запросов, их параметры, перфоманс, план и тд: http://localhost:8080/
