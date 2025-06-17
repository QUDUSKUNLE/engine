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
func (repo *Repository) CreateUser(ctx context.Context, user db.CreateUserParams) (*db.User, error) {
	createdUser, err := repo.database.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &db.User{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}, nil
}

// GetUser fetches a user by their unique ID.
func (repo *Repository) GetUser(ctx context.Context, id string) (*db.User, error) {
	return repo.database.GetUser(ctx, id)
}

// GetUserByEmail fetches a user by their email address.
func (repo *Repository) GetUserByEmail(ctx context.Context, email pgtype.Text) (*db.User, error) {
	response, err := repo.database.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// UpdateUser updates user fields that are provided in the params and returns the updated user row.
func (repo *Repository) UpdateUser(ctx context.Context, user db.UpdateUserParams) (*db.User, error) {
	updatedUser, err := repo.database.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &db.User{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}

// CreatePasswordResetToken creates a password reset token for a user.
func (repo *Repository) CreatePasswordResetToken(ctx context.Context, params db.CreatePasswordResetTokenParams) error {
	err := repo.database.CreatePasswordResetToken(ctx, params)
	return err
}

// GetPasswordResetToken retrieves a password reset token for a user.
func (repo *Repository) GetPasswordResetToken(ctx context.Context, token string) (*db.PasswordResetToken, error) {
	return repo.database.GetPasswordResetToken(ctx, token)
}

// MarkResetTokenUsed marks a password reset token as used.
func (repo *Repository) MarkResetTokenUsed(ctx context.Context, token string) error {
	return repo.database.MarkResetTokenUsed(ctx, token)
}

// UpdateUserPassword updates the password for a user.
func (repo *Repository) UpdateUserPassword(ctx context.Context, params db.UpdateUserPasswordParams) error {
	return repo.database.UpdateUserPassword(ctx, params)
}

// Email verification methods
func (repo *Repository) CreateEmailVerificationToken(ctx context.Context, arg db.CreateEmailVerificationTokenParams) (*db.EmailVerificationToken, error) {
	token, err := repo.database.CreateEmailVerificationToken(ctx, arg)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (repo *Repository) GetEmailVerificationToken(ctx context.Context, arg db.GetEmailVerificationTokenParams) (*db.EmailVerificationToken, error) {
	token, err := repo.database.GetEmailVerificationToken(ctx, arg)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (repo *Repository) MarkEmailVerificationTokenUsed(ctx context.Context, arg db.MarkEmailVerificationTokenUsedParams) error {
	return repo.database.MarkEmailVerificationTokenUsed(ctx, arg)
}

func (repo *Repository) MarkEmailAsVerified(ctx context.Context, email string) error {
	return repo.database.MarkEmailAsVerified(ctx, pgtype.Text{String: email, Valid: true})
}
