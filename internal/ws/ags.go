package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// AggregateTradeData содержит данные о агрегированных торгах
type AggregateTradeData struct {
	EventType     string `json:"e"` // Обычно "aggTrade"
	EventTime     int64  `json:"E"`
	Symbol        string `json:"s"`
	AggregateID   int64  `json:"a"`
	Price         string `json:"p"`
	Quantity      string `json:"q"`
	FirstTradeID  int64  `json:"f"`
	LastTradeID   int64  `json:"l"`
	TradeTime     int64  `json:"T"`
	IsMarketMaker bool   `json:"m"`
}

// InitAggregateTradeStream запускает горутину для получения обновлений агрегированных торгов
func InitAggregateTradeStream(symbol string) {
	go connectAggregateTradeStream(symbol)
}

// connectAggregateTradeStream устанавливает WebSocket-соединение с Binance для агрегированных торгов
func connectAggregateTradeStream(symbol string) {
	url := "wss://stream.binance.com:9443/ws/" + symbol + "@aggTrade"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("AggregateTrade WS error:", err)
		return
	}
	defer c.Close()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("AggregateTrade read error:", err)
			continue
		}

		var data AggregateTradeData
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Println("AggregateTrade unmarshal error:", err)
			continue
		}

		// Здесь вы можете обработать полученные данные
		log.Printf("Получены данные агрегированных торгов: %+v\n", data)
	}
}
