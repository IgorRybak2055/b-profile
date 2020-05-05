package services

import (
	"context"
)

type Profile interface {
	OrderProfile(ctx context.Context, iccids []string) ([][]string, error)
}
