package domain

var TripStatuses TripStatusCollection

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

func InitTripStatusCollection(
	canceled TripStatus,
	driverFound TripStatus,
	driverSearch TripStatus,
	ended TripStatus,
	onPosition TripStatus,
	started TripStatus,
) {
	TripStatuses = TripStatusCollection{
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
