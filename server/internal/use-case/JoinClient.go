package usecase

import (
	"talk/internal/core"
	. "talk/internal/lib/logger"
	. "talk/internal/models/messages"
	message_encoder "talk/internal/services/message_encoder"
)

type JoinClientUseCase struct {
	Hub            *core.Hub
	ShareRooms     ShareRoomsUseCase
	MessageEncoder message_encoder.MessageEncoder
}

func (uc *JoinClientUseCase) Execute(client *core.Client, dto JoinMessageDto) {
	room := uc.Hub.RoomsPool.FindByUuid(dto.RoomUuid)
	if room == nil {
		// TODO: Отправлять ошибку
		Log.Warn("[JoinClientUseCase] room not found", LogFields{
			"clientUuid": client.Uuid,
			"roomUuid":   dto.RoomUuid,
		})
		return
	}

	// TODO: Возможно лишняя проверка. Она есть в Room
	// if uc.Hub.RoomsPool.FindByClient(client) != nil {
	// 	// TODO: Отправлять ошибку
	// 	Log.Warn("[JoinClientUseCase] client already connected to room", LogFields{
	// 		"clientUuid": client.Uuid,
	// 		"roomUuid":   room.Uuid,
	// 	})
	// 	return
	// }

	room.Join(client)

	for _, roomClient := range room.Clients {
		if roomClient.Uuid == client.Uuid {
			continue
		}

		// Отправка сообщения существующему пользователю о новом пользователе
		roomClient.Outbound <- uc.MessageEncoder.BuildAddPeerMessage(AddPeerMessageDto{PeerID: client.Uuid, CreateOffer: false})

		// Отправка сообщения новому пользователю о существующем
		client.Outbound <- uc.MessageEncoder.BuildAddPeerMessage(AddPeerMessageDto{PeerID: roomClient.Uuid, CreateOffer: true})
	}

	uc.ShareRooms.Execute()
}
