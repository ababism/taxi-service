package models

import (
	"github.com/google/uuid"

	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

type User struct {
	Id       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Balance  float64   `db:"balance"`
}

func (u User) ToDomain() domain.User {
	return domain.User{
		Id:      u.Id,
		Name:    u.Username,
		Balance: u.Balance,
	}
}
