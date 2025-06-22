package usecase

import (
	core "talk/internal/core"
)

type FindClientByUuid struct {
	Hub *core.Hub
}

func (uc *FindClientByUuid) Execute(uuid string) *core.Client {

	for _, client := range uc.Hub.Clients {
		if client.Uuid == uuid {
			return client
		}
	}
	return nil
}
