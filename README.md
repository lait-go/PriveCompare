# PriceCompare — Сравнение цен в магазинах

Полноценное веб-приложение для поиска товаров и сравнения цен в разных магазинах.

## Быстрый старт

```bash
docker compose up --build
```

После запуска доступны:

| Сервис | URL |
|--------|-----|
| Frontend | http://localhost:5173 |
| Backend API | http://localhost:8080/api/v1 |
| Swagger | http://localhost:8080/swagger/index.html |

Миграции и тестовые данные загружаются автоматически при первом запуске.

## Демо-аккаунты

| Роль | Email | Пароль |
|------|-------|--------|
| Пользователь | user@pricecompare.ru | user123 |
| Администратор | admin@pricecompare.ru | admin123 |

## Функционал

- Поиск товаров (по названию, без учёта регистра)
- Фильтрация (категория, магазин, цена, бренд, скидки)
- Сортировка (цена, название, популярность, рейтинг магазина)
- Сравнение цен в таблице
- Корзина и оформление заказа (без оплаты)
- JWT-авторизация и история заказов
- Админ-панель (CRUD категорий, магазинов, товаров, цен, заказы)

## Технологии

**Backend:** Go, Gin, GORM, JWT, Swagger  
**Frontend:** Vue 3, Vite, Pinia, Vue Router, Axios  
**БД:** MySQL 8  
**Инфра:** Docker, Docker Compose

## Структура проекта

```
├── backend/
│   ├── cmd/server/          # Точка входа
│   ├── internal/
│   │   ├── config/          # Конфигурация
│   │   ├── database/        # Подключение к БД
│   │   ├── handlers/        # HTTP handlers
│   │   ├── middleware/      # JWT, CORS, Admin
│   │   ├── models/          # GORM модели
│   │   ├── repositories/    # Слой данных
│   │   ├── seed/            # Тестовые данные
│   │   └── services/        # Бизнес-логика
│   └── docs/                # Swagger
├── frontend/
│   └── src/
│       ├── components/
│       ├── layouts/
│       ├── pages/
│       ├── router/
│       ├── services/
│       └── stores/
├── docker-compose.yml
└── .env
```

## Конфигурация

Все настройки в файле `.env`:

```env
MYSQL_ROOT_PASSWORD=rootpassword
MYSQL_DATABASE=pricecompare
MYSQL_USER=priceuser
MYSQL_PASSWORD=pricepass
BACKEND_PORT=8080
FRONTEND_PORT=5173
JWT_SECRET=your-super-secret-jwt-key
```

## API Endpoints

| Метод | Путь | Описание |
|-------|------|----------|
| GET | /api/v1/home | Данные главной страницы |
| GET | /api/v1/products | Поиск и фильтрация товаров |
| GET | /api/v1/products/:id | Карточка товара |
| POST | /api/v1/auth/register | Регистрация |
| POST | /api/v1/auth/login | Вход |
| GET/POST | /api/v1/cart | Корзина |
| POST | /api/v1/orders | Оформление заказа |
| GET | /api/v1/orders/my | История заказов |
| * | /api/v1/admin/* | Админ CRUD |

Полная документация: http://localhost:8080/swagger/index.html

## Тестовые данные

При первом запуске создаются:
- 8 категорий
- 6 магазинов
- 120+ товаров
- цены в 3–6 магазинах на каждый товар
- демо-пользователь и администратор

## Остановка

```bash
docker compose down
```

Удаление данных БД:

```bash
docker compose down -v
```
