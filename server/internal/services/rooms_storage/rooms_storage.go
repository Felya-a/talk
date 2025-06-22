package room_storage

import (
	"sync"
	"talk/internal/core"

	. "talk/internal/lib/logger"
	in_memory_rooms_storage "talk/internal/services/rooms_storage/repository/in_memory"
	postgres_rooms_storage "talk/internal/services/rooms_storage/repository/postgres"

	"github.com/jmoiron/sqlx"
)

type RoomsStorage struct {
	mutex                sync.RWMutex
	liveRepository       RoomsRepository
	persistentRepository RoomsRepository
}

func NewRoomsStorage(db *sqlx.DB) *RoomsStorage {
	postgresRepository := postgres_rooms_storage.NewPostgresRoomsRepository(db)
	liveRepository := in_memory_rooms_storage.NewInMemoryRoomsRepository()

	// Инициализация in memory хранилища
	// Подгрузка комнат из бд
	persistentRooms, _ := postgresRepository.FindAll()
	for _, room := range persistentRooms {
		liveRepository.Save(room)
	}

	return &RoomsStorage{persistentRepository: postgresRepository, liveRepository: liveRepository}
}

func (rs *RoomsStorage) FindAll() ([]*core.Room, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return rs.liveRepository.FindAll()
}

func (rs *RoomsStorage) FindByUuid(uuid string) (*core.Room, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return rs.liveRepository.FindByUuid(uuid)
}

func (rs *RoomsStorage) FindByClient(client *core.Client) (*core.Room, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return rs.liveRepository.FindByClient(client)
}

func (rs *RoomsStorage) Save(room *core.Room) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	if err := rs.liveRepository.Save(room); err != nil {
		return err
	}

	if err := rs.persistentRepository.Save(room); err != nil {
		return err
	}

	Log.Info("[RoomsStorage] new room created", LogFields{
		"uuid": room.Uuid,
		"name": room.Name,
	})
	return nil
}

func (rs *RoomsStorage) Delete(room *core.Room) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	if err := rs.liveRepository.Delete(room); err != nil {
		return err
	}

	if err := rs.persistentRepository.Delete(room); err != nil {
		return err
	}

	Log.Info("[RoomsStorage] room deleted", LogFields{
		"uuid": room.Uuid,
		"name": room.Name,
	})
	return nil
}
