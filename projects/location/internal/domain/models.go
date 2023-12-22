package domain

import "github.com/google/uuid"

// Driver defines model for Driver.
type Driver struct {
	// Id Идентификатор водителя
	Id *uuid.UUID

	// Lat Latitude in decimal degrees
	Lat float32

	// Lng Longitude in decimal degrees
	Lng float32
}

// LatLngLiteral An object describing a specific location with Latitude and Longitude in decimal degrees.
type LatLngLiteral struct {
	// Lat Latitude in decimal degrees
	Lat float32

	// Lng Longitude in decimal degrees
	Lng float32
}

// GetDriversParams defines parameters for GetDrivers.
type GetDriversParams struct {
	// Lat Latitude in decimal degrees
	Lat float32 `form:"lat"`

	// Lng Longitude in decimal degrees
	Lng float32 `form:"lng"`

	// Radius of serach
	Radius float32 `form:"radius"`
}

// UpdateDriverLocationJSONRequestBody defines body for UpdateDriverLocation for application/json ContentType.
type UpdateDriverLocationJSONRequestBody = LatLngLiteral
