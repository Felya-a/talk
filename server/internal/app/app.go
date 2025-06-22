package app

import (
	"talk/internal/config"
	authService "talk/internal/services/auth"
	"talk/internal/transport/http"
	"talk/internal/transport/ws"

	"github.com/jmoiron/sqlx"
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
	db *sqlx.DB,
) *App {
	wsServer := ws.New(config.WebSocket.Port, db)
	httpServer := http.New(config.Http.Port, authService)

	return &App{
		WsServer:   wsServer,
		HttpServer: httpServer,
	}
}
