package http_handlers_v1

// type UserInfoResponseDto struct {
// 	ID    int64  `json:"id"`
// 	Email string `json:"email"`
// }

type AuthRequestDto struct {
	AccessToken       string
	RefreshToken      string
	AuthorizationCode string `json:"authorization_code"`
}

type AuthResponseDto struct {
	RedirectUrl string `json:"redirect_url"`
	UserData    `json:"user"`
}

type UserData struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
}

// type LoginRequestDto struct {
// 	Email    string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required,min=6,max=20"`
// }

// type LoginResponseDto struct {
// 	AuthorizationCode string `json:"authorization_code"`
// }

// type RegistrationRequestDto struct {
// 	Email    string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required,min=6,max=20"`
// }

// type RegistrationResponseDto struct {
// 	UserId int64 `json:"userid"`
// }

// type TokenRequestDto struct {
// 	AuthorizationCode string `json:"authorization_code" validate:"required"`
// }

// type TokenResponseDto struct {
// 	AccessToken  string `json:"access_token"`
// 	RefreshToken string `json:"refresh_token"`
// }

// type RefreshResponseDto struct {
// 	AccessToken  string `json:"access_token"`
// 	RefreshToken string `json:"refresh_token"`
// }
