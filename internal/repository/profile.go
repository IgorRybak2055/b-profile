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

// NewProfileRepository will create an object that represent the Profile interface
func NewAccountRepository(dbc *sqlx.DB) Profile {
	return &profile{dbc}
}

// CreateAccount use account data for registration new account in database.
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
