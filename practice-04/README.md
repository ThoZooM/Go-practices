
# Практическое занятие №4

## Маршрутизация с chi. Создание CRUD-сервиса «Список задач»

## 1. Структура проекта

```text
practice-04/
├── internal/
│   └── task/
│       ├── model.go
│       ├── repo.go
│       ├── handler.go
│       └── handler_test.go
├── pkg/
│   └── middleware/
│       ├── logger.go
│       └── cors.go
├── main.go
├── go.mod
├── go.sum
└── README.md
```

## 2. Описание приложения

Разработан HTTP-сервис для управления списком задач (ToDo).

Каждая задача содержит следующие поля:

| Поле | Тип | Описание |
|---|---|---|
| id | int64 | Уникальный идентификатор |
| title | string | Название задачи |
| done | bool | Статус выполнения |
| created_at | time.Time | Время создания |
| updated_at | time.Time | Время изменения |

Данные сохраняются в JSON-файл и восстанавливаются при повторном запуске приложения.

## 3. Маршрутизация

Для обработки HTTP-запросов используется библиотека `chi`.

Все операции с задачами доступны через версионированный API `/api/v1`.

| HTTP-метод | Маршрут | Описание |
|---|---|---|
| GET | /health | Проверка состояния сервера |
| GET | /api/v1/tasks | Получение списка задач |
| GET | /api/v1/tasks/{id} | Получение задачи по ID |
| POST | /api/v1/tasks | Создание задачи |
| PUT | /api/v1/tasks/{id} | Обновление задачи |
| DELETE | /api/v1/tasks/{id} | Удаление задачи |

Маршруты регистрируются в `Handler.Routes()` и подключаются в `main.go` через `chi.Router`.

Идентификатор задачи извлекается из URL с помощью `chi.URLParam()`.

## 4. Запуск приложения

```powershell
cd practice-04
```

Установить зависимости:

```powershell
go mod download
```

Запустить приложение:

```powershell
go run .
```

Сервер доступен по адресу:

http://localhost:8080

Проверка работоспособности:

```powershell
curl.exe -i http://localhost:8080/health
```

Ответ:

```text
HTTP/1.1 200 OK

OK
```

## 5. Примеры HTTP-запросов

### 5.1. Создание задачи

```powershell
$body = '{"title":"Изучить библиотеку chi"}'

Invoke-RestMethod `
    -Uri "http://localhost:8080/api/v1/tasks" `
    -Method POST `
    -ContentType "application/json; charset=utf-8" `
    -Body ([System.Text.Encoding]::UTF8.GetBytes($body))
```

Ответ: `201 Created`.

Пример JSON:

```json
{
  "id": 1,
  "title": "Изучить библиотеку chi",
  "done": false,
  "created_at": "2026-11-27T16:00:00Z",
  "updated_at": "2026-11-27T16:00:00Z"
}
```

### 5.2. Получение списка задач

```powershell
curl.exe -i http://localhost:8080/api/v1/tasks
```

Ответ: `200 OK`.

### 5.3. Получение задачи по ID

```powershell
curl.exe -i http://localhost:8080/api/v1/tasks/1
```

Ответ: `200 OK`.

### 5.4. Обновление задачи

```powershell
$body = @{
    title = "Изучить chi глубже"
    done = $true
} | ConvertTo-Json -Compress

Invoke-RestMethod `
    -Uri "http://localhost:8080/api/v1/tasks/1" `
    -Method PUT `
    -ContentType "application/json; charset=utf-8" `
    -Body ([System.Text.Encoding]::UTF8.GetBytes($body))
```

Ответ: `200 OK`.

### 5.5. Удаление задачи

```powershell
curl.exe -i -X DELETE http://localhost:8080/api/v1/tasks/1
```

Ответ: `204 No Content`.

## 6. Дополнительные задания

### 6.1. Валидация title

Реализована проверка длины названия задачи.

Допустимая длина: от 3 до 100 символов.

Пробелы по краям удаляются.

При нарушении ограничения возвращается статус `422 Unprocessable Entity`.

### 6.2. Пагинация

Реализована возможность получать задачи постранично.

Пример:

```powershell
curl.exe "http://localhost:8080/api/v1/tasks?page=2&limit=5"
```

Параметры:

- `page` — номер страницы (по умолчанию 1);
- `limit` — количество задач (по умолчанию 10, максимум 100).

### 6.3. Фильтрация по done

Реализована фильтрация задач по статусу выполнения.

Выполненные задачи:

```powershell
curl.exe "http://localhost:8080/api/v1/tasks?done=true"
```

Невыполненные задачи:

```powershell
curl.exe "http://localhost:8080/api/v1/tasks?done=false"
```

Фильтрацию можно совместить с пагинацией:

```powershell
curl.exe "http://localhost:8080/api/v1/tasks?done=true&page=1&limit=5"
```

### 6.4. Сохранение на диск

Добавлено сохранение задач в `tasks.json`.

При запуске приложение загружает ранее сохранённые задачи.

При создании, обновлении и удалении задачи JSON-файл перезаписывается.

### 6.5. Версионирование API

Реализована версия API `v1`.

Все основные маршруты задач расположены по адресу:

`/api/v1/tasks`

## 7. Middleware

### Logger

Реализовано логирование:

- HTTP-метода;
- запрошенного URL;
- времени обработки запроса.

### CORS

Добавлены заголовки:

- Access-Control-Allow-Origin
- Access-Control-Allow-Methods
- Access-Control-Allow-Headers

OPTIONS-запросы обрабатываются с ответом `204 No Content`.

Также используются стандартные middleware библиотеки chi:

- RequestID;
- Recoverer.

## 8. Обработка ошибок

API использует JSON для передачи ошибок.

Пример:

```json
{
  "error": "task not found"
}
```

Основные HTTP-статусы:

| Статус | Значение |
|---|---|
| 200 OK | Успешное получение или изменение |
| 201 Created | Задача создана |
| 204 No Content | Задача удалена |
| 400 Bad Request | Некорректный запрос или ID |
| 404 Not Found | Задача не найдена |
| 422 Unprocessable Entity | Недопустимая длина title |

## 9. Тестирование

Для автоматического тестирования используется пакет `net/http/httptest`.

Разработаны следующие тесты:

| Тест | Проверяемая операция | Ожидаемый результат |
|---|---|---|
| TestListTasks | Получение пустого списка | 200 |
| TestCreateTask | Создание задачи | 201 |
| TestGetTask | Получение задачи по ID | 200 |
| TestUpdateTask | Обновление задачи | 200 |
| TestDeleteTask | Удаление задачи | 204, затем 404 |
| TestInvalidTitle | Проверка длины title | 422 |
| TestInvalidID | Проверка некорректного ID | 400 |

Команда запуска:

```powershell
go test ./... -v
```

Для проверки кода:

```powershell
go fmt ./...
go vet ./...
go build .
```
