package messages

import "encoding/json"

// incoming
type CreateRoomMessageDto struct {
	RoomName string `json:"room_name" validate:"required,min=3,max=20"`
}

type JoinMessageDto struct {
	RoomUuid string `json:"room_uuid" validate:"required"`
}

type RelayIceMessageDto struct {
	PeerID       string          `json:"peer_id" validate:"required"`
	IceCandidate json.RawMessage `json:"ice_candidate" validate:"required"`
}

type RelaySdpMessageDto struct {
	PeerID             string          `json:"peer_id" validate:"required"`
	SessionDescription json.RawMessage `json:"session_description" validate:"required"`
}

// outgoing
type RemovePeerMessageDto struct {
	ClientUuid string
}

type PingMessageDto struct{}

type ClientInfoMessageDto struct {
	Uuid string
}

type AddPeerMessageDto struct {
	PeerID      string
	CreateOffer bool
}

type ShareRoomsMessageDto struct {
	Rooms []RoomInfo
}

type RoomInfo struct {
	Uuid    string
	Name    string
	Clients []string
}

type IceCandidateMessageDto struct {
	PeerID       string
	IceCandidate []byte
}

type SessionDescriptionMessageDto struct {
	PeerID             string
	SessionDescription []byte
}

type ErrorMessageDto struct {
	Err error
}

// bidirectional
