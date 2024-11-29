package ws

import (
	"fmt"
	. "talk/internal/models"
)

type MessageHandler interface {
	HandleMessage(client *Client, message ReceiveMessage)
}

type MessageRouter struct {
	handlers map[MessageType]MessageHandler
}

func NewMessageRouter() *MessageRouter {
	return &MessageRouter{
		handlers: make(map[MessageType]MessageHandler),
	}
}

func (r *MessageRouter) RegisterHandler(messageType MessageType, handler MessageHandler) {
	r.handlers[messageType] = handler
}

func (r *MessageRouter) RouteMessage(client *Client, message ReceiveMessage) error {
	handler, exists := r.handlers[message.Type]
	if !exists {
		fmt.Printf("no handler for message type: %s\n", message.Type)
		return fmt.Errorf("no handler for message type: %s", message.Type)
	}
	handler.HandleMessage(client, message)
	return nil
}
