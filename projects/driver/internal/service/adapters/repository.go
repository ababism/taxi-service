package adapters

import (
	"context"
	"github.com/google/uuid"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

type DriverRepository interface {
	GetTripByID(ctx context.Context, tripId uuid.UUID) (domain.Trip, error)
	ChangeTripStatus(ctx context.Context, tripId uuid.UUID, status domain.TripStatus) error
	CreateTrip(ctx context.Context, status domain.Trip) error
}
