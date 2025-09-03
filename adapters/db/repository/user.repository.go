package repository

import (
	"context"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/ports"
	"github.com/jackc/pgx/v5/pgtype"
)

// Ensure Repository implements ports.UserRepository
var _ ports.UserRepository = (*Repository)(nil)

func (repo *Repository) GetUsers(
	ctx context.Context,
	params db.GetUsersParams,
) ([]*db.User, error) {
	return repo.database.GetUsers(ctx, params)
}

func (repo *Repository) CreateUser(
	ctx context.Context,
	user db.CreateUserParams,
) (*db.User, error) {
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

func (repo *Repository) GetUser(
	ctx context.Context,
	id string,
) (*db.User, error) {
	return repo.database.GetUser(ctx, id)
}

func (repo *Repository) GetUserByEmail(
	ctx context.Context,
	email pgtype.Text,
) (*db.User, error) {
	response, err := repo.database.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (repo *Repository) ListManagersByadmin(
	ctx context.Context,
	arg db.ListUsersByAdminParams,
) ([]*db.ListUsersByAdminRow, error) {
	response, err := repo.database.ListUsersByAdmin(ctx, arg)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (repo *Repository) UpdateUser(
	ctx context.Context,
	user db.UpdateUserParams,
) (*db.User, error) {
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

func (repo *Repository) CreatePasswordResetToken(
	ctx context.Context,
	params db.CreatePasswordResetTokenParams,
) error {
	return repo.database.CreatePasswordResetToken(ctx, params)
}

func (repo *Repository) GetPasswordResetToken(
	ctx context.Context,
	token string,
) (*db.PasswordResetToken, error) {
	return repo.database.GetPasswordResetToken(ctx, token)
}

func (repo *Repository) MarkResetTokenUsed(
	ctx context.Context,
	token string,
) error {
	return repo.database.MarkResetTokenUsed(ctx, token)
}

func (repo *Repository) UpdateUserPassword(
	ctx context.Context,
	params db.UpdateUserPasswordParams,
) error {
	return repo.database.UpdateUserPassword(ctx, params)
}

func (repo *Repository) CreateEmailVerificationToken(
	ctx context.Context,
	arg db.CreateEmailVerificationTokenParams,
) (*db.EmailVerificationToken, error) {
	token, err := repo.database.CreateEmailVerificationToken(ctx, arg)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (repo *Repository) GetEmailVerificationToken(
	ctx context.Context,
	arg string,
) (*db.EmailVerificationToken, error) {
	token, err := repo.database.GetEmailVerificationToken(ctx, arg)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (repo *Repository) MarkEmailVerificationTokenUsed(
	ctx context.Context,
	arg string,
) error {
	return repo.database.MarkEmailVerificationTokenUsed(ctx, arg)
}

func (repo *Repository) MarkEmailAsVerified(
	ctx context.Context,
	email string,
) error {
	return repo.database.MarkEmailAsVerified(ctx, pgtype.Text{String: email, Valid: true})
}
