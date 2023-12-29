package models

import "gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/handler/generated"

type Error struct {
	Message string `json:"message"`
}

type GetDriversResponse struct {
	Drivers []generated.Driver `json:"drivers"`
}
