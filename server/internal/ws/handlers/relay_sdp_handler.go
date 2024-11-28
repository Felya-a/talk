package handlers

import (
	"encoding/json"
	"fmt"
	"talk/internal/lib/logger/sl"
	"talk/internal/models"
	. "talk/internal/models"
	. "talk/internal/ws"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type RelaySdpMessageDto struct {
	PeerID             string `json:"peer_id" validate:"required"`
	SessionDescription string `json:"session_description" validate:"required"`
}

type RelaySdpMessageHandler struct {
	RoomsPool *RoomsPool
	Hub       *Hub
}

func (h *RelaySdpMessageHandler) HandleMessage(client *Client, message Message) {
	var dto RelaySdpMessageDto

	if err := json.Unmarshal([]byte(message.Data), &dto); err != nil {
		fmt.Println("Ошибка при десериализации запроса ", sl.Err(err))
		return
	}

	if err := validator.New().Struct(dto); err != nil {
		fmt.Println("validation error", sl.Err(err))
		return
	}

	if err := uuid.Validate(dto.PeerID); err != nil {
		fmt.Println("Получен не uuid", err)
		return
	}

	room := h.RoomsPool.FindByClientId(client.Uuid)
	if room == nil {
		fmt.Println("Комната не найдена")
		return
	}

	targetClient := h.Hub.FindClientByUuid(uuid.MustParse(dto.PeerID))
	if !room.CheckExistUser(targetClient) {
		fmt.Println("Пользователи не состоят в одной комнате")
		return
	}

	if targetClient == nil {
		fmt.Println("Пользователь не найден")
		return
	}

	targetClient.Send <- models.Message{
		Type: MessageTypeSessionDescription,
		Data: dto.SessionDescription,
	}
}
