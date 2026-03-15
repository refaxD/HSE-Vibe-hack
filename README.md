# HSE-Vibe-hack — Trip Tuner

**Команда: Х Х и в продакшн**

Trip Tuner — приложение, которое **строит живые маршруты по текстовому описанию пользователя и его предпочтениям**.  
Пример: запрос «свидание, зайти в кафе, посмотреть фильм» превращается в маршрут вроде «парк Горького → кафе рядом с парком → кинотеатр → прогулка возле фонтанов».  
Пользователь может задавать **начальную и конечную точки**, бюджет, длительность, тип активности и другие параметры — система адаптирует маршрут под запрос.

Небольшая сводка проекта доступна в Google Docs (product overview), README фокусируется на технической стороне и локальном запуске.

---

## Зачем это нужно (продуктовый обзор)

- **Экономия времени на планирование**: вместо ручного поиска мест и маршрутов пользователь формулирует задачу одним-двумя предложениями, а маршрут подбирается автоматически.
- **Маршруты под конкретный контекст**: свидание, прогулка с друзьями, семейный выход, «быстро выдохнуть после работы», туристический маршрут в новом районе и т.д.
- **Учёт предпочтений и ограничений**: можно учитывать бюджет, тип мест (кафе, парки, бары, музеи), расстояние пешком/транспортом, время суток.
- **Освоение города**: Trip Tuner помогает «открывать» город заново — подсказывает новые места и связки, о которых пользователь сам бы не догадался.

---

## Функциональность

- **Карта и геолокация** — отображение карты (Yandex Maps), текущей позиции пользователя и точек маршрута.
- **Каталог мест и категорий** — места (кафе, рестораны, парки и т.д.) и события с привязкой к категориям; фильтрация и поиск.
- **Генерация маршрута по промпту** — текстовый запрос обрабатывается бэкендом (в т.ч. с использованием нейросети), в результате выдаётся последовательность мест/событий и сегменты маршрута.
- **Визуализация маршрута** — отрисовка маршрута и сегментов на карте, панель с информацией о точках и пути.
- **Авторизация** — регистрация, вход и профиль пользователя; персонализация на уровне аккаунта.
- **Мобильная сборка** — поддержка сборки под Android/iOS через Capacitor (режим PWA/мобильного приложения).

---

## Технологический стек

С учётом кода в репозитории и изначальной архитектуры проекта:

| Часть       | Технологии |
|------------|------------|
| Frontend   | Angular 19, TypeScript, Tailwind CSS, Yandex Maps API, Capacitor (Android/iOS) |
| Backend (текущая реализация) | Go 1.23, Fiber v2, Swagger (swaggo), MongoDB, Redis |
| Backend (изначальный дизайн) | Node.js (Koa, tsoa, TypeORM), MongoDB, RabbitMQ |
| Маршруты   | OpenRouteService (ORS) / собственный ORS-инстанс в Docker |
| Инфраструктура | Docker, Docker Compose |

---

## Установка и запуск

### Требования

- Node.js и npm (для frontend)
- Go 1.23+ (для backend)
- Docker и Docker Compose (для MongoDB, Redis и при необходимости ORS)
- Ключи: Yandex Map API, OpenRouteService (или локальный ORS), при использовании парсинга/LLM — API ключи (см. переменные окружения)

### 1. Клонирование репозитория

```bash
git clone <url-репозитория>
cd HSE-Vibe-hack
```

### 2. Backend (Go API)

```bash
cd backend
```

Создайте файл `.env` в каталоге `backend` (см. раздел «Переменные окружения»). Запустите MongoDB и Redis — например, через Docker Compose в корне проекта или только для backend:

```bash
# из каталога backend
docker compose up -d
```

Установите зависимости и запустите API:

```bash
go mod download
go run ./cmd/api
```

API будет доступен по адресу `http://localhost:3000`. Swagger UI: `http://localhost:3000/swagger/index.html`.

### 3. Frontend (Angular)

```bash
cd frontend
npm install
npm start
```

Приложение откроется по адресу `http://localhost:4200`. Для работы карты и маршрутов в `frontend/src/environments/environment.ts` должны быть заданы `apiBaseUrl` (адрес backend), `YandexMapKey` и при необходимости `RouteKey`/`RouteUrl` (OpenRouteService).

### 4. Запуск через Docker (вся инфраструктура)

В корне проекта:

```bash
# MongoDB, Redis, OpenRouteService
docker compose up -d
```

Backend и frontend в текущем виде в корневом `docker-compose.yaml` закомментированы; при необходимости их можно раскомментировать и донастроить (см. комментарии в файле).

### 5. OpenRouteService (опционально)

Для маршрутизации можно использовать либо внешний API (`https://api.openrouteservice.org`), либо локальный ORS:

- Локальный ORS: каталог `open-route-service` с отдельным `docker-compose.yml`.
- В основном проекте сервис `ors-app` уже описан в корневом `docker-compose.yaml`.

---

## Структура проекта

