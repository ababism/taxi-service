package locationclient

import (
	"errors"
	"github.com/google/uuid"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/repository/locationclient/generated"
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
	driversDomain := make([]domain.DriverLocation, len(dLocations))

	for i, driverLocation := range dLocations {
		dlDom, err := ToDriverLocationDomain(driverLocation)
		if err != nil {
			return nil, err
		}
		driversDomain[i] = *dlDom
	}

	return driversDomain, nil
}
