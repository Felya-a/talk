package core

import (
	"slices"
	"sync"

	. "talk/internal/lib/logger"
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

func (rp *RoomsPool) FindByUuid(uuid string) *Room {
	for _, room := range rp.rooms {
		if room.Uuid == uuid {
			return room
		}
	}
	return nil
}

func (rp *RoomsPool) FindByClient(client *Client) *Room {
	for _, room := range rp.rooms {
		for _, roomClient := range room.Clients {
			if roomClient.Uuid == client.Uuid {
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
	Log.Info("new room created", LogFields{
		"uuid": room.Uuid,
		"name": room.Name,
	})
}

func (rp *RoomsPool) Remove(room *Room) {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	index := slices.Index(rp.rooms, room)
	if index != -1 {
		rp.rooms[index] = rp.rooms[len(rp.rooms)-1]
		rp.rooms = rp.rooms[:len(rp.rooms)-1]
		Log.Info("room deleted from rooms pool", LogFields{"uuid": room.Uuid})
	} else {
		Log.Warn("room not found on delete from rooms pool", LogFields{"uuid": room.Uuid})
	}
}
