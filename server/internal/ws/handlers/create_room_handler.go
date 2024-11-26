package handlers

import (
	"encoding/json"
	"fmt"
	"talk/internal/lib/logger/sl"
	. "talk/internal/models"
	. "talk/internal/ws"

	"github.com/go-playground/validator"
)

type CreateRoomMessageDto struct {
	RoomName string `json:"room_name" validate:"required,min=3,max=20"`
}

type CreateRoomMessageHandler struct {
	RoomsPool *RoomsPool
}

func (h *CreateRoomMessageHandler) HandleMessage(client *Client, message Message) {
	var dto CreateRoomMessageDto

	if err := json.Unmarshal([]byte(message.Data), &dto); err != nil {
		fmt.Println("Ошибка при десериализации запроса ", sl.Err(err))
		return
	}

	if err := validator.New().Struct(dto); err != nil {
		fmt.Println("validation error", sl.Err(err))
		return
	}

	createdRoom := NewRoom(dto.RoomName)

	h.RoomsPool.Add(createdRoom)

	// Рассылка комнат
	client.Hub.ShareRooms()
}
