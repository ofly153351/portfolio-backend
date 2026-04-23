package upload

import "errors"

var (
	ErrFileRequired    = errors.New("file_required")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrUploadFailed    = errors.New("upload_failed")
	ErrInvalidFileType = errors.New("invalid_file_type")
)
