# Работа с BigQuery в облаке

## Перед началом
Для части запросов к Google Cloud или их API понадобится доступ из страны, где они доступны (vpn).

## Начало работы
Сначала google cloud попросит включить этот сервис и ввести нужные ему данные:
https://cloud.google.com/bigquery/docs/sandbox?hl=en

## (Опционально) авторизация с gcloud CLI
1. Установка gcloud-cli: https://cloud.google.com/sdk/docs/install-sdk
2. Инициализация gcloud-cli на локальной машине `gcloud init`
3. (Опционально) Установка авторизации по умолчанию на локальной машине `gcloud auth application-default login` , в таком случае токен не понадобится

Gcloud-cli также поставит BigQuery cli: https://cloud.google.com/bigquery/docs/bq-command-line-tool , 
через который можно обращаться к данным в проектах (он также работает с [эмулятором](../emulator/README.md)).

## Проект
Проекты содержат в себе **датасеты** (в контексте базы данных - аналог **схемы**), которые в свою очередь содержат **таблицы**.
BigQuery может оперировать и другими сущностями, которые мы рассматривать не будем. 
Подробнее: https://cloud.google.com/bigquery?#common-uses

Создать проект [можно](https://cloud.google.com/resource-manager/docs/creating-managing-projects#gcloud)
- Через консоль: https://console.cloud.google.com/cloud-resource-manager
- С помощью CLI: `gcloud projects create <id>` , в качестве id - уникальная строка от 6 до 30 символов

## Авторизация
Для обращения к Google cloud API в коде понадобится сервисный аккаунт и credentials файл для авторизации в нём.

Создать сервисный аккаунт можно двумя способами:
1. В консоли Google cloud
2. С помощью CLI

### Создание сервисного аккаунта и получение ключа в консоли Google cloud
1. Перейти в раздел Service accounts https://console.cloud.google.com/iam-admin/serviceaccounts
2. Выбрать проект
3. Нажать "+ Create service account"
4. Выбрать название и дать роли для BigQuery. Например `BigQuery admin` - для всех возможных действий.
5. Выбрать в списке аккаунтов у созданного аккаунта опцию "Manage keys"
6. Create new key -> JSON
7. Скачается json файл с токеном

### Создание сервисного аккаунта и получение ключа в gcloud CLI
Создать сервисный аккаунт 
```bash
gcloud iam service-accounts create <unique-name>
```

Посмотреть email сервисного аккаунта можно в консоли или с помощью `gcloud iam service-accounts list`

Дать доступ к BigQuery
```bash
gcloud projects add-iam-policy-binding <project-id> \
    --member=serviceAccount:<service-account-email> \
    --role=roles/bigquery.admin
```

И скачать файл с токеном:
```bash
gcloud iam service-accounts keys create <output-file>.json --iam-account=<service-account-email>
```

## Обращение к данным

### SQL 
BigQuery поддерживает google-sql запросы через [cli](https://cloud.google.com/bigquery/docs/bq-command-line-tool) и библиотеки

Библиотеки: https://cloud.google.com/bigquery/docs/reference/libraries

Справка google-SQL: https://cloud.google.com/bigquery/docs/reference/standard-sql/query-syntax

Пример работы sql запросов находится в `sql.go` для запуска примера:
```bash
go run . --project=<project-id> [--kind=sql] [--credentialsFile=<credentials-file>] 
```

По умолчанию `kind=sql`, `credentialsFile=credentials.json`.  Credentials file - токен с предыдущего шага. 

### GRPC
BigQuery API поддерживает [GRPC сервис](https://cloud.google.com/bigquery/docs/reference/storage/rpc) 
для потокового чтения в форматах Apache Arrow и Avro

API: https://cloud.google.com/bigquery/docs/reference/storage/

Библиотеки: https://cloud.google.com/bigquery/docs/reference/storage/libraries

Пример работы sql запросов находится в `grpc.go` для запуска примера:
```bash
go run . --project=<project-id> --kind=grpc [--credentialsFile=<credentials-file>] 
```
