package usecase

import (
	"talk/internal/core"
	. "talk/internal/lib/logger"
	. "talk/internal/models/messages"
	message_encoder "talk/internal/services/message_encoder"
)

type LeaveClientUseCase struct {
	Hub            *core.Hub
	ShareRooms     ShareRoomsUseCase
	MessageEncoder message_encoder.MessageEncoder
}

func (uc *LeaveClientUseCase) Execute(client *core.Client) {
	room := uc.Hub.RoomsPool.FindByClient(client)
	if room == nil {
		// TODO: Отправлять ошибку
		Log.Warn("room not found", nil)
		return
	}

	if err := room.Leave(client); err != nil {
		// TODO: Отправлять ошибку
		Log.Error("client not leave from room", Log.Err(err))
		return
	}

	uc.Hub.Broadcast <- uc.MessageEncoder.BuildRemovePeerMessage(RemovePeerMessageDto{ClientUuid: client.Uuid})

	uc.ShareRooms.Execute()
}
