package ws

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"talk/internal/config"
	"talk/internal/lib/logger/sl"
	"talk/internal/models"
	ws "talk/internal/ws"
	"talk/internal/ws/handlers"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsTransport struct {
	log      *slog.Logger
	wsServer *http.Server
	port     string
}

// Настройка WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Разрешаем все подключения (можно уточнить для безопасности)
		return true
	},
}

func New(
	log *slog.Logger,
	port string,
) *WsTransport {
	setGinMode()
	handler := gin.Default()

	roomsPool := ws.NewRoomsPool()

	hub := ws.NewHub(roomsPool)
	go hub.Run()

	router := ws.NewMessageRouter()
	router.RegisterHandler(models.MessageTypePing, &handlers.PingMessageHandler{})
	router.RegisterHandler(models.MessageTypeJoin, &handlers.JoinMessageHandler{RoomsPool: roomsPool})
	router.RegisterHandler(models.MessageTypeLeave, &handlers.LeaveMessageHandler{RoomsPool: roomsPool})
	router.RegisterHandler(models.MessageTypeCreateRoom, &handlers.CreateRoomMessageHandler{RoomsPool: roomsPool})

	handler.GET("/ws", func(ctx *gin.Context) {
		// Обновляем HTTP-соединение до WebSocket
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Warn("Ошибка обновления до WebSocket: %v", err)
			return
		}

		client := ws.NewClient(conn, hub)

		hub.Register <- client

		go client.WritePump()
		go client.ReadPump(router)
	})

	wsServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}

	return &WsTransport{
		log,
		wsServer,
		port,
	}
}

func (wst *WsTransport) MustRun() {
	if err := wst.run(); err != nil {
		panic(err)
	}
}

func (wst *WsTransport) run() error {
	const op = "ws.run"

	log := wst.log.With(
		slog.String("op", op),
		slog.String("port", wst.port),
	)

	log.Info("ws server is running", slog.String("addr", wst.wsServer.Addr))
	if err := wst.wsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("error on start ws server", sl.Err(err))
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func setGinMode() {
	var mode string

	switch config.Get().Env {
	case "local", "test":
		mode = "debug"
	case "stage", "prod":
		mode = "release"
	default:
		mode = "release"
	}

	gin.SetMode(mode)
}

func (a *WsTransport) Stop() {
	a.log.Info("stopping ws server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.wsServer.Shutdown(ctx); err != nil {
		a.log.Error("Server forced to shutdown: ", sl.Err(err))
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
