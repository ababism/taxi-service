package adapters

import (
	"context"
	"github.com/google/uuid"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

type DriverRepository interface {
	GetTripsByStatus(ctx context.Context, status domain.TripStatus) ([]domain.Trip, error)
	GetTripByID(ctx context.Context, tripId uuid.UUID) (*domain.Trip, error)
	UpdateTrip(ctx context.Context, tripId uuid.UUID, updatedTrip domain.Trip) error
	ChangeTripStatus(ctx context.Context, tripId uuid.UUID, status domain.TripStatus) error
	ChangeTripStatusAndDriver(ctx context.Context, tripId uuid.UUID, driverId string, status domain.TripStatus) error
	InsertTrip(ctx context.Context, status domain.Trip) error
}
