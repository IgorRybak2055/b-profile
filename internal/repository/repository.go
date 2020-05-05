package repository

import (
	"context"

	"github.com/IgorRybak2055/bamboo/internal/models"
)

type Profile interface {
	OrderProfile(ctx context.Context, profile models.Profile) error
}
