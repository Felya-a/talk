package message_encoder

import (
	. "talk/internal/models/messages"
)

type MessageEncoder interface {
	BuildPingMessage() TransmitMessage
	BuildAddPeerMessage(dto AddPeerMessageDto) TransmitMessage
	BuildShareRoomsMessage(dto ShareRoomsMessageDto) TransmitMessage
	BuildRemovePeerMessage(dto RemovePeerMessageDto) TransmitMessage
	BuildIceCandidateMessage(dto IceCandidateMessageDto) TransmitMessage
	BuildSessionDescriptionMessage(dto SessionDescriptionMessageDto) TransmitMessage
	BuildErrorMessage(dto ErrorMessageDto) TransmitMessage
	Encode(message TransmitMessage) (data []byte, err error)
}
