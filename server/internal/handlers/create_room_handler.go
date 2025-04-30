package handlers

import (
	"encoding/json"
	. "talk/internal/core"
	. "talk/internal/lib/logger"
	"talk/internal/models/errors"
	. "talk/internal/models/messages"
	usecase "talk/internal/use-case"

	"github.com/go-playground/validator"
)

type CreateRoomMessageHandler struct {
	CreateRoom usecase.CreateRoomUseCase
	ShareRooms usecase.ShareRoomsUseCase
}

func (h *CreateRoomMessageHandler) HandleMessage(client *Client, message ReceiveMessage) error {
	var dto CreateRoomMessageDto

	if err := json.Unmarshal([]byte(message.Data), &dto); err != nil {
		Log.Error("[CreateRoomMessageHandler] error on unmarshal message data", Log.Err(err))
		return errors.ErrValidation
	}

	if err := validator.New().Struct(dto); err != nil {
		Log.Error("[CreateRoomMessageHandler] validation error", Log.Err(err))
		return errors.ErrValidation
	}

	h.CreateRoom.Execute(dto.RoomName)

	return nil
}
