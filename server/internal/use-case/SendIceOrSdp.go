package usecase

import (
	"errors"
	"talk/internal/core"
	. "talk/internal/lib/logger"
	. "talk/internal/models/messages"
	message_encoder "talk/internal/services/message_encoder"
	rooms_storage "talk/internal/services/rooms_storage"
)

type SendIceOrSdpUseCase struct {
	RoomsStorage     *rooms_storage.RoomsStorage
	ShareRooms       ShareRoomsUseCase
	MessageEncoder   message_encoder.MessageEncoder
	FindClientByUuid FindClientByUuid
}

func (uc *SendIceOrSdpUseCase) Execute(client *core.Client, targetClientUuid string, data []byte, sendMessageType string) error {
	targetClient := uc.FindClientByUuid.Execute(targetClientUuid)
	if targetClient == nil {
		// TODO: Отправлять ошибку
		Log.Warn("client not found", LogFields{"clientUuid": targetClientUuid})
		return errors.New("client not found")
	}

	// TODO: Обрабатывать ошибку
	room, _ := uc.RoomsStorage.FindByClient(client)
	if room == nil {
		// TODO: Отправлять ошибку
		Log.Warn("[SendIceOrSdpUseCase] room not found", LogFields{"clientUuid": client.Uuid})
		return errors.New("room not found")
	}

	if sendMessageType == "ice" {
		targetClient.Outbound <- uc.MessageEncoder.BuildIceCandidateMessage(IceCandidateMessageDto{PeerID: client.Uuid, IceCandidate: data})
	} else if sendMessageType == "sdp" {
		targetClient.Outbound <- uc.MessageEncoder.BuildSessionDescriptionMessage(SessionDescriptionMessageDto{PeerID: client.Uuid, SessionDescription: data})
	} else {
		return errors.New("send message type not defined")
	}

	return nil
}
