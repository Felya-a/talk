package handlers

import (
	"fmt"
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

	room.Leave(client)

	client.Hub.ShareRooms()
}
