package health

import "context"

// Repository defines persistence behavior for the health module.
type Repository interface {
	Ping(ctx context.Context) error
}
