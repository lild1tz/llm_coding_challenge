
# 🚜 AgroAI: Hermes + Apollo + Superset + LangSmith

**AgroAI** — комплексное AI-решение для агропромышленности, разработанное в рамках **LLM Coding Challenge** для компании **АгроПрогресс**. Система включает два микросервиса — **Hermes** (на Go) и **Apollo** (на Python/FastAPI), а также полностью интегрирована с:

- **Apache Superset** — для визуальной аналитики результатов,
- **LangSmith** — для мониторинга и трассировки LLM-запросов и токенов.
---
## 🎯 Проблема

- Агрономы часто передают информацию о проделанной работе в виде свободных текстовых сообщений или фотографий таблиц. Ручная обработка таких данных занимает много времени, подвержена ошибкам и затрудняет оперативный анализ. Apollo решает эту проблему, автоматизируя процесс извлечения и структурирования ключевой информации.
---
## ✨ Возможности AgroAI

### 🛰 Hermes (мультиканальный шлюз)
- 📥 Приём сообщений агрономов из **Telegram** и **WhatsApp** (текст, фото, аудио)
- 🧹 Фильтрация и маршрутизация сообщений
- 🔁 Асинхронная передача данных в **Apollo**
- 💾 Сохранение оригиналов и результатов в **PostgreSQL**
- ☁️ Загрузка медиафайлов в **Google Drive**
- ⚙️ Работа с доверенными чатами (whitelist ID)

### 🤖 Apollo (AI-обработка сообщений)
- 🧠 **Извлечение структуры из текстов** — автоматическое построение таблиц из сообщений агрономов
- 🧾 **Классификация** — определение, является ли сообщение полевым отчётом
- 🖼️ **Обработка фото** — извлечение данных из изображений с таблицами
- 🎧 **Транскрибация аудио** — распознавание речи из голосовых сообщений
- 🧭 **Нормализация данных** — приведение к стандартным справочникам `cultures`, `operations`, `units`
- 🔁 **Гибкость** — работает с сообщениями различной структуры и полноты
- 📦 **Масштабируемость** — легко разворачивается в контейнерах (Docker), масштабируется горизонтально

### 📊 Superset (BI-аналитика)
- 📈 Интерактивные дашборды на основе данных из **PostgreSQL**
- 🔍 Анализ активности по операциям, культурам, подразделениям
- 🗓️ Историческая аналитика: динамика объёмов и валового сбора
- 🌍 Возможность добавления геопространственной аналитики (при наличии координат)

### 🔎 LangSmith (LLM-мониторинг)
- 🧵 Полная трассировка LLM-цепочек от **LangChain**
- 📊 Подсчёт и визуализация токенов (input/output)
- 🐞 Отладка промтов и цепочек обработки
- 📁 Разделение задач по проектам (например, классификация, разметка, нормализация)

---
## Сервисы и их роли

| Сервис            | Назначение                                                                                     |
|-------------------|-------------------------------------------------------------------------------------------------|
| **Hermes**        | Шлюз сообщений: приём, фильтрация, маршрутизация из WhatsApp / Telegram в дальнейшие модули.    |
| **Apollo**        | LLM‑сервис: извлечение сущностей, структурирование текстов, аудио, фото.                        |
| **PostgreSQL**    | Централизованное реляционное хранилище данных (метаданные, логи, справочники).                  |
| **Google Drive**  | Объектное хранилище обработанных файлов (изображения, документы, архивы отчётов).               |
| **Apache Superset** | Визуальная аналитика агроданных: дашборды, Ad‑hoc‑запросы, отчёты для менеджмента.            |
| **LangSmith**     | Мониторинг LLM‑цепочек (tracing), метрики качества, алерты на деградацию моделей.              |


## 📦 Компоненты

### 🔗 Связь между компонентами

- Hermes ↔️ Telegram/WhatsApp — получение сообщений
- Hermes ↔️ Apollo — обработка через REST API
- Hermes ↔️ PostgreSQL — хранение оригиналов и результатов
- Hermes ↔️ Google Drive — загрузка медиафайлов
- Superset ↔️ PostgreSQL — построение аналитических дашбордов

### 🛰️ Hermes (Go)

```Hermes``` — мультиканальный шлюз, обрабатывающий входящие сообщения агрономов.

