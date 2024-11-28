package handlers

import (
	"encoding/json"
	"fmt"
	"talk/internal/models"
	. "talk/internal/models"
	. "talk/internal/ws"
)

type LeaveMessageHandler struct {
	RoomsPool *RoomsPool
}

func (h *LeaveMessageHandler) HandleMessage(client *Client, message Message) {
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

	encodeData, err := json.Marshal([]map[string]interface{}{
		{"peerID": client.Uuid},
	})
	if err != nil {
		// TODO: Отправлять ошибку
		fmt.Println("ошибка при формировании сообщения: %w", err)
		return
	}

	room.Broadcast <- models.Message{
		Type: MessageTypeRemovePeer,
		Data: string(encodeData),
	}

	room.Leave(client)

	client.Hub.ShareRooms()
}
