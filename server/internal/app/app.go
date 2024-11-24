package app

import (
	"log/slog"
	"talk/internal/models"
	"talk/internal/transport/ws"
)

// Структура сервера WebSocket
type App struct {
	WsServer *ws.WsTransport
}

// Создание нового сервера
func New(
	log *slog.Logger,
	wsPort string,
) *App {
	hub := models.NewHub()
	go hub.Run()
	wsServer := ws.New(log, wsPort, hub)

	return &App{
		WsServer: wsServer,
	}
}
