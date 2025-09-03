package repository

import (
	"context"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/ports"
)

// Ensure Repository implements TestPriceRepository
var _ ports.TestPriceRepository = (*Repository)(nil)

func (repo *Repository) CreateTestPrice(
	ctx context.Context,
	price db.Create_Test_PriceParams,
) ([]*db.DiagnosticCentreTestPrice, error) {
	return repo.database.Create_Test_Price(ctx, price)
}
