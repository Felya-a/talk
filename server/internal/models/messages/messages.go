package messages

import "encoding/json"

type MessageType string

const (
	// incoming
	MessageTypeCreateRoom MessageType = "create-room"
	MessageTypeJoin       MessageType = "join"
	MessageTypeLeave      MessageType = "leave"
	MessageTypeRelaySdp   MessageType = "relay-sdp"
	MessageTypeRelayIce   MessageType = "relay-ice"
	MessageTypePong       MessageType = "pong"

	// outgoing
	MessageTypePing               MessageType = "ping"
	MessageTypeAddPeer            MessageType = "add-peer"
	MessageTypeShareRooms         MessageType = "share-rooms"
	MessageTypeRemovePeer         MessageType = "remove-peer"
	MessageTypeIceCandidate       MessageType = "ice-candidate"
	MessageTypeSessionDescription MessageType = "session-description"
	MessageTypeError              MessageType = "error"

	// bidirectional
)

type ReceiveMessage struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type TransmitMessage struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}
