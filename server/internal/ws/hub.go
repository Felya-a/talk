package ws

import (
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
	Broadcast  chan TransmitMessage
	RoomsPool  *RoomsPool
	mu         sync.Mutex
}

func NewHub(roomsPool *RoomsPool) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan TransmitMessage),
		RoomsPool:  roomsPool,
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			go hub.addClient(client)
		case client := <-hub.Unregister:
			go hub.removeClient(client)
		case message := <-hub.Broadcast:
			go hub.broadcastMessage(message)
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

	if _, ok := hub.Clients[client]; !ok {
		return
	}

	delete(hub.Clients, client)
	close(client.Send)

	room := hub.RoomsPool.FindByClientId(client.Uuid)
	if room != nil {
		room.Leave(client)
	}

	client.conn.Close()
	hub.ShareRooms()
	fmt.Printf("Пользователь %s отключился\n", client.Uuid)
}

// Широковещательная отправка сообщений
func (hub *Hub) broadcastMessage(message TransmitMessage) {
	fmt.Printf("Broadcast сообщение: %s\n", message)
	for client := range hub.Clients {
		client.Send <- message
	}
}

func (hub *Hub) ShareRooms() {
	rooms := hub.RoomsPool.FindAll()

	roomsForShareMessage := make([]map[string]interface{}, len(rooms))
	for i, room := range rooms {
		clientsForShareMessage := make([]map[string]interface{}, len(room.Clients))
		for j, client := range room.Clients {
			clientsForShareMessage[j] = map[string]interface{}{
				"uuid": client.Uuid,
			}
		}
		roomsForShareMessage[i] = map[string]interface{}{
			"uuid":    room.Uuid,
			"name":    room.Name,
			"clients": clientsForShareMessage,
		}
	}

	hub.Broadcast <- models.TransmitMessage{
		Type: MessageTypeShareRooms,
		Data: roomsForShareMessage,
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
