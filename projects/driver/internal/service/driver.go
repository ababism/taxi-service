package service

import (
	"context"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
)

var _ adapters.DriverService = driverService{}

type driverService struct {
	r adapters.UserRepository
}

func (d driverService) GetTripByID(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) (domain.Trip, error) {
	//TODO implement me
	panic("implement me")
}

func (d driverService) AcceptTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) (domain.Trip, error) {
	//TODO implement me
	panic("implement me")
}

func (d driverService) CancelTrip(ctx context.Context, diverId uuid.UUID, tripId openapi_types.UUID, reason *string) (domain.Trip, error) {
	//TODO implement me
	panic("implement me")
}

func (d driverService) StartTrip(ctx context.Context, diverId openapi_types.UUID, tripId openapi_types.UUID) (domain.Trip, error) {
	//TODO implement me
	panic("implement me")
}

func (d driverService) EndTrip(ctx context.Context, diverId openapi_types.UUID, tripId openapi_types.UUID) (domain.Trip, error) {
	//TODO implement me
	panic("implement me")
}

func (d driverService) GetTrips(ctx context.Context, diverId openapi_types.UUID) ([]domain.Trip, error) {
	//TODO implement me
	panic("implement me")
}

//func (s driverService) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
//	user, err := s.r.Get(ctx, id)
//	if err != nil {
//		return domain.User{}, domain.FilterErrors(err)
//	}
//	return user, nil
//}
//func (s driverService) Create(ctx context.Context, u domain.User) (domain.User, error) {
//	if len(u.Name) > 30 || len(u.Name) < 1 {
//		return domain.User{}, domain.ErrIncorrectBody
//	}
//	newUser, err := s.r.Create(ctx, u)
//	if err != nil {
//		return domain.User{}, domain.FilterErrors(err)
//	}
//	return newUser, nil
//}
//
//func (s driverService) GetBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
//	user, err := s.r.Get(ctx, userID)
//	if err != nil {
//		return 0, domain.FilterErrors(err)
//	}
//	return user.Balance, nil
//}
//
//func (s driverService) AddToBalance(ctx context.Context, userID uuid.UUID, amount float64) (float64, error) {
//	newBalance, err := s.r.AddToBalance(ctx, userID, amount)
//	if err != nil {
//		return 0, domain.FilterErrors(err)
//	}
//	return newBalance, nil
//}
//
//func (s driverService) TransferCurrency(ctx context.Context, senderID uuid.UUID, receiverID uuid.UUID, amount float64) (newSenderBalance float64, err error) {
//	if amount < 0 {
//		return 0, domain.ErrIncorrectBody
//	}
//	if senderID == receiverID {
//		return 0, domain.ErrAccessDenied
//	}
//	newBalance, err := s.r.TransferCurrency(ctx, senderID, receiverID, amount)
//	if err != nil {
//		return 0, domain.FilterErrors(err)
//	}
//	return newBalance, nil
//}

func NewDriverService(ur adapters.UserRepository) adapters.DriverService {
	return &driverService{r: ur}
}
