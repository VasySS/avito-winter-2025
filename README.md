## Запуск

Скопировать значения из `.env.example` в файл `.env` и выполнить команду:

```
make compose-up
```

Запуск юнит тестов:

```
make test
```

Запуск интеграционных + юнит тестов:

```
make test-full
```

## Задание

### API

1. Изменил `/api/sendCoin` на `/api/send-coin` для соответствия [RFC3986](https://datatracker.ietf.org/doc/html/rfc3986#section-6.2.2.1).
2. Изменил метод `/api/buy/{item}` с `GET` на `POST`, так как никакие данные по этому пути не возвращаются.
3. Изменил поле `type` в `InfoResponse/inventory` на `name`, так как у мерча нету типа.
