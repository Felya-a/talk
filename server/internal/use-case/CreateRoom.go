package usecase

import (
	core "talk/internal/core"
	message_encoder "talk/internal/services/message_encoder"
	rooms_storage "talk/internal/services/rooms_storage"

	"github.com/google/uuid"
)

type CreateRoomUseCase struct {
	RoomsStorage   *rooms_storage.RoomsStorage
	MessageEncoder message_encoder.MessageEncoder
	ShareRooms     ShareRoomsUseCase
}

func (uc *CreateRoomUseCase) Execute(roomName string) {
	room := core.NewRoom(
		core.RoomDto{
			Name: roomName,
			Uuid: uuid.New().String(),
		}, uc.MessageEncoder,
	)
	uc.RoomsStorage.Save(room)
	uc.ShareRooms.Execute()
}
