# Практическое занятие №3

## HTTP-сервер на Go. Обработка GET/POST-запросов

### Описание

Реализованы создание, получение, поиск, изменение статуса и удаление задач, обработка ошибок и логирование HTTP-запросов.

### Структура проекта

```text
practice-03/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers.go
│   │   ├── handlers_test.go
│   │   ├── middleware.go
│   │   ├── cors.go
│   │   └── responses.go
│   └── storage/
│       └── memory.go
├── go.mod
├── manage.ps1
├── requests.md
└── README.md
```

### HTTP-маршруты

| Метод | Маршрут | Описание | Статус |
|---|---|---|---|
| GET | `/health` | Проверка состояния сервера | 200 |
| GET | `/tasks` | Получение списка задач | 200 |
| GET | `/tasks?q=text` | Поиск задач по названию | 200 |
| POST | `/tasks` | Создание задачи | 201 |
| GET | `/tasks/{id}` | Получение задачи по ID | 200 |
| PATCH | `/tasks/{id}` | Отметка задачи как выполненной | 200 |
| DELETE | `/tasks/{id}` | Удаление задачи | 204 |

### Запуск

Для запуска из папки `practice-03`:

```powershell
go run ./cmd/server
```

По умолчанию сервер работает на порту `8080`.

### Изменение порта

```powershell
$env:PORT="9090"
go run ./cmd/server
```

При отсутствии переменной окружения `PORT` используется порт `8080`.

### Примеры запросов

Проверка сервера:

```powershell
curl.exe -i http://localhost:8080/health
```

Создание задачи:

```powershell
$body = '{"title":"Finish Go homework"}'
Invoke-RestMethod -Uri "http://localhost:8080/tasks" -Method POST -ContentType "application/json" -Body $body
```

Получение списка задач:

```powershell
curl.exe -i http://localhost:8080/tasks
```

Получение задачи по ID:

```powershell
curl.exe -i http://localhost:8080/tasks/1
```

Изменение статуса задачи:

```powershell
curl.exe -i -X PATCH http://localhost:8080/tasks/1
```

Удаление задачи:

```powershell
curl.exe -i -X DELETE http://localhost:8080/tasks/1
```

Полный набор тестовых запросов находится в файле `requests.md`.

### Дополнительные задания

Реализованы:

- Валидация длины названия задачи (3–140 символов).
- Изменение статуса задачи через PATCH.
- Удаление задач через DELETE.
- CORS-middleware.
- Логирование HTTP-запросов.
- Настройка порта через переменную окружения PORT.
- PowerShell-скрипт для запуска, сборки и тестирования.
- Автоматические тесты с использованием `httptest`.

### Сборка приложения

```powershell
go build -o ./bin/server.exe ./cmd/server
```

### Автоматическое тестирование

```powershell
go test -v ./...
```

Также предусмотрен скрипт `manage.ps1`:

```powershell
.\manage.ps1 run
.\manage.ps1 build
.\manage.ps1 test
```

### Особенности реализации

Данные сохраняются только в оперативной памяти и удаляются после перезапуска приложения.

Для обработки HTTP-запросов используется стандартная библиотека Go без сторонних веб-фреймворков.