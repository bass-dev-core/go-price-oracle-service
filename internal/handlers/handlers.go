// Файл: internal/handlers/handlers.go
package handlers

import (
	"context"
	"encoding/json"
	"go-price-oracle-service/internal/ws"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/adshao/go-binance/v2/futures"
)

var (
	data struct {
		Bid       float64
		Ask       float64
		Inventory int
		Another   string
	}
	mu sync.RWMutex
)

var client = futures.NewClient("", "")

func getBid() float64 {
	ticker, err := client.NewDepthService().Symbol("BCHUSDT").Limit(5).Do(context.Background())
	if err != nil {
		log.Println("Ошибка при получении bid:", err)
		return 0
	}
	if len(ticker.Bids) > 0 {
		price, err := strconv.ParseFloat(ticker.Bids[0].Price, 64)
		if err != nil {
			log.Println("Ошибка парсинга bid:", err)
			return 0
		}
		return price
	}
	return 0
}

func getAsk() float64 {
	ticker, err := client.NewDepthService().Symbol("BCHUSDT").Limit(5).Do(context.Background())
	if err != nil {
		log.Println("Ошибка при получении ask:", err)
		return 0
	}
	if len(ticker.Asks) > 0 {
		price, err := strconv.ParseFloat(ticker.Asks[0].Price, 64)
		if err != nil {
			log.Println("Ошибка парсинга ask:", err)
			return 0
		}
		return price
	}
	return 0
}

func getInventory() int {
	return data.Inventory
}

func getAnotherData() string {
	return data.Another
}

func GetParametersHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	response := map[string]interface{}{
		"bid":       getBid(),
		"ask":       getAsk(),
		"inventory": getInventory(),
		"another":   getAnotherData(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if bid, ok := payload["bid"].(float64); ok {
		data.Bid = bid
	}
	if ask, ok := payload["ask"].(float64); ok {
		data.Ask = ask
	}
	if inv, ok := payload["inventory"].(float64); ok {
		data.Inventory = int(inv)
	}
	if a, ok := payload["another"].(string); ok {
		data.Another = a
	}

	w.WriteHeader(http.StatusOK)
}

// Новый объединённый обработчик, который возвращает данные из WebSocket и REST
func GetCombinedPriceHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем WS-данные
	wsData := ws.CurrentPriceData()

	// Получаем REST-данные
	restData := map[string]interface{}{
		"bid":       getBid(),
		"ask":       getAsk(),
		"inventory": getInventory(),
		"another":   getAnotherData(),
	}

	// Объединяем оба набора данных
	response := map[string]interface{}{
		"ws":   wsData,
		"rest": restData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
