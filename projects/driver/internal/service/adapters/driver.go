package adapters

import (
	"context"
	"github.com/google/uuid"

	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

//type DriverService interface {
//	Create(ctx context.Context, user domain.User) (domain.User, error)
//	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
//	GetBalance(ctx context.Context, userID uuid.UUID) (float64, error)
//	AddToBalance(ctx context.Context, userID uuid.UUID, amount float64) (float64, error)
//	TransferCurrency(ctx context.Context, senderID uuid.UUID, receiverID uuid.UUID, amount float64) (newSenderBalance float64, err error)
//}

type DriverService interface {
	Get(ctx context.Context, userId uuid.UUID, tripId uuid.UUID) (domain.Trip, error)
}
