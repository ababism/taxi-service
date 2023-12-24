package adapters

import (
	"context"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

type KafkaClient interface {
	SendTripStatusUpdate(ctx context.Context, trip domain.Trip, commandType domain.CommandType, reason *string) error
}
