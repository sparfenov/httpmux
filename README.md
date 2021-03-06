# HTTP-мультиплексор
- приложение представляет собой http-сервер с одним хендлером
- хендлер на вход получает POST-запрос со списком url в json-формате
- сервер запрашивает данные по всем этим url и возвращает результат клиенту в json-формате
- если в процессе обработки хотя бы одного из url получена ошибка, обработка всего списка прекращается и клиенту возвращается текстовая ошибка

Ограничения:
- для реализации задачи следует использовать Go 1.13 или выше
- использовать можно только компоненты стандартной библиотеки Go
- сервер не принимает запрос если количество url в нем больше 20
- сервер не обслуживает больше чем 100 одновременных входящих http-запросов
- для каждого входящего запроса должно быть не больше 4 одновременных исходящих
- таймаут на запрос одного url - секунда
- обработка запроса может быть отменена клиентом в любой момент, это должно повлечь за собой остановку всех операций связанных с этим запросом
- сервис должен поддерживать 'graceful shutdown'
- результат должен быть выложен на github

## Установка
имея локально golang:

```bash 
$ go get github.com/sparfenov/httpmux
```

или через docker:

```bash 
$ git clone github.com/sparfenov/httpmux
$ cd httpmux
$ docker build -f deployment/Dockerfile -t httpmux:test .
$ docker run -p 8080:8080 httpmux:test
```

## Использование
Пример запроса:
```
curl --header "Content-Type: application/json" \
--request POST \
--data '{"urls": ["http://ifconfig.me/all","http://ifconfig.me/ua","http://ifconfig.co/ip"]}' \
http://localhost:8080/
```

## Разработка

### Установить golang, go linter, reflex (утилита для hot reload)
- install golang https://golang.org/doc/install
- install golangci-lint
```bash 
make install-golangci-lint
```
- install reflex 
```bash
- go get -u github.com/cespare/reflex
```

### Запустить dev сборку с hot reload
```bash
make dev
```

### Перед коммитом проверить линтером
перед запуском линтера код автоматически будет отформатирован через gofmt
```bash
make lint
```
