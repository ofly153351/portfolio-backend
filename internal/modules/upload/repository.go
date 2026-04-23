package upload

import "context"

// Repository defines persistence behavior for the upload module.
type Repository interface {
	Ping(ctx context.Context) error
}
