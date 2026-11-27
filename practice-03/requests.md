
# Практическое занятие №3 — тестирование API

Все запросы выполняются в Windows PowerShell.
Перед проверкой запустить сервер:

```powershell
go run ./cmd/server
```

## 1. Проверка состояния сервера

```powershell
curl.exe -i http://localhost:8080/health
```

Статус: 200 OK.

## 2. Создание задачи

```powershell
$body = @{title="Купить молоко"} | ConvertTo-Json -Compress
Invoke-WebRequest -Uri "http://localhost:8080/tasks" -Method POST -ContentType "application/json; charset=utf-8" -Body ([System.Text.Encoding]::UTF8.GetBytes($body))
```

Статус: 201 Created.

## 3. Получение списка задач

```powershell
curl.exe -i http://localhost:8080/tasks
```

Статус: 200 OK.

## 4. Фильтрация задач по названию

```powershell
curl.exe -i "http://localhost:8080/tasks?q=milk"
```

Статус: 200 OK. Возвращаются задачи с указанным текстом в названии.

## 5. Получение задачи по ID

```powershell
curl.exe -i http://localhost:8080/tasks/1
```

Статус: 200 OK, если задача существует.

## 6. Проверка некорректного ID

```powershell
curl.exe -i http://localhost:8080/tasks/abc
```

Статус: 400 Bad Request.

## 7. Проверка отсутствующей задачи

```powershell
curl.exe -i http://localhost:8080/tasks/9999
```

Статус: 404 Not Found.

## 8. Создание задачи без названия

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/tasks" -Method POST -ContentType "application/json" -Body "{}" -SkipHttpErrorCheck
```

Статус: 400 Bad Request.


---

## Дополнительные проверки

### 9. Проверка валидации названия

Название короче трёх символов:

```powershell
$body = '{"title":"Hi"}'
try {
    Invoke-RestMethod -Uri "http://localhost:8080/tasks" -Method POST -ContentType "application/json" -Body $body
} catch {
    $_.Exception.Response.StatusCode.value__
}
```

Ожидаемый результат: 422 Unprocessable Entity.

### 10. Изменение статуса задачи (PATCH)

Сначала необходимо создать задачу и узнать её ID.

```powershell
curl.exe -i -X PATCH http://localhost:8080/tasks/1
```

Ожидаемый результат: 200 OK, поле done = true.

### 11. Удаление задачи (DELETE)

```powershell
curl.exe -i -X DELETE http://localhost:8080/tasks/1
```

Ожидаемый результат: 204 No Content.

Повторный запрос GET /tasks/1 должен вернуть 404 Not Found.

### 12. Проверка CORS

```powershell
curl.exe -i -H "Origin: http://localhost:3000" http://localhost:8080/health
```

Ожидаемый заголовок: Access-Control-Allow-Origin: *

### 13. Проверка CORS preflight

```powershell
curl.exe -i -X OPTIONS http://localhost:8080/tasks -H "Origin: http://localhost:3000" -H "Access-Control-Request-Method: POST"
```

Ожидаемый результат: 204 No Content.

