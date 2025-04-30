package core

import (
	"slices"
	"sync"
	. "talk/internal/lib/logger"
	. "talk/internal/models/messages"
)

type Hub struct {
	mutex      sync.Mutex
	Clients    []*Client
	Register   chan Connection
	Unregister chan *Client
	Broadcast  chan TransmitMessage
	Router     *MessageRouter
	RoomsPool  *RoomsPool
}

func NewHub(
	router *MessageRouter,
	roomsPool *RoomsPool,
) *Hub {
	return &Hub{
		Clients:    make([]*Client, 0),
		Register:   make(chan Connection),
		Unregister: make(chan *Client),
		Broadcast:  make(chan TransmitMessage),
		Router:     router,
		RoomsPool:  roomsPool,
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case connection := <-hub.Register:
			go hub.addClient(connection)
		case client := <-hub.Unregister:
			go hub.removeClient(client)
		case message := <-hub.Broadcast:
			go hub.broadcastMessage(message)
		}
	}
}

func (hub *Hub) HandleMessage(client *Client, msg ReceiveMessage) error {
	return hub.Router.RouteMessage(client, msg)
}

// Добавление пользователя
func (hub *Hub) addClient(connection Connection) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	client := NewClient(connection, hub)
	go client.WritePump()
	go client.ReadPump(hub.Router)

	hub.Clients = append(hub.Clients, client)
	Log.Info("[Hub] client connected", LogFields{"clientUuid": client.Uuid})
}

// Удаление пользователя
func (hub *Hub) removeClient(client *Client) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	if !slices.Contains(hub.Clients, client) {
		Log.Warn("[Hub] client already deleted from hub", LogFields{"clientUuid": client.Uuid})
		return
	}

	index := slices.Index(hub.Clients, client)
	hub.Clients[index] = hub.Clients[len(hub.Clients)-1]
	hub.Clients = hub.Clients[:len(hub.Clients)-1]

	room := hub.RoomsPool.FindByClient(client)
	if room != nil {
		room.Leave(client)
	}

	client.Kill()

	Log.Info("[Hub] client disconnected", LogFields{"clientUuid": client.Uuid})
}

// Широковещательная отправка сообщений
func (hub *Hub) broadcastMessage(message TransmitMessage) {
	Log.Debug("[Hub] broadcast message", LogFields{"message": message})
	for _, client := range hub.Clients {
		client.Outbound <- message
	}
}
