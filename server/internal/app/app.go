package app

import (
	"log/slog"
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
	wsServer := ws.New(log, wsPort)

	return &App{
		WsServer: wsServer,
	}
}
