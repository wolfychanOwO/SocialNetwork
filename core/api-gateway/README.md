# API Gateway

## Зона ответственности

- Принимает запросы от пользователя.
- Проверяет аутентификацию (например, наличие и валидность токена) совместно с User Service.
- **Маршрутизизация** запросов к User Service, Post Service, Statistics Service
- Выступает **единой точкой входа** в архитектуру для внешних клиентов.

## Границы сервиса

- **Не** хранит основные бизнес-данные (пользователей, постов и т.д.) — для этого существуют соответствующие сервисы.
- **Не** реализует сложную бизнес-логику, а только проксирует/маршрутизирует запросы.
- Может иметь **минимальное** локальное состояние (например, кэш сессий), но все «тяжёлые» операции делегирует целевым сервисам.
- Все внешние клиенты обращаются к системе **только** через этот Gateway.

## Основные технологии

- **Python**
- **Web-фреймворк** (FastAPI)
- **Docker** (приложение упаковано в Docker-образ)
- (Необязательно) **Redis**, если нужно хранить короткоживущие данные.

