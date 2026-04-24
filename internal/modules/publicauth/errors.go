package publicauth

import "errors"

var (
	ErrTokenMissing = errors.New("public token missing")
	ErrTokenInvalid = errors.New("public token invalid")
	ErrTokenExpired = errors.New("public token expired")
)

