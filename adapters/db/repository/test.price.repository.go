package repository

import (
	"context"

	"github.com/medivue/adapters/db"
)

// CreateTestPrice implements the ports.TestPriceRepository interface.
// Update the signature and implementation as needed to match your interface definition.
func (r *Repository) CreateTestPrice(ctx context.Context, price db.Create_Test_PriceParams) ([]*db.DiagnosticCentreTestPrice, error) {
	// TODO: Implement the logic to create a test price in the database.
	return r.database.Create_Test_Price(ctx, price)
}
