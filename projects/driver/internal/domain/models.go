package domain

import (
	"github.com/google/uuid"
)

const (
	// TODO testing
	SearchRadius = 113_000
	ServiceName  = "mts-final-taxi/driver"
)

// LatLngLiteral An object describing a specific location with Latitude and Longitude in decimal degrees.
type LatLngLiteral struct {
	Lat float32
	Lng float32
}

type Money struct {
	// Amount expressed as a decimal number of major currency units
	Amount float64
	// Currency 3 letter currency code as defined by ISO-4217
	Currency string
}

// Trip defines model for Trip.
type Trip struct {
	DriverId *string
	From     *LatLngLiteral
	Id       uuid.UUID
	Price    *Money
	Status   *TripStatus
	To       *LatLngLiteral
}

type DriverLocation struct {
	DriverId    uuid.UUID
	Coordinates LatLngLiteral
}
