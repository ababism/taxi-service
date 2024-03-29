package kafkaconsumer

import "time"

//
//type Event struct {
//	RequestId              string    `json:"id"`
//	Source          string    `json:"source"`
//	Type            string    `json:"type"`
//	DataContentType string    `json:"datacontenttype"`
//	Time            time.Time `json:"time"`
//	Data            []byte    `json:"data"`
//}

type StatusUpdateEvent struct {
	RequestId       string    `json:"id"`
	Source          string    `json:"source"`
	Type            string    `json:"type"`
	DataContentType string    `json:"datacontenttype"`
	Time            time.Time `json:"time"`
	Data            TripData  `json:"data"`
}

type TripData struct {
	DriverId string `json:"driver_id,omitempty"`
	TripID   string `json:"trip_id"`
}

type CreatedTripEvent struct {
	ID              string          `json:"id"`
	Source          string          `json:"source"`
	Type            string          `json:"type"`
	DataContentType string          `json:"datacontenttype"`
	Time            time.Time       `json:"time"`
	Data            CreatedTripData `json:"data"`
}

type CreatedTripData struct {
	TripID   string        `json:"trip_id"`
	OfferID  string        `json:"offer_id"`
	Price    TripPrice     `json:"price"`
	From     LatLngLiteral `json:"from"`
	To       LatLngLiteral `json:"to"`
	Status   string        `json:"status"`
	DriverId string        `json:"driver_id,omitempty"`
}

// Тут было float64 по схеме kafka, но мне надо float32 по контракту location (кривая стркутура?)
type LatLngLiteral struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type TripPrice struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
