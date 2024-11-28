package ws

import (
	"fmt"
	"slices"
	. "talk/internal/models"

	"github.com/google/uuid"
)

type Room struct {
	Uuid      uuid.UUID
	Name      string
	Clients   []*Client
	Broadcast chan Message
}

func NewRoom(name string) *Room {
	return &Room{
		Uuid:      uuid.New(),
		Name:      name,
		Clients:   []*Client{},
		Broadcast: make(chan Message),
	}
}

func (r *Room) Run() {
	for message := range r.Broadcast {
		fmt.Println("Broadcast сообщение в комнате ", r.Uuid, message)
		for _, client := range r.Clients {
			client.Send <- message
		}
	}
}

func (r *Room) Join(client *Client) {
	index := slices.Index(r.Clients, client)
	if index != -1 {
		fmt.Println("Пользователь ", client.Uuid, " уже находится в комнате ", r.Uuid)
		return
	}
	r.Clients = append(r.Clients, client)
	fmt.Println("Пользователь ", client.Uuid, " добавлен в комнату ", r.Uuid)
}

func (r *Room) Leave(client *Client) {
	index := slices.Index(r.Clients, client)
	if index != -1 {
		r.Clients = append(r.Clients[:index], r.Clients[index+1:]...)
		fmt.Println("Пользователь ", client.Uuid, "исключен из комнаты", r.Uuid)
	} else {
		fmt.Println("Пользователь ", client.Uuid, "не найден в комнате", r.Uuid, "при исключении")
	}
}

func (r *Room) CheckExistUser(client *Client) bool {
	index := slices.Index(r.Clients, client)
	return index != -1
}
