package usecase

import (
	"talk/internal/core"
	. "talk/internal/lib/logger"
	. "talk/internal/models/messages"
	message_encoder "talk/internal/services/message_encoder"
	rooms_storage "talk/internal/services/rooms_storage"
)

type JoinClientUseCase struct {
	RoomsStorage   *rooms_storage.RoomsStorage
	ShareRooms     ShareRoomsUseCase
	MessageEncoder message_encoder.MessageEncoder
}

func (uc *JoinClientUseCase) Execute(client *core.Client, dto JoinMessageDto) {
	// TODO: Обрабатывать ошибку
	room, _ := uc.RoomsStorage.FindByUuid(dto.RoomUuid)
	if room == nil {
		// TODO: Отправлять ошибку
		Log.Warn("[JoinClientUseCase] room not found", LogFields{
			"clientUuid": client.Uuid,
			"roomUuid":   dto.RoomUuid,
		})
		return
	}

	// Переподключение пользователя из другой комнаты
	if foundRoom, _ := uc.RoomsStorage.FindByClient(client); foundRoom != nil {
		if foundRoom.Uuid == room.Uuid {
			Log.Info("client is trying to reconnect to his room", LogFields{"clientUuid": client.Uuid, "roomUuid": foundRoom.Uuid})
			return
		}

		Log.Info("client already belongs to one of the rooms", LogFields{"clientUuid": client.Uuid, "roomUuid": foundRoom.Uuid})
		if err := foundRoom.Leave(client); err != nil {
			Log.Error("error on leave from room", LogFields{"clientUuid": client.Uuid, "roomUuid": foundRoom.Uuid})
			// TODO: Отправлять ошибку
			return
		}
	}

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