**Возможности:**
- Поддержка **Telegram** и **WhatsApp**;
- Маршрутизация: Apollo API (текст, аудио, фото);
- Фильтрация и сохранение сообщений;
- Загрузка в Google Drive;
- Асинхронная обработка.

---
## 🗂️ Архитектура сервиса Hermes

```mermaid
flowchart LR
    %% ---------- Входящие каналы ----------
    subgraph Inbound
        direction TB
        WA[WhatsApp]
        TG[Telegram]
    end

    %% ---------- Hermes (Ingress) ----------
    H1["Hermes<br/>(ingress)"]

    WA --> H1
    TG --> H1

    %% ---------- Сториджи (Ingress) ----------
    subgraph Storage_In
        direction LR
        D1["Google&nbsp;Drive"]
        PG1[(PostgreSQL)]
    end
    H1 --> D1
    H1 --> PG1

    %% ---------- Apollo ----------
    AP["Apollo<br/>(LLM‑сервис)"]
    H1 -- "распознавание / флуд‑чек" --> AP

    %% ---------- LangSmith ----------
    LS["LangSmith<br/>(tracing)"]
    AP -. метрики .- LS

    %% ---------- Hermes (Egress) ----------
    H2["Hermes<br/>(egress)"]
    AP --> H2

    %% ---------- Сториджи (Egress) ----------
    subgraph Storage_Out
        direction LR
        D2["Google&nbsp;Drive"]
        PG2[(PostgreSQL)]
    end
    H2 --> D2
    H2 --> PG2

    %% ---------- Исходящие каналы ----------
    H2 --> WA2[WhatsApp]
    H2 --> TG2[Telegram]

    %% ---------- BI / аналитика ----------
    SUP["Apache Superset<br/>(dashboards)"]
    PG2 --> SUP

```
---

### 🗄️ Схема базы данных Hermes

```mermaid
erDiagram
    Worker {
        SERIAL worker_id PK
        varchar(1023) name
    }

    whatsapp {
        SERIAL id PK
        varchar(1023) whatsapp_id
        int  worker_id FK
    }

    telegram {
        SERIAL id PK
        varchar(1023) telegram_id
        int  worker_id FK
    }

    Chats {
        SERIAL id PK
        varchar(255) type
        varchar(1023) chat_name
    }

    ChatContext {
        SERIAL id PK
        int    chat_id   FK
        timestamp created_at
        varchar(1023) name
    }

    report {
        SERIAL report_id PK
        int    chat_context_id FK
        timestamp started_at
        timestamp last_updated_at
        timestamp finished_at
    }

    Message {
        SERIAL message_id PK
        int    worker_id FK
        int    chat_id   FK
        timestamp created_at
        text    content
    }

    images {
        SERIAL image_id PK
        int    message_id FK
        varchar(1023) image_url
    }

    tables {
        SERIAL id PK
        timestamp createdAt
        jsonb data
        int    message_id FK
    }

    Listener {
        SERIAL id PK
        int    worker_id FK
        int    chat_id   FK
    }

    Verbiage {
        SERIAL message_id PK
        int    worker_id FK
        int    chat_id   FK
        timestamp sended_at
        text   content
    }

    %% --- связи ---
    Worker   ||--|{ whatsapp : "has"
    Worker   ||--|{ telegram : "has"
    Worker   ||--|{ Message  : "writes"
    Worker   ||--|{ Listener : "listens"
    Worker   ||--|{ Verbiage : "sends"

    Chats    ||--|{ Message  : "contains"
    Chats    ||--|{ Listener : "observed by"
    Chats    ||--|{ Verbiage : "receives"
    Chats    ||--|{ ChatContext : "has contexts"

    ChatContext ||--|{ report : "aggregates"

    Message  ||--|{ images   : "includes"
    Message  ||--|{ tables   : "creates"
```
---
### 🤖 Apollo (Python + FastAPI)

```Apollo``` — интеллектуальный мультимодальный процессор для агросообщений.

**Возможности:**
- `process_message` — структурирование текста;
- `classify_message` — классификация по LLM;
- `process_photo` — OCR/LLM;
- `transcribe_audio` — транскрибация Whisper;
- `change_table` - исправление таблицы по запросу

**LLM-инфраструктура:**
- `ChatGPT-4.1`;
- `Whisper-1`;
- `sentence-transformers/all-MiniLM-L6-v2`;
- интеграция с LangChain и LangSmith.

---
## 🧪 API Эндпоинты (Apollo)

Все эндпоинты доступны по адресу: `http://localhost:8000`

