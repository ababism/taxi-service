package domain

import (
	"github.com/google/uuid"
)

const ServiceName = "mts-final-taxi/driver"

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
	Id     *uuid.UUID
	Price  *Money
	Status *TripStatus

	// To An object describing a specific location with Latitude and Longitude in decimal degrees.
	To *LatLngLiteral
}

type DriverLocation struct {
	DriverId    uuid.UUID
	Coordinates LatLngLiteral
}

// TripStatus defines model for Trip.Status.
type TripStatus string

type TripStatusCollection struct {
	canceled     TripStatus
	driverFound  TripStatus
	driverSearch TripStatus
	ended        TripStatus
	onPosition   TripStatus
	started      TripStatus
}

func NewTripStatusCollection(
	canceled TripStatus,
	driverFound TripStatus,
	driverSearch TripStatus,
	ended TripStatus,
	onPosition TripStatus,
	started TripStatus) TripStatusCollection {

	return TripStatusCollection{
		canceled:     canceled,
		driverFound:  driverFound,
		driverSearch: driverSearch,
		ended:        ended,
		onPosition:   onPosition,
		started:      started}
}

func (t *TripStatusCollection) GetCanceled() TripStatus {
	return t.canceled
}

func (t *TripStatusCollection) GetDriverFound() TripStatus {
	return t.driverFound
}

func (t *TripStatusCollection) GetDriverSearch() TripStatus {
	return t.driverSearch
}

func (t *TripStatusCollection) GetEnded() TripStatus {
	return t.ended
}

func (t *TripStatusCollection) GetOnPosition() TripStatus {
	return t.onPosition
}

func (t *TripStatusCollection) GetStarted() TripStatus {
	return t.started
}
