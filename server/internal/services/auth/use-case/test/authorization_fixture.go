package auth_test

import (
	models "talk/internal/services/auth/model"
	ssoFake "talk/internal/services/auth/sso/mock"
)

var FakeUser = &ssoFake.UserWithTokensModel{
	UserModel: models.UserModel{ID: 1, Email: "email@local.com"},
	Tokens: &models.JwtTokens{
		AccessJwtToken:  "valid_access_token",
		RefreshJwtToken: "valid_refresh_token",
	},
	AuthorizationCode: "valid_authorization_code",
}