---

### 📍 POST `/process_message`

🔹 **Описание:** Извлечение табличной структуры из текстового сообщения.

**Запрос:**
```json
{
  "message": "Пахота зяби под сою\nПо ПУ 7/1402\nОтд 17 7/141"
}
```
**Ответ:**
```json
{
  "table": [
    {
      "date": null,
      "division": "АОР",
      "operation": "Пахота",
      "culture": "Соя",
      "per_day": 7,
      "per_operation": 1402,
      "val_day": null,
      "val_beginning": null
    }
  ]
}
```

### 📍 POST `/classify_message`
🔹 **Описание**: Классификация сообщения — является ли оно отчётом об операции.

**Запрос:**
```json
{
  "message": "СЗР оз пш - 103/557 га"
}
```

**Ответ:**
```json
{
  "probability": 0.98,
  "prediction": 1
}
```
### 📍 POST `/process_photo`
🔹 **Описание**: Обработка фотографии (таблицы).

**Запрос:**
```json
{
  "photo": "BASE64_STRING",
  "type": "jpeg"
}
```



### 📍 POST `/transcribe_audio`
🔹 **Описание**: Распознавание речи из аудиофайла.

**Запрос:**
```json
{
  "audio": "BASE64_STRING",
  "type": "mp3"
}
```
**Ответ**:
```json
{
  "text": "Пример транскрибированного текста"
}
```

### 📍 POST `/change_table`
🔹 **Описание**: Изменение таблицы на основе пользовательских инструкций.

**Запрос:**
```json
{
  "table": [
    {
      "date": "17.04.2025",
      "division": "Юг",
      "operation": "Пахота",
      "culture": "Пшеница",
      "per_day": 50,
      "per_operation": 150,
      "val_day": null,
      "val_beginning": null
    }
  ],
  "message": "Изменить дату на 18.04.2025"
}
```
#### **Ответ**: аналогичен ответу /process_message.
---

## 🔄 Сценарий взаимодействия

- 👨‍🌾 Агроном отправляет 📩 сообщение (текст, фото, аудио).
- 📬 **Hermes**:
  - фильтрует, классифицирует;
  - сохраняет в БД, Google Drive;
  - направляет в **Apollo**.
- 🤖 **Apollo**:
  - обрабатывает через LLM;
  - извлекает данные и нормализует.
- 💾 Результаты сохраняются в PostgreSQL.
- 📊 **Superset** визуализирует данные в интерактивных дашбордах.
- 🔎 **LangSmith** позволяет отслеживать LLM-запросы и токены.

---

## ⚙️ Быстрый старт

```bash
git clone https://github.com/lild1tz/llm_coding_challenge.git
cd llm_coding_challenge/backend
docker-compose up --build -d
```

- Swagger Apollo: http://localhost:8000/docs  
- Superset: http://localhost:8088  
- LangSmith: https://smith.langchain.com 
---

## 🧰 Обслуживание

```bash
# Миграции
cd backend/hermes
migrate -path ./migrations -database "postgres://hermes:password@localhost:5432/hermes?sslmode=disable" up

# Логи
docker logs -f hermes
docker logs -f apollo

# Перегенерация QR-кода
rm -rf whatsapp-sessions
docker-compose restart hermes
```

---
### 💻 Coding AI Assistants

При создании и разработке сервиса **AgroAI Suite** активно использовались современные AI-инструменты для ускорения разработки, автодополнения кода и генерации решений:

- ⚡️ **Cursor** — AI-редактор на базе ChatGPT, использовался как основной IDE для Go и Python:
  - генерация и рефакторинг функций,
  - структурирование промтов и pydantic-моделей,
  - генерация документации и README.

- 🤖 **GitHub Copilot** — автодополнение на основе контекста:
  - помогал при написании FastAPI-эндпоинтов и SQL-запросов,
  - ускорял работу с типами данных и схемами.

---
### 👥 Команда проекта

| Имя | Роль |  Компания | Телеграм |
|-----|------|------------------------|----------|
| **Александр Иванов** | Team Lead, MLE | СБЕР | @lild1tz |
| **Валерий Хегай** | Go Developer | Яндекс, Backend Dev | @valerijhegaj |
| **Сергей Дементьев** | Product Analyst | Т‑Банк | @sdementev33 |

---

Сделано с ❤️ в рамках LLM Coding Challenge для АгроПрогресс.
