package handlers

import (
	. "talk/internal/core"
	. "talk/internal/lib/logger"
	. "talk/internal/models/messages"
)

type PingMessageHandler struct{}

func (h *PingMessageHandler) HandleMessage(client *Client, message ReceiveMessage) error {
	Log.Info("ping", LogFields{"clientUuid": client.Uuid})

	return nil
}
