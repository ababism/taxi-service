package models

import (
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/domain"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/handler/generated"
)

func ToDriverResponse(t domain.Driver) generated.Driver {
	id := t.Id.String()
	return generated.Driver{
		Id:  &id,
		Lat: t.Lat,
		Lng: t.Lng,
	}
}
