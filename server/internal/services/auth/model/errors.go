package auth_service

import (
	"errors"
)

// Было скопировано с sso
var (
	ErrUnathorizated = errors.New("user not authorized")
	// ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserAlreadyExists  = errors.New("user already exists")
	// ErrInvalidJwt         = errors.New("jwt is invalid")
	// ErrJwtExpired         = errors.New("token is expired")
)

var predefinedErrors = map[string]error{
	"ErrUnathorizated": errors.New("user not authorized"),
	// "ErrUserNotFound":       ErrUserNotFound,
	// "ErrInvalidCredentials": ErrInvalidCredentials,
	// "ErrUserAlreadyExists":  ErrUserAlreadyExists,
	// "ErrInvalidJwt":         ErrInvalidJwt,
	// "ErrJwtExpired":         ErrJwtExpired,
}

// IsDefinedError проверяет, принадлежит ли ошибка к предустановленным
func IsDefinedError(err error) bool {
	for _, predefinedError := range predefinedErrors {
		if errors.Is(err, predefinedError) {
			return true
		}
	}
	return false
}
