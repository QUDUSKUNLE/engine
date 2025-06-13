package ports

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/medicue/adapters/db"
)

type UserRepository interface {
	// User operations
	GetUser(ctx context.Context, id string) (*db.User, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (*db.User, error)
	GetUserByEmail(ctx context.Context, email pgtype.Text) (*db.User, error)
	GetUsers(ctx context.Context, arg db.GetUsersParams) ([]*db.User, error)
	UpdateUser(ctx context.Context, arg db.UpdateUserParams) (*db.User, error)

	// Password reset operations
	CreatePasswordResetToken(ctx context.Context, arg db.CreatePasswordResetTokenParams) error
	GetPasswordResetToken(ctx context.Context, token string) (*db.PasswordResetToken, error)
	MarkResetTokenUsed(ctx context.Context, id string) error
	UpdateUserPassword(ctx context.Context, arg db.UpdateUserPasswordParams) error

		CreateEmailVerificationToken(ctx context.Context, arg db.CreateEmailVerificationTokenParams) (db.EmailVerificationToken, error)
	GetEmailVerificationToken(ctx context.Context, arg db.GetEmailVerificationTokenParams) (db.EmailVerificationToken, error)
	MarkEmailVerificationTokenUsed(ctx context.Context, arg db.MarkEmailVerificationTokenUsedParams) error
	MarkEmailAsVerified(ctx context.Context, email string) error
}
