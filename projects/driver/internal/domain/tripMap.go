package domain

import (
	"github.com/google/uuid"
	"sync"
)

var IncomingTrips *TripMap

// TripMap Sync
type TripMap struct {
	trips map[*string]chan *uuid.UUID
	mu    sync.RWMutex
}

func InitTripMap() {
	IncomingTrips = &TripMap{
		trips: make(map[*string]chan *uuid.UUID),
	}
}

func init() {
	IncomingTrips = &TripMap{
		trips: make(map[*string]chan *uuid.UUID),
	}
}

func (tm *TripMap) AddTrip(driverID *string, tripIDChan chan *uuid.UUID) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.trips[driverID] = tripIDChan
}

func (tm *TripMap) GetTripChannel(driverID *string) (chan *uuid.UUID, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	channel, exists := tm.trips[driverID]
	return channel, exists
}

func (tm *TripMap) DeleteTripChannel(driverID *string) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	delete(tm.trips, driverID)
}
