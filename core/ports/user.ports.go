package ports

import (
	"context"

	"github.com/medicue/adapters/db"
)

type UserRepository interface {
	GetUsers(ctx context.Context, params db.GetUsersParams) ([]*db.User, error)
	CreateUser(ctx context.Context, user db.CreateUserParams) (*db.CreateUserRow, error)
	GetUser(ctx context.Context, id string) (*db.User, error)
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	UpdateUser(ctx context.Context, user db.UpdateUserParams) (*db.UpdateUserRow, error)
}
