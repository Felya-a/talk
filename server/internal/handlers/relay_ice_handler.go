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

type RelayIceMessageHandler struct {
	SendIceOrSdp usecase.SendIceOrSdpUseCase
}

func (h *RelayIceMessageHandler) HandleMessage(client *core.Client, message ReceiveMessage) error {
	var dto RelayIceMessageDto

	if err := json.Unmarshal([]byte(message.Data), &dto); err != nil {
		Log.Error("[RelayIceMessageHandler] error on unmarshal message data", Log.Err(err))
		return errors.ErrValidation
	}

	if err := validator.New().Struct(dto); err != nil {
		Log.Error("[RelayIceMessageHandler] validation error", Log.Err(err))
		return errors.ErrValidation
	}

	if err := uuid.Validate(dto.PeerID); err != nil {
		Log.Error("[RelayIceMessageHandler] uuid validation error", Log.Err(err))
		return errors.ErrValidation
	}

	h.SendIceOrSdp.Execute(client, dto.PeerID, dto.IceCandidate, "ice")

	return nil
}
