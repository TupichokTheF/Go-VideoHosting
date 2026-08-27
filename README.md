# Go Video Hosting

REST API видеохостинга на Go: регистрация и JWT-аутентификация, загрузка видео в S3-совместимое хранилище через presigned URL, выдача ссылок на просмотр.

Проект построен по принципам чистой архитектуры и DDD: доменный слой не зависит ни от базы данных, ни от HTTP, интерфейсы объявляются потребителем, а инфраструктура подставляется снаружи.

## Стек

| Компонент | Технология |
|---|---|
| Язык | Go 1.26 |
| HTTP-роутер | chi v5 |
| База данных | PostgreSQL 17 (pgx/v5) |
| Хранилище токенов | Redis 7 |
| Объектное хранилище | MinIO (S3-совместимое) |
| Аутентификация | JWT (golang-jwt/v5), bcrypt |
| Миграции | goose (встроены в бинарник) |
| Документация | Swagger (swaggo) |
| Конфигурация | cleanenv |
| Развёртывание | Docker, Docker Compose |

## Архитектура

```
cmd/
├── api/                   HTTP-сервер
└── migrate/               применение миграций

internal/
├── domain/                ядро: сущности, value objects, бизнес-правила
│   ├── user/              User, VO с валидацией, Repository, Hasher
│   └── video/             Video, машина состояний, VO, Repository
│
├── application/           сценарии использования
│   ├── services/          AuthService, UserService, VideoService
│   ├── ports/             интерфейсы инфраструктуры (JWT, Cache, Storage)
│   ├── dtos/              передача данных между слоями
│   └── errors/            ошибки уровня приложения
│
├── infrastructure/        реализации внешних зависимостей
│   ├── repositories/      PostgreSQL
│   ├── storage/           MinIO
│   ├── cache/             Redis
│   ├── security/          JWT-менеджер, bcrypt
│   └── database/          пулы подключений
│
└── presentation/          HTTP-слой
    ├── handlers/          обработчики запросов
    ├── routers/           сборка маршрутов через функциональные опции
    ├── middleware/        проверка авторизации
    ├── mappers/           схемы ↔ DTO, доменные ошибки → HTTP-статусы
    ├── schemas/           контракты запросов и ответов
    ├── ports/             интерфейсы сервисов
    └── context/           типобезопасные ключи контекста
```

Зависимости направлены внутрь: `presentation → application → domain`. Инфраструктура реализует интерфейсы, объявленные во внутренних слоях, — заменить PostgreSQL на другую СУБД или MinIO на AWS S3 можно, не трогая бизнес-логику.

**Доменный слой** защищает инварианты: поля сущностей неэкспортируемые, изменение состояния идёт только через методы, value objects нельзя сконструировать невалидными. Восстановление объектов из базы вынесено в отдельные функции (`Reconstitute`), которые намеренно обходят валидацию — данные уже прошли её при записи.

## Как устроена загрузка видео

Файл не проходит через API: сервис выдаёт временную подписанную ссылку, а клиент загружает видео напрямую в хранилище.

```
1. POST /video/add               клиент → API
                                 создаётся запись (status=draft)
                                 API → клиент: video_id + presigned URL

2. PUT <presigned_url>           клиент → MinIO (байты минуют API)

3. POST /video/{id}/complete     API проверяет наличие объекта в хранилище
                                 status: draft → uploaded
```

Так сервис не тратит память и время на передачу гигабайтных файлов, а хранилище масштабируется независимо от API.

Статусы видео образуют машину состояний, переходы валидируются в домене:

```
draft → uploaded → processing → ready
                 ↘  failed   ↗
```

## Быстрый старт

Требуется Docker и Docker Compose.

```bash
git clone https://github.com/TupichokTheF/Go-VideoHosting.git
cd Go-VideoHosting

cp .env.example .env
# заполнить переменные, сгенерировать ключи:
#   openssl rand -base64 32

docker compose up -d --build
```

Compose поднимет PostgreSQL, Redis и MinIO, дождётся их готовности, применит миграции и запустит API.

