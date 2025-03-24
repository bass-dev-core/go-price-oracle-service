// Файл: internal/ws/ws.go
package ws

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// PriceData содержит данные о markPrice и связанных параметрах
type PriceData struct {
	Symbol          string `json:"symbol"`
	EventTime       int64  `json:"E"`
	MarkPrice       string `json:"p"`
	IndexPrice      string `json:"i"`
	EstimatedSettle string `json:"P"`
	FundingRate     string `json:"r"`
	NextFundingTime int64  `json:"T"`
}

// markPriceMessage соответствует сообщению markPriceUpdate от Binance
type markPriceMessage struct {
	EventType       string `json:"e"` // Обычно "markPriceUpdate"
	EventTime       int64  `json:"E"`
	Symbol          string `json:"s"`
	MarkPrice       string `json:"p"`
	IndexPrice      string `json:"i"`
	EstimatedSettle string `json:"P"`
	FundingRate     string `json:"r"`
	NextFundingTime int64  `json:"T"`
}

var priceData = &PriceData{Symbol: "BCHUSDT"}
var mu sync.RWMutex

// InitPriceStream запускает горутину для получения markPrice-обновлений
func InitPriceStream() {
	go connectMarkPriceStream()
}

// connectMarkPriceStream устанавливает WebSocket-соединение с Binance Futures Testnet и обновляет priceData
func connectMarkPriceStream() {
	// Используем правильный URL для получения markPrice обновлений
	url := "wss://stream.binancefuture.com/ws/bchusdt@markPrice"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("markPrice WS error:", err)
		return
	}
	defer c.Close()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("markPrice read error:", err)
			continue
		}

		var data markPriceMessage
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Println("markPrice unmarshal error:", err)
			continue
		}

		mu.Lock()
		priceData.EventTime = data.EventTime
		priceData.MarkPrice = data.MarkPrice
		priceData.IndexPrice = data.IndexPrice
		priceData.EstimatedSettle = data.EstimatedSettle
		priceData.FundingRate = data.FundingRate
		priceData.NextFundingTime = data.NextFundingTime
		mu.Unlock()
	}
}

// GetPriceData отдает текущие данные по цене через HTTP
func GetPriceData(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(priceData)
}

// CurrentPriceData возвращает копию текущих данных
func CurrentPriceData() PriceData {
	mu.RLock()
	defer mu.RUnlock()
	return *priceData
}

// StartWebhookSender запускает горутину, которая каждые 3 секунды отправляет данные по webhookURL
func StartWebhookSender(webhookURL string) {
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for range ticker.C {
			// Получаем актуальные данные
			data := CurrentPriceData()
			payload, err := json.Marshal(data)
			if err != nil {
				log.Println("Webhook: ошибка маршалинга данных:", err)
				continue
			}

			// Отправляем HTTP POST запрос на webhookURL
			resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				log.Println("Webhook: ошибка отправки запроса:", err)
				continue
			}
			resp.Body.Close()
			log.Println("Webhook: данные отправлены успешно")
		}
	}()
}
