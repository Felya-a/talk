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

type RelayIceMessageDto struct {
	PeerID       string          `json:"peer_id" validate:"required"`
	IceCandidate json.RawMessage `json:"ice_candidate" validate:"required"`
}

type RelayIceMessageHandler struct {
	RoomsPool *RoomsPool
	Hub       *Hub
}

func (h *RelayIceMessageHandler) HandleMessage(client *Client, message ReceiveMessage) {
	var dto RelayIceMessageDto

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

	messageData := map[string]interface{}{
		"peer_id":       client.Uuid,
		"ice_candidate": dto.IceCandidate,
	}

	targetClient.Send <- models.TransmitMessage{
		Type: MessageTypeIceCandidate,
		Data: messageData,
	}
}
