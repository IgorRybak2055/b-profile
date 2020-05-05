package services

import (
	"context"
)

// Profile represents possible actions with profile.
type Profile interface {
	OrderProfile(ctx context.Context, iccids []string) ([][]string, error)
}
