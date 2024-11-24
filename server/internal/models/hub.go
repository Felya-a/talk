package models

import (
	"fmt"
	"sync"
)

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
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

// Добавление клиента
func (hub *Hub) addClient(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.Clients[client] = true
	fmt.Printf("Пользователь %s подключился\n", client.ID)
}

// Удаление клиента
func (hub *Hub) removeClient(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if _, ok := hub.Clients[client]; ok {
		delete(hub.Clients, client)
		close(client.Send)
		fmt.Printf("Пользователь %s отключился\n", client.ID)
	}
}

// Широковещательная отправка сообщений
func (hub *Hub) broadcastMessage(message Message) {
	fmt.Printf("Broadcast сообщение: %s\n", message)
	for client := range hub.Clients {
		select {
		case client.Send <- message:
			// Сообщение успешно отправлено
		}
	}
}
