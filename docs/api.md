
------------------------------------------------------------
aka API_DOCUMENTATION by Arlandaren/infoowner
------------------------------------------------------------

# Документация API

Это руководство описывает доступные эндпоинты, их параметры, методы HTTP, возможные коды ответа и примеры запросов/ответов.

## Содержимое
- [Общее API (Api)](#общее-api-api)
    - [GET /metrics](#get-metrics)
    - [GET /ping](#get-ping)
    - [POST /time/advance](#post-timeadvance)
    - [GET /content/moderate](#get-contentmoderate)
    - [GET /content/propose](#get-contentpropose)
    - [GET /content/file/upload](#get-contentfileupload)
- [Статистика (Statistics)](#статистика-statistics)
    - [GET /stats/campaigns/:campaign_id](#get-statscampaignscampaign_id)
    - [GET /stats/advertisers/:advertiser_id/campaigns](#get-statsadvertisersadvertiser_idcampaigns)
    - [GET /stats/campaigns/:campaign_id/daily](#get-statscampaignsdaily)
    - [GET /stats/advertisers/:advertiser_id/campaigns/daily](#get-statsadvertisersadvertiser_idcampaignsdaily)
- [Кампании (Campaigns)](#кампании-campaigns)
    - [POST /advertisers/:advertiser_id/campaigns](#post-advertisersadvertiser_idcampaigns)
    - [GET /advertisers/:advertiser_id/campaigns](#get-advertisersadvertiser_idcampaigns)
    - [GET /advertisers/:advertiser_id/campaigns/:campaign_id](#get-advertisersadvertiser_idcampaignscampaign_id)
    - [PUT /advertisers/:advertiser_id/campaigns/:campaign_id](#put-advertisersadvertiser_idcampaignscampaign_id)
    - [DELETE /advertisers/:advertiser_id/campaigns/:campaign_id](#delete-advertisersadvertiser_idcampaignscampaign_id)
- [Рекламодатели (Advertisers)](#рекламодатели-advertisers)
    - [GET /advertisers/:advertiser_id](#get-advertisersadvertiser_id)
    - [POST /advertisers/bulk](#post-advertisersbulk)
    - [POST /ml-scores](#post-ml-scores)
- [Объявления (Ads)](#объявления-ads)
    - [GET /ads](#get-ads)
    - [POST /ads/:ad_id/click](#post-adsad_idclick)
- [Клиенты (Clients)](#клиенты-clients)
    - [GET /clients/:client_id](#get-clientsclient_id)
    - [POST /clients/bulk](#post-clientsbulk)

---

![Схема архитектуры](arch.png)

---
## Общее API (Api)

Эндпоинты, предоставляющие специальные сервисы: метрики, пинг, управление датой, модерация и генерация контента, загрузка изоображений.

### GET /metrics

**Описание:**
Возвращает метрики Prometheus.

**Пример запроса:**
```
GET /metrics
```

**Ответ:**
Текстовый формат метрик для Prometheus.

---

### GET /ping

**Описание:**
Проверка работоспособности сервиса.

**Пример запроса:**
```
GET /ping
```

**Ответ (200):**
```json
{
  "message": "pong"
}
```

---

### POST /time/advance

**Описание:**
Сдвигает системную/рабочую дату (для тестирования или моделирования).

**Тело запроса (JSON):**
Пример:
```json
{
  "advance_by": 1,
  "unit": "day"
}
```

**Ответ (200):**
Возвращается отправленное тело запроса.
```json
{
  "advance_by": 1,
  "unit": "day"
}
```

**Коды ошибок:**
- 400 – Ошибка валидации или запроса.

---

### GET /content/moderate

**Описание:**
Модерирует переданный текст на наличие неподобающего содержания.

**Тело запроса (JSON):**
Пример:
```json
{
  "text": "Некоторый текст для проверки."
}
```

**Пример ответа (200):**
```json
{
  "result": true
}
```
Значение result:
- true – если текст не содержит запрещённого содержания (ответ "нет" с точки зрения модели).
- false – если текст содержит неподобающее содержание (ответ "да").

**Коды ошибок:**
- 400 – Неверный формат запроса.
- 500 – Модель недоступна/внутренняя ошибка.

---

### GET /content/propose

**Описание:**
Генерирует продающий вариант рекламного текста на основе названия рекламодателя и заголовка.

**Тело запроса (JSON):**
Пример:
```json
{
  "advertiser": "Компания А",
  "title": "Лучшие цены всегда"
}
```

**Ответ (200):**
```json
{
  "result": "Лучшее предложение от Компании А по супер цене!"
}
```

**Коды ошибок:**
- 400 – Неверный формат запроса.
- 500 – Модель недоступна/ошибка генерации.

---

### GET /content/file/upload

**Описание:**
Загружает пользовательские изоображения на сервер.

**Тело запроса (form-data):**
Пример:
```file
{
  "file": .gif/.png/.jpg
}
```

**Ответ (200):**
```json
{
  "Link": "http://minio:9000/images/R2.png"
}
```

**Коды ошибок:**
- 400 – Неверный формат запроса.
- 500 – Ошибка загрузки на сервер.

---

## Статистика (Statistics)

Эндпоинты статистики возвращают сводную информацию по кампаниям и рекламодателям.

### GET /stats/campaigns/:campaign_id

**Описание:**
Возвращает сводную статистику по конкретной кампании.

**Параметры URL:**
- `campaign_id` – UUID кампании.

**Пример запроса:**
```
GET /stats/campaigns/1e7d1c88-4c1b-4e7f-97d0-1234567890ab
```

**Ответ (200):**
```json
{
  "impressions_count": 1000,
  "clicks_count": 50,
  "conversion": 0.05,
  "spent_impressions": 950,
  "spent_clicks": 47,
  "spent_total": 200.0
}
```

**Коды ошибок:**
- 400 – Неверный формат campaign_id.
- 404 – Кампания не найдена.
- 500 – Внутренняя ошибка сервера.

---

### GET /stats/advertisers/:advertiser_id/campaigns

**Описание:**
Возвращает сводную статистику по всем кампаниям, принадлежащим рекламодателю.

**Параметры URL:**
- `advertiser_id` – UUID рекламодателя.

**Пример запроса:**
```
GET /stats/advertisers/3f1d9d22-7f37-4c2d-abc3-1234567890cd/campaigns
```

**Ответ (200):**
```json
{
  "impressions_count": 5000,
  "clicks_count": 250,
  "conversion": 0.05,
  "spent_impressions": 4800,
  "spent_clicks": 240,
  "spent_total": 1000.0
}
```

**Коды ошибок:**
- 400 – Неверный формат advertiser_id.
- 404 – Кампания(ы) не найдены.
- 500 – Внутренняя ошибка сервера.

---

### GET /stats/campaigns/:campaign_id/daily

**Описание:**
Возвращает ежедневную статистику по указанной кампании.

**Параметры URL:**
- `campaign_id` – UUID кампании.

**Пример запроса:**
```
GET /stats/campaigns/1e7d1c88-4c1b-4e7f-97d0-1234567890ab/daily
```

**Пример ответа (200):**
```json
[
  {
    "date": "2023-10-01",
    "impressions_count": 200,
    "clicks_count": 10,
    "conversion": 0.05,
    "spent_impressions": 190,
    "spent_clicks": 9,
    "spent_total": 40.0
  },
  {
    "date": "2023-10-02",
    "impressions_count": 300,
    "clicks_count": 15,
    "conversion": 0.05,
    "spent_impressions": 290,
    "spent_clicks": 14,
    "spent_total": 60.0
  }
]
```

**Коды ошибок:**
- 400 – Неверный формат campaign_id.
- 404 – Кампания не найдена.
- 500 – Внутренняя ошибка сервера.

---

### GET /stats/advertisers/:advertiser_id/campaigns/daily

**Описание:**
Возвращает ежедневную статистику по всем кампаниям рекламодателя.

**Параметры URL:**
- `advertiser_id` – UUID рекламодателя.

**Пример запроса:**
```
GET /stats/advertisers/3f1d9d22-7f37-4c2d-abc3-1234567890cd/campaigns/daily
```

**Пример ответа (200):**
```json
[
  {
    "date": "2023-10-01",
    "impressions_count": 500,
    "clicks_count": 25,
    "conversion": 0.05,
    "spent_impressions": 480,
    "spent_clicks": 24,
    "spent_total": 100.0
  }
]
```

**Коды ошибок:**
- 400 – Неверный формат advertiser_id.
- 404 – Кампания не найдена.
- 400 – Ошибка запроса.
- 500 – Внутренняя ошибка сервера.

---

## Кампании (Campaigns)

Эндпоинты для управления кампаниями. Все запросы выполняются внутри группы URL:
`/advertisers/:advertiser_id/campaigns`

### POST /advertisers/:advertiser_id/campaigns

**Описание:**
Создание новой кампании.

**Параметры URL:**
- `advertiser_id` – UUID рекламодателя.

**Тело запроса (JSON):**
Пример структуры (CampaignCreate):
```json
{
  "name": "Новая кампания",
  "budget": 1000,
  "start_date": "2023-10-05",
  "end_date": "2023-11-05",
  "...": "другие необходимые поля"
}
```

**Ответ (201):**
Ответ содержит данные созданной кампании:
```json
{
  "campaign_id": "1e7d1c88-4c1b-4e7f-97d0-1234567890ab",
  "advertiser_id": "3f1d9d22-7f37-4c2d-abc3-1234567890cd",
  "name": "Новая кампания",
  "budget": 1000,
  "start_date": "2023-10-05",
  "end_date": "2023-11-05"
  // ...
}
```

**Коды ошибок:**
- 400 – Неверный формат advertiser_id, ошибка валидации JSON или невалидные данные.
- 400 – Ошибки БД.

---

### GET /advertisers/:advertiser_id/campaigns

**Описание:**
Возвращает список кампаний для рекламодателя с поддержкой пагинации.

**Параметры URL:**
- `advertiser_id` – UUID рекламодателя.

**Параметры запроса:**
- `size` – Опционально. Количество кампаний на странице (по умолчанию 10).
- `page` – Опционально. Номер страницы (по умолчанию 1).

**Пример запроса:**
```
GET /advertisers/3f1d9d22-7f37-4c2d-abc3-1234567890cd/campaigns?size=10&page=1
```

**Ответ (200):**
```json
[
  {
    "campaign_id": "1e7d1c88-4c1b-4e7f-97d0-1234567890ab",
    "advertiser_id": "3f1d9d22-7f37-4c2d-abc3-1234567890cd",
    "name": "Кампания 1",
    // ...
  },
  {
    "campaign_id": "2f8e2d99-5d2e-4e8f-97f1-0987654321fe",
    "advertiser_id": "3f1d9d22-7f37-4c2d-abc3-1234567890cd",
    "name": "Кампания 2",
    // ...
  }
]
```

**Коды ошибок:**
- 400 – Неверный формат advertiser_id или неверные параметры пагинации.
- 400 – Ошибка запроса.

---

### GET /advertisers/:advertiser_id/campaigns/:campaign_id

**Описание:**
Возвращает информацию о конкретной кампании рекламодателя.

**Параметры URL:**
- `advertiser_id` – UUID рекламодателя.
- `campaign_id` – UUID кампании.

**Пример запроса:**
```
GET /advertisers/3f1d9d22-7f37-4c2d-abc3-1234567890cd/campaigns/1e7d1c88-4c1b-4e7f-97d0-1234567890ab
```

**Ответ (200):**
```json
{
  "campaign_id": "1e7d1c88-4c1b-4e7f-97d0-1234567890ab",
  "advertiser_id": "3f1d9d22-7f37-4c2d-abc3-1234567890cd",
  "name": "Кампания 1",
  // ...
}
```

**Коды ошибок:**
- 400 – Неверный формат advertiser_id или campaign_id.
- 404 – Кампания не найдена.
- 403 – Кампания не принадлежит указанному рекламодателю.

---

### PUT /advertisers/:advertiser_id/campaigns/:campaign_id

**Описание:**
Обновление данных кампании.

**Параметры URL:**
- `advertiser_id` – UUID рекламодателя.
- `campaign_id` – UUID кампании.

**Тело запроса (JSON):**
Пример структуры (CampaignUpdate):
```json
{
  "name": "Обновленная кампания",
  "budget": 1200,
  "...": "другие изменяемые поля"
}
```

**Ответ (200):**
```json
{
  "campaign_id": "1e7d1c88-4c1b-4e7f-97d0-1234567890ab",
  "advertiser_id": "3f1d9d22-7f37-4c2d-abc3-1234567890cd",
  "name": "Обновленная кампания",
  "budget": 1200
  // ...
}
```

**Коды ошибок:**
- 400 – Неверный формат advertiser_id или campaign_id, невалидные данные.
- 404 – Кампания не найдена.
- 403 – Кампания не принадлежит рекламодателю.
- 400/403 – Отказ действия на уровне бизнес-логики.

---

### DELETE /advertisers/:advertiser_id/campaigns/:campaign_id

**Описание:**
Удаление кампании.

**Параметры URL:**
- `advertiser_id` – UUID рекламодателя.
- `campaign_id` – UUID кампании.

**Пример запроса:**
```
DELETE /advertisers/3f1d9d22-7f37-4c2d-abc3-1234567890cd/campaigns/1e7d1c88-4c1b-4e7f-97d0-1234567890ab
```

**Ответ:**
Код: 204 No Content

**Коды ошибок:**
- 400 – Неверный формат advertiser_id или campaign_id.
- 404 – Кампания не найдена.
- 403 – Кампания не принадлежит рекламодателю.
- 400 – Ошибка запроса.

---

## Рекламодатели (Advertisers)

Эндпоинты для получения и массового обновления данных рекламодателей, а также загрузки ML-оценок.

### GET /advertisers/:advertiser_id

**Описание:**
Получение данных рекламодателя по ID.

**Параметры URL:**
- `advertiser_id` – UUID рекламодателя.

**Пример запроса:**
```
GET /advertisers/3f1d9d22-7f37-4c2d-abc3-1234567890cd
```

**Ответ (200):**
```json
{
  "advertiser_id": "3f1d9d22-7f37-4c2d-abc3-1234567890cd",
  "name": "Рекламодатель 1",
  // ...
}
```

**Коды ошибок:**
- 400 – Неверный формат advertiser_id.
- 404 – Рекламодатель не найден.

---

### POST /advertisers/bulk

**Описание:**
Массовое создание/обновление рекламодателей.

**Тело запроса (JSON):**
Пример:
```json
[
  {
    "advertiser_id": "3f1d9d22-7f37-4c2d-abc3-1234567890cd",
    "name": "Рекламодатель 1",
    "...": "другие поля"
  },
  {
    "advertiser_id": "4a2e9f33-8e47-4d3c-bcd4-0987654321ef",
    "name": "Рекламодатель 2",
    "...": "другие поля"
  }
]
```

**Ответ (201):**
Возвращается тело запроса с созданными/обновленными записями.

**Коды ошибок:**
- 400 – Ошибка валидации JSON или отсутствует advertiser_id.
- 400 – Ошибка запроса.

---

### POST /ml-scores

**Описание:**
Загрузка/обновление ML-оценки для рекламодателя.

**Тело запроса (JSON):**
Пример:
```json
{
  "client_id": "8cfa1b44-9eaa-4d7f-abbc-1234567890aa",
  "advertiser_id": "3f1d9d22-7f37-4c2d-abc3-1234567890cd",
  "score": 0.87
}
```

**Ответ:**
Код 200 OK при успешном обновлении.

**Коды ошибок:**
- 400 – Неверный формат client_id или advertiser_id, либо ошибка в теле запроса.

---

## Объявления (Ads)

Эндпоинты для получения объявления для клиента и регистрации кликов по объявлению.

### GET /ads

**Описание:**
Возвращает объявление для указанного клиента.

**Параметры запроса:**
- `client_id` – UUID клиента (передается как query параметр).

**Пример запроса:**
```
GET /ads?client_id=8cfa1b44-9eaa-4d7f-abbc-1234567890aa
```

**Ответ (200):**
```json
{
  "ad_id": "a1b2c3d4-1234-5678-9101-abcdefabcdef",
  "title": "Лучшее предложение",
  "content": "Описание объявления",
  // ...
}
```

**Коды ошибок:**
- 400 – Неверный формат client_id.
- 404 – Объявление не найдено.
- 500 – Внутренняя ошибка сервера.

---

### POST /ads/:ad_id/click

**Описание:**
Регистрирует клик по объявлению.

**Параметры URL:**
- `ad_id` – UUID объявления.

**Тело запроса (JSON):**
```json
{
  "client_id": "8cfa1b44-9eaa-4d7f-abbc-1234567890aa"
}
```

**Ответ:**
Код 204 No Content при успешной регистрации клика.

**Коды ошибок:**
- 400 – Неверный формат ad_id или client_id, либо невалидное тело запроса.
- 404 – Объявление не найдено.
- 400 – Прочие ошибки (например, бизнес-логика).

---

## Клиенты (Clients)

Эндпоинты для работы с клиентами.

### GET /clients/:client_id

**Описание:**
Возвращает данные клиента по его ID.

**Параметры URL:**
- `client_id` – UUID клиента.

**Пример запроса:**
```
GET /clients/8cfa1b44-9eaa-4d7f-abbc-1234567890aa
```

**Ответ (200):**
```json
{
  "client_id": "8cfa1b44-9eaa-4d7f-abbc-1234567890aa",
  "name": "Клиент 1",
  // ...
}
```

**Коды ошибок:**
- 400 – Неверный формат client_id.
- 404 – Клиент не найден.

---

### POST /clients/bulk

**Описание:**
Массовое создание/обновление клиентов.

**Тело запроса (JSON):**
Пример:
```json
[
  {
    "client_id": "8cfa1b44-9eaa-4d7f-abbc-1234567890aa",
    "name": "Клиент 1",
    "...": "другие поля"
  },
  {
    "client_id": "9dbf2c55-0ebb-4fc7-9a7d-0987654321bb",
    "name": "Клиент 2",
    "...": "другие поля"
  }
]
```

**Ответ (201):**
Возвращается тело запроса с созданными/обновленными клиентами.

**Коды ошибок:**
- 400 – Неверный формат JSON или отсутствует client_id.
- 400 – Ошибка запроса.

------------------------------------------------------------
Пояснения:

- Все UUID должны передаваться в корректном формате.
- При массовых операциях (bulk upsert) сервер возвращает отправленный запрос в случае успеха.
- Обработка ошибок производится согласно внутренней логике сервиса: 400 — ошибка валидации/запроса, 403 — отказ по правам, 404 — не найден, 500 — внутренняя ошибка сервера.

Это полная документация по перечисленным эндпоинтам.

# Postman коллекция

Для подробностей [клик](./collection.json).