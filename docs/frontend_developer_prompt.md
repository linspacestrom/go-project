# Prompt для Frontend-разработчика (Student & T)

Используй этот текст как полное ТЗ на фронтенд-интеграцию с текущим backend.

## Роль и цель

Ты Senior Frontend Developer.  
Нужно реализовать production-ready web-клиент для платформы Student & T (роли: `student`, `mentor`, `admin`) на основе реального backend-кода.

Требование: ориентируйся на код backend как source of truth, а не только на OpenAPI.

## Источники правды

1. Код роутов и хендлеров (`internal/transport/http/handlers/**`).
2. Сервисная логика (`internal/service/**`).
3. DTO/mapper (`internal/dto/**`, `internal/mapper/mapper.go`).
4. Миграции (`migrations/*.sql`) для бизнес-ограничений.
5. OpenAPI (`docs/openapi.yaml`) используй как вспомогательный источник.

Важно: `docs/openapi.yaml` покрывает только часть API и не содержит много рабочих ручек (mentor/admin blocks).

## Базовые правила интеграции

1. Base URL: `http://localhost:8080`.
2. Все защищенные ручки под `/api/v1/**` (кроме `/api/v1/auth/**`) требуют JWT в заголовке `Authorization: Bearer <access_token>`.
3. Access token TTL по умолчанию: `24h`, refresh token TTL: `720h` (30 дней).
4. Формат дат: отправляй RFC3339/ISO8601 (`2026-04-17T10:00:00Z`).
5. Формат списков:
   ```json
   {
     "items": [],
     "total_count": 0,
     "limit": 20,
     "offset": 0
   }
   ```
6. Пагинация по умолчанию: `limit=20`, max `limit=100` (кроме `/api/v1/mentor/slots`, там фиксировано `limit=100`).

## Унифицированная обработка ошибок на фронте

Поддержи оба формата, потому что backend отдает и так, и так:

1. Основной:
   ```json
   { "error": { "code": "INVALID_REQUEST", "message": "..." } }
   ```
2. Middleware/role-check:
   ```json
   { "error": "forbidden" }
   ```

Типовые коды: `INVALID_REQUEST`, `UNAUTHORIZED`, `FORBIDDEN`, `NOT_FOUND`, `CONFLICT`, `INTERNAL_ERROR`.

## Роли и доступ

1. `student`: все общие защищенные ручки (`/api/v1/me`, `/cities`, `/slots`, `/bookings`, `/mentor-requests`).
2. `mentor`: все как у student + `/api/v1/mentor/**`.
3. `admin`: все как у student + `/api/v1/admin/**`.

Навигация и UI должны строиться по `role` из токена/профиля.

## Справочник enum-значений

1. `role`: `student`, `mentor`, `admin`.
2. `slot.type`: `online`, `offline`.
3. `slot.status`: `draft`, `pending`, `active`, `canceled`, `completed`, `deleted`.
4. `booking.type`: `room_only`, `room_with_mentor`, `mentor_call`, `event_seat`.
5. `booking.status`: `active`, `canceled`, `completed`, `deleted`.
6. `mentor_request.request_type`: `category`, `skill`, `direct_mentor`, `other`.
7. `mentor_request.status`: `pending`, `in_review`, `approved`, `rejected`, `canceled`, `expired`.

## Полный список ручек и функциональность

### System / Docs

| Метод | Путь | Auth | Функциональность | Ответ |
|---|---|---|---|---|
| GET | `/heathz` | Нет | Liveness check (опечатка в пути именно `heathz`) | `{"status":"ok"}` |
| GET | `/readyz` | Нет | Readiness + доступность БД | `{"status":"ok","database":"connected\|disconnected"}` |
| GET | `/swagger` | Нет | Swagger UI | HTML |
| GET | `/swagger/openapi.yaml` | Нет | OpenAPI спецификация | YAML |

### Auth