| Сервис | Адрес |
|---|---|
| API | http://localhost:8080 |
| Swagger UI | http://localhost:8080/swagger/index.html |
| Консоль MinIO | http://localhost:9001 |

Логи и остановка:

```bash
docker compose logs -f api
docker compose down          # остановить
docker compose down -v       # остановить и удалить данные
```

## Локальный запуск без контейнера API

Удобно во время разработки — инфраструктура в Docker, приложение снаружи:

```bash
docker compose up -d postgres redis minio
go run ./cmd/migrate
go run ./cmd/api
```

Миграции можно применять и по частям:

```bash
go run ./cmd/migrate -cmd=status   # что применено
go run ./cmd/migrate -cmd=down     # откатить одну
```

Новая миграция:

```bash
go tool goose -dir migrations create create_comments sql
```

Обновить Swagger после изменения аннотаций:

```bash
go tool swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

## Конфигурация

Все параметры читаются из переменных окружения; `.env` используется при локальном запуске.

| Переменная | Описание |
|---|---|
| `HTTP_HOST`, `HTTP_PORT` | адрес HTTP-сервера |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | PostgreSQL |
| `REDIS_HOST`, `REDIS_PORT` | Redis |
| `MINIO_HOST`, `MINIO_PORT` | MinIO |
| `MINIO_USER`, `MINIO_PASSWORD` | доступ к MinIO |
| `MINIO_BUCKET` | имя бакета (по умолчанию `videos`) |
| `MINIO_TTL` | срок жизни presigned URL (по умолчанию `30m`) |
| `ACCESS_SECRET_KEY`, `REFRESH_SECRET_KEY` | ключи подписи JWT |
| `ACCESS_TOKEN_TTL`, `REFRESH_TOKEN_TTL` | время жизни токенов |
| `SWAGGER` | включить Swagger UI |

## API

Базовый путь — `/api/v1`.

| Метод | Путь | Доступ | Описание |
|---|---|---|---|
| `POST` | `/auth/register` | — | регистрация |
| `POST` | `/auth/login` | — | вход: access-токен в теле, refresh в cookie |
| `POST` | `/auth/refresh` | cookie | обновление access-токена |
| `POST` | `/auth/logout` | cookie | выход, отзыв refresh-токена |
| `GET` | `/user/me` | Bearer | профиль текущего пользователя |
| `POST` | `/video/add` | Bearer | создать запись и получить ссылку на загрузку |
| `POST` | `/video/{video_id}/complete` | Bearer | подтвердить загрузку файла |
| `GET` | `/video/get` | — | получить ссылку на просмотр |

Access-токен передаётся в заголовке `Authorization: Bearer <token>` и живёт минуты. Refresh-токен хранится в HttpOnly-cookie и в Redis, что позволяет отзывать сессию при выходе.

## Пример сценария

```bash
# Регистрация
curl -X POST localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"vasya","email":"vasya@example.com","password":"Passw0rd"}'

# Вход — refresh-cookie сохраняется в cookies.txt
curl -c cookies.txt -X POST localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"vasya","password":"Passw0rd"}'

TOKEN="<access_token из ответа>"

# Создание записи о видео
curl -X POST localhost:8080/api/v1/video/add \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Мой ролик","description":"Описание"}'

# Загрузка файла напрямую в хранилище
curl -X PUT "<upload_url>" --upload-file video.mp4

# Подтверждение загрузки
curl -X POST localhost:8080/api/v1/video/<video_id>/complete \
  -H "Authorization: Bearer $TOKEN"
```

## Планы

- [ ] Воркер транскодинга: ffmpeg, несколько разрешений, генерация превью
- [ ] Адаптивный стриминг (HLS)
- [ ] Уровни видимости видео: публичное, по ссылке, приватное
- [ ] Лента, поиск и пагинация
- [ ] Комментарии и счётчик просмотров
- [ ] Ротация refresh-токенов
- [ ] Тесты доменного слоя и интеграционные тесты
- [ ] CI на GitHub Actions

## Лицензия

MIT
