# Frontend (React + TypeScript + FSD)

## Стек

1. React 18
2. TypeScript
3. Vite
4. React Router
5. TanStack Query

## Структура

Использована FSD-структура:

1. `app` — провайдеры и роутинг
2. `shared` — API, утилиты, базовые UI-компоненты
3. `entities` — типы доменных сущностей
4. `features` — пользовательские фичи (auth, notifications, forms)
5. `widgets` — layout-уровень (sidebar/topbar/shell)
6. `pages` — страницы

## Запуск

1. Установить зависимости:
   - `npm install`
2. Создать `.env`:
   - `cp .env.example .env`
3. Запустить dev сервер:
   - `npm run dev`

## Переменные окружения

1. `VITE_MAIN_API_URL` — URL Main Service
2. `VITE_NOTIFICATION_API_URL` — URL Notification Service
3. `VITE_NOTIFICATION_REALTIME_URL` — URL realtime endpoint notification service

