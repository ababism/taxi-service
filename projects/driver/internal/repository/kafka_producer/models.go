package kafka_producer

import "time"

type CommandType string

type Data struct {
	TripID string  `json:"trip_id"`
	Reason *string `json:"reason,omitempty"`
}

type TripCommand struct {
	DriverId        string      `json:"id"`
	Source          string      `json:"source"`
	Type            CommandType `json:"type"`
	DataContentType string      `json:"datacontenttype"`
	Time            time.Time   `json:"time"`
	Data            Data        `json:"data"`
}
