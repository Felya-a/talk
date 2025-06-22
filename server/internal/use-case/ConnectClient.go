package usecase

import (
	"talk/internal/core"
	// . "talk/internal/lib/logger"
	auth_service "talk/internal/services/auth"
)

type ConnectClientUseCase struct {
	Hub         *core.Hub
	AuthService *auth_service.AuthService
}

func (uc *ConnectClientUseCase) Execute(connection core.Connection, accessToken string) {
	// if accessToken == "" {
	// 	Log.Info("[ConnectClientUseCase] access token is empty", nil)
	// }

	// _, _, _, err := uc.AuthService.Auth(context.Background(), accessToken, "", "")
	// if err != nil {
	// 	// TODO: Возвращать ошибку
	// 	Log.Error("[ConnectClientUseCase] error on auth", Log.Err(err))
	// 	return
	// }

	uc.Hub.Register <- connection
}