| Метод | Путь | Auth | Что делает | Request | Response | Важные правила |
|---|---|---|---|---|---|---|
| POST | `/api/v1/auth/register` | Нет | Регистрация пользователя | `email*`, `password* (min 8)`, `full_name*`, `birth_date?`, `university?`, `course?`, `degree_type?`, `role?` | `201` `{"user":UserResponse,"tokens":{"access_token","refresh_token"}}` | Если `role` пустой, backend ставит `student`. Student profile создается только если одновременно заданы `university`, `course>0`, `degree_type`. |
| POST | `/api/v1/auth/login` | Нет | Логин по email/password | `email*`, `password*` | `200` `{"access_token","refresh_token"}` | При неверных данных: `401 UNAUTHORIZED`. |
| POST | `/api/v1/auth/refresh` | Нет | Обновление пары токенов | `refresh_token*` | `200` новая пара токенов | Старый refresh отзывается (rotation). |
| POST | `/api/v1/auth/logout` | Нет | Logout (revoke refresh) | `refresh_token*` | `204` | Access token не blacklisted, только refresh revoke. |

### User / Catalog / Public Protected

| Метод | Путь | Role | Что делает | Вход | Выход | Важные правила |
|---|---|---|---|---|---|---|
| GET | `/api/v1/me` | Любая авторизованная | Профиль текущего пользователя | - | `UserResponse` | User берется из `user_id` claim токена. |
| PATCH | `/api/v1/me` | Любая авторизованная | Обновление профиля | `full_name?`, `birth_date?` | `UserResponse` | Частичный patch. |
| GET | `/api/v1/cities` | Любая авторизованная | Список городов | `limit?`, `offset?` | `ListResponse<CityResponse>` | default `limit=20`, max `100`. |
| GET | `/api/v1/cities/{cityId}/hubs` | Любая авторизованная | Хабы города | path `cityId`, `limit?`, `offset?` | `ListResponse<HubResponse>` | `cityId` должен быть UUID. |
| GET | `/api/v1/cities/{cityId}/rooms` | Любая авторизованная | Комнаты города | path `cityId`, `limit?`, `offset?` | `ListResponse<RoomResponse>` | `cityId` должен быть UUID. |
| GET | `/api/v1/slots` | Любая авторизованная | Поиск слотов | `city_id?`, `hub_id?`, `room_id?`, `mentor_id?`, `status?`, `type?`, `start_from?`, `end_to?`, `limit?`, `offset?` | `ListResponse<SlotResponse>` | Сортировка на backend фиксирована `start_at ASC`. Параметры `sort_by`/`order` есть в DTO, но backend их игнорирует. |
| GET | `/api/v1/slots/{id}` | Любая авторизованная | Слот по ID | path `id` | `SlotResponse` | `404 NOT_FOUND` если нет слота. |
| POST | `/api/v1/bookings` | Любая авторизованная | Создать бронь | `booking_type*`, `start_at*`, `end_at*`, `slot_id?`, `room_id?`, `meeting_url?`, `seat_number?`, `idempotency_key?` | `201 BookingResponse` | Проверки: `start<end`; город пользователя обязателен; при `slot_id` слот должен быть `active`; ограничение capacity для slot/room; при дублировании `idempotency_key` вернется существующая бронь (все равно статус ответа `201`). |
| DELETE | `/api/v1/bookings/{id}` | Любая авторизованная | Отмена своей брони | path `id` | `200 BookingResponse` | Отменяет только бронь текущего пользователя (`force=false`). |
| GET | `/api/v1/bookings` | Любая авторизованная | Список моих броней | `status?`, `city_id?`, `room_id?`, `slot_id?`, `mentor_id?`, `limit?`, `offset?` | `ListResponse<BookingResponse>` | Backend всегда фильтрует по текущему user. Сортировка фиксирована `start_at DESC`. `mentor_id` параметр сейчас не используется backend-логикой. |
| POST | `/api/v1/mentor-requests` | Любая авторизованная | Создать запрос на ментора | `request_type*`, `slot_id?`, `mentor_id?`, `skill_id?`, `comment?` | `201 MentorRequestResponse` | У текущего пользователя должен быть `city_id`, иначе `400`. Статус нового запроса всегда `pending`. |
| GET | `/api/v1/mentor-requests` | Любая авторизованная | Список mentor-запросов | `status?`, `skill_id?`, `city_id?`, `mentor_id?`, `limit?`, `offset?` | `ListResponse<MentorRequestResponse>` | Role-aware фильтрация: для `student` принудительно `mentee_id=current user`; для `mentor` принудительно `mentor_id=current user`; для `admin` фильтры как переданы. |

