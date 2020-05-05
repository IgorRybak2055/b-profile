package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/IgorRybak2055/bamboo/internal/models"
	"github.com/IgorRybak2055/bamboo/internal/storage"
)

type profile struct {
	dbc *sqlx.DB
}

// NewAccountRepository will create an object to access profile repository that represent the Profile interface
func NewAccountRepository(dbc *sqlx.DB) Profile {
	return &profile{dbc}
}

// OrderProfile added matchingID to repository and change status to 'Confirmed'.
func (p profile) OrderProfile(ctx context.Context, profile models.Profile) error {
	var query = `
UPDATE
	profiles
SET
	matchingid = $1,
    status = 'Confirmed'
WHERE 
	iccid = $2 AND status = 'Available'
`

	rows, err := p.dbc.ExecContext(ctx, query, profile.MatchingID, profile.ICCID)
	if err != nil {
		return storage.ErrDBInteracting
	}

	n, err := rows.RowsAffected()
	if err != nil {
		return storage.ErrDBInteracting
	}

	if n != 1 {
		return storage.ErrNotAvailableICCID
	}

	return nil
}
