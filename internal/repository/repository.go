package repository

import (
	"context"

	"github.com/IgorRybak2055/bamboo/internal/models"
)

// Profile represents possible database actions with profile.
type Profile interface {
	OrderProfile(ctx context.Context, profile models.Profile) error
}
