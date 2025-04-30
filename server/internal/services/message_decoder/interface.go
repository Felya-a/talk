package message_decoder

import (
	. "talk/internal/models/messages"
)

type MessageDecoder interface {
	Decode(data []byte) (ReceiveMessage, error)
}
