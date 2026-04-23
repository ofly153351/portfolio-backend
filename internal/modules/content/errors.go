package content

import "errors"

var (
	ErrInvalidLocale   = errors.New("invalid_locale")
	ErrVersionConflict = errors.New("version_conflict")
	ErrInvalidPayload  = errors.New("invalid_payload")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrNotFound        = errors.New("not_found")
)
