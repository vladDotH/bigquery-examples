# Spark

Apache Spark позволяет читать и записывать датафреймы из/в BigQuery.
Обратите внимание что работа (настройки и опции) в кластере Spark / Google cloud 
может отличаться от локальной работы. 
Подробнее: [github.com/GoogleCloudDataproc/spark-bigquery-connector](https://github.com/GoogleCloudDataproc/spark-bigquery-connector)

Здесь приведён пример подключения, чтения и записи датафреймов с помощью PySpark.

Для инициализации:
```
python3 -m venv env
source env/bin/activate
python3 -m pip install -r requirements.txt
```

Для авторизации в этой же папке должен быть json токен облака.

Для запуска записи в таблицу:
```
python3 write.py
```

Для запуска чтения из таблицу:
```
python3 read.py
```

Spark может долго прогреваться.