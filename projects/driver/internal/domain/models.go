package domain

import (
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type User struct {
	Id      uuid.UUID
	Name    string
	Balance float64
}

// LatLngLiteral An object describing a specific location with Latitude and Longitude in decimal degrees.
type LatLngLiteral struct {
	// Lat Latitude in decimal degrees
	Lat float32

	// Lng Longitude in decimal degrees
	Lng float32
}

// Money defines model for Money.
type Money struct {
	// Amount expressed as a decimal number of major currency units
	Amount float64

	// Currency 3 letter currency code as defined by ISO-4217
	Currency string
}

// Trip defines model for Trip.
type Trip struct {
	DriverId *string

	// From An object describing a specific location with Latitude and Longitude in decimal degrees.
	From   *LatLngLiteral
	Id     *openapi_types.UUID
	Price  *Money
	Status *TripStatus

	// To An object describing a specific location with Latitude and Longitude in decimal degrees.
	To *LatLngLiteral
}

// TripStatus defines model for Trip.Status.
type TripStatus string
