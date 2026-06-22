# TechZone

TechZone — учебный проект интернет-магазина электроники, разработанный на Go с использованием PostgreSQL и Docker.

Проект реализует полный цикл работы интернет-магазина: от регистрации пользователя и управления корзиной до оформления заказов и администрирования товаров.

---

## Функциональность

### Пользователи

* Регистрация аккаунта
* Авторизация через JWT
* Получение информации о текущем пользователе

### Товары

* Просмотр каталога товаров
* Просмотр товара по ID
* Создание товара (только для администратора)

### Корзина

* Добавление товаров в корзину
* Просмотр содержимого корзины
* Удаление товаров из корзины

### Заказы

* Создание заказа из корзины
* Просмотр списка заказов пользователя
* Просмотр заказа по ID
* Изменение статуса заказа (только для администратора)

### Дополнительно

* JWT-аутентификация
* Ролевая модель (client/admin)
* Транзакции при оформлении заказа
* Асинхронная обработка уведомлений через Worker Pool
* Unit-тесты бизнес-логики

---

## Технологический стек

### Backend

* Go
* PostgreSQL
* pgx
* JWT
* bcrypt

### Infrastructure

* Docker
* Docker Compose

### Frontend

* React
* Axios
* React Router

---

## Архитектура проекта

```text
techzone/
├── cmd/
│   └── server/
├── internal/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── service/
│   └── worker/
├── migrations/
├── pkg/
│   ├── jwt/
│   └── postgres/
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

Проект построен по принципам многослойной архитектуры:

* Handler — HTTP-обработчики
* Service — бизнес-логика
* Repository — работа с базой данных
* Middleware — авторизация и проверка ролей
* Worker Pool — асинхронная обработка задач

---

## Запуск проекта

### Клонирование репозитория

```bash
git clone <repository-url>
cd techzone
```

### Сборка контейнеров

```bash
docker compose build
```

### Запуск приложения

```bash
docker compose up
```

После запуска приложение будет доступно по адресу:

```text
http://localhost:8080
```

---

## Миграции

Применить миграции:

```bash
make migrate-up
```

Откатить последнюю миграцию:

```bash
make migrate-down
```

---

## Переменные окружения

```env
DB_URL=postgres://postgres:1234@postgres:5432/study
JWT_SECRET=your-secret-key
```

---

## Основные API

### Регистрация

```http
POST /register
```

```json
{
  "login": "admin",
  "email": "admin@mail.ru",
  "password": "123456"
}
```

### Авторизация

```http
POST /login
```

```json
{
  "login": "admin",
  "password": "123456"
}
```

### Получение товаров

```http
GET /products
```

### Создание товара

```http
POST /products
```

Требуется JWT-токен администратора.

```json
{
  "name": "iPhone 17",
  "description": "Apple smartphone",
  "price": 99990,
  "stock": 10
}
```

### Добавление товара в корзину

```http
POST /cart/items
```

### Удаление товара из корзины

```http
DELETE /cart/items/{item_id}
```

### Создание заказа

```http
POST /orders
```

### Получение заказов пользователя

```http
GET /orders
```

### Получение заказа по ID

```http
GET /orders/{id}
```

### Изменение статуса заказа

```http
PATCH /orders/{id}/status
```

Доступно только администраторам.

---

## Назначение роли администратора

Для выдачи роли администратора выполните SQL-запрос:

```sql
UPDATE users
SET role = 'admin'
WHERE login = 'admin';
```

---

## Будущие улучшения

* Swagger/OpenAPI документация
* Deploy в облако (Render)
* Kafka для асинхронной обработки событий
* Email-уведомления
* Интеграционные тесты
* CI/CD pipeline

```
```
