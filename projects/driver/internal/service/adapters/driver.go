package adapters

import (
	"context"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

type DriverService interface {
	GetTripByID(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) (domain.Trip, error)
	AcceptTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) (domain.Trip, error)
	CancelTrip(ctx context.Context, diverId uuid.UUID, tripId openapi_types.UUID, reason *string) (domain.Trip, error)
	StartTrip(ctx context.Context, diverId openapi_types.UUID, tripId openapi_types.UUID) (domain.Trip, error)
	EndTrip(ctx context.Context, diverId openapi_types.UUID, tripId openapi_types.UUID) (domain.Trip, error)
	GetTrips(ctx context.Context, diverId openapi_types.UUID) ([]domain.Trip, error)
}
