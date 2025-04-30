package direct_message_decoder

import (
	"encoding/json"
	. "talk/internal/lib/logger"
	. "talk/internal/models/messages"
)

type DirectMessageDecoder struct{}

func NewDirectMessageDecoder() *DirectMessageDecoder {
	return &DirectMessageDecoder{}
}

func (d *DirectMessageDecoder) Decode(data []byte) (ReceiveMessage, error) {
	var message ReceiveMessage
	if err := json.Unmarshal(data, &message); err != nil {
		Log.Error("[DirectMessageDecoder] invalid message format", LogFields{"error": err, "messageData": data})
		return ReceiveMessage{}, err
	}

	return message, nil
}
