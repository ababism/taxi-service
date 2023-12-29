package domain

import (
	"github.com/google/uuid"
	"sync"
)

var AvailableTripEvents *TripEventMap

// TripEventMap is safe for concurrent use by multiple goroutines
type TripEventMap struct {
	trips map[string]chan uuid.UUID
	mu    sync.RWMutex
}

func InitTripMap() {
	AvailableTripEvents = &TripEventMap{
		trips: make(map[string]chan uuid.UUID),
	}
}

func init() {
	AvailableTripEvents = &TripEventMap{
		trips: make(map[string]chan uuid.UUID),
	}
}

func (tm *TripEventMap) AddTripChannel(driverID string, tripEvent chan uuid.UUID) func() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.trips[driverID] = tripEvent

	return func() {
		close(tripEvent)
	}
}

func (tm *TripEventMap) SendTrip(driverID string, tripId uuid.UUID) bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	channel, ok := tm.trips[driverID]
	if !ok {
		return false
	}

	select {
	case channel <- tripId:
		return true
	default:
		return false
	}
}

func (tm *TripEventMap) GetTripChannel(driverID string) (chan uuid.UUID, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	channel, exists := tm.trips[driverID]
	return channel, exists
}

func (tm *TripEventMap) DeleteTripChannel(driverID string) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	channel, exists := tm.trips[driverID]
	if exists {
		close(channel)
		delete(tm.trips, driverID)
	}
}
