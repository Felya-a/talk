package adapter

import (
	"io"
	"sync"
	"time"

	. "talk/internal/lib/logger"
	. "talk/internal/models/errors"
	. "talk/internal/models/messages"
	. "talk/internal/services/message_decoder"
	. "talk/internal/services/message_encoder"
	usecase "talk/internal/use-case"

	"github.com/gorilla/websocket"
)

var closeCodes = []int{1000, 1001, 1002, 1003, 1005, 1006, 1007, 1008, 1009, 1010, 1011, 1012, 1013, 1015}

type WebSocketConnection struct {
	mutex          sync.Mutex
	conn           *websocket.Conn
	messageEncoder MessageEncoder
	messageDecoder MessageDecoder
	closeCh        chan struct{}
	shareRooms     usecase.ShareRoomsUseCase
}

func NewWebSocketConnection(
	conn *websocket.Conn,
	messageEncoder MessageEncoder,
	messageDecoder MessageDecoder,
	shareRooms usecase.ShareRoomsUseCase,
) *WebSocketConnection {
	connection := &WebSocketConnection{
		conn:           conn,
		messageEncoder: messageEncoder,
		messageDecoder: messageDecoder,
		closeCh:        make(chan struct{}),
		shareRooms:     shareRooms,
	}

	go connection.Ping()

	return connection
}

func (w *WebSocketConnection) Receive() (ReceiveMessage, error) {
	_, messageBytes, err := w.conn.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, closeCodes...) {
			Log.Info("[WebSocketConnection] close connection", LogFields{"closeMessage": err.Error()})
			return ReceiveMessage{}, ErrCloseConnection
		}
		if err == io.EOF {
			Log.Warn("[WebSocketConnection] connection is close", Log.Err(err))
		}
		return ReceiveMessage{}, err
	}

	message, err := w.messageDecoder.Decode(messageBytes)
	if err != nil {
		Log.Error("[WebSocketConnection] invalid message format", LogFields{"error": err, "messageBytes": messageBytes})
		return ReceiveMessage{}, err
	}

	return message, nil
}

func (w *WebSocketConnection) Send(message TransmitMessage, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if err != nil {
		message = w.messageEncoder.BuildErrorMessage(ErrorMessageDto{Err: err})
	}

	encodeMessage, err := w.messageEncoder.Encode(message)
	if err != nil {
		Log.Error("[WebSocketConnection] error on marshal message", Log.Err(err))
		return
	}

	err = w.conn.WriteMessage(websocket.TextMessage, encodeMessage)
	if err != nil {
		Log.Error("[WebSocketConnection] error on send message", Log.Err(err))
	}
}

func (w *WebSocketConnection) Close() {
	select {
	// Насколько помню нужно для предотвращения отправки пинга отключенному пользователю
	case w.closeCh <- struct{}{}:
		close(w.closeCh) // Освобождение памяти. Не уверен что нужно, но пусть будет
	default:
	}

	err := w.conn.Close()
	if err != nil {
		Log.Error("[WebSocketConnection] error on close connection", Log.Err(err))
	}

	w.shareRooms.Execute()
}

func (w *WebSocketConnection) Ping() {
	w.Send(w.messageEncoder.BuildPingMessage(), nil)

	ticker := time.NewTicker(time.Second * 15)
	for {
		select {
		case <-w.closeCh:
			return
		case <-ticker.C:
			w.Send(w.messageEncoder.BuildPingMessage(), nil)
		}
	}
}
