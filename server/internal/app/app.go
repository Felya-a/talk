package app

import (
	"talk/internal/config"
	authService "talk/internal/services/auth"
	"talk/internal/transport/http"
	"talk/internal/transport/ws"
)

// Структура сервера WebSocket
type App struct {
	WsServer   *ws.WsTransport
	HttpServer *http.HttpTransport
}

// Создание нового сервера
func New(
	config config.Config,
	authService *authService.AuthService,
) *App {
	wsServer := ws.New(config.WebSocket.Port)
	httpServer := http.New(config.Http.Port, authService)

	return &App{
		WsServer:   wsServer,
		HttpServer: httpServer,
	}
}
