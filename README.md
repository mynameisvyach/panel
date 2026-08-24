# 3x-ui: Happ fork patch

Комплект файлов для форка актуальной ветки `main` проекта MHSanaei/3x-ui.

## Что добавлено

1. Необязательное поле **«Описание сервера (Happ)»** в окне:
   `Подключения → Изменить подключение → Основное`.
2. Если описание пустое, `serverDescription` не добавляется.
3. Если описание заполнено, оно обрезается по краям, кодируется в Base64 и добавляется:
   - в VLESS, Trojan, Shadowsocks и Hysteria как query-параметр `serverDescription`;
   - в `vmess://` JSON как `meta.serverDescription`.
4. При открытии обычной HTTPS-ссылки подписки в браузере стандартная страница через 200 мс открывает:
   `happ://add/<исходная HTTPS-ссылка>`.
5. Запрос самого Happ по HTTPS-ссылке продолжает получать обычное тело подписки, поэтому цикл редиректов не возникает.
6. Во вкладку `Настройки панели → Подписка → Happ` добавлены отдельные настройки расширенного объявления, окончания подписки, скрытия серверных настроек, обновления, закрепления, фрагментации, шумов, ping и сортировки.

Включённые параметры добавляются директивами `#...` только в обычное тело `/sub/...` перед Base64-кодированием. Отдельные `/json/...` и Clash-подписки не изменялись.

## Куда копировать

Скопируйте всё содержимое этого архива, кроме `README.md`, в корень своего форка 3x-ui с сохранением структуры каталогов и заменой файлов.

Например:

```text
3x-ui-happ-fork-files/
├── frontend/
│   ├── public/openapi.json
│   └── src/...
└── internal/...
```

Каталоги `frontend` и `internal` должны попасть прямо в корень репозитория 3x-ui.

## Основные изменённые файлы

```text
frontend/src/pages/sub/SubPage.tsx
frontend/src/pages/settings/HappOptionsForm.tsx
frontend/src/pages/settings/SubscriptionGeneralTab.tsx
frontend/src/pages/inbounds/form/InboundFormModal.tsx
frontend/src/schemas/forms/inbound-form.ts
frontend/src/lib/xray/inbound-form-adapter.ts
frontend/src/models/dbinbound.ts
internal/database/model/model.go
internal/sub/service.go
internal/sub/happ_server_description_test.go
internal/sub/happ_options.go
internal/sub/happ_options_test.go
internal/web/translation/ru-RU.json
internal/web/translation/en-US.json
```

Также включены обновлённые генерируемые TypeScript/OpenAPI-файлы.

## Сборка

Требования актуальной ветки на момент подготовки комплекта:

- Go согласно `go.mod` проекта;
- Node.js 24 согласно `.nvmrc`;
- npm.

Из корня форка:

```bash
npm --prefix frontend ci
make build
```

Для полной проверки перед публикацией:

```bash
make verify
```

## Обновление базы

Отдельную SQL-миграцию выполнять не требуется. Существующий `AutoMigrate` при запуске добавит колонку `server_description` в таблицу inbound.

Перед заменой рабочего бинарника рекомендуется сохранить резервную копию базы данных панели.

## Проверка

### Описание сервера

1. Откройте `Подключения`.
2. Выберите `Изменить подключение`.
3. Во вкладке `Основное` заполните `Описание сервера (Happ)`.
4. Сохраните подключение.
5. Обновите обычную подписку в Happ.

Если поле оставить пустым, ссылка останется без `serverDescription`.

### Редирект

Откройте обычную ссылку вида:

```text
https://example.com/sub/ID
```

в браузере. Стандартная страница автоматически попытается открыть:

```text
happ://add/https://example.com/sub/ID
```

Если браузер блокирует автоматический запуск внешнего приложения, на стандартной странице остаётся штатная кнопка Happ для ручного открытия.
