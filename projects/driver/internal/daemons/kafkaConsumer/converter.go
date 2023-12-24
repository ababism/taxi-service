package kafkaconsumer

import (
	"github.com/google/uuid"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

func ParseUUID(strUUID string) (uuid.UUID, error) {
	return uuid.Parse(strUUID)
}

func ToDomainTrip(trip CreatedTripData) (domain.Trip, error) {
	tripId, err := ParseUUID(trip.TripID)
	if err != nil {
		return domain.Trip{}, err
	}
	tripStatus := domain.TripStatus(trip.Status)

	return domain.Trip{
		DriverId: nil,
		From: &domain.LatLngLiteral{
			Lat: trip.From.Lat,
			Lng: trip.From.Lng,
		},
		Id: &tripId,
		Price: &domain.Money{
			Amount:   trip.Price.Amount,
			Currency: trip.Price.Currency,
		},
		Status: &tripStatus,
		To: &domain.LatLngLiteral{
			Lat: trip.To.Lat,
			Lng: trip.To.Lng,
		},
	}, nil
}
