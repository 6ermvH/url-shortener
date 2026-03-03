# URL Shortener

Сервис сокращения ссылок на Go с поддержкой in-memory и PostgreSQL хранилищ.

## Архитектура

Слои: `handler → service → repository`. Каждый слой зависит только от следующего через интерфейс.

## API

**POST /**
```json
// request
{ "url": "https://example.com" }

// response 201
{ "shortUrl": "aBcD1234XY" }
```

**GET /{short}**
```json
// response 200
{ "originalUrl": "https://example.com" }

// response 404
{ "error": "url not found" }
```

## Генерация короткой ссылки

1. Берём SHA-256 хеш от исходного URL
2. Берём первые 8 байт хеша
3. Кодируем в base63 с длиной 10 символов

Одинаковый URL всегда даёт одинаковую короткую ссылку (идемпотентность).
## Конфигурация

Параметры задаются через флаги или переменные окружения (флаги имеют приоритет):

| Флаг       | Env       | По умолчанию | Описание                            |
| ---------- | --------- | ------------ | ----------------------------------- |
| `-addr`    | `ADDR`    | `:8080`      | Адрес сервера                       |
| `-storage` | `STORAGE` | `memory`     | Тип хранилища: `postgres or memory` |
| `-dsn`     | `DSN`     | `empty`      | PostgreSQL DSN (при `postgres`)     |

## Сборка и запуск

### Локально

```bash
go build -o url-shortener ./cmd

# in-memory
./url-shortener --storage=memory

./url-shortener --storage=postgres --dsn="postgres://user:pass@localhost:5432/db?sslmode=disable"
```

### Docker Compose

```bash
docker compose up --build
```

Сервис будет доступен на `http://localhost:8080`.

## Юнит-тесты

```bash
go install go.uber.org/mock/mockgen@latest
go generate ./...

go test ./...
```
## E2E-тесты

Тесты проверяют сервис через реальные HTTP-запросы. Используют отдельный docker-compose стек с PostgreSQL.

```bash
docker compose -f docker-compose.e2e.yaml up --build -d

go test -tags e2e ./tests/e2e/...
```

По умолчанию тесты обращаются на `http://localhost:8081`.
## Нагрузочное тестирование

Использована библиотека `locust` которая позволяет конфигурировать нагрузочное тестирование на сервис
```bash
pip install locust

locust -f tests/load/locustfile.py --host=http://localhost:8080
```

### Результаты

**1 пользователь, PostgreSQL**

![1 пользователь, PostgreSQL](metrics/1USER_1MIN_POSTGRES.png)

**500 пользователей, In-Memory**

![500 пользователей, In-Memory](metrics/500USER_1MIN_MEMORY.png)

**500 пользователей, PostgreSQL**

![500 пользователей, PostgreSQL](metrics/500USER_1MIN_POSTGRES.png)

## Линтинг

```bash
golangci-lint run
```

Конфигурация в `.golangci.yaml`: включены все линтеры кроме `depguard`, `wsl`, `nolintlint`. Активирован билд-тег `e2e`.

## Проблемы и решения

### Пакет base63

Стандартный `encoding/base64` использует алфавит из 64 символов. Символ `+` и `/` не URL-safe, а `=` нужен для паддинга. Написал собственный пакет `pkg/base63` по аналогии с `encoding/base64`: тип `Encoding` с методами `Encode`, `EncodeToString`, `EncodedLen`. Паддинг не нужен, так как длина выходной строки фиксирована.
### Миграции через golang-migrate

Миграции хранятся в `migrations/*.sql` и встраиваются в бинарник через `go:embed`. Применяются автоматически при старте сервиса с PostgreSQL. Повторное применение безопасно — `ErrNoChange` обрабатывается отдельно.

## Допущения

- **Коллизии игнорируются.** Вероятность совпадения первых 8 байт SHA-256 для двух разных URL 1/2^64.
- **Оригинальный URL уникален.** Два одинаковых URL дадут одну короткую ссылку.
- **Сервис не возвращает 301/302.** Клиент получает JSON с оригинальным URL и сам выполняет переход. Это упрощает тестирование и позволяет использовать API программно.
- **Нет аутентификации.** Сервис полностью открыт — подходит для внутреннего использования.
- **Нет TTL.** Ссылки хранятся бессрочно.
- **Бизнес логика в сервисе.** Длина коротких ссылок, их формат определён в `service` слое, что даёт гибкость в изменении сервиса
