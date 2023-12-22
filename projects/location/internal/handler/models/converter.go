package models

import (
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/handler/generated"
)

func ToDriverResponse(t domain.Driver) generated.Driver {
	id := t.Id.String()
	return generated.Driver{
		Id:  &id,
		Lat: t.Lat,
		Lng: t.Lng,
	}
}