### Admin (`role=admin`)

| Метод | Путь | Что делает | Request | Response | Важные правила |
|---|---|---|---|---|---|
| POST | `/api/v1/admin/mentors/register_mentor` | Создать mentor-пользователя админом | `email*`, `password*`, `full_name*`, `city_id*`, `description?`, `title?` | `201 UserResponse` | Создает пользователя с ролью mentor + mentor profile. Токены не выдаются. |
| POST | `/api/v1/admin/cities` | Создать город | `name*`, `is_active?` | `201 CityResponse` | `is_active` default `true`. |
| PATCH | `/api/v1/admin/cities/{id}` | Обновить город | path `id`, `name?`, `is_active?` | `200 CityResponse` | `404` если city не найден. |
| POST | `/api/v1/admin/hubs` | Создать хаб | `city_id*`, `name*`, `address*`, `is_active?` | `201 HubResponse` | `is_active` default `true`. |
| PATCH | `/api/v1/admin/hubs/{id}` | Обновить хаб | path `id`, `name?`, `address?`, `is_active?` | `200 HubResponse` | `404` если hub не найден. |
| POST | `/api/v1/admin/rooms` | Создать комнату | `hub_id*`, `name*`, `capacity* (min 1)`, `description?`, `room_type?`, `is_active?` | `201 RoomResponse` | `is_active` default `true`. |
| PATCH | `/api/v1/admin/rooms/{id}` | Обновить комнату | path `id`, `name?`, `description?`, `room_type?`, `capacity?`, `is_active?` | `200 RoomResponse` | `404` если room не найден. |
| DELETE | `/api/v1/admin/rooms/{id}` | Удалить комнату (soft) | path `id` | `204` | Фактически `is_active=false`. |
| PATCH | `/api/v1/admin/users/{id}/city` | Привязать пользователя к городу | path `id`, body `city_id*` | `204` | Нужен валидный user и city. |
| GET | `/api/v1/admin/analytics/business` | Бизнес-метрики | - | `200` counters | Возвращает агрегаты: approvals/rejections/active_bookings/etc. |
| GET | `/api/v1/admin/analytics/technical` | Тех-метрики | - | `200` counters | Возвращает `booking_conflict_count`, `failed_outbox_events`, `outbox_lag_count`. |

### Mentor (`role=mentor`)

| Метод | Путь | Что делает | Request | Response | Важные правила |
|---|---|---|---|---|---|
| POST | `/api/v1/mentor/skills/subscribe` | Подписка ментора на skill | `skill_id* (uuid)` | `204` | Upsert: повторный запрос безопасен. |
| DELETE | `/api/v1/mentor/skills/subscribe/{skillId}` | Отписка от skill | path `skillId` | `204` | Делает подписку неактивной. |
| GET | `/api/v1/mentor/slots` | Слоты текущего ментора | - | `ListResponse<SlotResponse>` | Фильтр по текущему mentor; backend фиксирует `limit=100`, `offset=0`. |
| POST | `/api/v1/mentor/slots` | Создать слот ментора | `city_id*`, `type*`, `start_at*`, `end_at*`, `room_id?`, `meeting_url?`, `status?`, `capacity?` | `201 SlotResponse` | `status` default `active`; `start<end`; `offline` требует `room_id`; `online` требует отсутствия `room_id`; backend требует `meeting_url` в любом случае; проверяется конфликт по времени для ментора/комнаты. |
| PATCH | `/api/v1/mentor/slots/{id}` | Частичное обновление слота | path `id`, `room_id?`, `start_at?`, `end_at?`, `meeting_url?`, `status?`, `capacity?` | `200 SlotResponse` | Сейчас backend не проверяет ownership и почти не валидирует patch (ограничь это в UI). |
| DELETE | `/api/v1/mentor/slots/{id}` | Удалить слот (soft) | path `id` | `204` | Ставит статус `deleted`. |
| GET | `/api/v1/mentor/requests` | Список запросов ментора | query как у `/api/v1/mentor-requests` | `ListResponse<MentorRequestResponse>` | Для mentor всегда фильтр по `mentor_id=current user`. |
| POST | `/api/v1/mentor/requests/{id}/approve` | Одобрить mentor request | path `id` | `200 MentorRequestResponse` | Ставит `status=approved`, проставляет `mentor_id=current user`. |
| POST | `/api/v1/mentor/requests/{id}/reject` | Отклонить mentor request | path `id` | `200 MentorRequestResponse` | Ставит `status=rejected`, проставляет `mentor_id=current user`. |

