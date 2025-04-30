package usecase

import (
	core "talk/internal/core"
	message_encoder "talk/internal/services/message_encoder"
)

type CreateRoomUseCase struct {
	Hub            *core.Hub
	MessageEncoder message_encoder.MessageEncoder
	ShareRooms     ShareRoomsUseCase
}

func (uc *CreateRoomUseCase) Execute(roomName string) {
	room := core.NewRoom(roomName, uc.MessageEncoder)
	uc.Hub.RoomsPool.Add(room)
	uc.ShareRooms.Execute()
}
