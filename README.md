# 🔐 KVStore API

Простое и быстрое HTTP хранилище ключ-значение, написанное на Go. Использует Tarantool как хранилище, с модульной архитектурой и автогенерацией Swagger-документации.

---

## 🚀 Возможности

- ⚙️ REST-интерфейс (`GET`, `POST`, `DELETE`)
- 🧱 Чистая архитектура: handler → service → repository → storage
- 🗃️ Поддержка Tarantool в качестве хранилища
- 🧪 Покрытие тестами с использованием `gomock` и `testify`
- 📚 Swagger-документация (`swag init`)

---

## 🌐 Доступность API и документации

- API доступно по: `http://91.135.156.60/api/v1`
- Swagger UI с документацией доступен по адресу:
  `http://91.135.156.60/swagger/index.html`

> **Примечание:** В продакшене приложение слушает порт 80, поэтому доступ осуществляется без указания порта.

---

🛠️ Быстрый старт с Makefile
| Команда        | Описание                              |
| -------------- | ------------------------------------- |
| `make up`      | Собрать и запустить контейнеры в фоне |
| `make down`    | Остановить и удалить контейнеры       |
| `make restart` | Перезапустить контейнеры              |
| `make test`    | Запустить все тесты с покрытием       |

