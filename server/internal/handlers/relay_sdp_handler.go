package handlers

import (
	"encoding/json"
	"talk/internal/core"
	. "talk/internal/lib/logger"
	"talk/internal/models/errors"
	. "talk/internal/models/messages"
	usecase "talk/internal/use-case"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type RelaySdpMessageHandler struct {
	SendIceOrSdp usecase.SendIceOrSdpUseCase
}

func (h *RelaySdpMessageHandler) HandleMessage(client *core.Client, message ReceiveMessage) error {
	var dto RelaySdpMessageDto

	if err := json.Unmarshal([]byte(message.Data), &dto); err != nil {
		Log.Error("[RelaySdpMessageHandler] error on unmarshal message", Log.Err(err))
		return errors.ErrValidation
	}

	if err := validator.New().Struct(dto); err != nil {
		Log.Error("[RelaySdpMessageHandler] validation error", Log.Err(err))
		return errors.ErrValidation
	}

	if err := uuid.Validate(dto.PeerID); err != nil {
		Log.Error("[RelaySdpMessageHandler] uuid validation error", Log.Err(err))
		return errors.ErrValidation
	}

	h.SendIceOrSdp.Execute(client, dto.PeerID, dto.SessionDescription, "sdp")

	return nil
}
