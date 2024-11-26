package handlers

import (
	"fmt"
	. "talk/internal/models"
	. "talk/internal/ws"
)

type PingMessageHandler struct{}

func (h *PingMessageHandler) HandleMessage(client *Client, message Message) {
	fmt.Println("PING HANDLER")
	fmt.Println(message)
}
