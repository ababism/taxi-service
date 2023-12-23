package adapters

import (
	"context"
	"github.com/google/uuid"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

type DriverService interface {
	GetTripByID(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) (*domain.Trip, error)
	AcceptTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) error
	CancelTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID, reason *string) error
	StartTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) error
	EndTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) error
	GetTrips(ctx context.Context, driverId uuid.UUID) ([]domain.Trip, error)
}
