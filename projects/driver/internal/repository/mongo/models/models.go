package models

import (
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

// MongoMoney defines the MongoDB representation of Money.
type MongoMoney struct {
	Amount   float64 `bson:"amount"`
	Currency string  `bson:"currency"`
}

// MongoTrip defines the MongoDB representation of Trip.
type MongoTrip struct {
	DriverId *string             `bson:"driver_id,omitempty"`
	From     *MongoLatLngLiteral `bson:"from"`
	Id       *uuid.UUID          `bson:"trip_id"`
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

func ToMongoTripModel(trip domain.Trip) MongoTrip {
	// Convert the domain.Trip to its MongoDB representation
	mongoTrip := MongoTrip{
		DriverId: trip.DriverId,
		From: &MongoLatLngLiteral{
			Lat: trip.From.Lat,
			Lng: trip.From.Lng,
		},
		Id: trip.Id,
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
