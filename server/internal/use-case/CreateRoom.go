package usecase

import (
	core "talk/internal/core"
	message_encoder "talk/internal/services/message_encoder"
	rooms_storage "talk/internal/services/rooms_storage"
)

type CreateRoomUseCase struct {
	RoomsStorage   *rooms_storage.RoomsStorage
	MessageEncoder message_encoder.MessageEncoder
	ShareRooms     ShareRoomsUseCase
}

func (uc *CreateRoomUseCase) Execute(roomName string) {
	room := core.NewRoom(roomName, uc.MessageEncoder)
	uc.RoomsStorage.Save(room)
	uc.ShareRooms.Execute()
}
