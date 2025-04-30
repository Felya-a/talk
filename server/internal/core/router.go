package core

import (
	"context"
	"fmt"
	. "talk/internal/lib/logger"
	. "talk/internal/models/messages"

	"golang.org/x/sync/errgroup"
)

type MessageRouter struct {
	handlers map[MessageType][]MessageHandler
}

func NewMessageRouter() *MessageRouter {
	return &MessageRouter{
		handlers: make(map[MessageType][]MessageHandler),
	}
}

func (r *MessageRouter) RegisterHandler(messageType MessageType, handler MessageHandler) {
	if r.handlers[messageType] == nil {
		r.handlers[messageType] = make([]MessageHandler, 0)
	}

	r.handlers[messageType] = append(r.handlers[messageType], handler)
}

func (r *MessageRouter) RouteMessage(client *Client, message ReceiveMessage) error {
	handlers, exists := r.handlers[message.Type]
	if !exists {
		Log.Warn("no handler for message type: %s\n", LogFields{"messageType": message.Type})
		return fmt.Errorf("no handler for message type: %s", message.Type)
	}

	g, _ := errgroup.WithContext(context.Background())

	for _, handler := range handlers {
		g.Go(func() error { return handler.HandleMessage(client, message) })
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
