# 💡 Go Price Oracle Service

A Go-based server for retrieving and updating parameters (Bid, Ask, Inventory, etc.) via REST API and Webhook. Used in market-making and DEX integration projects.

---

## 🛠 Installation

```bash
git clone https://your-repo-url.git
cd go-price-oracle-service
make install
```

---

## 🧾 Creating the `.env` File

Create a `.env` file in the project root with the following content:

```env
BINANCE_API_KEY=your_api_key
BINANCE_API_SECRET=your_secret_key
```

---

## 🧪 Running

### 🔹 Locally

```bash
make build
./go-price-oracle-service
```

or:

```bash
make run
```

---

## 🚀 Deployment (Docker)

```bash
make deploy
```

Service will be available on port `6942`.

---

## 🌐 Endpoints

| Method | URL               | Description                      |
|--------|-------------------|----------------------------------|
| GET    | `http://localhost:6942/api/price` | Get current parameters           |
| POST   | `http://localhost:3000/webhook`        | Update parameters via Webhook    |
