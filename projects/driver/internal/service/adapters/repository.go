package adapters

import (
	"context"
	"github.com/google/uuid"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

type DriverRepository interface {
	GetTripByID(ctx context.Context, tripId uuid.UUID) (domain.Trip, error)
	GetTrips(ctx context.Context, driverId uuid.UUID) ([]domain.Trip, error)
	// Only to change to DRIVER_FOUND
	ChangeTripStatus(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID, status domain.TripStatus) error
}
