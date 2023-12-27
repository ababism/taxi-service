package adapters

import (
	"context"
	"github.com/google/uuid"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

type DriverService interface {
	GetTripsByStatus(ctx context.Context, status domain.TripStatus) ([]domain.Trip, error)
	GetTripByID(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) (*domain.Trip, error)
	InsertTrip(ctx context.Context, trip domain.Trip) error
	UpdateTrip(ctx context.Context, tripId uuid.UUID, updatedTrip domain.Trip) error
	UpdateTripStatus(ctx context.Context, tripId uuid.UUID, status domain.TripStatus) error
	UpdateTripStatusAndDriver(ctx context.Context, tripId uuid.UUID, driverId string, status domain.TripStatus) error
	AcceptTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) error
	CancelTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID, reason *string) error
	StartTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) error
	EndTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) error
	GetTrips(ctx context.Context, driverId uuid.UUID) ([]domain.Trip, error)
	GetDrivers(ctx context.Context, driverLocation domain.LatLngLiteral, radius float32) ([]domain.DriverLocation, error)
}
