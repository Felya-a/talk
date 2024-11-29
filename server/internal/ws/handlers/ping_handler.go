package handlers

import (
	"fmt"
	. "talk/internal/models"
	. "talk/internal/ws"
)

type PingMessageHandler struct {
	Hub *Hub
}

func (h *PingMessageHandler) HandleMessage(client *Client, message ReceiveMessage) {
	fmt.Println("Ping from ", client.Uuid)
	h.Hub.ShareRooms() // DEBUG ONLY
}
