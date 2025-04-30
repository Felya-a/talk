package core

import (
	"errors"
	. "talk/internal/lib/logger"
	. "talk/internal/models/errors"
	. "talk/internal/models/messages"

	"github.com/google/uuid"
)

type Client struct {
	Uuid       string
	Connection Connection
	Hub        *Hub
	Outbound   chan TransmitMessage
}

func NewClient(
	connection Connection,
	hub *Hub,
) *Client {
	return &Client{
		// TODO: получать uuid из sso
		Uuid:       uuid.New().String(),
		Connection: connection,
		Hub:        hub,
		Outbound:   make(chan TransmitMessage),
	}
}

func (client *Client) Kill() {
	close(client.Outbound)
	client.Connection.Close()
}

// Забирает сообщения из WebSocket-соединения и отправляет их в hub (концентратор).
func (client *Client) ReadPump(router *MessageRouter) {
	defer func() {
		client.Hub.Unregister <- client
	}()

	for {
		message, err := client.Connection.Receive()
		if err != nil {
			if errors.Is(err, ErrCloseConnection) {
				Log.Info("[Client] close connection", LogFields{"clientUuid": client.Uuid})
				return
			}
			Log.Error("[Client] error on read message from connection", LogFields{"clientUuid": client.Uuid, "error": err})
			continue
		}

		if err := client.Hub.HandleMessage(client, message); err != nil {
			client.Connection.Send(TransmitMessage{}, err)
		}
	}
}

// Перенаправляет сообщения из hub обратно в WebSocket-соединение.
func (client *Client) WritePump() {
	for {
		message, ok := <-client.Outbound
		if !ok {
			Log.Info("[Client] client channel was closed", LogFields{"clientUuid": client.Uuid})
			return
		}
		Log.Info("[Client] message to send", LogFields{
			"clientUuid":  client.Uuid,
			"messageType": message.Type,
		})

		client.Connection.Send(message, nil)
	}
}
