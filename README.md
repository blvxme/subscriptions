Простой сервис подписок, работающий по gRPC.

# Использование

1. `git clone git@github.com:blvxme/subscriptions.git`, `cd subscriptions`
2. Сборка: `go build -o subscriptions cmd/subscriptions/main.go`
3. Установка переменных окружения:
    - Порт, на котором будет работать сервер: `export SUBSCRIPTIONS_PORT=8080` (если не устанавливать, по умолчанию
      будет использоваться порт 8080)
4. Запуск: `./subscriptions`. После этого сервер начнет принимать соединения на установленном порту.

### Работа с сервисом

- Конфигурация и запуск сервера реализованы в `./internal/server/server.go`
- Обработка клиентов реализована в `./internal/handler/handler.go`
- Protobuf-схема, а также сгенерированный код на Go находится в `./internal/proto/`

Для проверки работы можно использовать [`grpcurl`](https://github.com/fullstorydev/grpcurl).

- Для подписки на события можно воспользоваться командой:

```shell
grpcurl                                                    \
    -import-path /path/to/cloned/repository/internal/proto \
    -proto pubsub.proto                                    \
    -d '{"key": "subject_name"}'                           \
    -plaintext                                             \
    localhost:<PORT NUMBER> PubSub.Subscribe
```

- Для публикации событий можно воспользоваться командой:

```shell
grpcurl                                                    \
    -import-path /path/to/cloned/repository/internal/proto \
    -proto pubsub.proto                                    \
    -d '{"key": "subject_name", "data": "hello"}'          \
    -plaintext                                             \
    localhost:<PORT NUMBER> PubSub.Publish
```

- Для остановки сервера можно послать ему сигнал SIGINT (Ctrl+C) или SIGTERM. После этого сервер очистит ресурсы (
  закроет SubPub и остановит gRPC сервер) и завершится (используется graceful shutdown).
