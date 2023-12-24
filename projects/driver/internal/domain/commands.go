package domain

const (
	Source                        = "/driver"
	TripCommandAccept CommandType = "trip.command.accept"
	TripCommandCancel CommandType = "trip.command.cancel"
	TripCommandEnd    CommandType = "trip.command.end"
	TripCommandStart  CommandType = "trip.command.start"
)

type CommandType string
