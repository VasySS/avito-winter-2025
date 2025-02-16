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

Запуск нагрузочного теста:

```
make test-load
```

## Задание

### Изменения в API

1. Изменил `/api/sendCoin` на `/api/send-coin` для соответствия [RFC3986](https://datatracker.ietf.org/doc/html/rfc3986#section-6.2.2.1).
2. Изменил метод `/api/buy/{item}` с `GET` на `POST`, так как никакие данные по этому пути не возвращаются.
3. Изменил поле `type` в `InfoResponse/inventory` на `name`, так как у мерча нету типа.

### Тесты

Интеграционные тесты для всех сценариев, покрытие:

```
ok      github.com/VasySS/avito-winter-2025/internal/usecase/merch      0.010s  coverage: 85.3% of statements
ok      github.com/VasySS/avito-winter-2025/internal/usecase/auth       0.212s  coverage: 80.0% of statements
```

Результаты нагрузочного тестирования (основное ограничение из-за вычислительной сложности bcrypt):

```
checks.........................: 99.91% 14588 out of 14601
data_received..................: 4.1 MB 137 kB/s
data_sent......................: 4.0 MB 134 kB/s
http_req_blocked...............: avg=20.67µs  min=1.72µs  med=7.94µs   max=6.02ms   p(90)=16.03µs  p(95)=18.2µs
http_req_connecting............: avg=10.59µs  min=0s      med=0s       max=5.94ms   p(90)=0s       p(95)=0s
✓ http_req_duration..............: avg=40.42ms  min=2.28ms  med=11.4ms   max=364.76ms p(90)=111.48ms p(95)=139.24ms
{ expected_response:true }...: avg=40.44ms  min=2.28ms  med=11.39ms  max=364.76ms p(90)=111.49ms p(95)=139.27ms
✓ http_req_failed................: 0.11%  13 out of 11811
http_req_receiving.............: avg=82.96µs  min=12.6µs  med=70.81µs  max=3.72ms   p(90)=137.51µs p(95)=163.87µs
http_req_sending...............: avg=38.53µs  min=5.34µs  med=30.7µs   max=6.12ms   p(90)=62.25µs  p(95)=73.83µs
http_req_tls_handshaking.......: avg=0s       min=0s      med=0s       max=0s       p(90)=0s       p(95)=0s
http_req_waiting...............: avg=40.3ms   min=2.22ms  med=11.29ms  max=364.54ms p(90)=111.35ms p(95)=139.04ms
http_reqs......................: 11811  392.801227/s
iteration_duration.............: avg=160.07ms min=12.79ms med=129.41ms max=610.96ms p(90)=298.78ms p(95)=345.43ms
iterations.....................: 3007   100.004512/s
vus............................: 1      min=1              max=47
vus_max........................: 50     min=50             max=50
```
