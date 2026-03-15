# DelayedNotifier — отложенные уведомления через очередь

## Описание

Сервис принимает запросы на создание уведомлений, складывает их в очередь через RabbitMQ и отправляет в указанное время. При ошибке отправки повторяет попытку с экспоненциальной задержкой. Поддерживает отправку через Telegram и Log (для тестирования).

## Технологии

- **Go 1.25** — язык разработки
- **Gin** — HTTP-фреймворк
- **RabbitMQ** — очередь сообщений для фоновой обработки
- **Redis** — хранилище уведомлений и кэш статусов
- **Telegram Bot API** — канал доставки уведомлений
- **wbf** — внутренняя библиотека (rabbitmq, redis, logger, config, retry)
- **Docker Compose** — оркестрация инфраструктуры

## Архитектура
```
Клиент → HTTP API → Logic → Redis + RabbitMQ
                                      ↓
                                   Worker (Consumer из wbf)
                                      ↓
                                Ожидание scheduled_at
                                      ↓
                                Проверка статуса в Redis
                                      ↓
                                Sender (Telegram / Log)
                                      ↓
                                Обновление статуса в Redis
```

- **API** — принимает HTTP-запросы, возвращает JSON
- **Logic** — валидация, создание уведомлений, публикация в RabbitMQ
- **Storage** — CRUD в Redis (JSON-сериализация)
- **Worker** — handler для wbf Consumer, обрабатывает уведомления из очереди
- **Sender** — интерфейс + реализации (Telegram, Log)

## Структура проекта
```
├── cmd/main.go                     — точка входа
├── config/
│   ├── config.yaml                 — конфигурация (не в git)
│   └── config.yaml.example         — шаблон конфигурации
├── docker-compose.yml              — RabbitMQ + Redis
├── web/index.html                  — UI
├── internal/
│   ├── api/                        — HTTP-хендлеры и маршруты
│   ├── config/                     — чтение и валидация конфига
│   ├── interfaces/                 — интерфейсы storage
│   ├── logic/                      — бизнес-логика
│   ├── models/                     — структуры данных
│   ├── sender/                     — интерфейс и реализации каналов доставки
│   ├── storage/                    — работа с Redis
│   └── worker/                     — обработчик очереди RabbitMQ
```

## API

| Метод | Путь | Описание |
|-------|------|----------|
| POST | /notify | Создать уведомление |
| GET | /notify/:id | Получить статус уведомления |
| DELETE | /notify/:id | Отменить уведомление |

### POST /notify

Запрос:
```json
{
  "message": "Напоминание о встрече",
  "channel": "telegram",
  "recipient": "123456789",
  "delay_sec": 60
}
```
`channel` — `telegram` или `log`. `scheduled_at` (RFC3339) можно указать вместо `delay_sec`.

### Ретраи

При ошибке отправки воркер увеличивает `retry_count` и возвращает сообщение в очередь. wbf Consumer выполняет повторную обработку с экспоненциальной задержкой через `ConsumingStrat`. Максимум 5 попыток, после чего статус — `failed`.

## Запуск
```bash
# 1. Поднять RabbitMQ и Redis
make start

# 2. Скопировать и заполнить конфиг
cp config/config.yaml config/config.yaml
# Указать bot_token и default_chat_id

# 3. Запустить
make run
```

## UI

Открыть http://localhost:8080/ в браузере.

## Тесты
```bash
go test ./...
```
