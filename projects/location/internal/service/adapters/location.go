package adapters

import (
	"context"
	"github.com/google/uuid"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/domain"
)

type LocationService interface {
	GetDrivers(c context.Context, lat float32, lng float32, radius float32) ([]domain.Driver, error)
	UpdateDriverLocation(c context.Context, driverId uuid.UUID, lat float32, lng float32) error
}
