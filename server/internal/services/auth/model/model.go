package auth_service

type JwtTokens struct {
	AccessJwtToken  string
	RefreshJwtToken string
}

type UserModel struct {
	ID    int64
	Email string
}
