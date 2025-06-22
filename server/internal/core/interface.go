package core

import (
	messages "talk/internal/models/messages"
)

type Connection interface {
	Receive() (messages.ReceiveMessage, error)
	Send(message messages.TransmitMessage, err error)
	Close()
}

type MessageHandler interface {
	HandleMessage(client *Client, message messages.ReceiveMessage) error
}

type RoomsStorage interface {
	FindByClient(client *Client) (*Room, error)
}
