# test-service

HTTP-приложение на Go с 4 обработчиками: `/one`, `/two`, `/three`, `/four`.

Каждый обработчик принимает в теле запроса список downstream-целей и опциональный payload, вызывает все цели и возвращает результат.

## Конфигурация

| Переменная | Описание | По умолчанию |
|---|---|---|
| `PORT` | Порт сервера | `8080` |

## Формат запроса

```json
{
  "downstream": [
    {"url": "http://localhost:9001/endpoint", "method": "POST"},
    {"url": "http://localhost:9002/api", "method": "GET"}
  ],
  "payload": {
    "key": "value"
  }
}
```

- `downstream` — обязательный массив целей
- `payload` — опционально, пробрасывается телом в каждый downstream-вызов
- `method` — опционально, по умолчанию `POST`

## Пример

```bash
curl -X POST http://localhost:8080/one \
  -H 'Content-Type: application/json' \
  -d '{
    "downstream": [
      {"url": "http://localhost:9001/echo"},
      {"url": "http://localhost:9002/hello", "method": "GET"}
    ],
    "payload": {"hello": "world"}
  }'
```

Ответ:

```json
{
  "handler": "one",
  "results": [
    {"url": "http://localhost:9001/echo", "status": 200, "body": "..."},
    {"url": "http://localhost:9002/hello", "status": 200, "body": "..."}
  ]
}
```

## Запуск

```bash
go run .
```

## Структура проекта

```
.
├── main.go                 # точка входа
├── config/config.go        # чтение env
├── handlers/handlers.go    # HTTP-обработчики
├── downstream/downstream.go # HTTP-клиент для вызовов
└── README.md
```
