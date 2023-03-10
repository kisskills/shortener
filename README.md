# Cервис, предоставляющий API по созданию сокращённых ссылок.
___
## Функционал

Сервис принимает следующие запросы по http и grpc:
1. Метод Post,который сохраняет оригинальный URL в базе и возвращает сокращённый.
2. Метод Get, который принимает сокращённый URL и возвращает оригинальный.
___
## Реализация

Сервис работает через GRPC. Распространяется в виде Docker-образа.
Можно запустить сервис и вне контейнера (командой `go run cmd/main.go`). Он будет запущен с in-memory хранилищем.
Если создать env переменную `STORAGE_MODE`, и задать ей значение `postgres`, можно будет настроить подключение к PostgreSQL.

Сервис принимает http-запросы на порту 8080, rpc-запросы на порту 6565
___
## Запуск
```shell
# Клонируем репозиторий
> git clone https://github.com/kisskills/shortener
```
```shell
# Генерируем proto
> make gen-grpc
```
```shell
# Запуск с PostgreSQL
> make postgres
```
```shell
# Запуск с in-memory решением
> make inmemory
```
___
## Запросы
### Общий вид

`GET /get/{short_link}`
```
POST /create
body: {"link": "https://github.com"}
```

### Примеры
* Получение оригинальной ссылки
```shell
curl -X 'GET' \
  'http://localhost:8080/get/sHkrVHsldk' \
  -H 'accept: application/json'
```
* Создание короткой ссылки
```shell
curl -X 'POST' \
  'http://localhost:8080/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "link": "https://github.com"
}'
```
___
## Прочее
```shell
# Для завершения
> make clean
# Для тестирования
> make test-cover
```
Для создания коротких ссылок был выбран алгоритм хэширования xxHash версии xxh3.
Основывался на скорости работы и количестве коллизий (Рассматривался так же вариант от Google - city64)

> https://github.com/Cyan4973/xxHash

> https://github.com/Cyan4973/xxHash/wiki/Collision-ratio-comparison
