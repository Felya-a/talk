package ws

import (
	"encoding/json"
	"fmt"
	"sync"
	"talk/internal/models"
	. "talk/internal/models"

	"github.com/google/uuid"
)

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	RoomsPool  *RoomsPool
	mu         sync.Mutex
}

func NewHub(roomsPool *RoomsPool) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
		RoomsPool:  roomsPool,
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.addClient(client)
		case client := <-hub.Unregister:
			hub.removeClient(client)
		case message := <-hub.Broadcast:
			hub.broadcastMessage(message)
		}
	}
}

// Добавление пользователя
func (hub *Hub) addClient(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	hub.Clients[client] = true
	fmt.Printf("Пользователь %s подключился\n", client.Uuid)
}

// Удаление пользователя
func (hub *Hub) removeClient(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if _, ok := hub.Clients[client]; ok {
		delete(hub.Clients, client)
		close(client.Send)
		client.conn.Close()
		fmt.Printf("Пользователь %s отключился\n", client.Uuid)
	}
}

// Широковещательная отправка сообщений
func (hub *Hub) broadcastMessage(message Message) {
	fmt.Printf("Broadcast сообщение: %s\n", message)
	for client := range hub.Clients {
		client.Send <- message
	}
}

func (hub *Hub) ShareRooms() {
	rooms := hub.RoomsPool.FindAll()

	var roomsForShareMessage []map[string]interface{}
	for _, room := range rooms {
		clientsForShareMessage := []map[string]interface{}{}
		for _, client := range room.Clients {
			clientsForShareMessage = append(clientsForShareMessage, map[string]interface{}{
				"id": client.Uuid,
			})
		}
		roomsForShareMessage = append(roomsForShareMessage, map[string]interface{}{
			"id":      room.Uuid,
			"name":    room.Name,
			"clients": clientsForShareMessage,
		})
	}

	encodeData, err := json.Marshal(roomsForShareMessage)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	hub.Broadcast <- models.Message{
		Type: MessageTypeShareRooms,
		Data: string(encodeData),
	}
}

func (hub *Hub) FindClientByUuid(uuid uuid.UUID) *Client {
	for client := range hub.Clients {
		if client.Uuid == uuid {
			return client
		}
	}
	return nil
}
