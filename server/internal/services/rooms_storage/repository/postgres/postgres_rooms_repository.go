package postgres_rooms_storage

import (
	"database/sql"
	"errors"
	"talk/internal/core"
	. "talk/internal/lib/logger"
	. "talk/internal/services/message_encoder"

	"github.com/jmoiron/sqlx"
)

type PostgresRoomsRepository struct {
	db             *sqlx.DB
	messageEncoder MessageEncoder
}

func NewPostgresRoomsRepository(
	db *sqlx.DB,
	messageEncoder MessageEncoder,
) *PostgresRoomsRepository {
	return &PostgresRoomsRepository{db, messageEncoder}
}

func (r *PostgresRoomsRepository) FindAll() ([]*core.Room, error) {
	var roomsDtos = make([]core.RoomDto, 0)

	err := r.db.Select(&roomsDtos, `
		select
			uuid,
			name
		from room
	`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var rooms = make([]*core.Room, 0)
	for _, roomDto := range roomsDtos {
		rooms = append(rooms, core.NewRoom(roomDto, r.messageEncoder))
	}

	return rooms, nil
}

func (r *PostgresRoomsRepository) FindByUuid(uuid string) (*core.Room, error) {
	var room *core.Room

	err := r.db.Get(&room, `
		select
			uuid as "Uuid",
			name as "Name"
		from room
		where uuid = $1
	`, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return room, nil
		}
		return room, err
	}

	return room, nil
}

func (r *PostgresRoomsRepository) FindByClient(*core.Client) (*core.Room, error) {
	return nil, errors.New("not implemented") // имплиментируется только в in memory хранилище
}

func (r *PostgresRoomsRepository) Save(room *core.Room) error {
	_, err := r.db.Exec(`
		insert into room (
			uuid,
			name
		) values (
			$1,
			$2
		)
	`, room.Uuid, room.Name)
	if err != nil {
		Log.Error("[PostgresRoomsRepository] error on save", Log.Err(err))
		return err
	}

	return nil
}

func (r *PostgresRoomsRepository) Delete(room *core.Room) error {
	_, err := r.db.Exec(`
		delete from room
		where uuid = $1
	`, room.Uuid)
	if err != nil {
		Log.Error("[PostgresRoomsRepository] error on delete", Log.Err(err))
		return err
	}

	return nil
}