## Форматы response-моделей (кратко)

1. `UserResponse`: `id`, `full_name`, `birth_date?`, `email`, `role`, `city_id?`, `is_active`, `created_at`, `updated_at`.
2. `CityResponse`: `id`, `name`, `is_active`.
3. `HubResponse`: `id`, `city_id`, `name`, `address`, `is_active`.
4. `RoomResponse`: `id`, `hub_id`, `name`, `description?`, `room_type?`, `capacity`, `is_active`.
5. `SlotResponse`: `id`, `mentor_id`, `room_id?`, `city_id`, `type`, `start_at`, `end_at`, `meeting_url?`, `status`, `capacity?`.
6. `BookingResponse`: `id`, `slot_id?`, `room_id?`, `user_id`, `booking_type`, `status`, `start_at`, `end_at`, `meeting_url?`, `seat_number?`, `created_at`, `updated_at`.
7. `MentorRequestResponse`: `id`, `slot_id?`, `mentor_id?`, `mentee_id`, `city_id`, `skill_id?`, `request_type`, `status`, `comment?`.

## Что реализовать во фронтенде

1. Auth-модуль: register/login/logout/refresh + хранение токенов + автоподновление access token.
2. Role-based routing:
   - student: профиль, каталог, слоты, бронирования, mentor requests.
   - mentor: всё student + mentor skills/slots/requests.
   - admin: всё student + управление city/hub/room/users + analytics.
3. Единый API-client:
   - общий обработчик ошибок (оба формата error),
   - автоматическая подстановка `Authorization`,
   - retry через `/auth/refresh` при `401` (однократно).
4. UX-валидации на клиенте до запроса:
   - корректные UUID,
   - `start_at < end_at`,
   - минимальные обязательные поля по каждой форме,
   - в бронировании требуй хотя бы `slot_id` или `room_id`, даже если backend технически допускает отсутствие обоих.
5. Таблицы/фильтры со state URL query для list-ручек.
6. Пагинация и статусы загрузки/пустых состояний/ошибок для каждой страницы.

## Известные backend-нюансы (учесть в UI/клиенте)

1. OpenAPI не полон и не отражает все реальные маршруты.
2. `sort_by` и `order` в некоторых list DTO сейчас не используются backend.
3. `mentor_id` query в `/api/v1/bookings` сейчас не применяется.
4. Для mentor slot create backend требует `meeting_url`, даже если DTO делает его optional.
5. Проверка владения при update/delete mentor slot на backend нестрогая: не давай пользователю редактировать чужие записи в UI.
6. `/heathz` написан с опечаткой и должен использоваться именно в таком виде.

## Definition of Done

1. Все перечисленные ручки интегрированы.
2. Реализованы экраны по ролям и скрытие недоступных разделов.
3. Ошибки и edge-cases обработаны (включая `409 CONFLICT` сценарии бронирований/слотов).
4. Реализован стабильный auth flow с refresh токеном.
5. Все даты отображаются в локальном часовом поясе пользователя, в API отправляются в RFC3339 UTC.

