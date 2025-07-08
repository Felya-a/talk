package core

import (
	"errors"
	"slices"
	"sync"
	. "talk/internal/lib/logger"
	. "talk/internal/models/messages"
	. "talk/internal/services/message_encoder"
)

type RoomDto struct {
	Uuid string
	Name string
}

type Room struct {
	RoomDto
	Clients []*Client

	mutex          sync.Mutex
	broadcast      chan TransmitMessage
	messageEncoder MessageEncoder
}

func NewRoom(
	dto RoomDto,
	messageEncoder MessageEncoder,
) *Room {
	var room = &Room{
		RoomDto: RoomDto{
			Uuid: dto.Uuid,
			Name: dto.Name,
		},
		Clients:        []*Client{},
		broadcast:      make(chan TransmitMessage),
		messageEncoder: messageEncoder,
	}

	go room.Run()

	return room
}

func (r *Room) Run() {
	go r.serveBroadcast()
}

func (r *Room) Join(client *Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.checkExistUserLocked(client) {
		Log.Warn("client already is located in room", LogFields{"clientUuid": client.Uuid, "roomUuid": r.Uuid})
		return
	}
	r.Clients = append(r.Clients, client)
	Log.Info("client added in room", LogFields{"clientUuid": client.Uuid, "roomUuid": r.Uuid})
}

func (r *Room) Leave(client *Client) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	index := slices.Index(r.Clients, client)
	if index != -1 {
		r.Clients[index] = r.Clients[len(r.Clients)-1]
		r.Clients = r.Clients[:len(r.Clients)-1]
		Log.Info("client excluded from room", LogFields{"clientUuid": client.Uuid, "roomUuid": r.Uuid})
	} else {
		Log.Info("client not found in room", LogFields{"clientUuid": client.Uuid, "roomUuid": r.Uuid})
		return errors.New("client not found in room")
	}

	r.broadcast <- r.messageEncoder.BuildRemovePeerMessage(
		RemovePeerMessageDto{ClientUuid: client.Uuid},
	)

	return nil
}

func (r *Room) serveBroadcast() {
	for message := range r.broadcast {
		Log.Debug("broadcast message in room", LogFields{"roomUuid": r.Uuid, "message": message})
		for _, client := range r.Clients {
			client.Outbound <- message
		}
	}
}

func (r *Room) CheckExistUser(client *Client) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.checkExistUserLocked(client)
}

func (r *Room) checkExistUserLocked(client *Client) bool {
	return slices.Contains(r.Clients, client)
}
