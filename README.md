# Task Service

Сервис для управления задачами с HTTP API на Go.

## Требования

- Go `1.23+`
- Docker и Docker Compose

## Быстрый запуск через Docker Compose

```bash
docker compose up --build
```

После запуска сервис будет доступен по адресу `http://localhost:8080`.

Если `postgres` уже запускался ранее со старой схемой, пересоздай volume:

```bash
docker compose down -v
docker compose up --build
```

Причина в том, что SQL-файл из `migrations/0001_create_tasks.up.sql` монтируется в `docker-entrypoint-initdb.d` и применяется только при инициализации пустого data volume.

## Swagger

Swagger UI:

```text
http://localhost:8080/swagger/
```

OpenAPI JSON:

```text
http://localhost:8080/swagger/openapi.json
```

## API

Базовый префикс API:

```text
/api/v1
```

Основные маршруты:

- `POST /api/v1/tasks`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/{id}`
- `PUT /api/v1/tasks/{id}`
- `DELETE /api/v1/tasks/{id}`


## Test 
Все изменения, которые были внесены в проект:
Сервис теперь поддерживает рекуррентные задачи с новыми полями recurrence_type и recurrence_rule. recurrence_type может принимать значения "daily", "monthly", "specific", "even", "odd", а recurrence_rule задаёт правило повторения: например, "5,15,25" для дней месяца или "2026-04-06,2026-04-07" для конкретных дат.
В OpenAPI/Swagger добавлены новые схемы CreateTaskRequest и UpdateTaskRequest, включающие эти поля, а также исправлена схема ошибки Error (добавлена запятая после схемы Task). Все маршруты API (POST, GET, PUT, DELETE) обновлены для поддержки рекуррентных задач.
Сервис включает периодический воркер StartRecurrenceWorker, который проверяет задачи с установленной периодичностью и создаёт новые задачи согласно правилам. Валидация входных данных обновлена: проверяется корректность status и recurrence_type, а также обязательность поля title.
SQL-схема базы данных дополнена колонками recurrence_type и recurrence_rule и индексом по recurrence_type. Docker Compose монтирует SQL-файл миграции, который применяется при инициализации пустого volume.
Swagger UI доступен по адресу http://localhost:8080/swagger/, OpenAPI JSON — http://localhost:8080/swagger/openapi.json.
