# TechZone

Учебный проект интернет-магазина электроники, разработанный мной на Go.

## Возможности:
- Регистрация пользователей,
- Авторизация через JWT
- Просмотр товаров
- Добавление товаров (только администратор)
- Корзина покупок
- Создание заказов
- Просмотр заказов пользователя
- Изменение статуса заказа (только администратор)

PostgreSQL в качестве базы данных,
Docker и Docker Compose

## Стек технологий:
### Backend
- Go
- PostgreSQL
- pgx
- JWT
- bcrypt
- Инфраструктура
- Docker
- Docker Compose

## Особенности

- JWT-аутентификация
- Ролевая модель (client/admin)
- Транзакции при создании заказа
- Docker Compose для развертывания
- PostgreSQL в качестве СУБД

## Архитектура проекта:
```
techzone/
├── cmd/
│   └── server/
├── internal/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   └── service/
├── migrations/
├── pkg/
│   ├── jwt/
│   └── postgres/
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

- Запуск через Docker
- Сборка:
docker compose build
- Запуск:
docker compose up

- После запуска приложение будет доступно по адресу:
http://localhost:8080

## Миграции
- Применить миграции:
make migrate-up

- Откатить последнюю миграцию:
make migrate-down

## Основные API:

- Регистрация:
```
POST /register
{
  "login": "admin",
  "email": "admin@mail.ru",
  "password": "123456",
}
```

- Авторизация:

```
POST /login
{
  "login": "admin",
  "password": "123456"
}
```

- Получение списка товаров:

```
GET /products
```

- Создание товара:

```
POST /products
(Требуется JWT-токен администратора.)
{
  "name": "iPhone 17",
  "description": "Apple smartphone",
  "price": 99990,
  "stock": 10
}
```

- Добавление товара в корзину:
```
POST /cart/items
```

- Удалить товар из корзины по ID:

```
DELETE /cart/items/{item_id}
```

- Создание заказа:

```
POST /orders
```

- Получение заказов пользователя:

```
GET /orders
```

- Получение заказа по ID:

```
GET /orders/{id}
```

- Изменение статуса заказа:

```
PATCH /orders/{id}/status
(Доступно только администратору.)
```

- Для назначения роли администратора, использовать команду в бд:

```
UPDATE users
SET role = 'admin'
WHERE login = 'admin';
```

- Присутствуют переменные окружения
