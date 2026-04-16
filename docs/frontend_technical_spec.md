# ТЗ для Frontend (дополнение к backend ТЗ “Студент и Т”)

Документ описывает требования к фронтенду для интеграции с уже существующим backend-проектом.

## 1. Цель

Реализовать web-клиент платформы “Студент и Т” с поддержкой ролей `student`, `mentor`, `admin`, интеграцией с:

1. `Main Service` (основные бизнес-функции).
2. `Notification Service` (уведомления, unread/read, realtime).

## 2. Границы frontend-реализации

В scope фронтенда:

1. Полный UX для сценариев из backend ТЗ.
2. Интеграция с REST API обеих сервисов.
3. Realtime уведомления через WS/SSE.
4. Ролевой доступ к разделам и действиям.
5. Валидация пользовательского ввода.
6. Пагинация, фильтрация, сортировка, состояния загрузки/ошибок.

Out of scope:

1. Реализация Notification Service на backend.
2. Изменение существующей архитектуры Main Service.

## 3. Архитектурные требования к frontend

1. Разделить API-клиенты:
   - `mainApiClient` для `Main Service`.
   - `notificationApiClient` для `Notification Service`.
2. Единый auth-модуль:
   - login/register/refresh/logout.
   - автоподстановка `Authorization: Bearer <token>`.
   - retry один раз через refresh на `401`.
3. Единый error-handler:
   - поддержка форматов:
     - `{ "error": { "code": "...", "message": "..." } }`
     - `{ "error": "forbidden" }`
4. Единая модель пагинации:
   - `limit`, `offset`, `total_count`.
5. Typed контракты:
   - отдельные типы request/response/filter/list для всех сущностей.

## 4. Роли и матрица доступа в UI

1. `student`:
   - профиль, каталог город/хаб/комнаты;
   - слоты;
   - бронирования;
   - mentor requests;
   - уведомления.
2. `mentor`:
   - всё `student` +
   - управление своими слотами;
   - подписки на skills;
   - обработка mentor requests.
3. `admin`:
   - всё `student` +
   - админка городов/хабов/комнат;
   - регистрация ментора;
   - смена города пользователя;
   - бизнес/тех аналитика.

Frontend обязан скрывать недоступные разделы по роли и не показывать запрещенные действия.

## 5. Основные экраны

1. Auth:
   - регистрация;
   - логин;
   - выход.
2. Общие:
   - `Мой профиль` (`GET/PATCH /api/v1/me`);
   - `Города` → `Хабы` → `Комнаты`;
   - `Слоты` (поиск/фильтрация);
   - `Мои бронирования`;
   - `Mentor Requests`;
   - `Уведомления`.
3. Mentor:
   - `Мои слоты` (CRUD);
   - `Подписки на skills`;
   - `Запросы ко мне` (approve/reject).
4. Admin:
   - `Города` (create/update);
   - `Хабы` (create/update);
   - `Комнаты` (create/update/delete);
   - `Регистрация ментора`;
   - `Смена города пользователя`;
   - `Аналитика`.

## 6. API интеграция (frontend view)

### 6.1 Main Service (обязательный минимум)

1. Auth:
   - `POST /api/v1/auth/register`
   - `POST /api/v1/auth/login`
   - `POST /api/v1/auth/refresh`
   - `POST /api/v1/auth/logout`
2. User/Public:
   - `GET /api/v1/me`
   - `PATCH /api/v1/me`
   - `GET /api/v1/cities`
   - `GET /api/v1/cities/{cityId}/hubs`
   - `GET /api/v1/cities/{cityId}/rooms`
   - `GET /api/v1/slots`
   - `GET /api/v1/slots/{id}`
   - `POST /api/v1/bookings`
   - `DELETE /api/v1/bookings/{id}`
   - `GET /api/v1/bookings`
   - `POST /api/v1/mentor-requests`
   - `GET /api/v1/mentor-requests`
3. Mentor:
   - `POST /api/v1/mentor/skills/subscribe`
   - `DELETE /api/v1/mentor/skills/subscribe/{skillId}`
   - `GET /api/v1/mentor/slots`
   - `POST /api/v1/mentor/slots`
   - `PATCH /api/v1/mentor/slots/{id}`
   - `DELETE /api/v1/mentor/slots/{id}`
   - `GET /api/v1/mentor/requests`
   - `POST /api/v1/mentor/requests/{id}/approve`
   - `POST /api/v1/mentor/requests/{id}/reject`
