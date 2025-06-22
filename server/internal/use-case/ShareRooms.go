package usecase

import (
	core "talk/internal/core"
	events "talk/internal/core/events"
	. "talk/internal/models/messages"
	message_encoder "talk/internal/services/message_encoder"
	rooms_storage "talk/internal/services/rooms_storage"
)

type ShareRoomsUseCase struct {
	Hub            *core.Hub
	RoomsStorage   *rooms_storage.RoomsStorage
	MessageEncoder message_encoder.MessageEncoder
}

func (uc *ShareRoomsUseCase) Execute() {
	// TODO: Обрабатывать ошибку
	rooms, _ := uc.RoomsStorage.FindAll()

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

func (uc *ShareRoomsUseCase) ExecuteEvent(event core.Event) {
	uc.Execute()

	eventData, ok := event.(events.ClientConnectedEvent)
	if !ok {
		return
	}

	for _, client := range uc.Hub.Clients {
		if client.Uuid == eventData.ClientUuid {
			client.Outbound <- uc.MessageEncoder.BuildClientInfoMessage(ClientInfoMessageDto{Uuid: eventData.ClientUuid})
		}
	}
}
