package events

type ClientConnectedEvent struct {
	ClientUuid string
}

func (e ClientConnectedEvent) Name() string {
	return "client.connected"
}
