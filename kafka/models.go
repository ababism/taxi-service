package main

import (
	"time"
)

const (
	Source                        = "/client"
	TripCommandAccept CommandType = "trip.event.accept"
	TripCommandCancel CommandType = "trip.event.cancel"
	TripCommandEnd    CommandType = "trip.event.end"
	TripCommandStart  CommandType = "trip.event.start"
	TripCommandCreate CommandType = "trip.event.create"
)

type CommandType string

type Data struct {
	TripID   string  `json:"trip_id"`
	DriverId string  `json:"driver_id,omitempty"`
	Reason   *string `json:"reason,omitempty"`
}

type TripCommand struct {
	RequestId       string      `json:"id"`
	Source          string      `json:"source"`
	Type            CommandType `json:"type"`
	DataContentType string      `json:"datacontenttype"`
	Time            time.Time   `json:"time"`
	Data            Data        `json:"data"`
}

type CreatedTripCommand struct {
	RequestId       string          `json:"id"`
	Source          string          `json:"source"`
	Type            string          `json:"type"`
	DataContentType string          `json:"datacontenttype"`
	Time            time.Time       `json:"time"`
	Data            CreatedTripData `json:"data"`
}

type CreatedTripData struct {
	TripID  string        `json:"trip_id"`
	OfferID string        `json:"offer_id"`
	Price   TripPrice     `json:"price"`
	From    LatLngLiteral `json:"from"`
	To      LatLngLiteral `json:"to"`
	Status  string        `json:"status"`
}

type LatLngLiteral struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type TripPrice struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func getCommandType(cmd int) CommandType {
	switch cmd {
	case 0:
		return TripCommandCreate
	case 1:
		return TripCommandAccept
	case 2:
		return TripCommandCancel
	case 3:
		return TripCommandEnd
	}
	return TripCommandStart
}
