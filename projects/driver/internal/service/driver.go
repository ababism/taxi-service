package service

import (
	"context"
	"github.com/google/uuid"

	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
)

var _ adapters.DriverService = userService{}

type userService struct {
	r adapters.UserRepository
}

func (s userService) Create(ctx context.Context, u domain.User) (domain.User, error) {
	if len(u.Name) > 30 || len(u.Name) < 1 {
		return domain.User{}, domain.ErrIncorrectBody
	}
	newUser, err := s.r.Create(ctx, u)
	if err != nil {
		return domain.User{}, domain.FilterErrors(err)
	}
	return newUser, nil
}

func (s userService) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.r.Get(ctx, id)
	if err != nil {
		return domain.User{}, domain.FilterErrors(err)
	}
	return user, nil
}

func (s userService) GetBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	user, err := s.r.Get(ctx, userID)
	if err != nil {
		return 0, domain.FilterErrors(err)
	}
	return user.Balance, nil
}

func (s userService) AddToBalance(ctx context.Context, userID uuid.UUID, amount float64) (float64, error) {
	newBalance, err := s.r.AddToBalance(ctx, userID, amount)
	if err != nil {
		return 0, domain.FilterErrors(err)
	}
	return newBalance, nil
}

func (s userService) TransferCurrency(ctx context.Context, senderID uuid.UUID, receiverID uuid.UUID, amount float64) (newSenderBalance float64, err error) {
	if amount < 0 {
		return 0, domain.ErrIncorrectBody
	}
	if senderID == receiverID {
		return 0, domain.ErrAccessDenied
	}
	newBalance, err := s.r.TransferCurrency(ctx, senderID, receiverID, amount)
	if err != nil {
		return 0, domain.FilterErrors(err)
	}
	return newBalance, nil
}

func NewDriverService(ur adapters.UserRepository) adapters.DriverService {
	return &userService{r: ur}
}
