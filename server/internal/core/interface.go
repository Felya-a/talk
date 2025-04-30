package core

import (
	. "talk/internal/models/messages"
)

type Connection interface {
	Receive() (ReceiveMessage, error)
	Send(message TransmitMessage, err error)
	Close()
}

type MessageHandler interface {
	HandleMessage(client *Client, message ReceiveMessage) error
}
