package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   uuid.UUID
	Hub  *Hub
	Conn *websocket.Conn
	Send chan (Message)
}
