package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/juju/zaputil/zapctx"
	"github.com/pkg/errors"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var _ adapters.DriverService = driverService{}

type driverService struct {
	r              adapters.DriverRepository
	locationClient adapters.LocationClient
	tripStatuses   domain.TripStatusCollection
}

// GetTripsByStatus returns all trips with the given status
func (s driverService) GetTripsByStatus(ctx context.Context, status domain.TripStatus) ([]domain.Trip, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: GetTripsByStatus")
	defer span.End()

	trips, err := s.r.GetTripsByStatus(newCtx, status)
	if err != nil {
		logger.Error("driver-service: get trips by status from repository", zap.Error(err))
		return nil, err
	}

	return trips, nil
}

func (s driverService) GetTripByID(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) (*domain.Trip, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: GetTripByID")
	defer span.End()

	// err if trip driver != nil and driver != driverId
	trip, err := s.r.GetTripByID(newCtx, tripId)
	if err != nil {
		logger.Error("driver-service: get trip from repository")
		return nil, domain.FilterErrors(err)
	}
	if trip.DriverId != nil && *trip.DriverId != driverId.String() {
		return nil, errors.Wrap(domain.ErrAccessDenied, "trip driver id does not match passed id")
	}
	return trip, err
}

// UpdateTrip updates all fields of the given trip in the service layer
func (s driverService) UpdateTrip(ctx context.Context, tripId uuid.UUID, updatedTrip domain.Trip) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: UpdateTrip")
	defer span.End()

	// Call the repository method to update the trip
	err := s.r.UpdateTrip(newCtx, tripId, updatedTrip)
	if err != nil {
		logger.Error("driver-service: update trip in repository", zap.Error(err))
		return domain.FilterErrors(err)
	}

	return nil
}

func (s driverService) AcceptTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: AcceptTrip")
	defer span.End()

	ctx = zapctx.WithLogger(newCtx, logger)

	trip, err := s.r.GetTripByID(newCtx, tripId)
	if err != nil {
		logger.Error("can't get trip from repository")
		return domain.FilterErrors(err)
	}

	if trip.Status == nil || *trip.Status != s.tripStatuses.GetDriverSearch() {
		return errors.Wrap(domain.ErrAccessDenied, "trip doesn't need driver")
	}
	// TODO Kafka <- acceptTrip

	return nil
}

func (s driverService) CancelTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID, reason *string) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: CancelTrip")
	defer span.End()

	ctx = zapctx.WithLogger(newCtx, logger)

	trip, err := s.r.GetTripByID(newCtx, tripId)
	if err != nil {
		logger.Error("can't get trip from repository")
		return domain.FilterErrors(err)
	}
	if trip.Status == nil || *trip.Status == s.tripStatuses.GetDriverSearch() {
		return errors.Wrap(domain.ErrAccessDenied, "trip hasn't connected with driver yet")
	}
	if trip.DriverId != nil && *trip.DriverId != driverId.String() {
		return errors.Wrap(domain.ErrAccessDenied, "trip driver id does not match passed id")
	}
	// TODO Kafka <- cancel (with Reason)

	return nil
}

func (s driverService) StartTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: StartTrip")
	defer span.End()

	ctx = zapctx.WithLogger(newCtx, logger)

	trip, err := s.r.GetTripByID(newCtx, tripId)
	if err != nil {
		logger.Error("can't get trip from repository")
		return domain.FilterErrors(err)
	}
	// TODO ask about *trip.Status != s.tripStatuses.GetOnPosition()
	if trip.Status == nil || *trip.Status != s.tripStatuses.GetOnPosition() || *trip.Status != s.tripStatuses.GetDriverFound() {
		return errors.Wrap(domain.ErrAccessDenied, "trip hasn't connected with driver yet")
	}
	if trip.DriverId != nil && *trip.DriverId != driverId.String() {
		return errors.Wrap(domain.ErrAccessDenied, "trip driver id does not match passed id")
	}

	// TODO Kafka <- start trip

	return nil
}

func (s driverService) EndTrip(ctx context.Context, driverId uuid.UUID, tripId uuid.UUID) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: EndTrip")
	defer span.End()

	ctx = zapctx.WithLogger(newCtx, logger)

	trip, err := s.r.GetTripByID(newCtx, tripId)
	if err != nil {
		logger.Error("can't get trip from repository")
		return domain.FilterErrors(err)
	}
	// TODO REDO ask about *trip.Status != s.tripStatuses.GetOnPosition()
	if trip.Status == nil || *trip.Status != s.tripStatuses.GetStarted() {
		return errors.Wrap(domain.ErrAccessDenied, "trip hasn't connected with driver yet")
	}
	if trip.DriverId != nil && *trip.DriverId != driverId.String() {
		return errors.Wrap(domain.ErrAccessDenied, "trip driver id does not match passed id")
	}
	// TODO Kafka <- End trip

	return nil
}

// Long poll
func (s driverService) GetTrips(ctx context.Context, driverId uuid.UUID) ([]domain.Trip, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	newCtx, span := tr.Start(ctx, "driver.service: GetTrips")
	defer span.End()

	ctx = zapctx.WithLogger(newCtx, logger)

	trips := make([]domain.Trip, 0)
	// fot loop for 5-10 seconds with timeout:
	//
	// 		span.AddEvent("Message from Kafka")
	//		TripCreated <- Kafka
	//
	//		drivers := s.locationClient.GetDrivers(TripCreated.userLocation, radius)
	//
	//      for range drivers {
	//			if dr.id == id {
	//					trips = trips.append(TripCreated)
	//			{
	//      {
	//

	// TODO DELETE?
	//trips, err := s.locationClient.GetDrivers(newCtx, driverId)
	//if err != nil {
	//	logger.Error("driver-service: EndTrip")
	//	return nil, domain.FilterErrors(err)
	//}

	return trips, nil
}

func NewDriverService(driverRepository adapters.DriverRepository, locationClient adapters.LocationClient, tsc domain.TripStatusCollection) adapters.DriverService {
	return &driverService{
		r:              driverRepository,
		locationClient: locationClient,
		tripStatuses:   tsc,
	}
}
