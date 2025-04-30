package handlers

import (
	. "talk/internal/core"
	. "talk/internal/lib/logger"
	. "talk/internal/models/messages"
	usecase "talk/internal/use-case"
)

type LeaveMessageHandler struct {
	LeaveClient usecase.LeaveClientUseCase
}

func (h *LeaveMessageHandler) HandleMessage(client *Client, message ReceiveMessage) error {
	Log.Debug("[LeaveMessageHandler]", LogFields{"clientUuid": client.Uuid})
	h.LeaveClient.Execute(client)

	return nil
}
