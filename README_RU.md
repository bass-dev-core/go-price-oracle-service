# 💡 Go Price Oracle Service

Сервер на Go для получения и обновления параметров (Bid, Ask, Inventory и др.) через REST API и Webhook. Используется в проектах с маркет-мейкерами и DEX-интеграциями.

---

## 🛠 Установка

```bash
git clone https://your-repo-url.git
cd go-price-oracle-service
make install
```

---

## 🧾 Создание файла `.env`

Создайте `.env` в корне проекта со следующими переменными:

```env
BINANCE_API_KEY=your_api_key
BINANCE_API_SECRET=your_secret_key
```

---

## 🧪 Запуск

### 🔹 Локально

```bash
make build
./go-price-oracle-service
```

или:

```bash
make run
```

---

## 🚀 Развертывание (Docker)

```bash
make deploy
```

Сервис будет доступен по порту `6942`.

---

## 🌐 Эндпоинты

| Метод | URL               | Описание                        |
|-------|-------------------|---------------------------------|
| GET    | `http://localhost:6942/api/price` | Получение текущих котировок           |
| POST   | `http://localhost:3000/webhook`        | Обновление котировок чере webhook    |
