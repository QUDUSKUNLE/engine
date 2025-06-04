package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)
// GetUsers retrieves a paginated list of users from the database.
func (u *Repository) GetUsers(ctx context.Context, params GetUsersParams) ([]*User, error) {
	return u.database.GetUsers(ctx, params)
}

// CreateUser inserts a new user into the database and returns the created user row.
func (u *Repository) CreateUser(ctx context.Context, user CreateUserParams) (*CreateUserRow, error) {
	return u.database.CreateUser(ctx, user)
}

// GetUser fetches a user by their unique ID.
func (u *Repository) GetUser(ctx context.Context, id string) (*User, error) {
	return u.database.GetUser(ctx, id)
}

// GetUserByEmail fetches a user by their email address.
func (u *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	emailText := pgtype.Text{String: email, Valid: true}
	return u.database.GetUserByEmail(ctx, emailText)
}

// UpdateUser updates user fields that are provided in the params and returns the updated user row.
func (u *Repository) UpdateUser(ctx context.Context, user UpdateUserParams) (*UpdateUserRow, error) {
	updatedUser, err := u.database.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}
