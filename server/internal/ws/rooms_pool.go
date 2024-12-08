package ws

import (
	"fmt"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type RoomsPool struct {
	rooms []*Room
	mu    sync.Mutex
}

func NewRoomsPool() *RoomsPool {
	return &RoomsPool{
		rooms: []*Room{},
	}
}

func (rp *RoomsPool) FindAll() []*Room {
	return rp.rooms
}

func (rp *RoomsPool) FindByUuid(uuid uuid.UUID) *Room {
	for _, room := range rp.rooms {
		if room.Uuid == uuid {
			return room
		}
	}
	return nil
}

func (rp *RoomsPool) FindByClientId(clientUuid uuid.UUID) *Room {
	for _, room := range rp.rooms {
		for _, client := range room.Clients {
			if client.Uuid == clientUuid {
				return room
			}
		}
	}
	return nil
}

func (rp *RoomsPool) Add(room *Room) {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	rp.rooms = append(rp.rooms, room)
	fmt.Println("Создана новая комната")
	fmt.Println("UUID: ", room.Uuid)
	fmt.Println("Name: ", room.Name)
}

func (rp *RoomsPool) Remove(room *Room) {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	index := slices.Index(rp.rooms, room)

	if index != -1 {
		rp.rooms = append(rp.rooms[:index], rp.rooms[index+1:]...)
		fmt.Println("Комната ", room.Uuid, "удалена")
	} else {
		fmt.Println("Комната ", room.Uuid, "не найдена при удалении")
	}

	// Отключение пользователей
	for _, client := range room.Clients {
		client.Hub.Unregister <- client
	}

	room = nil
}
