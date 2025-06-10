package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/ports"
)

// Ensure Repository implements ports.UserRepository
var _ ports.UserRepository = (*Repository)(nil)

// GetUsers retrieves a paginated list of users from the database.
func (repo *Repository) GetUsers(ctx context.Context, params db.GetUsersParams) ([]*db.User, error) {
	return repo.database.GetUsers(ctx, params)
}

// CreateUser inserts a new user into the database and returns the created user row.
func (repo *Repository) CreateUser(ctx context.Context, user db.CreateUserParams) (*db.CreateUserRow, error) {
	return repo.database.CreateUser(ctx, user)
}

// GetUser fetches a user by their unique ID.
func (repo *Repository) GetUser(ctx context.Context, id string) (*db.User, error) {
	return repo.database.GetUser(ctx, id)
}

// GetUserByEmail fetches a user by their email address.
func (repo *Repository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	emailText := pgtype.Text{String: email, Valid: true}
	return repo.database.GetUserByEmail(ctx, emailText)
}

// UpdateUser updates user fields that are provided in the params and returns the updated user row.
func (repo *Repository) UpdateUser(ctx context.Context, user db.UpdateUserParams) (*db.UpdateUserRow, error) {
	updatedUser, err := repo.database.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}
