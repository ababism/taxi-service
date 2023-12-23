package domain

import (
	"errors"
)

var (
	ErrTokenInvalid  = errors.New("token invalid")
	ErrIncorrectBody = errors.New("incorrect json body")
	ErrInternal      = errors.New("server internal error")
	ErrNotFound      = errors.New("not found")
	ErrAccessDenied  = errors.New("access denied")
	ErrAlreadyExists = errors.New("element already exists")
	ErrBadUUID       = errors.New("could not get correct uuid")
	ErrBadLatitude   = errors.New("latitude outside the boundaries [90, 90]")
	ErrBadLongitude  = errors.New("longitude outside the boundaries [180, 180]")
)

func FilterErrors(err error) error {
	switch {
	case errors.Is(err, ErrTokenInvalid),
		errors.Is(err, ErrIncorrectBody),
		errors.Is(err, ErrNotFound),
		errors.Is(err, ErrAccessDenied),
		errors.Is(err, ErrAlreadyExists),
		errors.Is(err, ErrBadLatitude),
		errors.Is(err, ErrBadLongitude),
		errors.Is(err, ErrBadUUID):
		return err
	default:
		return ErrInternal
	}
}
