package domain

import (
	"github.com/google/uuid"
)

type User struct {
	Id      uuid.UUID
	Name    string
	Balance float64
}
