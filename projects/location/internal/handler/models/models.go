package models

import "gitlab/ArtemFed/mts-final-taxi/projects/location/internal/handler/generated"

type Error struct {
	Message string `json:"message"`
}

type GetDriversResponse struct {
	drivers []generated.Driver
}
