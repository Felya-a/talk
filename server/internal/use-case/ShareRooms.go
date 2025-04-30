package usecase

import (
	core "talk/internal/core"
	. "talk/internal/models/messages"
	message_encoder "talk/internal/services/message_encoder"
)

type ShareRoomsUseCase struct {
	Hub            *core.Hub
	MessageEncoder message_encoder.MessageEncoder
}

func (uc *ShareRoomsUseCase) Execute() {
	rooms := uc.Hub.RoomsPool.FindAll()

	roomsInfoDto := make([]RoomInfo, len(rooms))
	for i, room := range rooms {
		clientsUuids := make([]string, len(room.Clients))
		for i, client := range room.Clients {
			clientsUuids[i] = client.Uuid
		}

		roomsInfoDto[i] = RoomInfo{
			Uuid:    room.Uuid,
			Name:    room.Name,
			Clients: clientsUuids,
		}
	}

	uc.Hub.Broadcast <- uc.MessageEncoder.BuildShareRoomsMessage(ShareRoomsMessageDto{Rooms: roomsInfoDto})
}
