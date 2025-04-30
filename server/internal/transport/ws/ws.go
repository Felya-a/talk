package ws

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"talk/internal/adapter"
	"talk/internal/config"
	"talk/internal/core"
	"talk/internal/handlers"
	. "talk/internal/lib/logger"
	. "talk/internal/models/messages"
	message_decoder "talk/internal/services/message_decoder/direct"
	message_encoder "talk/internal/services/message_encoder/direct"
	usecase "talk/internal/use-case"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsTransport struct {
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
	port string,
) *WsTransport {
	setGinMode()
	handler := gin.Default()

	roomsPool := core.NewRoomsPool()
	router := core.NewMessageRouter()
	hub := core.NewHub(router, roomsPool)
	go hub.Run()

	messageEncoder := message_encoder.NewDirectMessageEncoder()
	messageDecoder := message_decoder.NewDirectMessageDecoder()

	// Use-case
	shareRooms := usecase.ShareRoomsUseCase{Hub: hub, MessageEncoder: messageEncoder}
	joinClient := usecase.JoinClientUseCase{Hub: hub, ShareRooms: shareRooms, MessageEncoder: messageEncoder}
	leaveClient := usecase.LeaveClientUseCase{Hub: hub, ShareRooms: shareRooms, MessageEncoder: messageEncoder}
	createRoom := usecase.CreateRoomUseCase{Hub: hub, ShareRooms: shareRooms, MessageEncoder: messageEncoder}
	sendIceOrSdp := usecase.SendIceOrSdpUseCase{Hub: hub, ShareRooms: shareRooms}

	// Регистрация обработчиков сообщений
	router.RegisterHandler(MessageTypePing, &handlers.PingMessageHandler{})
	router.RegisterHandler(MessageTypeJoin, &handlers.JoinMessageHandler{JoinClient: joinClient})
	router.RegisterHandler(MessageTypeLeave, &handlers.LeaveMessageHandler{LeaveClient: leaveClient})
	router.RegisterHandler(MessageTypeCreateRoom, &handlers.CreateRoomMessageHandler{CreateRoom: createRoom, ShareRooms: shareRooms})
	router.RegisterHandler(MessageTypeRelaySdp, &handlers.RelaySdpMessageHandler{SendIceOrSdp: sendIceOrSdp})
	router.RegisterHandler(MessageTypeRelayIce, &handlers.RelayIceMessageHandler{SendIceOrSdp: sendIceOrSdp})

	handler.GET("/ws", func(ctx *gin.Context) {
		// Обновляем HTTP-соединение до WebSocket
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			Log.Warn("error upgrade http to websocket", Log.Err(err))
			return
		}

		webSocketConnection := adapter.NewWebSocketConnection(
			conn,
			messageEncoder,
			messageDecoder,
			shareRooms,
		)

		hub.Register <- webSocketConnection
	})

	wsServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}

	return &WsTransport{
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
	if err := wst.wsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		Log.Error("error on start ws server", Log.Err(err))
		return fmt.Errorf("error on start ws server")
	}
	Log.Info("ws server is running", LogFields{"port": wst.port, "addr": wst.wsServer.Addr})

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
	Log.Info("stopping ws server", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.wsServer.Shutdown(ctx); err != nil {
		Log.Error("Server forced to shutdown: ", Log.Err(err))
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
