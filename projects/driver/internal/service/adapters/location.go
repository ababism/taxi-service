package adapters

import (
	"context"

	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

type LocationClient interface {
	GetDrivers(ctx context.Context, centerLocation domain.LatLngLiteral, radius float32) ([]domain.DriverLocation, error)
}
