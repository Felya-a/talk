package handlers

import (
	"fmt"
	"talk/internal/models"
	. "talk/internal/models"
	. "talk/internal/ws"
)

type LeaveMessageHandler struct {
	RoomsPool *RoomsPool
}

func (h *LeaveMessageHandler) HandleMessage(client *Client, message ReceiveMessage) {
	room := h.RoomsPool.FindByClientId(client.Uuid)
	if room == nil {
		// TODO: Отправлять ошибку
		fmt.Println("Комната не найдена")
		return
	}

	if h.RoomsPool.FindByClientId(client.Uuid) == nil {
		// TODO: Отправлять ошибку
		fmt.Println("Пользователь не находится в этой комнате")
		return
	}

	room.Leave(client)

	messageData := map[string]interface{}{
		"peerID": client.Uuid,
	}

	room.Broadcast <- models.TransmitMessage{
		Type: MessageTypeRemovePeer,
		Data: messageData,
	}

	client.Hub.ShareRooms()
}
