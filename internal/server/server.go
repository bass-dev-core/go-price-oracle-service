// Файл: internal/server/server.go
package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go-price-oracle-service/internal/handlers"
	"go-price-oracle-service/internal/ws"
)

func StartServer() {
	r := mux.NewRouter()

	// Ваши маршруты, например, для API
	r.HandleFunc("/api/price", handlers.GetCombinedPriceHandler).Methods("GET")

	// Запуск WebSocket стрима
	ws.InitPriceStream()

	// Запуск webhook отправки данных каждые 3 секунды
	ws.StartWebhookSender("http://127.0.0.1:3000/webhook")

	log.Println("Starting server on :6942")
	http.ListenAndServe(":6942", r)
}
