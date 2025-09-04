package ports

import (
	"context"

	"github.com/diagnoxix/adapters/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository interface {
	// User operations
	GetUser(
		ctx context.Context,
		arg string,
	) (*db.User, error)

	CreateUser(
		ctx context.Context,
		arg db.CreateUserParams,
	) (*db.User, error)

	GetUserByEmail(
		ctx context.Context,
		email pgtype.Text,
	) (*db.User, error)

	GetUsers(
		ctx context.Context,
		arg db.GetUsersParams,
	) ([]*db.User, error)

	UpdateUser(
		ctx context.Context,
		arg db.UpdateUserParams,
	) (*db.User, error)

	ListManagersByadmin(
		ctx context.Context,
		arg db.ListUsersByAdminParams,
	) ([]*db.ListUsersByAdminRow, error)

	// Password reset operations
	CreatePasswordResetToken(
		ctx context.Context,
		arg db.CreatePasswordResetTokenParams,
	) error

	GetPasswordResetToken(
		ctx context.Context,
		arg string,
	) (*db.PasswordResetToken, error)

	MarkResetTokenUsed(
		ctx context.Context,
		arg string,
	) error

	UpdateUserPassword(
		ctx context.Context,
		arg db.UpdateUserPasswordParams,
	) error

	CreateEmailVerificationToken(
		ctx context.Context,
		arg db.CreateEmailVerificationTokenParams,
	) (*db.EmailVerificationToken, error)

	GetEmailVerificationToken(
		ctx context.Context,
		rg string,
	) (*db.EmailVerificationToken, error)

	MarkEmailVerificationTokenUsed(
		ctx context.Context,
		arg string,
	) error

	MarkEmailAsVerified(
		ctx context.Context,
		email string,
	) error
}
