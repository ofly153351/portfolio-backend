package auth

import "context"

// Repository defines persistence behavior for the auth module.
type Repository interface {
	Ping(ctx context.Context) error
}
