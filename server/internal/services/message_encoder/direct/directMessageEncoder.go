package direct_message_encoder

import (
	"encoding/json"
	"fmt"
	. "talk/internal/lib/logger"
	. "talk/internal/models/messages"
	. "talk/internal/services/message_encoder"
	"time"
)

type DirectMessageEncoder struct{}

func NewDirectMessageEncoder() MessageEncoder {
	return &DirectMessageEncoder{}
}

func (b *DirectMessageEncoder) BuildRemovePeerMessage(dto RemovePeerMessageDto) TransmitMessage {
	messageData := map[string]interface{}{
		"peer_id": dto.ClientUuid,
	}

	return TransmitMessage{
		Type: MessageTypeRemovePeer,
		Data: messageData,
	}
}

func (b *DirectMessageEncoder) BuildShareRoomsMessage(dto ShareRoomsMessageDto) TransmitMessage {
	messageData := make([]map[string]interface{}, len(dto.Rooms))

	for i, room := range dto.Rooms {
		messageData[i] = map[string]interface{}{
			"uuid":    room.Uuid,
			"name":    room.Name,
			"clients": room.Clients,
		}
	}

	return TransmitMessage{
		Type: MessageTypeShareRooms,
		Data: messageData,
	}
}

func (b *DirectMessageEncoder) BuildAddPeerMessage(dto AddPeerMessageDto) TransmitMessage {
	messageData := map[string]interface{}{
		"peer_id":      dto.PeerID,
		"create_offer": dto.CreateOffer,
	}

	return TransmitMessage{
		Type: MessageTypeAddPeer,
		Data: messageData,
	}
}

func (b *DirectMessageEncoder) BuildPingMessage() TransmitMessage {
	return TransmitMessage{
		Type: MessageTypePing,
		Data: fmt.Sprint(time.Now().Unix()),
	}
}

func (b *DirectMessageEncoder) BuildSessionDescriptionMessage(dto SessionDescriptionMessageDto) TransmitMessage {
	var sessionDescription map[string]interface{}
	err := json.Unmarshal([]byte(dto.SessionDescription), &sessionDescription)
	if err != nil {
		// TODO: возвращать ошибку
	}

	messageData := map[string]interface{}{
		"peer_id":             dto.PeerID,
		"session_description": sessionDescription,
	}

	return TransmitMessage{
		Type: MessageTypeSessionDescription,
		Data: messageData,
	}
}

func (b *DirectMessageEncoder) BuildIceCandidateMessage(dto IceCandidateMessageDto) TransmitMessage {
	var iceCandidate map[string]interface{}
	err := json.Unmarshal([]byte(dto.IceCandidate), &iceCandidate)
	if err != nil {
		// TODO: возвращать ошибку
	}

	messageData := map[string]interface{}{
		"peer_id":       dto.PeerID,
		"ice_candidate": iceCandidate,
	}

	return TransmitMessage{
		Type: MessageTypeIceCandidate,
		Data: messageData,
	}
}

func (b *DirectMessageEncoder) BuildClientInfoMessage(dto ClientInfoMessageDto) TransmitMessage {
	messageData := map[string]interface{}{
		"uuid": dto.Uuid,
	}

	return TransmitMessage{
		Type: MessageTypeClientInfo,
		Data: messageData,
	}
}

func (b *DirectMessageEncoder) BuildErrorMessage(dto ErrorMessageDto) TransmitMessage {
	messageData := map[string]interface{}{
		"error": dto.Err.Error(),
	}

	return TransmitMessage{
		Type: MessageTypeError,
		Data: messageData,
	}
}

func (b *DirectMessageEncoder) Encode(message TransmitMessage) (data []byte, err error) {
	encodeMessage, err := json.Marshal(message)
	if err != nil {
		Log.Error("[DirectMessageEncoder] error on marshal message", Log.Err(err))
		return nil, err
	}

	return encodeMessage, nil
}
