package adapters

import (
	"context"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/domain"
)

type LocationRepository interface {
	GetDrivers(c context.Context, lat float32, lng float32, radius float32) ([]domain.Driver, error)
	UpdateDriverLocation(c context.Context, driverId openapitypes.UUID, lat float32, lng float32) error
}
