package handlers

import (
	"fmt"
	"log"
	"net/http"
	"talk/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Настройка WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Разрешаем все подключения (можно уточнить для безопасности)
		return true
	},
}

func HandleConnections(hub *models.Hub) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("Новое подключение")

		// Обновляем HTTP-соединение до WebSocket
		ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Printf("Ошибка обновления до WebSocket: %v", err)
			return
		}

		clientId, err := uuid.NewV6()
		if err != nil {
			fmt.Println("Ошибка при создании uuid пользователя")
			return
		}

		client := &models.Client{
			ID:   clientId,
			Conn: ws,
			Hub:  hub,
			Send: make(chan models.Message),
		}
		hub.Register <- client

		ws.SetCloseHandler(func(code int, text string) error {
			fmt.Println("Close handler", code, text)
			hub.Unregister <- client
			return nil
		})

		// Если в течении 8 часов от клиента не придёт ни одного сообщения то соединение разорвётся
		ws.SetReadDeadline(time.Now().Add(8 * time.Hour))
		// Если в течении 8 часов клиенту не будет отправлено ни одного сообщения то соединение разорвётся
		ws.SetWriteDeadline(time.Now().Add(8 * time.Hour))

		// Читаем сообщения от клиента
		for {
			_, msgBytes, err := ws.ReadMessage()
			if err != nil {
				log.Printf("Ошибка чтения сообщения: %v", err)
				break
			}
			fmt.Println("Сообщение от пользователя: ", string(msgBytes))

			msg := models.Message{
				Type: "message",
				Data: string(msgBytes),
			}

			// // Отправляем сообщение в бродкаст канал
			hub.Broadcast <- msg
		}
	}
}
