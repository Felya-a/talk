package models

import "encoding/json"

type MessageType string

const (
	MessageTypePing               MessageType = "ping"
	MessageTypePong               MessageType = "pong"
	MessageTypeCreateRoom         MessageType = "create-room"
	MessageTypeJoin               MessageType = "join"
	MessageTypeLeave              MessageType = "leave"
	MessageTypeShareRooms         MessageType = "share-rooms"
	MessageTypeAddPeer            MessageType = "add-peer"
	MessageTypeRemovePeer         MessageType = "remove-peer"
	MessageTypeRelaySdp           MessageType = "relay-sdp"
	MessageTypeRelayIce           MessageType = "relay-ice"
	MessageTypeIceCandidate       MessageType = "ice-candidate"
	MessageTypeSessionDescription MessageType = "session-description"
)

type ReceiveMessage struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type TransmitMessage struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}
