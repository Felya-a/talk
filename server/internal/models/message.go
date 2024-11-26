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

type Message struct {
	Type MessageType `json:"type"`
	Data string      `json:"data"`
}

func (msg *Message) ToJson() ([]byte, error) {
	json, err := json.Marshal(msg)
	return json, err
}