```
HSE-Vibe-hack/
├── backend/                    # Go API
│   ├── cmd/api/                 # Точка входа приложения
│   ├── internal/
│   │   ├── delivery/http/       # HTTP-обработчики, роутер, middleware (auth)
│   │   ├── domain/              # Доменные модели и интерфейсы репозиториев
│   │   ├── repository/          # MongoDB и Redis реализации
│   │   │   ├── mongo/           # user, category, place, event_place, path, path_segment
│   │   │   └── redis/           # session
│   │   └── usecase/             # Бизнес-логика: user, category, place, path, algorithm (генерация маршрута)
│   ├── pkg/config/              # Загрузка конфигурации и env
│   ├── docs/                    # Swagger (docs.go, swagger.json, swagger.yaml)
│   ├── Dockerfile
│   └── docker-compose.yml       # API + MongoDB + Redis
├── frontend/                    # Angular-приложение (trip-tuner)
│   ├── src/
│   │   ├── app/
│   │   │   ├── core/            # Сервисы (api, map, notifications), модели, guards, компоненты (bottom-bar, slide-categories)
│   │   │   ├── features/home/   # Главный экран с картой и маршрутом
│   │   │   ├── pop-up-menus/   # Компоненты: информация о точке, о маршруте, уведомления
│   │   │   └── generated/      # OpenAPI-клиент и модели
│   │   └── environments/       # environment.ts, environment.prod.ts
│   ├── Dockerfile
│   └── package.json
├── open-route-service/          # Конфигурация и docker-compose для локального ORS
└── docker-compose.yaml         # MongoDB, Redis, ORS (backend/frontend опционально)
```

---

## Архитектура и логика работы

```mermaid
flowchart LR
  subgraph client [Клиент]
    UI[Angular SPA]
  end
  subgraph backend [Backend]
    Handlers[HTTP Handlers]
    Usecases[Usecases]
    Repos[Repositories]
  end
  subgraph storage [Хранилище]
    Mongo[(MongoDB)]
    Redis[(Redis)]
  end
  subgraph external [Внешние сервисы]
    ORS[OpenRouteService]
    LLM[LLM API]
  end
  UI --> Handlers
  Handlers --> Usecases
  Usecases --> Repos
  Usecases --> ORS
  Usecases --> LLM
  Repos --> Mongo
  Repos --> Redis
```

- **Клиент (Angular):** запросы к API через сгенерированный OpenAPI-клиент; карта (Yandex Maps), отображение точек и маршрута, поиск мест, ввод промпта для маршрута.
- **Backend (Go / Node.js):** HTTP‑API для аутентификации, категорий, мест, событий и маршрутов (`/path/create`, получение списка и сегментов). В изначальном дизайне использовались Koa + tsoa + TypeORM и RabbitMQ для фоновой обработки; в текущей реализации в репозитории используется Go + Fiber.
- **Генерация маршрута:** пользовательский промпт передаётся в нейросеть, которая разбирает текст на категории, фиксированные точки и события. После этого маршрут подбирается алгоритмом (итерационный отжиг) с учётом минимизации времени/расстояния/стоимости прогулки и ограничений пользователя (бюджет, длительность и т.д.). Для расчёта расстояний и времени используются данные ORS.
- **Данные:** пользователи, категории, места (в том числе спаршенные из open‑source БД Mos.ru), события и маршруты (path, path_segment) хранятся в MongoDB; сессии и временные данные — в Redis (и/или брокере сообщений RabbitMQ в изначальной версии).

---

## Переменные окружения

### Backend (`backend/.env`)

| Переменная        | Пример / значение по умолчанию | Назначение |
|-------------------|---------------------------------|------------|
| `PORT`            | `3000`                          | Порт HTTP-сервера |
| `MONGO_URI`       | `mongodb://user:pass@mongo:27017/apidb?authSource=admin` | Строка подключения к MongoDB |
| `MONGO_USER`      | —                               | Пользователь MongoDB (для docker-compose) |
| `MONGO_PASSWORD`  | —                               | Пароль MongoDB |
| `MONGO_DB`        | `apidb`                         | Имя БД |
| `REDIS_HOST`      | `localhost`                     | Хост Redis |
| `REDIS_PORT`      | `6379`                          | Порт Redis |
| `REDIS_TTL`       | `86400`                         | TTL сессии в секундах |
| `ROUTER_URL`      | `https://api.openrouteservice.org` | URL OpenRouteService |
| `ROUTER_KEY`      | —                               | Ключ API OpenRouteService |
| `OPENAI_API_KEY`  | —                               | Ключ OpenAI (для разбора промпта и генерации маршрута) |
| `MOS_DATA_KEY`    | —                               | Ключ Mos Data (если используется) |
| `TIMEPAD`         | —                               | Ключ Timepad (если используется) |
| `PARSE_ENTITIES`  | `false`                         | Запускать ли парсер сущностей при старте |

### Frontend

В `frontend/src/environments/environment.ts` (и `environment.prod.ts` для продакшена) задаются:

- `apiBaseUrl` — базовый URL backend (например, `http://localhost:3000`);
- `YandexMapKey` — ключ Yandex Maps API;
- `RouteKey`, `RouteUrl` — ключ и URL OpenRouteService.

---

## Тестирование

- **Backend:** тесты можно запускать стандартными средствами Go, например: `go test ./...` из каталога `backend`.
- **Frontend:** в проекте настроены Jasmine и Karma; запуск тестов через Angular CLI: `npm test` (из каталога `frontend`).

---

## Ограничения и дальнейшие планы

- В корневом `docker-compose.yaml` сервисы backend и frontend закомментированы — для полного запуска «одной командой» их нужно включить и при необходимости донастроить env и порты.
- Генерация маршрута по промпту зависит от наличия ключей внешних API (OpenRouteService, при необходимости OpenAI и др.); без них часть сценариев может быть недоступна.
- Возможные улучшения: расширение набора категорий и источников событий, офлайн-режим, доработка мобильных сборок (Capacitor), добавление тестов и CI.

---

## Команда

**Х Х и в продакшн** — проект для HSE VibeHack 15.03.
