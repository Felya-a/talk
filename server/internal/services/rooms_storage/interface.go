package room_storage

import (
	core "talk/internal/core"
)

type RoomsRepository interface {
	FindAll() ([]*core.Room, error)
	FindByUuid(uuid string) (*core.Room, error)
	FindByClient(*core.Client) (*core.Room, error)
	Save(*core.Room) error
	Delete(*core.Room) error
}
