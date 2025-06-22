package core

import (
	"slices"
	"sync"
	events "talk/internal/core/events"
	. "talk/internal/lib/logger"
	. "talk/internal/models/messages"
	"talk/internal/utils"
)

type Hub struct {
	mutex    sync.Mutex
	Clients  []*Client
	EventBus *EventBus

	Register   chan Connection
	Unregister chan *Client
	Broadcast  chan TransmitMessage

	Router *MessageRouter

	RoomsStorage RoomsStorage
}

func NewHub(
	router *MessageRouter,
	roomsStorage RoomsStorage,
) *Hub {
	return &Hub{
		Clients:      make([]*Client, 0),
		Register:     make(chan Connection),
		Unregister:   make(chan *Client),
		Broadcast:    make(chan TransmitMessage),
		Router:       router,
		RoomsStorage: roomsStorage,
		EventBus:     NewEventBus(),
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
	hub.EventBus.Publish(events.ClientConnectedEvent{ClientUuid: client.Uuid})
}

// Удаление пользователя
func (hub *Hub) removeClient(client *Client) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	if !slices.Contains(hub.Clients, client) {
		Log.Warn("[Hub] client already deleted from hub", LogFields{"clientUuid": client.Uuid})
		return
	}

	// TODO: проверить работу
	hub.Clients = utils.RemoveSliceElement(hub.Clients, client)

	room, err := hub.RoomsStorage.FindByClient(client)
	if room != nil && err == nil {
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
