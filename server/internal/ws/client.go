package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"talk/internal/lib/logger/sl"
	. "talk/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Время, разрешённое на запись сообщения клиенту
	writeWait = 10 * time.Second

	// Время, разрешённое на получение следующего сообщения pong от клиента
	pongWait = 10 * time.Second

	// Интервал отправки ping-сообщений клиенту. Должен быть меньше pongWait
	pingPeriod = (pongWait * 9) / 10

	// Максимальный размер сообщения в байтах, разрешённый от клиента
	maxMessageSize = 512
)

type Client struct {
	Uuid uuid.UUID
	Hub  *Hub
	conn *websocket.Conn
	Send chan (Message)
}

func NewClient(
	conn *websocket.Conn,
	hub *Hub,
) *Client {
	return &Client{
		Uuid: uuid.New(),
		conn: conn,
		Hub:  hub,
		Send: make(chan Message),
	}
}

// Забирает сообщения из WebSocket-соединения и отправляет их в hub (концентратор).
func (client *Client) ReadPump(router *MessageRouter) {
	defer func() {
		client.Hub.Unregister <- client
	}()

	client.conn.SetReadLimit(maxMessageSize)

	for {
		_, messageBytes, err := client.conn.ReadMessage()
		fmt.Println("Новое сообщение от пользователя " + client.Uuid.String())
		fmt.Println(string(messageBytes))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("IsUnexpectedCloseError: %v", err)
				return
			}
			fmt.Println("Неивестная ошибка при получении сообщения ", sl.Err(err))
			time.Sleep(2 * time.Second)
			continue
		}

		// Десериализация сообщения
		var message Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		router.RouteMessage(client, message)
	}
}

// Перенаправляет сообщения из hub обратно в WebSocket-соединение.
func (client *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Hub.Unregister <- client
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				fmt.Println("Канал пользователя ", client.Uuid, " был закрыт")
				return
			}
			fmt.Println("Новое сообщение в канале пользователя: ", message)

			// Сериализация сообщения
			encodeMessage, err := message.ToJson()
			if err != nil {
				fmt.Println("Ошибка формирования сообщения ", sl.Err(err))
			}

			if err := client.conn.WriteMessage(websocket.TextMessage, encodeMessage); err != nil {
				fmt.Println("Ошибка при отправке сообщения", sl.Err(err))
				continue
			}
		case <-ticker.C:
			unixTimeMillis := time.Now().UnixNano() / int64(time.Millisecond) // Время в миллисекундах

			pingMessage := &Message{
				Type: MessageTypePing,
				Data: strconv.FormatInt(unixTimeMillis, 10),
			}

			encodeMessage, err := pingMessage.ToJson()
			if err != nil {
				fmt.Println("Ошибка формирования сообщения ", sl.Err(err))
			}

			if err := client.conn.WriteMessage(websocket.TextMessage, encodeMessage); err != nil {
				fmt.Println("Ошибка при отправке ping сообщения", sl.Err(err))
				return
			}
		}
	}
}
