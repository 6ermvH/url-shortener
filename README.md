# URL Shortener

Тестовое задание в `OZON Банк` на позицию Стажёр GO
Сервис сокращения ссылок на Go с поддержкой in-memory и PostgreSQL хранилищ.

## Архитектура

Слои: `handler -> service -> repository`.
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

## Реляционная схема данных

```sql
CREATE TABLE IF NOT EXISTS links (
    short_url    TEXT PRIMARY KEY,
    original_url TEXT NOT NULL UNIQUE
);
```

## Генерация короткой ссылки

1. Берём SHA-256 хеш от исходного URL
2. Берём первые 8 байт хеша
3. Кодируем в base63 с длиной 10 символов
## Сборка и запуск

### Локально

```bash
go build -o url-shortener ./cmd

# in-memory
./url-shortener --addr=:8080 --storage=memory

# PostgreSQL
./url-shortener --addr=:8080 --storage=postgres \
	--dsn="postgres://user:pass@localhost:5432/db?sslmode=disable"

# Set migration layer
./url-shortener --addr=:8080 --storage=postgres \
	--dsn="postgres://user:pass@localhost:5432/db?sslmode=disable" \ 
	--migrate-version=1 
```

### Docker Compose

```bash
docker compose up --build
```

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

## Проблемы и решения

### Пакет base63

Стандартный `encoding/base64` использует алфавит из 64 символов. Символ `+` и `/` не URL-safe, а `=` нужен для паддинга. Написал собственный пакет `pkg/base63` по аналогии с `encoding/base64`: тип `Encoding` с методами `Encode`, `EncodeToString`, `EncodedLen`. Паддинг не нужен, так как длина выходной строки фиксирована.
### Миграции через golang-migrate

Миграции хранятся в `migrations/*.sql` и встраиваются в бинарник через `go:embed`. Применяются автоматически при старте сервиса с PostgreSQL. Повторное применение безопасно — `ErrNoChange` обрабатывается отдельно.

### Получение оригинальной ссылки из базы

Для получения оригинальной ссылки из БД используется простой `SELECT` запрос.
т.к на каждую оригинальную ссылку гарантируется только одна короткая ссылка, то оба поля помечены как `UNIQUE`, так-же добавлен хэш-индекс на `short_url`
## Допущения

- **Коллизии игнорируются.** Вероятность совпадения первых 8 байт SHA-256 для двух разных URL 1/2^64.
- **Оригинальный URL уникален.** Два одинаковых URL дадут одну короткую ссылку.
- **Сервис не возвращает 3XX.** Клиент получает JSON с оригинальным URL и сам выполняет переход. Это упрощает тестирование и позволяет использовать API программно.
- **Нет аутентификации.** Сервис полностью открыт, подходит для внутреннего использования.
- **Нет Очищения памяти.** Ссылки хранятся бессрочно.
- **Бизнес логика в сервисе.** Длина коротких ссылок, их формат определён в `service` слое, что даёт гибкость в изменении сервиса
- Решено было не использовать `cobra + viper` для реализации CLI, для такого маленького проекта это слишком
