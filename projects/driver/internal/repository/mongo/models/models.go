package models

import (
	"errors"
	"github.com/google/uuid"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
)

const (
	DriverCollectionName = "driver"
)

// MongoLatLngLiteral defines the MongoDB representation of LatLngLiteral.
type MongoLatLngLiteral struct {
	Lat float32 `bson:"lat"`
	Lng float32 `bson:"lng"`
}

func (lll *MongoLatLngLiteral) ToDomain() *domain.LatLngLiteral {
	return &domain.LatLngLiteral{
		Lat: lll.Lat,
		Lng: lll.Lng,
	}
}

// MongoMoney defines the MongoDB representation of Money.
type MongoMoney struct {
	Amount   float64 `bson:"amount"`
	Currency string  `bson:"currency"`
}

func (price *MongoMoney) ToDomain() *domain.Money {
	return &domain.Money{
		Amount:   price.Amount,
		Currency: price.Currency,
	}
}

// MongoTrip defines the MongoDB representation of Trip.
type MongoTrip struct {
	DriverId *string             `bson:"driver_id,omitempty"`
	From     *MongoLatLngLiteral `bson:"from"`
	Id       string              `bson:"trip_id"`
	Price    *MongoMoney         `bson:"price,omitempty"`
	Status   MongoTripStatus     `bson:"status"`
	To       *MongoLatLngLiteral `bson:"to"`
}

// MongoDriverLocation defines the MongoDB representation of DriverLocation.
type MongoDriverLocation struct {
	DriverId    uuid.UUID          `bson:"driver_id"`
	Coordinates MongoLatLngLiteral `bson:"coordinates"`
}

// TripStatus defines model for Trip.Status.
type MongoTripStatus string

func ToMongoStatusModel(trip domain.TripStatus) MongoTripStatus {
	return MongoTripStatus(trip)
}

func ToDomainTripModel(trip MongoTrip) (*domain.Trip, error) {
	//if trip.Id == nil {
	//	return nil, errors.New("trip id is nil")
	//}
	tripId, err := uuid.Parse(trip.Id)
	if err != nil {
		return nil, errors.New("can't parse trip id to uuid")
	}
	status := domain.TripStatus(trip.Status)
	t := &domain.Trip{
		DriverId: trip.DriverId,
		From:     trip.From.ToDomain(),
		Id:       tripId,
		Price:    trip.Price.ToDomain(),
		Status:   &status,
		To:       trip.To.ToDomain(),
	}
	return t, nil
}
func ToMongoTripModel(trip domain.Trip) MongoTrip {
	// Convert the domain.Trip to its MongoDB representation
	tripId := trip.Id.String()
	mongoTrip := MongoTrip{
		DriverId: trip.DriverId,
		From: &MongoLatLngLiteral{
			Lat: trip.From.Lat,
			Lng: trip.From.Lng,
		},
		Id: tripId,
		Price: &MongoMoney{
			Amount:   trip.Price.Amount,
			Currency: trip.Price.Currency,
		},
		Status: ToMongoStatusModel(*trip.Status),
		To: &MongoLatLngLiteral{
			Lat: trip.To.Lat,
			Lng: trip.To.Lng,
		},
	}
	return mongoTrip
}