4. Admin:
   - `POST /api/v1/admin/mentors/register_mentor`
   - `POST /api/v1/admin/cities`
   - `PATCH /api/v1/admin/cities/{id}`
   - `POST /api/v1/admin/hubs`
   - `PATCH /api/v1/admin/hubs/{id}`
   - `POST /api/v1/admin/rooms`
   - `PATCH /api/v1/admin/rooms/{id}`
   - `DELETE /api/v1/admin/rooms/{id}`
   - `PATCH /api/v1/admin/users/{id}/city`
   - `GET /api/v1/admin/analytics/business`
   - `GET /api/v1/admin/analytics/technical`

### 6.2 Notification Service (frontend обязан интегрировать)

1. `GET /api/v1/notifications`
2. `GET /api/v1/notifications/unread`
3. `PATCH /api/v1/notifications/{id}/read`
4. `PATCH /api/v1/notifications/read-all`
5. WS/SSE endpoint для realtime доставки.

## 7. Валидации на клиенте

1. UUID-поля должны валидироваться до отправки.
2. Email и пароль по правилам backend.
3. Все `start_at/end_at` с правилом `start_at < end_at`.
4. Для `offline slot`: обязателен `room_id`.
5. Для `online slot`: `room_id` отсутствует.
6. Для slot в UI обязательно требовать `meeting_url`.
7. Для бронирования валидировать обязательные поля по типу брони.
8. Для mentor request валидировать `request_type`.

## 8. Обработка ошибок и UX-статусов

1. `400 INVALID_REQUEST`: показать сообщение поля/формы.
2. `401 UNAUTHORIZED`: пробовать refresh, при провале logout + redirect на login.
3. `403 FORBIDDEN`: экран/тост “Недостаточно прав”.
4. `404 NOT_FOUND`: экран “Сущность не найдена”.
5. `409 CONFLICT`: отдельные дружелюбные тексты для конфликтов бронирования/слотов.
6. `500 INTERNAL_ERROR`: безопасное сообщение + retry action.

Обязательно:

1. Loading state для всех network-экранов.
2. Empty state для всех list-экранов.
3. Retry-кнопка при падении запроса.

## 9. Фильтрация, сортировка, пагинация

Для list-экранов поддержать:

1. фильтры по `city/hub/room/mentor/skill/category/status/type/date-range`;
2. `limit/offset`;
3. сохранение фильтров в URL query;
4. восстановление состояния фильтров при перезагрузке страницы.

Если backend в некоторых ручках игнорирует сортировку, UI должен:

1. не обещать пользователю недоступную сортировку;
2. показывать только поддерживаемые backend режимы.

## 10. Уведомления и realtime

1. Реализовать центр уведомлений:
   - список;
   - счетчик unread;
   - mark as read;
   - mark all as read.
2. Realtime:
   - подключение к WS/SSE после успешного логина;
   - reconnection policy;
   - dedup новых уведомлений по `id`;
   - обновление unread-бейджа без перезагрузки.
3. Уведомления типа alert выделять визуально.

## 11. Набор пользовательских сценариев (acceptance)

1. Student:
   - регистрируется/логинится;
   - видит ресурсы своего города;
   - создает и отменяет бронь;
   - создает mentor request;
   - получает уведомление о статусе.
2. Mentor:
   - подписывается на skill;
   - получает уведомление о подходящем request;
   - принимает/отклоняет request;
   - управляет своими слотами.
3. Admin:
   - создает город/хаб/комнату;
   - регистрирует ментора;
   - меняет город пользователю;
   - видит аналитику.

## 12. Тестирование frontend

Минимум:

1. Unit-тесты:
   - маппинг DTO в view-model;
   - форматтеры дат и валидаторы;
   - auth/refresh flow helper.
2. Integration/API tests:
   - обработка `401 -> refresh -> retry`;
   - обработка `409 conflict`.
3. E2E:
   - student booking happy-path;
   - mentor approve/reject flow;
   - admin CRUD и аналитика;
   - уведомления unread/read и realtime.

## 13. Нефункциональные требования

1. Безопасное хранение токенов согласно стандарту проекта.
2. Логи клиентских ошибок (без утечки персональных данных).
3. Доступность интерфейса: базовые требования a11y.
4. Производительность:
   - lazy-loading списков;
   - батч-обновление уведомлений;
   - оптимизация повторных запросов.

## 14. Definition of Done (frontend)

1. Все обязательные экраны и API-интеграции реализованы.
2. Ролевой доступ и UI-ограничения работают корректно.
3. Уведомления работают через Notification Service + realtime.
4. Конфликтные сценарии бронирований корректно отображаются пользователю.
5. Добавлены автотесты на ключевые пользовательские сценарии.
6. Документация по env и запуску frontend обновлена.

