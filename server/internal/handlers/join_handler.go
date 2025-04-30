package handlers

import (
	"encoding/json"
	. "talk/internal/core"
	. "talk/internal/lib/logger"
	"talk/internal/models/errors"
	. "talk/internal/models/messages"
	usecase "talk/internal/use-case"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type JoinMessageHandler struct {
	JoinClient usecase.JoinClientUseCase
}

func (h *JoinMessageHandler) HandleMessage(client *Client, message ReceiveMessage) error {
	var dto JoinMessageDto

	if err := json.Unmarshal([]byte(message.Data), &dto); err != nil {
		Log.Error("[JoinMessageHandler] error on unmarshal data", Log.Err(err))
		return errors.ErrValidation
	}

	if err := validator.New().Struct(dto); err != nil {
		Log.Error("[JoinMessageHandler] validation error", Log.Err(err))
		return errors.ErrValidation
	}

	if err := uuid.Validate(dto.RoomUuid); err != nil {
		Log.Error("[JoinMessageHandler] uuid is not valid", Log.Err(err))
		return errors.ErrValidation
	}

	h.JoinClient.Execute(client, dto)

	return nil
}
