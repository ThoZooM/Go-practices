# Go-practices
1 semestr, Programming_Go
# Практическое занятие №1

## Установка и настройка окружения Go

### Описание
Сервер поддерживает следующие маршруты:

| Маршрут | Описание |
|---|---|
| `/hello` | Возвращает текст `Hello, world!` |
| `/user` | Возвращает JSON с уникальным UUID пользователя |
| `/health` | Возвращает статус сервера и текущее время в формате RFC3339 |

Для генерации UUID используется библиотека `github.com/google/uuid`.

### Структура проекта

```text
practice-01/
├── cmd/
│   └── server/
│       └── main.go
├── go.mod
├── go.sum
└── README.md
```

### Запуск

Из папки `practice-01`:

```powershell
go mod download
go run ./cmd/server
```

По умолчанию сервер запускается на порту `8080`.

### Настройка порта

Для изменения порта используется переменная окружения `APP_PORT`.

Пример для PowerShell:

```powershell
$env:APP_PORT="8081"
go run ./cmd/server
```

### Сборка приложения

```powershell
go build -o helloapi.exe ./cmd/server
.\helloapi.exe
```

### Проверка и форматирование кода

```powershell
go fmt ./...
go vet ./...
```

### Дополнительные задания

Реализованы:
- Настройка порта через переменную окружения `APP_PORT`.
- Маршрут `/health` с текущим временем в формате RFC3339.