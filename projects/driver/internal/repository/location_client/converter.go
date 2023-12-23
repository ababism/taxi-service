package location_client

import (
	"errors"
	"github.com/google/uuid"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/location_client/generated"
)

func ToDriverLocationDomain(r generated.Driver) (*domain.DriverLocation, error) {
	if r.Id == nil {
		return nil, errors.New("driver id is nil")
	}
	driverId, err := uuid.Parse(*r.Id)
	if err != nil {
		return nil, err
	}
	res := domain.DriverLocation{
		DriverId: driverId,
		Coordinates: domain.LatLngLiteral{
			Lat: r.Lat,
			Lng: r.Lng,
		},
	}
	return &res, nil
}

func ToDriverLocationsDomain(dLocations []generated.Driver) ([]domain.DriverLocation, error) {
	tripsResponse := make([]domain.DriverLocation, len(dLocations))

	for _, driverLocation := range dLocations {
		l, err := ToDriverLocationDomain(driverLocation)
		if err != nil {
			return nil, err
		}
		tripsResponse = append(tripsResponse, *l)
	}

	return tripsResponse, nil
}
