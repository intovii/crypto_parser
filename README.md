# crypto_parser
## Описание
Микросервис, который собирает, хранит и отображает стоимости криптовалют.
Сбор данных производится каждые 5 секунд. 
Источник данных: `https://api.binance.com/api/v3/ticker/price?symbol=`

## API Endpoints (Задание)
`/currency/add` - Добавление криптовалюты в список наблюдения
Добавление криптовалюты в список наблюдения подразумевает что мы собираем и записываем цену криптовалюты в локальную базу раз в N секунд. 

`/currency/remove` - Удаление криптовалюты из списка наблюдения
Удаление из этого списка останавливает сбор цены. Для получения актуальной цены можно воспользоваться открытыми API.

`/currency/price` - Получение цены криптовалюты
Должна быть возможность получить цену конкретной криптовалюты в конкретный момент времени из локальной базы. 
Пример: запрос {"coin": "BTC", "timestamp": 1736500490}, сервис должен вернуть стоимость BTC в момент времени 1736500490. Если не удалось найти стоимость в конкретный момент времени, возвращаем стоимость в ближайший к запрошенному моменту времени.

`/swagger/index.html` - Документация Swagger.

## Инструкция по запуску
Склонировать репозиторий и перейти в рабочую директорию.
``` bash
git clone https://github.com/intovii/crypto_parser.git
cd crypto_parser
```
Запуск проекта при помощи Docker Compose.
``` bash
docker compose up -d
```
Остановка проекта и удаление соответсвующего тома.
``` bash
docker compose down -v
```

## Документация OpenAPI
Для получения детальной информации по запосам, ответам и параметрам перейдите по ссылке `http://localhost:3000/swagger/index.html` после запуска проекта.

## Тестирование
Postman.
Импортируйте `TestCryptoParser.postman_collection.json`.
- 4 Запроса:
    - запрос к API binance:
        `BinanceAPI` для быстрой и удобной проверки;
    - запросы к разработанному API:
        `CurrencyAdd` для создания валютной пары `BTCUSDT`;
        `CurrencyRemove` для удаления валютной пары `BTCUSDT`;
        `CurrencyPrice` для получения цены по валютной паре `BTCUSDT` в момент времени `1737392194`.