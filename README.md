# TechZone

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)](https://react.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)
[![Swagger](https://img.shields.io/badge/OpenAPI-Swagger-85EA2D?logo=swagger&logoColor=black)](https://swagger.io/)
[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Live Demo](https://img.shields.io/badge/Live-Demo-2563EB?style=for-the-badge)](https://techzone-phi.vercel.app)

TechZone — учебный проект интернет-магазина электроники, разработанный на **Go** с использованием многослойной архитектуры.

Основная цель проекта — не просто реализовать CRUD-приложение, а показать построение полноценного backend-сервиса с разделением ответственности между слоями, работой с транзакциями, кэшированием, авторизацией, асинхронной обработкой событий и современным процессом разработки.

## Демо

| Сервис | Ссылка |
|--------|---------|
| 🛍️ Frontend | https://techzone-phi.vercel.app |
| 📚 Swagger | https://techzone-production-e050.up.railway.app/swagger/index.html |
| ⚙️ Backend API | https://techzone-production-e050.up.railway.app/ |


---

## Содержание

- [Возможности](#возможности)
- [Используемые технологии](#используемые-технологии)
- [Архитектура](#архитектура)
- [Структура проекта](#структура-проекта)
- [Быстрый запуск](#быстрый-запуск)
- [Переменные окружения](#переменные-окружения)
- [API](#api)
- [Swagger](#swagger)
- [Тестирование](#тестирование)
- [Roadmap](#roadmap)

---

# Возможности

### Пользователи

- регистрация аккаунта;
- авторизация через JWT;
- получение информации о текущем пользователе;
- разделение ролей (`client` и `admin`).

### Каталог товаров

- просмотр каталога;
- просмотр товара по ID;
- создание товара;
- изменение товара;
- скрытие и восстановление товара;
- хранение изображений по URL.

### Корзина

- добавление товара;
- просмотр корзины;
- удаление товара;
- оформление заказа.

### Заказы

- создание заказа из корзины;
- просмотр списка заказов;
- просмотр конкретного заказа;
- изменение статуса заказа;
- автоматическое уменьшение остатков товаров.

### Оплата

- mock payment gateway;
- идемпотентность платежей;
- изменение статуса оплаты.

### Уведомления

- асинхронная обработка уведомлений;
- Worker Pool;
- подготовленная Kafka Event Driven архитектура.

---

# Используемые технологии

| Категория | Технологии |
|-----------|------------|
| Backend | Go, net/http |
| Database | PostgreSQL, pgx |
| Cache | Redis |
| Messaging | Kafka (Franz-Go) |
| Authentication | JWT |
| Password Hashing | bcrypt |
| Documentation | Swagger/OpenAPI |
| Frontend | React, React Router, Axios |
| Infrastructure | Docker, Docker Compose |
| Deploy | Railway, Vercel |
| Testing | Unit Tests, Integration Tests |

---

# Архитектура

Проект построен по принципам Layered Architecture.

```text
                    React Frontend
                           │
                           ▼
                    HTTP REST API
                           │
                           ▼
                      HTTP Handlers
                           │
                           ▼
                     Business Logic
                        (Services)
                           │
                           ▼
                    Repository Layer
                           │
          ┌────────────────┴────────────────┐
          ▼                                 ▼
    PostgreSQL                         Redis Cache
          │
          ▼
     Transactions
          │
          ▼
      Kafka Events
          │
          ▼
      Worker Pool
          │
          ▼
     Notifications
```

Каждый слой отвечает только за свою область ответственности:

- **Handler** — обработка HTTP-запросов;
- **Service** — бизнес-логика;
- **Repository** — работа с PostgreSQL;
- **Middleware** — авторизация и проверка ролей;
- **Redis** — кэширование каталога товаров;
- **Worker Pool** — асинхронная обработка уведомлений;
- **Kafka** — публикация событий (архитектура подготовлена для работы с брокером сообщений).

---

# Структура проекта

```text
techzone/
│
├── cmd/
│   └── server/
│
├── frontend/
│
├── internal/
│   ├── config/
│   ├── event/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── payment/
│   ├── repository/
│   ├── service/
│   └── tests/
│
├── migrations/
│
├── pkg/
│   ├── jwt/
│   ├── kafka/
│   ├── postgres/
│   └── redis/
│
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

---

# Быстрый запуск

## Клонирование проекта

```bash
git clone https://github.com/USERNAME/techzone.git
cd techzone
```

## Запуск через Docker

```bash
docker compose up --build
```

После запуска будут доступны:

| Сервис | Адрес |
|--------|--------|
| Backend | http://localhost:8080 |
| Frontend | http://localhost:5173 |
| Swagger | http://localhost:8080/swagger/index.html |

---

# Миграции

Применить миграции

```bash
make migrate-up
```

Откатить последнюю миграцию

```bash
make migrate-down
```

---

# Переменные окружения

```env
DB_URL=postgres://postgres:1234@postgres:5432/study?sslmode=disable

JWT_SECRET=your-secret

REDIS_ADDR=localhost:6379

KAFKA_BROKERS=localhost:9092
```

# API

Ниже представлены основные HTTP-эндпоинты приложения.

| Метод | Endpoint | Описание |
|--------|----------|----------|
| POST | `/register` | Регистрация пользователя |
| POST | `/login` | Авторизация |
| GET | `/me` | Информация о текущем пользователе |
| GET | `/products` | Получить каталог товаров |
| GET | `/products/{id}` | Получить товар по ID |
| POST | `/products` | Создать товар *(admin)* |
| PUT | `/products/{id}` | Обновить товар *(admin)* |
| PATCH | `/products/{id}/status` | Скрыть / восстановить товар *(admin)* |
| POST | `/cart/items` | Добавить товар в корзину |
| GET | `/cart` | Получить корзину |
| DELETE | `/cart/items/{id}` | Удалить товар из корзины |
| POST | `/orders` | Создать заказ |
| GET | `/orders` | Получить список заказов |
| GET | `/orders/{id}` | Получить заказ по ID |
| PATCH | `/orders/{id}/status` | Изменить статус заказа *(admin)* |
| POST | `/payments` | Выполнить оплату |
| GET | `/notifications` | Получить уведомления |

---

# Swagger

Проект содержит автоматически генерируемую документацию Swagger.

После запуска приложения она доступна по адресу

```
http://localhost:8080/swagger/index.html
```

Swagger позволяет:

- просматривать все доступные эндпоинты;
- выполнять запросы непосредственно из браузера;
- изучать модели запросов и ответов;
- тестировать API без использования Postman.

---

# Авторизация

После успешной авторизации сервер возвращает JWT-токен.

Все защищённые маршруты используют заголовок

```http
Authorization: Bearer <token>
```

Маршруты администратора дополнительно проверяют роль пользователя.

---

# Кэширование

Для уменьшения нагрузки на базу данных используется Redis.

В настоящий момент кэшируется:

- список товаров;
- автоматическая инвалидизация кэша при создании;
- инвалидизация после изменения товара;
- инвалидизация после скрытия товара.

---

# Асинхронная обработка

При оформлении заказа сервис публикует события.

В проекте используется следующая схема:

```
Создание заказа
        │
        ▼
 Publish Event
        │
        ▼
     Kafka
        │
        ▼
 Consumer
        │
        ▼
 Worker Pool
        │
        ▼
 Notification Service
```

При отсутствии Kafka приложение автоматически переключается на локальный Publisher, что позволяет запускать проект без брокера сообщений.

---

# Транзакции

Создание заказа выполняется внутри транзакции PostgreSQL.

В одну транзакцию входят:

- создание заказа;
- создание позиций заказа;
- уменьшение остатков товаров;
- очистка корзины.

Если одна из операций завершается ошибкой, изменения откатываются.

---

# Тестирование

Проект содержит unit- и integration-тесты.

Запуск unit-тестов

```bash
make unit
```

Запуск integration-тестов

```bash
make integration
```

---

# Разработка

Основные команды Makefile

```bash
make run
```

Запуск сервера.

```bash
make build
```

Сборка приложения.

```bash
make fmt
```

Форматирование кода.

```bash
make vet
```

Статический анализ.

```bash
make lint
```

Проверка линтером.

```bash
make migrate-up
```

Применить миграции.

```bash
make migrate-down
```

Откатить миграции.

---

# Назначение администратора

После регистрации пользователя роль администратора можно выдать SQL-запросом.

```sql
UPDATE users
SET role = 'admin'
WHERE login = 'admin';
```

---

# Деплой

Проект развёрнут в облаке.

| Компонент | Платформа |
|-----------|-----------|
| Backend | Railway |
| Frontend | Vercel |


---

# Скриншоты

## Интерфейс

### Каталог

![Catalog](docs/catalog.png)

### Корзина

![Cart](docs/cart.png)

### Профиль

![Profile](docs/profile.png)

### Админ-панель

![Admin](docs/admin.png)

### Swagger

![Swagger](docs/swagger.png)
```

---

# Roadmap

Планируемые улучшения проекта.

- [x] JWT-аутентификация
- [x] Role Based Access Control
- [x] Управление товарами
- [x] Корзина
- [x] Заказы
- [x] Redis Cache
- [x] Worker Pool
- [x] Docker
- [x] Swagger
- [x] React Frontend
- [x] CI/CD
- [ ] Загрузка изображений
- [ ] Поиск товаров
- [ ] Фильтрация каталога
- [ ] WebSocket-уведомления
- [ ] Интеграция с реальным платёжным сервисом

---

# О проекте

TechZone создавался как учебный проект для практики разработки backend-приложений на Go.

Основной акцент сделан на архитектуре приложения, разделении ответственности между слоями, работе с PostgreSQL, Redis, транзакциями, JWT-аутентификации и асинхронной обработке событий. Проект постепенно развивается и используется как площадка для изучения современных подходов к разработке серверных приложений.

---

# Лицензия

Проект распространяется по лицензии **MIT**.

При желании вы можете использовать его в качестве основы для собственных учебных проектов.