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

type JoinMessageDto struct {
	RoomUuid string `json:"room_uuid" validate:"required"`
}

type JoinMessageHandler struct {
	RoomsPool *RoomsPool
}

func (h *JoinMessageHandler) HandleMessage(client *Client, message Message) {
	var dto JoinMessageDto

	if err := json.Unmarshal([]byte(message.Data), &dto); err != nil {
		fmt.Println("Ошибка при десериализации запроса ", sl.Err(err))
		return
	}

	if err := validator.New().Struct(dto); err != nil {
		fmt.Println("validation error", sl.Err(err))
		return
	}

	if err := uuid.Validate(dto.RoomUuid); err != nil {
		fmt.Println("Получен не uuid", err)
		return
	}

	room := h.RoomsPool.FindByUuid(uuid.MustParse(dto.RoomUuid))
	if room == nil {
		// TODO: Отправлять ошибку
		fmt.Println("Комната ", dto.RoomUuid, " не найдена")
		return
	}

	if h.RoomsPool.FindByClientId(client.Uuid) != nil {
		// TODO: Отправлять ошибку
		fmt.Println("Пользователь уже находится в этой комнате")
		return
	}

	for _, roomClient := range room.Clients {
		// Отправка сообщения новому пользователю о существующих
		if err := sendMessage(roomClient, client.Uuid, false); err != nil {
			fmt.Println("Ошибка:", sl.Err(err))
			return
		}

		// Отправка сообщения существующему пользователю о новом пользователе
		if err := sendMessage(client, roomClient.Uuid, true); err != nil {
			fmt.Println("Ошибка:", sl.Err(err))
			return
		}
	}

	room.Join(client)
	client.Hub.ShareRooms()
}

func sendMessage(receiver *Client, peerID uuid.UUID, createOffer bool) error {
	dataForMessage := []map[string]interface{}{
		{
			"peerID":      peerID,
			"createOffer": createOffer,
		},
	}

	encodeData, err := json.Marshal(dataForMessage)
	if err != nil {
		return fmt.Errorf("ошибка при формировании сообщения: %w", err)
	}

	message := models.Message{
		Type: MessageTypeAddPeer,
		Data: string(encodeData),
	}

	receiver.Send <- message
	return nil
}
