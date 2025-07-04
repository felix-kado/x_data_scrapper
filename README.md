# x_data_scrapper — Twitter Data Aggregation API

REST API сервис на Go для агрегации данных пользователей Twitter через Twitter API v2. Использует chi router, чистую архитектуру и легко расширяется.

## Возможности
- Получение публичного профиля пользователя по username
- Получение твитов пользователя (с лимитом)
- Агрегация engagement-метрик
- Построение подграфа (seed users)
- Удобные JSON-ответы и обработка ошибок

## Быстрый старт

1. Получите Bearer Token для Twitter API v2
2. Создайте `.env` (или используйте `.env-example`):

```env
TWITTER_BEARER_TOKEN=ваш_токен_от_Twitter
PORT=3000
LOG_LEVEL=info
SHUTDOWN_TIMEOUT_SEC=15
```

3. Запустите сервер:
```sh
go run ./cmd/server
```

## Примеры запросов

**Профиль пользователя:**
```sh
curl -s http://localhost:3000/users/TwitterDev | jq
```

**Последние твиты:**
```sh
curl -s "http://localhost:3000/users/TwitterDev/tweets?limit=100" | jq
```

**Агрегированные метрики:**
```sh
curl -s http://localhost:3000/users/TwitterDev/metrics | jq
```

**Построение подграфа:**
```sh
curl -X POST http://localhost:3000/expand \
  -H "Content-Type: application/json" \
  -d '{"seed_ids": ["2244994945"], "depth": 1, "direction": 0, "collect_tweets": false, "collect_metrics": true}' | jq
```

**Проверка здоровья:**
```sh
curl -i http://localhost:3000/health
```

## Переменные окружения
- `TWITTER_BEARER_TOKEN` — Bearer-токен Twitter API (обязателен)
- `PORT` — порт HTTP-сервера (по умолчанию 3000)
- `LOG_LEVEL` — уровень логирования (`info`, `debug`, ...)
- `SHUTDOWN_TIMEOUT_SEC` — таймаут graceful shutdown (сек)

## Архитектура
- `internal/model` — структуры данных (User, Tweet, Edge, Graph, Metrics)
- `internal/twitter` — клиент Twitter API
- `internal/service` — бизнес-логика (UserService, ExpandService)
- `internal/handler` — HTTP-обработчики
- `internal/server` — роутинг
- `cmd/server` — точка входа, DI, запуск сервера
- `internal/util` — утилиты (JSON-ответы)
