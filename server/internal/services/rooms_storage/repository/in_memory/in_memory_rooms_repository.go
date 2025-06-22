package in_memory_rooms_storage

import (
	"slices"
	. "talk/internal/core"
	. "talk/internal/lib/logger"
)

type InMemoryRoomsRepository struct {
	rooms []*Room
}

func NewInMemoryRoomsRepository() *InMemoryRoomsRepository {
	return &InMemoryRoomsRepository{
		rooms: make([]*Room, 0),
	}
}

func (r *InMemoryRoomsRepository) FindAll() ([]*Room, error) {
	return r.rooms, nil
}

func (r *InMemoryRoomsRepository) FindByUuid(uuid string) (*Room, error) {
	for _, room := range r.rooms {
		if room.Uuid == uuid {
			return room, nil
		}
	}
	return nil, nil
}

func (r *InMemoryRoomsRepository) FindByClient(client *Client) (*Room, error) {
	for _, room := range r.rooms {
		for _, roomClient := range room.Clients {
			if roomClient.Uuid == client.Uuid {
				return room, nil
			}
		}
	}
	return nil, nil
}

func (r *InMemoryRoomsRepository) Save(room *Room) error {
	r.rooms = append(r.rooms, room)
	return nil
}

func (r *InMemoryRoomsRepository) Delete(room *Room) error {
	index := slices.Index(r.rooms, room)
	if index != -1 {
		r.rooms[index] = r.rooms[len(r.rooms)-1]
		r.rooms = r.rooms[:len(r.rooms)-1]
		Log.Info("[InMemoryRoomsRepository] room deleted from rooms pool", LogFields{"uuid": room.Uuid})
	} else {
		Log.Warn("[InMemoryRoomsRepository] room not found on delete from rooms pool", LogFields{"uuid": room.Uuid})
	}

	return nil
}
